package dashboard

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/dashboard/controller"
	"github.com/royhairul/live-studio-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.DashboardController) {
	// TODO: define routes
	route := router.Group("/dashboard")
	route.Use(middleware.RequireRoles("superadmin"))
	route.Use(middleware.TenantMiddleware())

	{
		route.GET("", controller.Dashboard)
	}
}
