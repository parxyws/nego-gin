package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"gorm.io/gorm"
)

type RolePermissionRepositoryImpl struct {
	DB *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) user.RolePermissionRepository {
	return &RolePermissionRepositoryImpl{DB: db}
}

func (r *RolePermissionRepositoryImpl) CreateRolePermission(ctx context.Context, entity *domain.RolePermission) (*domain.RolePermission, error) {
	dbTx := r.DB.WithContext(ctx)
	err := dbTx.Transaction(func(tx *gorm.DB) error {
		dbRes := tx.Where("role_id = ? AND permission_id", entity.RoleID, entity.PermissionID).FirstOrCreate(entity)
		if dbRes.Error != nil {
			return dbRes.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *RolePermissionRepositoryImpl) DeleteRolePermission(ctx context.Context, entity *domain.RolePermission) error {
	dbTx := r.DB.WithContext(ctx)
	err := dbTx.Transaction(func(tx *gorm.DB) error {
		dbRes := tx.Delete(entity, "role_id = ? AND permission_id = ?", entity.RoleID, entity.PermissionID)
		if dbRes.Error != nil {
			return dbRes.Error
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *RolePermissionRepositoryImpl) ReadRolePermissionByPermissionId(ctx context.Context, entity *domain.RolePermission) (*domain.RolePermission, error) {
	rolePermission := new(domain.RolePermission)
	dbTx := r.DB.WithContext(ctx)
	if err := dbTx.Where("permission_id = ?", entity.PermissionID).First(rolePermission).Error; err != nil {
		return nil, err
	}

	return rolePermission, nil
}

func (r *RolePermissionRepositoryImpl) ReadRolePermissionByRoleId(ctx context.Context, entity *domain.RolePermission) (*domain.RolePermission, error) {
	rolePermission := new(domain.RolePermission)
	dbTx := r.DB.WithContext(ctx)
	if err := dbTx.Where("role_id = ?", entity.RoleID).First(rolePermission).Error; err != nil {
		return nil, err
	}

	return rolePermission, nil
}

func (r *RolePermissionRepositoryImpl) ReadAllRolePermission(ctx context.Context, entity *domain.RolePermission) ([]domain.RolePermission, error) {
	var rolePermissions []domain.RolePermission
	dbTx := r.DB.WithContext(ctx)
	if err := dbTx.Find(&rolePermissions).Error; err != nil {
		return nil, err
	}

	return rolePermissions, nil
}

func (r *RolePermissionRepositoryImpl) ReadAllRolePermissionByRoleId(ctx context.Context, entity *domain.RolePermission) ([]domain.RolePermission, error) {
	var rolePermissions []domain.RolePermission
	dbTx := r.DB.WithContext(ctx)
	if err := dbTx.Where("role_id", entity.RoleID).Find(&rolePermissions).Error; err != nil {
		return nil, err
	}

	return rolePermissions, nil
}
