package host

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/host/controller"
)

func RegisterRouter(router *gin.RouterGroup, controller controller.HostController) {
	routes := router.Group("/host")
	{
		routes.GET("", controller.FindAll)
		routes.POST("", controller.Create)

		routes.GET("/group-by-studio", controller.FindAllGroupedByStudio)

		routes.GET("/:id", controller.FindByID)
		routes.PUT("/:id", controller.Update)
		routes.DELETE("/:id", controller.Delete)
	}
}
