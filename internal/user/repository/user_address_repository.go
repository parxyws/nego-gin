package repository

import (
	"context"
	"fmt"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"gorm.io/gorm"
)

type UserAddressRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserAddressRepository(db *gorm.DB) user.UserAddressRepository {
	return &UserAddressRepositoryImpl{DB: db}
}

func (u *UserAddressRepositoryImpl) CreateUserAddress(ctx context.Context, entity *domain.UserAddress) (*domain.UserAddress, error) {
	tx := u.DB.WithContext(ctx)
	err := tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(entity).Error; err != nil {
			return fmt.Errorf("UserAddressRepositoryImpl.CreateUserAddress - %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (u *UserAddressRepositoryImpl) DeleteUserAddress(ctx context.Context, entity []domain.UserAddress) error {
	tx := u.DB.WithContext(ctx)
	err := tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(entity).Error; err != nil {
			return fmt.Errorf("UserAddressRepositoryImpl.DeleteUserAddress - %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (u *UserAddressRepositoryImpl) UpdateUserAddress(ctx context.Context, entity *domain.UserAddress) (*domain.UserAddress, error) {
	tx := u.DB.WithContext(ctx)
	err := tx.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(entity).Error; err != nil {
			return fmt.Errorf("UserAddressRepositoryImpl.UpdateUserAddress - %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (u *UserAddressRepositoryImpl) ReadUserAddress(ctx context.Context, entity *domain.UserAddress) (*domain.UserAddress, error) {
	_address := new(domain.UserAddress)
	tx := u.DB.WithContext(ctx)
	if err := tx.Take(_address, "address_id = ? AND user_id = ?", entity.AddressID, entity.UserID).Error; err != nil {
		return nil, fmt.Errorf("UserAddressRepositoryImpl.ReadUserAddress - %w", err)
	}

	return _address, nil
}

func (u *UserAddressRepositoryImpl) ReadAllUserAddress(ctx context.Context, entity *domain.UserAddress) ([]domain.UserAddress, error) {
	var addresses []domain.UserAddress
	tx := u.DB.WithContext(ctx)
	if err := tx.Find(&addresses, "user_id = ?", entity.UserID).Error; err != nil {
		return nil, fmt.Errorf("UserAddressRepositoryImpl.ReadAllUserAddress - %w", err)
	}

	return addresses, nil
}
