package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
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
)

type AuthServiceImpl struct {
	cfg                *config.Config
	userRepository     user.UserRepository
	userRoleRepository user.UserRoleRepository
	roleRepository     user.RoleRepository
	rdb                *redis.Client
	jwtMiddleware      *jwt.GinJWTMiddleware
}

func NewAuthService(cfg *config.Config, userRepository user.UserRepository, userRoleRepository user.UserRoleRepository, roleRepository user.RoleRepository, rdb *redis.Client) user.AuthService {
	return &AuthServiceImpl{cfg: cfg, userRepository: userRepository, userRoleRepository: userRoleRepository, roleRepository: roleRepository, rdb: rdb}
}

func (a *AuthServiceImpl) Register(ctx context.Context, entity *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(entity.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	id := ulid.MustNew(ulid.Now(), ulid.Monotonic(rand.Reader, 0))

	_user := &domain.User{
		UserID:       id.String(),
		Username:     entity.Username,
		PasswordHash: string(password),
		Email:        entity.Email,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = a.userRepository.CreateUser(ctx, _user)
	if err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	otp, err := util.GenerateRandomInteger()
	if err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	identifier := base64.StdEncoding.EncodeToString([]byte(_user.Email))

	if err := a.rdb.Set(ctx, "otp"+identifier, otp, 10*time.Minute).Err(); err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	if err := a.rdb.Set(ctx, "user-ref"+identifier, _user.Email, 10*time.Minute).Err(); err != nil {
		return nil, fmt.Errorf("AuthService.Register - %w", err)
	}

	return &dto.UserRegisterResponse{
		ReferenceID: identifier,
	}, nil
}

func (a *AuthServiceImpl) ValidateUser(ctx context.Context, entity *dto.UserValidateAccRequest) error {
	token := strings.TrimSpace(entity.ReferenceID)

	val, err := a.rdb.Get(ctx, "otp"+token).Result()
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser - %w", err)
	}

	if val != entity.OTP {
		return fmt.Errorf("AuthService.ValidateUser - invalid token")
	}

	email, err := a.rdb.Get(ctx, "user-ref"+token).Result()
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser - %w", err)
	}

	_request := &domain.User{
		Email: email,
	}

	result, err := a.userRepository.ReadByEmail(ctx, _request)
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser - %w", err)
	}

	_, err = a.userRepository.UpdateUser(ctx, &domain.User{UserID: result.UserID, IsVerified: sql.NullTime{Time: time.Now(), Valid: true}})
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser - %w", err)
	}

	_role, err := a.roleRepository.ReadAllRole(ctx)
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser - %w", err)
	}
	index := make(map[string]domain.Role)
	for _, role := range _role {
		index[role.RoleName] = role
	}

	_, err = a.userRoleRepository.CreateUserRole(ctx, &domain.UserRole{UserID: result.UserID, RoleID: index["customer"].RoleID})
	if err != nil {
		return fmt.Errorf("AuthService.ValidateUser - %w", err)
	}

	return nil
}

func (a *AuthServiceImpl) Login(ctx context.Context, entity *dto.UserLoginRequest) (*dto.UserResponse, error) {
	_request := &domain.User{
		Username:     entity.Username,
		PasswordHash: entity.Password,
	}

	result, err := a.userRepository.ReadByUsername(ctx, _request)
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

	var roles []string
	for i, v := range role {
		roles[i] = v.Role.RoleName
	}

	payload := &middleware.JwtPayload{
		ID:       result.UserID,
		Username: result.Username,
		Email:    result.Email,
		Role:     roles,
	}

	token, err := a.jwtMiddleware.TokenGenerator(ctx, payload)
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
			AccessToken:  token.AccessToken,
			RefreshToken: token.AccessToken,
		},
	}, nil
}
