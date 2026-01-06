package user

import (
	"context"

	"github.com/parxyws/nego-gin/internal/user/domain/dto"
)

type AuthService interface {
	Register(ctx context.Context, entity *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error)
	ValidateUser(ctx context.Context, entity *dto.UserValidateAccRequest) error
	Login(ctx context.Context, entity *dto.UserLoginRequest) (*dto.UserResponse, error)
	RefreshToken(ctx context.Context)
	ForgotPassword(ctx context.Context)
	ResetPassword(ctx context.Context)
	ResendVerification(ctx context.Context)
}

type RoleService interface {
	RegisterRole(ctx context.Context, entity *dto.RoleRegisterRequest) (*dto.RoleResponse, error)
	CheckDeleteRole(ctx context.Context, entity []dto.RoleDeleteRequest) ([]dto.RoleResponse, error)
	DeleteRole(ctx context.Context, entity []dto.RoleDeleteRequest) error
	GetAllRoles(ctx context.Context) ([]dto.RoleResponse, error)
}

type UserService interface {
	ReadCurrentUser(ctx context.Context)
	ReadAllUser(ctx context.Context)
}
