package service

import (
	"context"

	"github.com/parxyws/nego-gin/config"
	"github.com/parxyws/nego-gin/internal/middleware"
	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"github.com/parxyws/nego-gin/internal/user/domain/dto"
	"github.com/parxyws/nego-gin/pkg/database/aws"
	"github.com/sirupsen/logrus"
)

type UserServiceImpl struct {
	cfg            *config.Config
	userRepository user.UserRepository
}

func NewUserService(cfg *config.Config, userRepository user.UserRepository) user.UserService {
	return &UserServiceImpl{cfg: cfg, userRepository: userRepository}
}

func (u *UserServiceImpl) GetCurrentUser(ctx context.Context, entity middleware.JwtPayload) (*dto.UserResponse, error) {
	useRequest := &domain.User{
		UserID:   entity.ID,
		Username: entity.Username,
		Email:    entity.Email,
	}

	result, err := u.userRepository.ReadById(ctx, useRequest)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "UserService.GetCurrentUser", "user_id": entity.ID}).Errorf("failed to read user: %v", err)
		return nil, err
	}

	return &dto.UserResponse{
		UserID:      result.UserID,
		Username:    result.Username,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Avatar:      result.AvatarURL,
		Email:       result.Email,
		PhoneNumber: result.Phone,
	}, nil
}

func (u *UserServiceImpl) UpdateCurrentUser(ctx context.Context, entity middleware.JwtPayload, request *dto.UpdateCurrentUserRequest) (*dto.UserResponse, error) {
	userRequest := &domain.User{
		UserID:    entity.ID,
		Username:  request.Username,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Phone:     request.PhoneNumber,
	}

	result, err := u.userRepository.UpdateUser(ctx, userRequest)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "UserService.UpdateCurrentUser", "user_id": entity.ID}).Errorf("failed to update user: %v", err)
		return nil, err
	}

	logrus.WithFields(logrus.Fields{"function": "UserService.UpdateCurrentUser", "user_id": entity.ID}).Info("User updated successfully")

	return &dto.UserResponse{
		UserID:      result.UserID,
		Username:    result.Username,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Avatar:      result.AvatarURL,
		Email:       result.Email,
		PhoneNumber: result.Phone,
	}, nil
}

func (u *UserServiceImpl) UpdateAvatar(ctx context.Context, entity middleware.JwtPayload, avatar *aws.UploadInput) (*dto.UserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *UserServiceImpl) DeleteCurrentUser(ctx context.Context, entity middleware.JwtPayload) error {
	userRequest := &domain.User{
		UserID:   entity.ID,
		Username: entity.Username,
		Email:    entity.Email,
	}

	err := u.userRepository.DeleteUser(ctx, userRequest)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "UserService.DeleteCurrentUser", "user_id": entity.ID}).Errorf("failed to delete user: %v", err)
		return err
	}

	logrus.WithFields(logrus.Fields{"function": "UserService.DeleteCurrentUser", "user_id": entity.ID}).Info("User deleted successfully")
	return nil
}

func (u *UserServiceImpl) GetUser(ctx context.Context, request *dto.GetUserProfileRequest) (*dto.UserProfileResponse, error) {
	userRequest := &domain.User{
		UserID: request.UserID,
	}

	result, err := u.userRepository.ReadById(ctx, userRequest)
	if err != nil {
		logrus.WithFields(logrus.Fields{"function": "UserService.GetUser", "user_id": request.UserID}).Errorf("failed to read user profile: %v", err)
		return nil, err
	}

	return &dto.UserProfileResponse{
		UserID:    result.UserID,
		Username:  result.Username,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Avatar:    result.AvatarURL,
		Email:     result.Email,
	}, nil
}
