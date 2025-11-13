package finance

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/finance/controller"
	"github.com/royhairul/live-studio-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.FinanceController) {
	route := router.Group("/finance")
	route.Use(middleware.RequireRoles("superadmin", "admin"))
	route.Use(middleware.TenantMiddleware())

	{
		route.GET("", controller.FindAll)
	}
}
