package user

import "github.com/gin-gonic/gin"

type AuthController interface {
	Register(c *gin.Context)
	ValidateUser(c *gin.Context)
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
	ForgotPassword(c *gin.Context)
	ResetPassword(c *gin.Context)
	ResendVerification(c *gin.Context)
}

type UserController interface {
	GetCurrentUser(c *gin.Context)
	UpdateCurrentUser(c *gin.Context)
	UpdateAvatar(c *gin.Context)
	DeleteCurrentUser(c *gin.Context)
	GetUserProfile(c *gin.Context)

	GetUserListOfAddresses(c *gin.Context)
	CreateUserAddress(c *gin.Context)
	UpdateUserAddress(c *gin.Context)
	DeleteUserAddress(c *gin.Context)
	SetDefaultAddress(c *gin.Context)
}

type AdminController interface {
	GetListOfRoles(c *gin.Context)
	CreateRole(c *gin.Context)
	DeleteRole(c *gin.Context)
	GetListOfPermissions(c *gin.Context)
	AssignRoleToUser(c *gin.Context)
	RemoveRoleFromUser(c *gin.Context)
}
