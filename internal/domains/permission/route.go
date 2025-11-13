package permission

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/permission/controller"
	"github.com/royhairul/live-studio-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.PermissionController) {
	// TODO: define routes
	route := router.Group("/permission")
	route.Use(middleware.RequireRoles("superadmin"))

	{
		route.GET("", controller.FindAll)
		route.GET("/grouped", controller.FindAllGrouped)
		route.POST("", controller.Create)
		route.GET("/:id", controller.FindByID)
		route.PUT("/:id", controller.Update)
		route.DELETE("/:id", controller.Delete)
	}
}
