package repository

import (
	"context"
	"fmt"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"gorm.io/gorm"
)

type UserAddressRepository struct {
	DB *gorm.DB
}

func NewUserAddressRepository(db *gorm.DB) user.UserAddressRepository {
	return &UserAddressRepository{DB: db}
}

func (u *UserAddressRepository) CreateUserAddress(ctx context.Context, entity *domain.UserAddress) (*domain.UserAddress, error) {
	tx := u.DB.WithContext(ctx)
	if err := tx.Create(entity).Error; err != nil {
		return nil, fmt.Errorf("UserAddressRepository.CreateUserAddress - %w", err)
	}

	return entity, nil
}

func (u *UserAddressRepository) DeleteUserAddress(ctx context.Context, entity *domain.UserAddress) error {
	tx := u.DB.WithContext(ctx)
	if err := tx.Delete(entity, "address_id = ? AND user_id = ?", entity.AddressID, entity.UserID).Error; err != nil {
		return fmt.Errorf("UserAddressRepository.DeleteUserAddress - %w", err)
	}

	return nil
}

func (u *UserAddressRepository) UpdateUserAddress(ctx context.Context, entity *domain.UserAddress) (*domain.UserAddress, error) {
	tx := u.DB.WithContext(ctx)
	if err := tx.Model(&domain.UserAddress{}).Where("address_id = ? AND user_id = ?", entity.AddressID, entity.UserID).Updates(entity).Error; err != nil {
		return nil, fmt.Errorf("UserAddressRepository.UpdateUserAddress - %w", err)
	}

	return entity, nil
}

func (u *UserAddressRepository) ReadUserAddress(ctx context.Context, entity *domain.UserAddress) (*domain.UserAddress, error) {
	foundAddress := new(domain.UserAddress)
	tx := u.DB.WithContext(ctx)
	if err := tx.Take(foundAddress, "address_id = ? AND user_id = ?", entity.AddressID, entity.UserID).Error; err != nil {
		return nil, fmt.Errorf("UserAddressRepository.ReadUserAddress - %w", err)
	}

	return foundAddress, nil
}

func (u *UserAddressRepository) ReadAllUserAddress(ctx context.Context, entity *domain.UserAddress) ([]domain.UserAddress, error) {
	var addresses []domain.UserAddress
	tx := u.DB.WithContext(ctx)
	if err := tx.Find(&addresses, "user_id = ?", entity.UserID).Error; err != nil {
		return nil, fmt.Errorf("UserAddressRepository.ReadAllUserAddress - %w", err)
	}

	return addresses, nil
}
