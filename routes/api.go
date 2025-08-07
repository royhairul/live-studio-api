package routes

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var RouteModule = fx.Module(
	"router",
	fx.Provide(SetupRouter),
)

func GroupAPI(router *gin.Engine) *gin.RouterGroup {
	return router.Group("/api")
}

func SetupRouter() *gin.Engine {
	route := gin.Default()

	route.Group("/api")

	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:4173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Role"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// User Management Routes
	RegisterSuperAdminRoutes(route)
	return route
}
