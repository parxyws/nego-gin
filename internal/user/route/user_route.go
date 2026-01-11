package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/config"
)

func UserRoute(route *gin.Engine, cfg *config.Config) {
	user := route.Group("/users")
	{
		user.GET("/me")
		user.PUT("/me")
		user.PATCH("/me/avatar")
		user.DELETE("/me")
		user.GET("/:userId")
		user.GET("/:userId/ratings")

		address := user.Group("/me/addresses")
		{
			address.GET("/")
			address.POST("/")
			address.PUT("/:addressId")
			address.DELETE("/:addressId")
			address.PATCH("/:addressId/default")
		}
	}
}
