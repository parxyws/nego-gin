package repository

import (
	"context"
	"fmt"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"gorm.io/gorm"
)

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) user.RoleRepository {
	return &RoleRepository{DB: db}
}

func (r *RoleRepository) CreateRole(ctx context.Context, entity *domain.Role) (*domain.Role, error) {
	tx := r.DB.WithContext(ctx)
	result := tx.Where("role_name = ?", entity.RoleName).FirstOrCreate(entity)

	if result.Error != nil {
		return nil, fmt.Errorf("RoleRepository.CreateRole - %w", result.Error)
	}

	return entity, nil
}

func (r *RoleRepository) DeleteRole(ctx context.Context, entity []domain.Role) error {
	tx := r.DB.WithContext(ctx)
	if err := tx.Unscoped().Delete(entity).Error; err != nil {
		return fmt.Errorf("RoleRepository.DeleteRole - %w", err)
	}

	return nil
}

func (r *RoleRepository) ReadAllRole(ctx context.Context) ([]domain.Role, error) {
	var roles []domain.Role
	tx := r.DB.WithContext(ctx)

	if err := tx.Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("RoleRepository.ReadAllRole - %w", err)
	}

	return roles, nil
}

func (r *RoleRepository) ReadAllRoleById(ctx context.Context, roleId []int32) ([]domain.Role, error) {
	var roles []domain.Role
	tx := r.DB.WithContext(ctx)
	if err := tx.Where("role_id = ?", roleId).Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("RoleRepository.ReadAllRoleById - %w", err)
	}

	return roles, nil
}
