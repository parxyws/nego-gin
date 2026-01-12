package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/user"
)

func UserRoute(route *gin.Engine, controller user.UserController) {
	users := route.Group("/users")
	{
		users.GET("/me", controller.GetCurrentUser)
		users.PUT("/me")
		users.PATCH("/me/avatar")
		users.DELETE("/me")
		users.GET("/:userId")
		users.GET("/:userId/ratings")

		address := users.Group("/me/addresses")
		{
			address.GET("/")
			address.POST("/")
			address.PUT("/:addressId")
			address.DELETE("/:addressId")
			address.PATCH("/:addressId/default")
		}
	}
}
