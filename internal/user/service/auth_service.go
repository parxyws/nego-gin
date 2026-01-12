package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/oklog/ulid/v2"
	"github.com/parxyws/nego-gin/config"
	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"github.com/parxyws/nego-gin/internal/user/domain/dto"
	"github.com/parxyws/nego-gin/middleware"
	"github.com/parxyws/nego-gin/pkg/util"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type AuthServiceImpl struct {
	cfg                *config.Config
	userRepository     user.UserRepository
	userRoleRepository user.UserRoleRepository
	roleRepository     user.RoleRepository
	rdb                *redis.Client
	jwtMiddleware      *jwt.GinJWTMiddleware
	mailDialer         *gomail.Dialer
}

func NewAuthService(cfg *config.Config, userRepository user.UserRepository, userRoleRepository user.UserRoleRepository, roleRepository user.RoleRepository, rdb *redis.Client, mailDialer *gomail.Dialer) user.AuthService {
	return &AuthServiceImpl{cfg: cfg, userRepository: userRepository, userRoleRepository: userRoleRepository, roleRepository: roleRepository, rdb: rdb, mailDialer: mailDialer}
}

func (a *AuthServiceImpl) RefreshToken(ctx context.Context, entity *dto.JwtToken, payload *middleware.JwtPayload) (*dto.UserResponse, error) {
	request := &domain.User{
		UserID: payload.ID,
	}
	result, err := a.userRepository.ReadByIdMinimal(ctx, request)
	token, err := a.jwtMiddleware.TokenGeneratorWithRevocation(ctx, payload, entity.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		UserID:    result.UserID,
		Username:  result.Username,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Email:     result.Email,
		Jwt: dto.JwtToken{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		},
	}, nil
}

func (a *AuthServiceImpl) ForgotPassword(ctx context.Context, request *dto.ForgotPasswordRequest) error {
	userRequest := &domain.User{Email: request.Email}
	foundUser, err := a.userRepository.ReadByEmail(ctx, userRequest)
	if err != nil {
		return util.ErrNotFound
	}

	otp, err := util.GenerateRandomInteger(6)
	if err != nil {
		return fmt.Errorf("AuthService.ForgotPassword - %w", err)
	}

	if err := a.rdb.Set(ctx, "forgot-otp:"+request.Email, otp, 15*time.Minute).Err(); err != nil {
		return fmt.Errorf("AuthService.ForgotPassword - %w", err)
	}

	mailMessage, err := util.GenerateOTPMailMessage(a.cfg, foundUser, otp)
	if err != nil {
		return fmt.Errorf("AuthService.ForgotPassword - %w", err)
	}

	// Async email sending
	go func() {
		_ = a.mailDialer.DialAndSend(mailMessage)
	}()

	return nil
}

func (a *AuthServiceImpl) ResetPassword(ctx context.Context, request *dto.ResetPasswordRequest) error {
	otpValue, err := a.rdb.Get(ctx, "forgot-otp:"+request.Email).Result()
	if errors.Is(err, redis.Nil) {
		return util.NewAppError(http.StatusBadRequest, "OTP expired or invalid", nil)
	}
	if err != nil {
		return fmt.Errorf("AuthService.ResetPassword - %w", err)
	}

	if otpValue != request.OTP {
		return util.NewAppError(http.StatusBadRequest, "Invalid OTP", nil)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("AuthService.ResetPassword - %w", err)
	}

	foundUser, err := a.userRepository.ReadByEmail(ctx, &domain.User{Email: request.Email})
	if err != nil {
		return util.ErrNotFound
	}

	_, err = a.userRepository.UpdateUser(ctx, &domain.User{
		UserID:       foundUser.UserID,
		PasswordHash: string(passwordHash),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		return fmt.Errorf("AuthService.ResetPassword - %w", err)
	}

	_ = a.rdb.Del(ctx, "forgot-otp:"+request.Email).Err()

	return nil
}

func (a *AuthServiceImpl) ResendVerification(ctx context.Context, entity *dto.UserRegisterResponse) error {
	var userData *domain.User
	token := strings.TrimSpace(entity.ReferenceID)

	data, err := a.rdb.Get(ctx, "user-ref"+token).Result()
	if err := json.Unmarshal([]byte(data), &userData); err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("AuthService.ResendVerification - %w", err)
	}

	if err := a.rdb.Del(ctx, "otp"+token).Err(); err != nil {
		return fmt.Errorf("AuthService.ValidateUser - %w", err)
	}

	otp, err := util.GenerateRandomInteger(6)
	if err != nil {
		return fmt.Errorf("AuthService.ResendVerification - %w", err)
	}

	if err := a.rdb.Set(ctx, "otp"+token, otp, 10*time.Minute).Err(); err != nil {
		return fmt.Errorf("AuthService.Register - %w", err)
	}

	mailMessage, err := util.GenerateOTPMailMessage(a.cfg, userData, otp)
	if err != nil {
		return fmt.Errorf("AuthService.ResendVerification - %w", err)
	}

	go func() {
		_ = a.mailDialer.DialAndSend(mailMessage)
	}()

	return nil
}

func (a *AuthServiceImpl) Register(ctx context.Context, entity *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(entity.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	id := ulid.MustNew(ulid.Now(), ulid.Monotonic(rand.Reader, 0))

	newUser := &domain.User{
		UserID:       id.String(),
		Username:     entity.Username,
		PasswordHash: string(password),
		Email:        entity.Email,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	registeredUser, err := a.userRepository.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	otp, err := util.GenerateRandomInteger(6)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	identifier := base64.StdEncoding.EncodeToString([]byte(newUser.UserID))

	if err := a.rdb.Set(ctx, "otp"+identifier, otp, 10*time.Minute).Err(); err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	userData, err := json.Marshal(registeredUser)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}
	if err := a.rdb.Set(ctx, "user-ref"+identifier, userData, 10*time.Minute).Err(); err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	mailMessage, err := util.GenerateOTPMailMessage(a.cfg, newUser, otp)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	go func() {
		_ = a.mailDialer.DialAndSend(mailMessage)
	}()

	return &dto.UserRegisterResponse{
		ReferenceID: identifier,
	}, nil
}

func (a *AuthServiceImpl) ValidateUser(ctx context.Context, entity *dto.UserValidateAccRequest) error {
	referenceToken := strings.TrimSpace(entity.ReferenceID)

	otpValue, err := a.rdb.Get(ctx, "otp"+referenceToken).Result()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("AuthService.ValidateUser - reference token expired or invalid")
	}
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser - %w", err)
	}
	if otpValue != entity.OTP {
		return fmt.Errorf("AuthService.ValidateUser - invalid token")
	}

	userData, err := a.rdb.Get(ctx, "user-ref"+referenceToken).Result()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("AuthService.ValidateUser - reference token expired or invalid")
	}
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser - failed to get user data: %w", err)
	}

	// Unmarshal JSON string to User struct
	var _user *domain.User
	if err := json.Unmarshal([]byte(userData), &_user); err != nil {
		return fmt.Errorf("AuthService.ValidateUser - failed to unmarshal user data: %w", err)
	}

	result, err := a.userRepository.ReadByEmail(ctx, _user)
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser - %w", err)
	}

	_, err = a.userRepository.UpdateUser(ctx, &domain.User{UserID: result.UserID, IsVerified: sql.NullTime{Time: time.Now(), Valid: true}})
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser.UpdateUser - %w", err)
	}

	allRoles, err := a.roleRepository.ReadAllRole(ctx)
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser.ReadAllRole - %w", err)
	}
	roleIndex := make(map[string]domain.Role)
	for _, role := range allRoles {
		roleIndex[role.RoleName] = role
	}

	_, err = a.userRoleRepository.CreateUserRole(ctx, &domain.UserRole{UserID: result.UserID, RoleID: roleIndex["customer"].RoleID})
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser.CreateUserRole - %w", err)
	}

	if err := a.rdb.Del(ctx, "otp"+referenceToken).Err(); err != nil {
		return fmt.Errorf("AuthService.ValidateUser - %w", err)
	}

	if err := a.rdb.Del(ctx, "user-ref"+referenceToken).Err(); err != nil {
		return fmt.Errorf("AuthService.ValidateUser - %w", err)
	}

	return nil
}

func (a *AuthServiceImpl) Login(ctx context.Context, entity *dto.UserLoginRequest) (*dto.UserResponse, error) {
	loginRequest := &domain.User{
		Username:     entity.Username,
		PasswordHash: entity.Password,
	}

	result, err := a.userRepository.ReadByUsername(ctx, loginRequest)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Login - %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.PasswordHash), []byte(entity.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, fmt.Errorf("AuthService.Login - invalid password")
		}

		return nil, fmt.Errorf("AuthService.Login - invalid password")
	}

	userRole := &domain.UserRole{
		UserID: result.UserID,
	}

	role, err := a.userRoleRepository.ReadUserRoleByUserID(ctx, userRole)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Login - %w", err)
	}

	var roleNames []string
	for _, userRole := range role {
		roleNames = append(roleNames, userRole.Role.RoleName)
	}

	jwtPayload := &middleware.JwtPayload{
		ID:       result.UserID,
		Username: result.Username,
		Email:    result.Email,
		Role:     roleNames,
	}

	authToken, err := a.jwtMiddleware.TokenGenerator(ctx, jwtPayload)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Login - %w", err)
	}

	return &dto.UserResponse{
		UserID:    result.UserID,
		Username:  result.Username,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Email:     result.Email,
		Jwt: dto.JwtToken{
			AccessToken:  authToken.AccessToken,
			RefreshToken: authToken.AccessToken,
		},
	}, nil
}
