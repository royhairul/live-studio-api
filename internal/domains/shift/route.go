package shift

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/shift/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.ShiftController) {
	// TODO: define routes
	shiftRouter := router.Group("/shift")
	{
		shiftRouter.GET("", controller.FindAll)
		shiftRouter.POST("", controller.Create)
		shiftRouter.GET("/:id", controller.FindByID)
		shiftRouter.PUT("/:id", controller.Update)
		shiftRouter.DELETE("/:id", controller.Delete)
	}
}
