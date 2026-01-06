package service

import (
	"context"
	"fmt"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"github.com/parxyws/nego-gin/internal/user/domain/dto"
)

type RoleServiceImpl struct {
	roleRepository     user.RoleRepository
	userRoleRepository user.UserRoleRepository
}

func NewRoleService(roleRepository user.RoleRepository, userRoleRepository user.UserRoleRepository) user.RoleService {
	return &RoleServiceImpl{roleRepository: roleRepository, userRoleRepository: userRoleRepository}
}

func (r *RoleServiceImpl) RegisterRole(ctx context.Context, entity *dto.RoleRegisterRequest) (*dto.RoleResponse, error) {
	_role := &domain.Role{
		RoleName:    entity.RoleName,
		Description: entity.Description,
	}

	res, err := r.roleRepository.CreateRole(ctx, _role)
	if err != nil {
		return nil, fmt.Errorf("RoleService.RegisterRole - %w", err)
	}

	return &dto.RoleResponse{
		ID:          res.ID,
		RoleName:    res.RoleName,
		Description: res.Description,
		CreatedAt:   res.CreatedAt,
	}, nil

}

func (r *RoleServiceImpl) CheckDeleteRole(ctx context.Context, entity []dto.RoleDeleteRequest) ([]dto.RoleResponse, error) {
	roleIDs := make([]int32, len(entity))
	for i, v := range entity {
		roleIDs[i] = v.ID
	}

	countMap, err := r.userRoleRepository.CountUserRoleByRoleID(ctx, roleIDs)
	if err != nil {
		return nil, fmt.Errorf("RoleService.CheckDeleteRole - CountUserRoleByRoleID: %w", err)
	}

	roles, err := r.roleRepository.ReadAllRoleById(ctx, roleIDs)
	if err != nil {
		return nil, fmt.Errorf("RoleService.CheckDeleteRole - FindRolesByIDs: %w", err)
	}

	responses := make([]dto.RoleResponse, 0, len(roles))
	for _, role := range roles {
		userCount := countMap[role.ID]

		response := dto.RoleResponse{
			ID:          role.ID,
			RoleName:    role.RoleName,
			Description: role.Description,
			UserCount:   userCount,
			IsUsed:      userCount > 0,
			CreatedAt:   role.CreatedAt,
		}
		responses = append(responses, response)
	}

	return responses, nil
}

func (r *RoleServiceImpl) DeleteRole(ctx context.Context, entity []dto.RoleDeleteRequest) error {

	_role := make([]domain.Role, len(entity))
	for i, v := range entity {
		_role[i] = domain.Role{
			ID:       v.ID,
			RoleName: v.RoleName,
		}
	}

	if err := r.roleRepository.DeleteRole(ctx, _role); err != nil {
		return fmt.Errorf("RoleService.DeleteRole - %w", err)
	}
	return nil

}

func (r *RoleServiceImpl) GetAllRoles(ctx context.Context) ([]dto.RoleResponse, error) {
	res, err := r.roleRepository.ReadAllRole(ctx)
	if err != nil {
		return nil, fmt.Errorf("RoleService.GetAllRoles - %w", err)
	}

	roles := make([]dto.RoleResponse, len(res))
	for i, role := range res {
		roles[i] = dto.RoleResponse{
			ID:          role.ID,
			RoleName:    role.RoleName,
			Description: role.Description,
			CreatedAt:   role.CreatedAt,
		}
	}

	return roles, nil
}
