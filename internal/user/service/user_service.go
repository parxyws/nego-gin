package service

import (
	"context"

	"github.com/parxyws/nego-gin/config"
	"github.com/parxyws/nego-gin/internal/user"
)

type UserServiceImpl struct {
	cfg            *config.Config
	userRepository user.UserRepository
}

func NewUserService(cfg *config.Config, userRepository user.UserRepository) user.UserService {
	return &UserServiceImpl{cfg: cfg, userRepository: userRepository}
}

func (u *UserServiceImpl) ReadCurrentUser(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (u *UserServiceImpl) ReadAllUser(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}
