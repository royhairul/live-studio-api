package performa

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/domains/performa/controller"
	"github.com/royhairul/live-studio-api/internal/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, controller controller.PerformaController) {
	// TODO: define routes
	route := router.Group("/performa")
	route.Use(middleware.RequireRoles("superadmin", "admin"))

	{
		route.GET("/host", controller.GetHosts)
		route.GET("/host/:id", controller.GetHostByID)

		route.GET("/account", controller.GetAccounts)

		route.GET("/studio", controller.GetStudios)
		route.GET("/studio/:id", controller.GetStudioByID)
	}
}
