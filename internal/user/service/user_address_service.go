package service

import (
	"context"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"github.com/parxyws/nego-gin/internal/user/domain/dto"
	"github.com/parxyws/nego-gin/pkg/logger"
)

type UserAddressService struct {
	userAddressRepository user.UserAddressRepository
}

func NewUserAddressService(userAddressRepository user.UserAddressRepository) user.UserAddressService {
	return &UserAddressService{userAddressRepository: userAddressRepository}
}

func (u *UserAddressService) CreateUserAddress(ctx context.Context, entity *dto.UserAddressCreateRequest) (*dto.UserAddressResponse, error) {
	log := logger.WithCtx(ctx, "service", "UserAddressService.CreateUserAddress").WithField("user_id", entity.UserId)
	request := &domain.UserAddress{
		UserID:        entity.UserId,
		AddressType:   entity.AddressType,
		RecipientName: entity.RecipientName,
		AddressLine1:  entity.AddressLine1,
		AddressLine2:  entity.AddressLine2,
		City:          entity.City,
		StateProvince: entity.StateProvince,
		PostalCode:    entity.PostalCode,
		CountryCode:   entity.CountryCode,
		Phone:         entity.Phone,
		IsDefault:     entity.IsDefault,
	}

	result, err := u.userAddressRepository.CreateUserAddress(ctx, request)
	if err != nil {
		log.Errorf("failed to create address: %v", err)
		return nil, err
	}

	log.WithField("address_id", result.AddressID).Info("Address created successfully")

	return &dto.UserAddressResponse{
		AddressID:     result.AddressID,
		UserID:        result.UserID,
		AddressType:   result.AddressType,
		RecipientName: result.RecipientName,
		AddressLine1:  result.AddressLine1,
		AddressLine2:  result.AddressLine2,
		City:          result.City,
		StateProvince: result.StateProvince,
		PostalCode:    result.PostalCode,
		CountryCode:   result.CountryCode,
		Phone:         result.Phone,
		IsDefault:     result.IsDefault,
	}, nil
}

func (u *UserAddressService) UpdateUserAddress(ctx context.Context, entity *dto.UserAddressUpdateRequest) (*dto.UserAddressResponse, error) {
	log := logger.WithCtx(ctx, "service", "UserAddressService.UpdateUserAddress").WithField("user_id", entity.UserId)
	request := &domain.UserAddress{
		UserID:        entity.UserId,
		AddressType:   entity.AddressType,
		RecipientName: entity.RecipientName,
		AddressLine1:  entity.AddressLine1,
		AddressLine2:  entity.AddressLine2,
		City:          entity.City,
		StateProvince: entity.StateProvince,
		PostalCode:    entity.PostalCode,
		CountryCode:   entity.CountryCode,
		Phone:         entity.Phone,
		IsDefault:     entity.IsDefault,
	}

	result, err := u.userAddressRepository.UpdateUserAddress(ctx, request)
	if err != nil {
		log.Errorf("failed to update address: %v", err)
		return nil, err
	}

	log.WithField("address_id", result.AddressID).Info("Address updated successfully")

	return &dto.UserAddressResponse{
		AddressID:     result.AddressID,
		UserID:        result.UserID,
		AddressType:   result.AddressType,
		RecipientName: result.RecipientName,
		AddressLine1:  result.AddressLine1,
		AddressLine2:  result.AddressLine2,
		City:          result.City,
		StateProvince: result.StateProvince,
		PostalCode:    result.PostalCode,
		CountryCode:   result.CountryCode,
		Phone:         result.Phone,
		IsDefault:     result.IsDefault,
	}, nil
}

func (u *UserAddressService) DeleteUserAddress(ctx context.Context, addressId int32, userId string) error {
	log := logger.WithCtx(ctx, "service", "UserAddressService.DeleteUserAddress").WithField("user_id", userId).WithField("address_id", addressId)
	request := &domain.UserAddress{
		UserID:    userId,
		AddressID: addressId,
	}

	if err := u.userAddressRepository.DeleteUserAddress(ctx, request); err != nil {
		log.Errorf("failed to delete address: %v", err)
		return err
	}

	log.Info("Address deleted successfully")
	return nil
}

func (u *UserAddressService) ListUserAddress(ctx context.Context, userId string) ([]dto.UserAddressResponse, error) {
	log := logger.WithCtx(ctx, "service", "UserAddressService.ListUserAddress").WithField("user_id", userId)
	request := &domain.UserAddress{
		UserID: userId,
	}

	result, err := u.userAddressRepository.ReadAllUserAddress(ctx, request)
	if err != nil {
		log.Errorf("failed to list addresses: %v", err)
		return nil, err
	}

	addresses := make([]dto.UserAddressResponse, len(result))
	for i, address := range result {
		addresses[i] = dto.UserAddressResponse{
			AddressID:     address.AddressID,
			UserID:        address.UserID,
			AddressType:   address.AddressType,
			RecipientName: address.RecipientName,
			AddressLine1:  address.AddressLine1,
			AddressLine2:  address.AddressLine2,
			City:          address.City,
			StateProvince: address.StateProvince,
			PostalCode:    address.PostalCode,
			CountryCode:   address.CountryCode,
			Phone:         address.Phone,
			IsDefault:     address.IsDefault,
		}
	}

	return addresses, nil
}
