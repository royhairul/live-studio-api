package live

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/live/controller"
	"github.com/royhairul/live-studio-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.LiveController) {
	route := router.Group("/live")

	{
		// Preview: WebSocket, ongoing sessions only. Auth is handled inside the
		// handler via ?token= because browsers cannot set headers on a
		// WebSocket handshake.
		route.GET("/preview", controller.GetLive)
		route.GET("/preview/detail/:id/:sessionId", controller.GetLiveDetail)
	}

	// History: plain REST over the lives table, so it uses header auth.
	history := route.Group("")
	history.Use(middleware.RequireRoles("superadmin", "admin"))
	history.Use(middleware.TenantMiddleware())

	{
		history.GET("/history", controller.GetStoredHistory)
		history.GET("/history/:id", controller.GetStoredHistory)

		// Pull from Shopee into the database. Idempotent — re-running refreshes
		// figures that settle after a stream ends instead of duplicating rows.
		history.POST("/history/sync", controller.SyncAllHistory)
		history.POST("/history/:id/sync", controller.SyncHistory)
	}
}
