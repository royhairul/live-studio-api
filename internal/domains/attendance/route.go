package attendance

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/attendance/controller"
)

func RegisterRouter(router *gin.RouterGroup, controller controller.AttendanceController) {
	routes := router.Group("/attendance")
	{
		routes.GET("", controller.FindAll)
		routes.GET("/unchecked-out", controller.FindUncheckedOut)
		routes.POST("/check-in", controller.CheckIn)
		routes.POST("/check-out", controller.CheckOut)
	}
}
