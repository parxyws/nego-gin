package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/config"
)

func AdminRoute(router *gin.Engine, cfg *config.Config) {
	admin := router.Group("/admin")
	{
		role := admin.Group("/role")
		{
			role.GET("/")
			role.POST("/")
			role.PUT("/:roleId")
			role.DELETE("/:roleId")
		}

		admin.GET("permissions")
		admin.POST("/users/:userId/roles")
		admin.DELETE("/users/:userId/roles/:roleId")
	}
}
