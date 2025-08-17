package performa

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/performa/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.PerformaController) {
	// TODO: define routes
	route := router.Group("/performa")
	{
		route.GET("/host", controller.GetHosts)
		route.GET("/host/:id", controller.GetHostByID)
		route.GET("/studio", controller.GetStudios)
		route.GET("/studio/:id", controller.GetStudioByID)
	}
}
