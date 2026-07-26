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

		// One ended session's detail (overview + products) as a single GET —
		// the REST twin of /preview/detail. Coexists with /history/:id/sync
		// because that one is POST-only.
		history.GET("/history/:id/:sessionId", controller.GetStoredHistoryDetail)

		// Pull from Shopee into the database. Idempotent — re-running refreshes
		// figures that settle after a stream ends instead of duplicating rows.
		history.POST("/history/sync", controller.SyncAllHistory)
		history.POST("/history/:id/sync", controller.SyncHistory)

		// Backfill a multi-month range for one account. liveList/v2 caps a single
		// call at a 30-day timeDim, so this walks the range in 30-day windows.
		// Span: ?startDate=&endDate= or ?months= (back from endDate, default 3).
		history.POST("/history/:id/sync/range", controller.SyncHistoryRange)
	}
}
