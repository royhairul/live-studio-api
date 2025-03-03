package routes

import (
	"live-studio-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()

	route.GET("/api/products/all", controllers.ProductIndex)
	route.POST("/api/products/create", controllers.ProductCreate)

	return route
}