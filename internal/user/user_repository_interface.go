package user

import (
	"context"

	"github.com/parxyws/nego-gin/internal/user/domain"
)

type RoleRepository interface {
	CreateRole(ctx context.Context, entity *domain.Role) (*domain.Role, error)
	DeleteRole(ctx context.Context, entity []domain.Role) error
	ReadAllRole(ctx context.Context) ([]domain.Role, error)
	ReadAllRoleById(ctx context.Context, roleId []int32) ([]domain.Role, error)
}

type PermissionRepository interface {
	CreatePermission(ctx context.Context)
	DeletePermission(ctx context.Context)
	ReadAllPermission(ctx context.Context)
	ReadAllPermissionById(ctx context.Context)
}

type RolePermissionRepository interface {
	CreateRolePermission(ctx context.Context)
	DeleteRolePermission(ctx context.Context)
	ReadAllRolePermission(ctx context.Context)
	ReadAllRolePermissionByRoleId(ctx context.Context)
}

type UserRepository interface {
	CreateUser(ctx context.Context, entity *domain.User) (*domain.User, error)
	UpdateUser(ctx context.Context, entity *domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, entity *domain.User) error

	ReadByUsername(ctx context.Context, entity *domain.User) (*domain.User, error)
	ReadByEmail(ctx context.Context, entity *domain.User) (*domain.User, error)
	ReadById(ctx context.Context, entity *domain.User) (*domain.User, error)
	ReadAllByRoles(ctx context.Context, id string, sortOrderDesc bool, limit int, createdAt string) ([]domain.User, error)
}

type UserRoleRepository interface {
	CreateUserRole(ctx context.Context, entity *domain.UserRole) (*domain.UserRole, error)
	DeleteUserRole(ctx context.Context, entity []domain.UserRole) error
	ReadAllUserRole(ctx context.Context) ([]domain.UserRole, error)
	CountUserRoleByRoleID(ctx context.Context, entity []int32) (map[int32]int64, error)
	ReadUserRoleByUserID(ctx context.Context, entity *domain.UserRole) ([]domain.UserRole, error)
	//ReadAllUserRole()
}

type UserAddressRepository interface {
	CreateUserAddress(ctx context.Context, entity *domain.UserAddress) (*domain.UserAddress, error)
	DeleteUserAddress(ctx context.Context, entity []domain.UserAddress) error
	UpdateUserAddress(ctx context.Context, entity *domain.UserAddress) (*domain.UserAddress, error)
	ReadUserAddress(ctx context.Context, entity *domain.UserAddress) (*domain.UserAddress, error)
	ReadAllUserAddress(ctx context.Context, entity *domain.UserAddress) ([]domain.UserAddress, error)
}
