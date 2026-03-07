package service

import (
	"context"

	"github.com/parxyws/nego-gin/internal/user"
)

type PermissionService struct {
	PermissionRepository user.PermissionRepository
}

func NewPermissionService(permissionRepository user.PermissionRepository) user.PermissionService {
	return &PermissionService{PermissionRepository: permissionRepository}
}

func (p *PermissionService) RegisterPermission(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (p *PermissionService) DeletePermission(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (p *PermissionService) GetAllPermission(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}
