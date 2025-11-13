package user

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/user/controller"
	"github.com/royhairul/live-studio-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.UserController) {
	route := router.Group("/users")
	route.Use(middleware.TenantMiddleware())
	{
		route.GET("", controller.GetAll)
		route.POST("", controller.Create)
		route.GET("/:id", controller.GetByID)
		route.PUT("/:id", controller.Update)
		route.DELETE("/:id", controller.Delete)

	}

	superadmin := router.Group("/superadmin")
	superadmin.Use(middleware.RequireRoles("superadmin"))
	superadmin.Use(middleware.TenantMiddleware())
	{
		superadmin.GET("/user", controller.GetAll)
		superadmin.POST("/user", controller.Create)
		superadmin.GET("/user/:id", controller.Create)
		superadmin.PUT("user/:id", controller.Update)
		superadmin.DELETE("/user/:id", controller.Delete)
	}
}
