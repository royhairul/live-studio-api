package schedule

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {
	controller := ProvideScheduleController()

	// TODO: define routes
	route := router.Group("/schedule")
	{
		route.GET("", controller.FindAll)
		route.POST("", controller.Create)
		route.GET("/:id", controller.FindByID)
		route.PUT("/:id", controller.Update)
		route.DELETE("/:id", controller.Delete)

		route.GET("/scheduled", controller.FindByShiftAndDate)
	}
}
