package route

import (
	"github.com/gin-gonic/gin"
	"github.com/parxyws/nego-gin/internal/user"
)

func AuthRoute(router *gin.RouterGroup, controller user.AuthController) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/verify", controller.ValidateUser)
		auth.POST("/resend-verification", controller.ResendVerification)
		auth.POST("/login", controller.Login)
		auth.POST("/forgot-password", controller.ForgotPassword)
		auth.POST("/reset-password", controller.ResetPassword)
		auth.POST("/refresh-token", controller.RefreshToken)
		//auth.POST("/logout")
	}
}
