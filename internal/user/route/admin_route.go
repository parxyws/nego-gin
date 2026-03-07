package route

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/user"
)

func AdminRoute(router *gin.RouterGroup, controller user.AdminController, authMiddleware *jwt.GinJWTMiddleware) {
	admin := router.Group("/admin")
	admin.Use(authMiddleware.MiddlewareFunc())
	{
		// Roles management
		roles := admin.Group("/roles")
		{
			roles.GET("", controller.GetListOfRoles)
			roles.POST("", controller.CreateRole)
			roles.PUT("/:roleId", controller.UpdateRole)
			roles.DELETE("/:roleId", controller.DeleteRole)
		}

		// Permissions
		admin.GET("/permissions", controller.GetListOfPermissions)

		// User role assignment
		admin.POST("/users/:userId/roles", controller.AssignRoleToUser)
		admin.DELETE("/users/:userId/roles/:roleId", controller.RemoveRoleFromUser)
	}
}
