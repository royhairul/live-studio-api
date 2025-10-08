package live

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/live/controller"
)

func RegisterRouter(router *gin.RouterGroup, controller controller.LiveController) {
	routes := router.Group("/live")
	{
		routes.GET("/shopee", controller.GetLive)
		routes.GET("/shopee/:id/:sessionId", controller.GetLiveDetail)
	}
}
