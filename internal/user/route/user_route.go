package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/user"
)

func UserRoute(router *gin.RouterGroup, controller user.UserController) {
	users := router.Group("/users")
	{
		users.GET("/me", controller.GetCurrentUser)
		users.PUT("/me", controller.UpdateCurrentUser)
		users.PATCH("/me/avatar", controller.UpdateAvatar)
		users.DELETE("/me", controller.DeleteCurrentUser)
		users.GET("/:userId", controller.GetUserProfile)
		users.GET("/:userId/ratings", controller.GetUserSellerRatings)

		address := users.Group("/me/addresses")
		{
			address.GET("/", controller.GetUserListOfAddresses)
			address.POST("/", controller.CreateUserAddress)
			address.PUT("/:addressId", controller.UpdateUserAddress)
			address.DELETE("/:addressId", controller.DeleteUserAddress)
			address.PATCH("/:addressId/default", controller.SetDefaultAddress)
		}
	}
}
