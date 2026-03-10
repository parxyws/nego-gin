package route

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/user"
)

func AuthRoute(router *gin.RouterGroup, controller user.AuthController, authMiddleware *jwt.GinJWTMiddleware) {
	auth := router.Group("/auth")
	{
		// Public endpoints
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
		auth.POST("/refresh-token", controller.RefreshToken) // Public — uses refresh token, not access token
		auth.POST("/forgot-password", controller.ForgotPassword)
		auth.POST("/reset-password", controller.ResetPassword)
		auth.POST("/verify-email", controller.VerifyEmail)

		// Protected endpoints
		protected := auth.Group("/")
		protected.Use(authMiddleware.MiddlewareFunc())
		{
			protected.POST("/logout", authMiddleware.LogoutHandler)
			protected.POST("/resend-verification", controller.ResendVerification)
		}
	}
}
