package user

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/user/controller"
	"github.com/royhairul/live-studio-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.UserController) {
	routes := router.Group("/users")
	{
		routes.GET("", controller.GetAll)
		routes.POST("", controller.Create)
		routes.GET("/:id", controller.GetByID)
		routes.PUT("/:id", controller.Update)
		routes.DELETE("/:id", controller.Delete)

	}

	superadmin := router.Group("/superadmin")
	superadmin.Use(middleware.RequireRoles("superadmin"))
	{
		superadmin.GET("/user", controller.GetAll)
		superadmin.POST("/user", controller.Create)
		superadmin.GET("/user/:id", controller.Create)
		superadmin.PUT("user/:id", controller.Update)
		superadmin.DELETE("/user/:id", controller.Delete)
	}
}
