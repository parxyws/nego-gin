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
	"github.com/parxyws/nego-gin/internal/middleware"
	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"github.com/parxyws/nego-gin/internal/user/domain/dto"
	"github.com/parxyws/nego-gin/pkg/util"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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

func NewAuthService(cfg *config.Config, userRepository user.UserRepository, userRoleRepository user.UserRoleRepository, roleRepository user.RoleRepository, rdb *redis.Client, mailDialer *gomail.Dialer, jwtMiddleware *jwt.GinJWTMiddleware) user.AuthService {
	return &AuthServiceImpl{cfg: cfg, userRepository: userRepository, userRoleRepository: userRoleRepository, roleRepository: roleRepository, rdb: rdb, mailDialer: mailDialer, jwtMiddleware: jwtMiddleware}
}

func (s *AuthServiceImpl) RefreshToken(ctx context.Context, tokenReq *dto.JwtToken, payload *middleware.JwtPayload) (*dto.UserResponse, error) {
	request := &domain.User{
		UserID: payload.ID,
	}
	user, err := s.userRepository.ReadByIdMinimal(ctx, request)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.RefreshToken", "user_id": payload.ID}).Errorf("failed to read user: %v", err)
		return nil, err
	}

	token, err := s.jwtMiddleware.TokenGeneratorWithRevocation(ctx, payload, tokenReq.RefreshToken)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.RefreshToken", "user_id": payload.ID}).Errorf("failed to generate token: %v", err)
		return nil, err
	}

	logrus.WithFields(logrus.Fields{"function": "AuthService.RefreshToken", "user_id": user.UserID}).Info("Token refreshed successfully")

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
		logrus.WithFields(logrus.Fields{"function": "AuthService.ForgotPassword", "email": request.Email}).Warnf("user not found: %v", err)
		return util.ErrNotFound
	}

	otp, err := util.GenerateRandomInteger(6)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ForgotPassword", "email": request.Email}).Errorf("failed to generate otp: %v", err)
		return fmt.Errorf("Generate OTP failed: %w", err)
	}

	if err := s.rdb.Set(ctx, "forgot-otp:"+request.Email, otp, 15*time.Minute).Err(); err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ForgotPassword", "email": request.Email}).Errorf("failed to store otp in redis: %v", err)
		return fmt.Errorf("Store OTP failed: %w", err)
	}

	mailMessage, err := util.GenerateOTPMailMessage(s.cfg, foundUser, otp)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ForgotPassword", "email": request.Email}).Errorf("failed to generate mail message: %v", err)
		return fmt.Errorf("Generate mail message failed: %w", err)
	}

	// Async email sending
	go func() {
		if err := s.mailDialer.DialAndSend(mailMessage); err != nil {
			logrus.WithFields(logrus.Fields{"function": "AuthService.ForgotPassword.Async", "email": request.Email}).Errorf("failed to send email: %v", err)
		} else {
			logrus.WithFields(logrus.Fields{"function": "AuthService.ForgotPassword.Async", "email": request.Email}).Info("Forgot password email sent successfully")
		}
	}()

	return nil
}

func (s *AuthServiceImpl) ResetPassword(ctx context.Context, request *dto.ResetPasswordRequest) error {
	otpValue, err := s.rdb.Get(ctx, "forgot-otp:"+request.Email).Result()
	if errors.Is(err, redis.Nil) {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ResetPassword", "email": request.Email}).Warn("otp expired or invalid")
		return util.NewAppError(http.StatusBadRequest, "Validation failed: OTP expired or invalid", nil)
	}
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ResetPassword", "email": request.Email}).Errorf("failed to get otp from redis: %v", err)
		return fmt.Errorf("Retrieve OTP failed: %w", err)
	}

	if otpValue != request.OTP {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ResetPassword", "email": request.Email}).Warn("invalid otp version")
		return util.NewAppError(http.StatusBadRequest, "Validation failed: invalid OTP", nil)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ResetPassword", "email": request.Email}).Errorf("failed to hash password: %v", err)
		return fmt.Errorf("Hash password failed: %w", err)
	}

	foundUser, err := s.userRepository.ReadByEmail(ctx, &domain.User{Email: request.Email})
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ResetPassword", "email": request.Email}).Warnf("user not found: %v", err)
		return util.ErrNotFound
	}

	_, err = s.userRepository.UpdateUser(ctx, &domain.User{
		UserID:       foundUser.UserID,
		PasswordHash: string(passwordHash),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ResetPassword", "email": request.Email}).Errorf("failed to update user password: %v", err)
		return fmt.Errorf("Update user password failed: %w", err)
	}

	_ = s.rdb.Del(ctx, "forgot-otp:"+request.Email).Err()
	logrus.WithFields(logrus.Fields{"function": "AuthService.ResetPassword", "email": request.Email}).Info("Password reset successfully")

	return nil
}

func (s *AuthServiceImpl) ResendVerification(ctx context.Context, resp *dto.UserRegisterResponse) error {
	var userDataModel *domain.User
	token := strings.TrimSpace(resp.ReferenceID)

	userData, err := s.rdb.Get(ctx, "user-ref"+token).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ResendVerification", "ref_id": token}).Errorf("failed to get user data from redis: %v", err)
		return fmt.Errorf("Retrieve user data failed: %w", err)
	}

	if err := json.Unmarshal([]byte(userData), &userDataModel); err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ResendVerification", "ref_id": token}).Errorf("failed to unmarshal user data: %v", err)
		return err
	}

	if err := s.rdb.Del(ctx, "otp"+token).Err(); err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ResendVerification", "ref_id": token}).Errorf("failed to delete old otp from redis: %v", err)
		return fmt.Errorf("Delete old OTP failed: %w", err)
	}

	otp, err := util.GenerateRandomInteger(6)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ResendVerification", "ref_id": token}).Errorf("failed to generate otp: %v", err)
		return fmt.Errorf("Generate OTP failed: %w", err)
	}

	if err := s.rdb.Set(ctx, "otp"+token, otp, 10*time.Minute).Err(); err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ResendVerification", "ref_id": token}).Errorf("failed to set otp in redis: %v", err)
		return fmt.Errorf("Store OTP failed: %w", err)
	}

	mailMessage, err := util.GenerateOTPMailMessage(s.cfg, userDataModel, otp)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ResendVerification", "ref_id": token}).Errorf("failed to generate mail message: %v", err)
		return fmt.Errorf("Generate mail message failed: %w", err)
	}

	go func() {
		if err := s.mailDialer.DialAndSend(mailMessage); err != nil {
			logrus.WithFields(logrus.Fields{"function": "AuthService.ResendVerification.Async", "ref_id": token}).Errorf("failed to send email: %v", err)
		} else {
			logrus.WithFields(logrus.Fields{"function": "AuthService.ResendVerification.Async", "ref_id": token}).Info("Verification email resent successfully")
		}
	}()

	return nil
}

func (s *AuthServiceImpl) Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.Register", "email": req.Email}).Errorf("failed to hash password: %v", err)
		return nil, fmt.Errorf("Hash password failed: %w", err)
	}

	id := ulid.MustNew(ulid.Now(), ulid.Monotonic(rand.Reader, 0))

	newUser := &domain.User{
		UserID:       id.String(),
		Username:     req.Username,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: string(password),
		Email:        req.Email,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	registeredUser, err := s.userRepository.CreateUser(ctx, newUser)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.Register", "email": req.Email}).Errorf("failed to create user in repo: %v", err)
		return nil, fmt.Errorf("Create user in repository failed: %w", err)
	}

	otp, err := util.GenerateRandomInteger(6)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.Register", "email": req.Email}).Errorf("failed to generate otp: %v", err)
		return nil, fmt.Errorf("Generate OTP failed: %w", err)
	}

	identifier := base64.StdEncoding.EncodeToString([]byte(newUser.UserID))

	if err := s.rdb.Set(ctx, "otp"+identifier, otp, 10*time.Minute).Err(); err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.Register", "email": req.Email, "ref_id": identifier}).Errorf("failed to set otp in redis: %v", err)
		return nil, fmt.Errorf("Store OTP failed: %w", err)
	}

	userData, err := json.Marshal(registeredUser)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.Register", "email": req.Email, "ref_id": identifier}).Errorf("failed to marshal user data: %v", err)
		return nil, fmt.Errorf("Marshal user data failed: %w", err)
	}
	if err := s.rdb.Set(ctx, "user-ref"+identifier, userData, 10*time.Minute).Err(); err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.Register", "email": req.Email, "ref_id": identifier}).Errorf("failed to set user-ref in redis: %v", err)
		return nil, fmt.Errorf("Store user reference failed: %w", err)
	}

	mailMessage, err := util.GenerateOTPMailMessage(s.cfg, newUser, otp)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.Register", "email": req.Email, "ref_id": identifier}).Errorf("failed to generate mail message: %v", err)
		return nil, fmt.Errorf("Generate mail message failed: %w", err)
	}

	go func() {
		if err := s.mailDialer.DialAndSend(mailMessage); err != nil {
			logrus.WithFields(logrus.Fields{"function": "AuthService.Register.Async", "email": req.Email, "ref_id": identifier}).Errorf("failed to send email: %v", err)
		} else {
			logrus.WithFields(logrus.Fields{"function": "AuthService.Register.Async", "email": req.Email, "ref_id": identifier}).Info("registration otp email sent successfully")
		}
	}()

	logrus.WithFields(logrus.Fields{"function": "AuthService.Register", "email": req.Email, "ref_id": identifier}).Info("User registered successfully")

	return &dto.UserRegisterResponse{
		ReferenceID: identifier,
	}, nil
}

func (s *AuthServiceImpl) ValidateUser(ctx context.Context, req *dto.UserValidateAccRequest) error {
	referenceToken := strings.TrimSpace(req.ReferenceID)

	otpValue, err := s.rdb.Get(ctx, "otp"+referenceToken).Result()
	if errors.Is(err, redis.Nil) {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ValidateUser", "ref_id": referenceToken}).Warn("reference token expired or invalid")
		return fmt.Errorf("AuthService.ValidateUser - reference token expired or invalid")
	}
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ValidateUser", "ref_id": referenceToken}).Errorf("failed to get otp from redis: %v", err)
		return fmt.Errorf("Retrieve OTP failed: %w", err)
	}
	if otpValue != req.OTP {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ValidateUser", "ref_id": referenceToken}).Warn("invalid token attempt")
		return fmt.Errorf("AuthService.ValidateUser - invalid token")
	}

	userData, err := s.rdb.Get(ctx, "user-ref"+referenceToken).Result()
	if errors.Is(err, redis.Nil) {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ValidateUser", "ref_id": referenceToken}).Warn("user reference token expired or invalid")
		return fmt.Errorf("AuthService.ValidateUser - reference token expired or invalid")
	}
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ValidateUser", "ref_id": referenceToken}).Errorf("failed to get user data from redis: %v", err)
		return fmt.Errorf("AuthService.ValidateUser - failed to get user data: %w", err)
	}

	// Unmarshal JSON string to User struct
	var unmarshalUser *domain.User
	if err := json.Unmarshal([]byte(userData), &unmarshalUser); err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ValidateUser", "ref_id": referenceToken}).Errorf("failed to unmarshal user data: %v", err)
		return fmt.Errorf("AuthService.ValidateUser - failed to unmarshal user data: %w", err)
	}

	user, err := s.userRepository.ReadByEmail(ctx, unmarshalUser)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ValidateUser", "email": unmarshalUser.Email}).Errorf("failed to read user from repo: %v", err)
		return fmt.Errorf("Retrieve user failed: %w", err)
	}

	_, err = s.userRepository.UpdateUser(ctx, &domain.User{UserID: user.UserID, IsVerified: sql.NullTime{Time: time.Now(), Valid: true}})
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ValidateUser", "user_id": user.UserID}).Errorf("failed to update user verification status: %v", err)
		return fmt.Errorf("AuthService.ValidateUser.UpdateUser - %w", err)
	}

	roles, err := s.roleRepository.ReadAllRole(ctx)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ValidateUser", "user_id": user.UserID}).Errorf("failed to read roles: %v", err)
		return fmt.Errorf("AuthService.ValidateUser.ReadAllRole - %w", err)
	}
	roleMap := make(map[string]domain.Role)
	for _, role := range roles {
		roleMap[role.RoleName] = role
	}

	_, err = s.userRoleRepository.CreateUserRole(ctx, &domain.UserRole{UserID: user.UserID, RoleID: roleMap["customer"].RoleID})
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ValidateUser", "user_id": user.UserID}).Errorf("failed to create user role: %v", err)
		return fmt.Errorf("AuthService.ValidateUser.CreateUserRole - %w", err)
	}

	if err := s.rdb.Del(ctx, "otp"+referenceToken).Err(); err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ValidateUser", "ref_id": referenceToken}).Errorf("failed to delete otp from redis: %v", err)
		return fmt.Errorf("Delete OTP failed: %w", err)
	}

	if err := s.rdb.Del(ctx, "user-ref"+referenceToken).Err(); err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.ValidateUser", "ref_id": referenceToken}).Errorf("failed to delete user-ref from redis: %v", err)
		return fmt.Errorf("Delete user reference failed: %w", err)
	}

	logrus.WithFields(logrus.Fields{"function": "AuthService.ValidateUser", "user_id": user.UserID}).Info("User verified successfully")

	return nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserResponse, error) {
	loginRequest := &domain.User{
		Username:     req.Username,
		PasswordHash: req.Password,
	}

	user, err := s.userRepository.ReadByUsername(ctx, loginRequest)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.Login", "username": req.Username}).Warnf("user not found: %v", err)
		return nil, fmt.Errorf("Retrieve user by username failed: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.Login", "username": req.Username}).Warn("invalid password attempt")
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, fmt.Errorf("AuthService.Login - invalid password")
		}

		return nil, fmt.Errorf("AuthService.Login - invalid password")
	}

	if !user.IsVerified.Valid {
		logrus.WithFields(logrus.Fields{"function": "AuthService.Login", "username": req.Username}).Warn("unverified user login attempt")
		return nil, fmt.Errorf("AuthService.Login - user is not verified")
	}

	user.LastLogin = sql.NullTime{Time: time.Now(), Valid: true}

	userLogin, err := s.userRepository.UpdateSingleColumnUser(ctx, user, "last_login")
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.Login", "username": req.Username}).Errorf("failed to update last login: %v", err)
		return nil, fmt.Errorf("Update last login failed: %w", err)
	}

	userRoleReq := &domain.UserRole{
		UserID: userLogin.UserID,
	}

	roles, err := s.userRoleRepository.ReadUserRoleByUserID(ctx, userRoleReq)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.Login", "username": req.Username, "user_id": userLogin.UserID}).Errorf("failed to read user roles: %v", err)
		return nil, fmt.Errorf("Retrieve user roles failed: %w", err)
	}

	var roleNames []string
	for _, ur := range roles {
		roleNames = append(roleNames, ur.Role.RoleName)
	}

	payload := &middleware.JwtPayload{
		ID:       userLogin.UserID,
		Username: userLogin.Username,
		Email:    userLogin.Email,
		Role:     roleNames,
	}

	token, err := s.jwtMiddleware.TokenGenerator(ctx, payload)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "AuthService.Login", "username": req.Username, "user_id": userLogin.UserID}).Errorf("failed to generate jwt: %v", err)
		return nil, fmt.Errorf("Generate JWT failed: %w", err)
	}

	logrus.WithFields(logrus.Fields{"function": "AuthService.Login", "username": userLogin.Username, "user_id": userLogin.UserID}).Info("User logged in successfully")

	return &dto.UserResponse{
		UserID:    userLogin.UserID,
		Username:  userLogin.Username,
		FirstName: userLogin.FirstName,
		LastName:  userLogin.LastName,
		Email:     userLogin.Email,
		Jwt: dto.JwtToken{
			AccessToken:  token.AccessToken,
			RefreshToken: token.AccessToken,
		},
	}, nil
}
