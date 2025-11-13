package attendance

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/controller"
	"github.com/royhairul/live-studio-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.AttendanceController) {
	route := router.Group("/attendance")
	route.Use(middleware.RequireRoles("superadmin", "admin"))
	route.Use(middleware.TenantMiddleware())

	{
		route.GET("", controller.FindAll)
		route.GET("/unchecked-out", controller.FindUncheckedOut)
		route.POST("/check-in", controller.CheckIn)
		route.POST("/check-out", controller.CheckOut)
	}
}
