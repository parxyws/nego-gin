package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"gorm.io/gorm"
)

type RolePermissionRepository struct {
	DB *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) user.RolePermissionRepository {
	return &RolePermissionRepository{DB: db}
}

func (r *RolePermissionRepository) CreateRolePermission(ctx context.Context, entity *domain.RolePermission) (*domain.RolePermission, error) {
	dbTx := r.DB.WithContext(ctx)
	dbRes := dbTx.Where("role_id = ? AND permission_id", entity.RoleID, entity.PermissionID).FirstOrCreate(entity)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}

	return entity, nil
}

func (r *RolePermissionRepository) DeleteRolePermission(ctx context.Context, entity *domain.RolePermission) error {
	dbTx := r.DB.WithContext(ctx)
	dbRes := dbTx.Delete(entity, "role_id = ? AND permission_id = ?", entity.RoleID, entity.PermissionID)
	if dbRes.Error != nil {
		return dbRes.Error
	}

	return nil
}

func (r *RolePermissionRepository) ReadRolePermissionByPermissionId(ctx context.Context, entity *domain.RolePermission) (*domain.RolePermission, error) {
	rolePermission := new(domain.RolePermission)
	dbTx := r.DB.WithContext(ctx)
	if err := dbTx.Where("permission_id = ?", entity.PermissionID).First(rolePermission).Error; err != nil {
		return nil, err
	}

	return rolePermission, nil
}

func (r *RolePermissionRepository) ReadRolePermissionByRoleId(ctx context.Context, entity *domain.RolePermission) (*domain.RolePermission, error) {
	rolePermission := new(domain.RolePermission)
	dbTx := r.DB.WithContext(ctx)
	if err := dbTx.Where("role_id = ?", entity.RoleID).First(rolePermission).Error; err != nil {
		return nil, err
	}

	return rolePermission, nil
}

func (r *RolePermissionRepository) ReadAllRolePermission(ctx context.Context, entity *domain.RolePermission) ([]domain.RolePermission, error) {
	var rolePermissions []domain.RolePermission
	dbTx := r.DB.WithContext(ctx)
	if err := dbTx.Find(&rolePermissions).Error; err != nil {
		return nil, err
	}

	return rolePermissions, nil
}

func (r *RolePermissionRepository) ReadAllRolePermissionByRoleId(ctx context.Context, entity *domain.RolePermission) ([]domain.RolePermission, error) {
	var rolePermissions []domain.RolePermission
	dbTx := r.DB.WithContext(ctx)
	if err := dbTx.Where("role_id", entity.RoleID).Find(&rolePermissions).Error; err != nil {
		return nil, err
	}

	return rolePermissions, nil
}
