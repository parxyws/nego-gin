package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain/dto"
	"github.com/parxyws/nego-gin/pkg/helper"
)

type AdminControllerImpl struct {
	userService user.UserService
	roleService user.RoleService
}

func NewAdminController(userService user.UserService) user.AdminController {
	return &AdminControllerImpl{userService: userService}
}

func (a *AdminControllerImpl) GetListOfRoles(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := a.roleService.GetAllRoles(ctx)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "failed to get roles", nil)
	}

	helper.Success(c, http.StatusOK, "success to get roles list", result)
}

func (a *AdminControllerImpl) CreateRole(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body := new(dto.RoleRegisterRequest)
	if err := c.ShouldBind(body); err != nil {
		helper.Error(c, http.StatusBadRequest, "failed to bind body", err)
	}

	result, err := a.roleService.RegisterRole(ctx, body)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "failed to register role", err)
	}

	helper.Success(c, http.StatusOK, "success to register role", result)
}

func (a *AdminControllerImpl) UpdateRole(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AdminControllerImpl) DeleteRole(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AdminControllerImpl) GetListOfPermissions(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AdminControllerImpl) AssignRoleToUser(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AdminControllerImpl) RemoveRoleFromUser(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
