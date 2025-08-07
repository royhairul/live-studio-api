package studio

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/studio/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.StudioController) {
	// TODO: define routes
	studioRouter := router.Group("/studio")
	{
		studioRouter.GET("", controller.FindAll)
		studioRouter.POST("", controller.Create)
		studioRouter.GET("/:id", controller.FindByID)
		studioRouter.PUT("/:id", controller.Update)
		studioRouter.DELETE("/:id", controller.Delete)
	}
}
