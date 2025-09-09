package target

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/target/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.TargetController) {
	// TODO: define routes
	targetRouter := router.Group("/target")
	{
		targetRouter.GET("", controller.FindAll)
		targetRouter.POST("", controller.Create)
		targetRouter.GET("/:id", controller.FindByID)
		targetRouter.PUT("/:id", controller.Update)
		targetRouter.DELETE("/:id", controller.Delete)
	}
}
