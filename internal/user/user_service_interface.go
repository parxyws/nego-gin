package user

import (
	"context"

	"github.com/parxyws/nego-gin/internal/user/domain/dto"
	"github.com/parxyws/nego-gin/middleware"
	"github.com/parxyws/nego-gin/pkg/database/aws"
)

type AuthService interface {
	Register(ctx context.Context, entity *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error)
	ValidateUser(ctx context.Context, entity *dto.UserValidateAccRequest) error
	Login(ctx context.Context, entity *dto.UserLoginRequest) (*dto.UserResponse, error)
	RefreshToken(ctx context.Context, entity *dto.JwtToken, payload *middleware.JwtPayload) (*dto.UserResponse, error)
	ForgotPassword(ctx context.Context, request *dto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, request *dto.ResetPasswordRequest) error
	ResendVerification(ctx context.Context, entity *dto.UserRegisterResponse) error
}

type PermissionService interface {
	RegisterPermission(ctx context.Context)
	DeletePermission(ctx context.Context)
	GetAllPermission(ctx context.Context)
}

type RoleService interface {
	RegisterRole(ctx context.Context, entity *dto.RoleRegisterRequest) (*dto.RoleResponse, error)
	CheckDeleteRole(ctx context.Context, entity []dto.RoleDeleteRequest) ([]dto.RoleResponse, error)
	DeleteRole(ctx context.Context, entity []dto.RoleDeleteRequest) error
	GetAllRoles(ctx context.Context) ([]dto.RoleResponse, error)
}

type UserService interface {
	GetCurrentUser(ctx context.Context, entity middleware.JwtPayload) (*dto.UserResponse, error)
	UpdateCurrentUser(ctx context.Context, entity middleware.JwtPayload, request *dto.UpdateCurrentUserRequest) (*dto.UserResponse, error)
	UpdateAvatar(ctx context.Context, entity middleware.JwtPayload, avatar *aws.UploadInput) (*dto.UserResponse, error)
	DeleteCurrentUser(ctx context.Context, entity middleware.JwtPayload) error
	GetUser(ctx context.Context, request *dto.GetUserProfileRequest) (*dto.UserProfileResponse, error)
}

type UserAddressService interface {
	CreateUserAddress(ctx context.Context, entity *dto.UserAddressCreateRequest) (*dto.UserAddressResponse, error)
	UpdateUserAddress(ctx context.Context, entity *dto.UserAddressUpdateRequest) (*dto.UserAddressResponse, error)
	DeleteUserAddress(ctx context.Context, addressId int32, userId string) error
	ListUserAddress(ctx context.Context, userId string) ([]dto.UserAddressResponse, error)
}
