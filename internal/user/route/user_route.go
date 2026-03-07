package route

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/user"
)

func UserRoute(router *gin.RouterGroup, controller user.UserController, authMiddleware *jwt.GinJWTMiddleware) {
	users := router.Group("/users")
	{
		// Public endpoints — visible to everyone
		users.GET("/:userId", controller.GetUserProfile)
		users.GET("/:userId/ratings", controller.GetUserSellerRatings)

		// Protected endpoints — require authentication
		protected := users.Group("/me")
		protected.Use(authMiddleware.MiddlewareFunc())
		{
			protected.GET("", controller.GetCurrentUser)
			protected.PUT("", controller.UpdateCurrentUser)
			protected.PATCH("/avatar", controller.UpdateAvatar)
			protected.DELETE("", controller.DeleteCurrentUser)

			// Addresses sub-group
			addresses := protected.Group("/addresses")
			{
				addresses.GET("", controller.ListUserAddress)
				addresses.POST("", controller.CreateUserAddress)
				addresses.PUT("/:addressId", controller.UpdateUserAddress)
				addresses.DELETE("/:addressId", controller.DeleteUserAddress)
				addresses.PATCH("/:addressId/default", controller.SetDefaultAddress)
			}
		}
	}
}
