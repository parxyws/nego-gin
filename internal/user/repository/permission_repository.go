package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	DB *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) user.PermissionRepository {
	return &PermissionRepository{DB: db}
}

func (p *PermissionRepository) CreatePermission(ctx context.Context, entity *domain.Permission) (*domain.Permission, error) {
	dbTx := p.DB.WithContext(ctx)
	if dbRes := dbTx.Where("name = ?", entity.Name).FirstOrCreate(entity); dbRes.Error != nil {
		return nil, dbRes.Error
	}

	return entity, nil
}

func (p *PermissionRepository) DeletePermission(ctx context.Context, entity *domain.Permission) error {
	dbTx := p.DB.WithContext(ctx)
	if dbRes := dbTx.Delete(entity, "permission_id", entity.PermissionID); dbRes.Error != nil {
		return dbRes.Error
	}

	return nil
}

func (p *PermissionRepository) ReadPermissionById(ctx context.Context, entity *domain.Permission) (*domain.Permission, error) {
	dbTx := p.DB.WithContext(ctx)
	permission := new(domain.Permission)
	if err := dbTx.Where("permission_id = ?", entity.PermissionID).First(&permission).Error; err != nil {
		return nil, err
	}

	return permission, nil
}

func (p *PermissionRepository) ReadAllPermission(ctx context.Context, entity *domain.Permission) ([]domain.Permission, error) {
	var permissions []domain.Permission

	dbTx := p.DB.WithContext(ctx)
	dbRes := dbTx.Find(&permissions)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}

	return permissions, nil
}
