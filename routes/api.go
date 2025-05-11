package routes

import (
	"live-studio-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()

	// User Routes
	route.POST("/api/login", controllers.Login)
	route.POST("/api/register", controllers.Register)

	// Product Routes
	route.GET("/api/products/all", controllers.ProductIndex)
	route.POST("/api/products/create", controllers.ProductCreate)

	// Akun Routes
	route.GET("/api/akun/all", controllers.AkunIndex)
	route.POST("/api/akun/create", controllers.AkunCreate)

	// Host Routes
	route.GET("/api/host/all", controllers.HostIndex)
	route.GET("/api/host/:id", controllers.HostShow)
	route.POST("/api/host/create", controllers.HostCreate)
	route.PUT("/api/host/:id/update", controllers.HostUpdate)
	route.DELETE("/api/host/:id/delete", controllers.HostDelete)

	// Studio Routes
	route.GET("/api/studio/all", controllers.StudioIndex)
	route.POST("/api/studio/create", controllers.StudioCreate)

	return route
}