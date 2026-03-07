package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/user"
	"github.com/parxyws/nego-gin/internal/user/domain/dto"
	"github.com/parxyws/nego-gin/pkg/helper"
)

type AdminController struct {
	userService user.UserService
	roleService user.RoleService
}

func NewAdminController(userService user.UserService) user.AdminController {
	return &AdminController{userService: userService}
}

func (a *AdminController) GetListOfRoles(c *gin.Context) {

	ctx, cancel := helper.GetContext(c)
	defer cancel()

	result, err := a.roleService.GetAllRoles(ctx)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, "Retrieve roles failed", nil)
	}

	helper.Success(c, http.StatusOK, "Roles list fetched successfully", result)
}

func (a *AdminController) CreateRole(c *gin.Context) {
	ctx, cancel := helper.GetContext(c)
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

func (a *AdminController) UpdateRole(c *gin.Context) {
	ctx, cancel := helper.GetContext(c)
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

func (a *AdminController) DeleteRole(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AdminController) GetListOfPermissions(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AdminController) AssignRoleToUser(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AdminController) RemoveRoleFromUser(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
