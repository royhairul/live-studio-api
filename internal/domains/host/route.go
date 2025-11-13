package host

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/host/controller"
	"github.com/royhairul/live-studio-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.HostController) {
	route := router.Group("/host")
	route.Use(middleware.RequireRoles("superadmin", "admin"))
	route.Use(middleware.TenantMiddleware())

	{
		route.GET("", controller.FindAll)
		route.POST("", controller.Create)

		route.GET("/group-by-studio", controller.FindAllGroupedByStudio)

		route.GET("/:id", controller.FindByID)
		route.PUT("/:id", controller.Update)
		route.DELETE("/:id", controller.Delete)
	}
}
