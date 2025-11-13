package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/auth/controller"
	"github.com/royhairul/live-studio-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.AuthController) {
	route := router.Group("/auth")

	{
		route.POST("/login", controller.Login)
		route.POST("/register", controller.Register)
		route.POST("/forgot-password", controller.ForgotPassword)
		route.POST("/verify-otp", controller.VerifyOtp)
		route.POST("/reset-password", controller.ResetPassword)

		route.GET("/me", middleware.RequireRoles("superadmin", "admin", "host"), controller.Me)
	}
}
