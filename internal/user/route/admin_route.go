package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/user"
)

func AdminRoute(router *gin.RouterGroup, controller user.AdminController) {
	admin := router.Group("/admin")
	{
		role := admin.Group("/role")
		{
			role.GET("/", controller.GetListOfRoles)
			role.POST("/", controller.CreateRole)
			//role.PUT("/:roleId")
			role.DELETE("/:roleId", controller.DeleteRole)
		}

		admin.GET("permissions", controller.GetListOfPermissions)
		admin.POST("/users/:userId/roles", controller.AssignRoleToUser)
		admin.DELETE("/users/:userId/roles/:roleId", controller.RemoveRoleFromUser)
	}
}
