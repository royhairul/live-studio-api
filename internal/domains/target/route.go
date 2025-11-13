package target

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/target/controller"
	"github.com/royhairul/live-studio-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.TargetController) {
	// TODO: define route
	route := router.Group("/target")
	route.Use(middleware.RequireRoles("superadmin", "admin"))
	route.Use(middleware.TenantMiddleware())

	{
		route.GET("", controller.FindAll)
		route.POST("", controller.Create)
		route.GET("/:id", controller.FindByID)
		route.PUT("/:id", controller.Update)
		route.DELETE("/:id", controller.Delete)
	}
}
