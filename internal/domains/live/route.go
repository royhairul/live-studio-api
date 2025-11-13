package live

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/live/controller"
	"github.com/royhairul/live-studio-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.LiveController) {
	route := router.Group("/live")
	route.Use(middleware.RequireRoles("superadmin", "admin"))
	route.Use(middleware.TenantMiddleware())

	{
		route.GET("/shopee", controller.GetLive)
		route.GET("/shopee/:id/:sessionId", controller.GetLiveDetail)
	}
}
