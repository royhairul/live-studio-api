package dashboard

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/dashboard/controller"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.DashboardController) {
	// TODO: define routes
	dashboardRouter := router.Group("/dashboard")
	{
		dashboardRouter.GET("", controller.Dashboard)
	}
}
