package service

import (
	"context"

	"github.com/parxyws/nego-gin/internal/user"
)

type PermissionServiceImpl struct {
	PermissionRepository user.PermissionRepository
}

func NewPermissionService(permissionRepository user.PermissionRepository) user.PermissionService {
	return &PermissionServiceImpl{PermissionRepository: permissionRepository}
}

func (p *PermissionServiceImpl) RegisterPermission(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (p *PermissionServiceImpl) DeletePermission(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (p *PermissionServiceImpl) GetAllPermission(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}
