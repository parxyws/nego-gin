package route

import (
	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/user"
)

func AuthRoute(router *gin.RouterGroup, controller user.AuthController, authMiddleware *jwt.GinJWTMiddleware) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/verify-email", controller.VerifyEmail)
		auth.POST("/resend-verification", controller.ResendVerification)
		auth.POST("/login", controller.Login)
		auth.POST("/forgot-password", controller.ForgotPassword)
		auth.POST("/reset-password", controller.ResetPassword)

		protected := router.Group("")
		protected.Use(authMiddleware.MiddlewareFunc())
		protected.POST("/refresh-token", controller.RefreshToken)
		protected.POST("/logout", authMiddleware.LogoutHandler)
	}
}
