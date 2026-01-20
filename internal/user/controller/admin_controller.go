package controller

import (
	"context"
	"net/http"
	"strconv"
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
		helper.Error(c, http.StatusInternalServerError, "Retrieve roles failed", nil)
	}

	helper.Success(c, http.StatusOK, "Roles list fetched successfully", result)
}

func (a *AdminControllerImpl) CreateRole(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body := new(dto.RoleRegisterRequest)
	if err := c.ShouldBind(body); err != nil {
		helper.Error(c, http.StatusBadRequest, "Parse request body failed", err)
	}

	result, err := a.roleService.RegisterRole(ctx, body)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Register role failed", err)
	}

	helper.Success(c, http.StatusOK, "Role registered successfully", result)
}

func (a *AdminControllerImpl) UpdateRole(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	param := c.Param("roleId")
	body := new(dto.RoleUpdateRequest)

	paramInt, err := strconv.ParseInt(param, 10, 32)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, "Validation failed: invalid role identifier", err)
	}
	body.RoleId = int32(paramInt)
	if err := c.ShouldBind(body); err != nil {
		helper.Error(c, http.StatusBadRequest, "Parse request body failed", err)
	}

	result, err := a.roleService.UpdateRole(ctx, body)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Update role failed", err)
	}

	helper.Success(c, http.StatusOK, "Role updated successfully", result)
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
