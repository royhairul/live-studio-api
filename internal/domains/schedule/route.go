package schedule

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/schedule/controller"
	"github.com/royhairul/live-studio-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.ScheduleController) {
	// TODO: define routes
	route := router.Group("/schedule")
	route.Use(middleware.RequireRoles("superadmin"))
	route.Use(middleware.TenantMiddleware())

	{
		route.GET("", controller.FindAll)
		route.POST("", controller.Create)
		route.GET("/:id", controller.FindByID)
		route.PUT("/:id", controller.Update)
		route.DELETE("/:id", controller.Delete)

		route.GET("/scheduled", controller.FindByShiftAndDate)
	}
}
