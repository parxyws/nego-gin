package repository

import (
	"context"
	"fmt"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"gorm.io/gorm"
)

type UserRoleRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) user.UserRoleRepository {
	return &UserRoleRepositoryImpl{DB: db}
}

func (u *UserRoleRepositoryImpl) CreateUserRole(ctx context.Context, entity *domain.UserRole) (*domain.UserRole, error) {
	tx := u.DB.WithContext(ctx)
	err := tx.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(entity)

		if result.Error != nil {
			return fmt.Errorf("UserRoleRepository.CreateUserRole %w", result.Error)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (u *UserRoleRepositoryImpl) DeleteUserRole(ctx context.Context, entity []domain.UserRole) error {
	tx := u.DB.WithContext(ctx)
	err := tx.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&entity)

		if result.Error != nil {
			return fmt.Errorf("UserRoleRepository.DeleteUserRole %w", result.Error)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (u *UserRoleRepositoryImpl) ReadAllUserRole(ctx context.Context) ([]domain.UserRole, error) {
	var userRoles []domain.UserRole
	tx := u.DB.WithContext(ctx)
	if err := tx.Find(&userRoles).Error; err != nil {
		return nil, fmt.Errorf("UserRoleRepository.ReadAllUserRole %w", err)
	}

	return userRoles, nil
}

func (u *UserRoleRepositoryImpl) ReadUserRoleByUserID(ctx context.Context, entity *domain.UserRole) ([]domain.UserRole, error) {
	var userRoles []domain.UserRole
	tx := u.DB.WithContext(ctx)

	if err := tx.Where("user_id = ?", entity.UserID).Preload("Role").Find(&userRoles).Error; err != nil {
		return nil, fmt.Errorf("UserRoleRepository.ReadUserRoleByUserID %w", err)
	}

	return userRoles, nil
}

func (u *UserRoleRepositoryImpl) CountUserRoleByRoleID(ctx context.Context, roleIDs []int32) (map[int32]int64, error) {
	if len(roleIDs) == 0 {
		return map[int32]int64{}, nil
	}

	type RoleCount struct {
		RoleID int32 `gorm:"column:role_id"`
		Count  int64 `gorm:"column:count"`
	}

	var results []RoleCount
	tx := u.DB.WithContext(ctx)

	err := tx.Model(&domain.UserRole{}).
		Select("role_id, COUNT(*) as count").
		Where("role_id IN ?", roleIDs).
		Group("role_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	countMap := make(map[int32]int64)
	for _, result := range results {
		countMap[result.RoleID] = result.Count
	}

	return countMap, nil
}
