package service

import (
	"context"
	"fmt"

	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain"
	"github.com/parxyws/nego-gin/internal/user/domain/dto"
	"github.com/parxyws/nego-gin/pkg/logger"
)

type RoleService struct {
	roleRepository     user.RoleRepository
	userRoleRepository user.UserRoleRepository
}

func NewRoleService(roleRepository user.RoleRepository, userRoleRepository user.UserRoleRepository) user.RoleService {
	return &RoleService{roleRepository: roleRepository, userRoleRepository: userRoleRepository}
}

func (r *RoleService) UpdateRole(ctx context.Context, entity *dto.RoleUpdateRequest) (*dto.RoleResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RoleService) RegisterRole(ctx context.Context, entity *dto.RoleRegisterRequest) (*dto.RoleResponse, error) {
	log := logger.WithCtx(ctx, "service", "RoleService.RegisterRole").WithField("role_name", entity.RoleName)
	newRole := &domain.Role{
		RoleName:    entity.RoleName,
		Description: entity.Description,
	}

	res, err := r.roleRepository.CreateRole(ctx, newRole)
	if err != nil {
		log.Errorf("failed to create role: %v", err)
		return nil, fmt.Errorf("Create role failed: %w", err)
	}

	log.WithField("role_id", res.RoleID).Info("Role registered successfully")

	return &dto.RoleResponse{
		ID:          res.RoleID,
		RoleName:    res.RoleName,
		Description: res.Description,
		CreatedAt:   res.CreatedAt,
	}, nil

}

func (r *RoleService) CheckDeleteRole(ctx context.Context, entity []dto.RoleDeleteRequest) ([]dto.RoleResponse, error) {
	log := logger.WithCtx(ctx, "service", "RoleService.CheckDeleteRole")
	roleIDs := make([]int32, len(entity))
	for idx, deleteRequest := range entity {
		roleIDs[idx] = deleteRequest.RoleId
	}

	countMap, err := r.userRoleRepository.CountUserRoleByRoleID(ctx, roleIDs)
	if err != nil {
		log.Errorf("failed to count user roles: %v", err)
		return nil, fmt.Errorf("Count user roles failed: %w", err)
	}

	roles, err := r.roleRepository.ReadAllRoleById(ctx, roleIDs)
	if err != nil {
		log.Errorf("failed to read roles by ids: %v", err)
		return nil, fmt.Errorf("Retrieve roles by IDs failed: %w", err)
	}

	responses := make([]dto.RoleResponse, 0, len(roles))
	for _, role := range roles {
		userCount := countMap[role.RoleID]

		response := dto.RoleResponse{
			ID:          role.RoleID,
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

func (r *RoleService) DeleteRole(ctx context.Context, entity []dto.RoleDeleteRequest) error {
	log := logger.WithCtx(ctx, "service", "RoleService.DeleteRole")

	rolesToDelete := make([]domain.Role, len(entity))
	for idx, deleteRequest := range entity {
		rolesToDelete[idx] = domain.Role{
			RoleID:   deleteRequest.RoleId,
			RoleName: deleteRequest.RoleName,
		}
	}

	if err := r.roleRepository.DeleteRole(ctx, rolesToDelete); err != nil {
		log.Errorf("failed to delete roles: %v", err)
		return fmt.Errorf("Delete roles failed: %w", err)
	}

	log.Info("Roles deleted successfully")
	return nil

}

func (r *RoleService) GetAllRoles(ctx context.Context) ([]dto.RoleResponse, error) {
	log := logger.WithCtx(ctx, "service", "RoleService.GetAllRoles")
	res, err := r.roleRepository.ReadAllRole(ctx)
	if err != nil {
		log.Errorf("failed to read all roles: %v", err)
		return nil, fmt.Errorf("Retrieve all roles failed: %w", err)
	}

	roles := make([]dto.RoleResponse, len(res))
	for i, role := range res {
		roles[i] = dto.RoleResponse{
			ID:          role.RoleID,
			RoleName:    role.RoleName,
			Description: role.Description,
			CreatedAt:   role.CreatedAt,
		}
	}

	return roles, nil
}
