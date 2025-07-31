package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/auth/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.AuthController) {
	routes := router.Group("/auth")
	{
		routes.POST("/login", controller.Login)
		routes.POST("/register", controller.Register)
		routes.POST("/forgot-password", controller.ForgotPassword)
		routes.POST("/verify-otp", controller.VerifyOtp)
		routes.POST("/reset-password", controller.ResetPassword)
	}
}
