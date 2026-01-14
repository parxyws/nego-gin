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

func (s *AuthServiceImpl) RefreshToken(ctx context.Context, tokenReq *dto.JwtToken, payload *middleware.JwtPayload) (*dto.UserResponse, error) {
	request := &domain.User{
		UserID: payload.ID,
	}
	user, err := s.userRepository.ReadByIdMinimal(ctx, request)
	token, err := s.jwtMiddleware.TokenGeneratorWithRevocation(ctx, payload, tokenReq.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Jwt: dto.JwtToken{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		},
	}, nil
}

func (s *AuthServiceImpl) ForgotPassword(ctx context.Context, request *dto.ForgotPasswordRequest) error {
	userRequest := &domain.User{Email: request.Email}
	foundUser, err := s.userRepository.ReadByEmail(ctx, userRequest)
	if err != nil {
		return util.ErrNotFound
	}

	otp, err := util.GenerateRandomInteger(6)
	if err != nil {
		return fmt.Errorf("AuthService.ForgotPassword - %w", err)
	}

	if err := s.rdb.Set(ctx, "forgot-otp:"+request.Email, otp, 15*time.Minute).Err(); err != nil {
		return fmt.Errorf("AuthService.ForgotPassword - %w", err)
	}

	mailMessage, err := util.GenerateOTPMailMessage(s.cfg, foundUser, otp)
	if err != nil {
		return fmt.Errorf("AuthService.ForgotPassword - %w", err)
	}

	// Async email sending
	go func() {
		_ = s.mailDialer.DialAndSend(mailMessage)
	}()

	return nil
}

func (s *AuthServiceImpl) ResetPassword(ctx context.Context, request *dto.ResetPasswordRequest) error {
	otpValue, err := s.rdb.Get(ctx, "forgot-otp:"+request.Email).Result()
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

	foundUser, err := s.userRepository.ReadByEmail(ctx, &domain.User{Email: request.Email})
	if err != nil {
		return util.ErrNotFound
	}

	_, err = s.userRepository.UpdateUser(ctx, &domain.User{
		UserID:       foundUser.UserID,
		PasswordHash: string(passwordHash),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		return fmt.Errorf("AuthService.ResetPassword - %w", err)
	}

	_ = s.rdb.Del(ctx, "forgot-otp:"+request.Email).Err()

	return nil
}

func (s *AuthServiceImpl) ResendVerification(ctx context.Context, resp *dto.UserRegisterResponse) error {
	var userDataModel *domain.User
	token := strings.TrimSpace(resp.ReferenceID)

	userData, err := s.rdb.Get(ctx, "user-ref"+token).Result()
	if err := json.Unmarshal([]byte(userData), &userDataModel); err != nil {
		return err
	}

	if err != nil {
		return fmt.Errorf("AuthService.ResendVerification - %w", err)
	}

	if err := s.rdb.Del(ctx, "otp"+token).Err(); err != nil {
		return fmt.Errorf("AuthService.ValidateUser - %w", err)
	}

	otp, err := util.GenerateRandomInteger(6)
	if err != nil {
		return fmt.Errorf("AuthService.ResendVerification - %w", err)
	}

	if err := s.rdb.Set(ctx, "otp"+token, otp, 10*time.Minute).Err(); err != nil {
		return fmt.Errorf("AuthService.Register - %w", err)
	}

	mailMessage, err := util.GenerateOTPMailMessage(s.cfg, userDataModel, otp)
	if err != nil {
		return fmt.Errorf("AuthService.ResendVerification - %w", err)
	}

	go func() {
		_ = s.mailDialer.DialAndSend(mailMessage)
	}()

	return nil
}

func (s *AuthServiceImpl) Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	id := ulid.MustNew(ulid.Now(), ulid.Monotonic(rand.Reader, 0))

	newUser := &domain.User{
		UserID:       id.String(),
		Username:     req.Username,
		PasswordHash: string(password),
		Email:        req.Email,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	registeredUser, err := s.userRepository.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	otp, err := util.GenerateRandomInteger(6)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	identifier := base64.StdEncoding.EncodeToString([]byte(newUser.UserID))

	if err := s.rdb.Set(ctx, "otp"+identifier, otp, 10*time.Minute).Err(); err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	userData, err := json.Marshal(registeredUser)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}
	if err := s.rdb.Set(ctx, "user-ref"+identifier, userData, 10*time.Minute).Err(); err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	mailMessage, err := util.GenerateOTPMailMessage(s.cfg, newUser, otp)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	go func() {
		_ = s.mailDialer.DialAndSend(mailMessage)
	}()

	return &dto.UserRegisterResponse{
		ReferenceID: identifier,
	}, nil
}

func (s *AuthServiceImpl) ValidateUser(ctx context.Context, req *dto.UserValidateAccRequest) error {
	referenceToken := strings.TrimSpace(req.ReferenceID)

	otpValue, err := s.rdb.Get(ctx, "otp"+referenceToken).Result()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("AuthService.ValidateUser - reference token expired or invalid")
	}
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser - %w", err)
	}
	if otpValue != req.OTP {
		return fmt.Errorf("AuthService.ValidateUser - invalid token")
	}

	userData, err := s.rdb.Get(ctx, "user-ref"+referenceToken).Result()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("AuthService.ValidateUser - reference token expired or invalid")
	}
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser - failed to get user data: %w", err)
	}

	// Unmarshal JSON string to User struct
	var unmarshaledUser *domain.User
	if err := json.Unmarshal([]byte(userData), &unmarshaledUser); err != nil {
		return fmt.Errorf("AuthService.ValidateUser - failed to unmarshal user data: %w", err)
	}

	user, err := s.userRepository.ReadByEmail(ctx, unmarshaledUser)
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser - %w", err)
	}

	_, err = s.userRepository.UpdateUser(ctx, &domain.User{UserID: user.UserID, IsVerified: sql.NullTime{Time: time.Now(), Valid: true}})
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser.UpdateUser - %w", err)
	}

	roles, err := s.roleRepository.ReadAllRole(ctx)
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser.ReadAllRole - %w", err)
	}
	roleMap := make(map[string]domain.Role)
	for _, role := range roles {
		roleMap[role.RoleName] = role
	}

	_, err = s.userRoleRepository.CreateUserRole(ctx, &domain.UserRole{UserID: user.UserID, RoleID: roleMap["customer"].RoleID})
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser.CreateUserRole - %w", err)
	}

	if err := s.rdb.Del(ctx, "otp"+referenceToken).Err(); err != nil {
		return fmt.Errorf("AuthService.ValidateUser - %w", err)
	}

	if err := s.rdb.Del(ctx, "user-ref"+referenceToken).Err(); err != nil {
		return fmt.Errorf("AuthService.ValidateUser - %w", err)
	}

	return nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserResponse, error) {
	loginRequest := &domain.User{
		Username:     req.Username,
		PasswordHash: req.Password,
	}

	user, err := s.userRepository.ReadByUsername(ctx, loginRequest)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Login - %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, fmt.Errorf("AuthService.Login - invalid password")
		}

		return nil, fmt.Errorf("AuthService.Login - invalid password")
	}

	userRoleReq := &domain.UserRole{
		UserID: user.UserID,
	}

	roles, err := s.userRoleRepository.ReadUserRoleByUserID(ctx, userRoleReq)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Login - %w", err)
	}

	var roleNames []string
	for _, ur := range roles {
		roleNames = append(roleNames, ur.Role.RoleName)
	}

	payload := &middleware.JwtPayload{
		ID:       user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Role:     roleNames,
	}

	token, err := s.jwtMiddleware.TokenGenerator(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Login - %w", err)
	}

	return &dto.UserResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Jwt: dto.JwtToken{
			AccessToken:  token.AccessToken,
			RefreshToken: token.AccessToken,
		},
	}, nil
}
