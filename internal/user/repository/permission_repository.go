package repository

import (
	"context"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"gorm.io/gorm"
)

type PermissionRepositoryImpl struct {
	DB *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) user.PermissionRepository {
	return &PermissionRepositoryImpl{DB: db}
}

func (p *PermissionRepositoryImpl) CreatePermission(ctx context.Context, entity *domain.Permission) (*domain.Permission, error) {
	dbTx := p.DB.WithContext(ctx)
	err := dbTx.Transaction(func(tx *gorm.DB) error {
		if dbRes := tx.Where("name = ?", entity.Name).FirstOrCreate(entity); dbRes.Error != nil {
			return dbRes.Error
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (p *PermissionRepositoryImpl) DeletePermission(ctx context.Context, entity *domain.Permission) error {
	dbTx := p.DB.WithContext(ctx)
	err := dbTx.Transaction(func(tx *gorm.DB) error {
		if dbRes := tx.Delete(entity, "permission_id", entity.PermissionID); dbRes.Error != nil {
			return dbRes.Error
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *PermissionRepositoryImpl) ReadPermissionById(ctx context.Context, entity *domain.Permission) (*domain.Permission, error) {
	dbTx := p.DB.WithContext(ctx)
	permission := new(domain.Permission)
	if err := dbTx.Where("permission_id = ?", entity.PermissionID).First(&permission).Error; err != nil {
		return nil, err
	}

	return permission, nil
}

func (p *PermissionRepositoryImpl) ReadAllPermission(ctx context.Context, entity *domain.Permission) ([]domain.Permission, error) {
	var permissions []domain.Permission

	dbTx := p.DB.WithContext(ctx)
	dbRes := dbTx.Find(&permissions)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}

	return permissions, nil
}
