package accountads

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/accountads/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.AccountadsController) {
	// TODO: define routes
	route := router.Group("/accountads")
	{
		route.GET("", controller.FindAll)
		route.POST("", controller.Create)
		route.GET("/:id", controller.FindByID)
		route.PUT("/:id", controller.Update)
		route.DELETE("/:id", controller.Delete)
	}
}
