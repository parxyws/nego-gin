package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/user"
)

type AdminControllerImpl struct {
	userService user.UserService
	roleService user.RoleService
}

func NewAdminController(userService user.UserService) user.AdminController {
	return &AdminControllerImpl{userService: userService}
}

func (a *AdminControllerImpl) GetListOfRoles(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (a *AdminControllerImpl) CreateRole(c *gin.Context) {
	//TODO implement me
	panic("implement me")
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
