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
	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:4173", "http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Role"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	route.Group("/api")

	// Serve OpenAPI JSON
	route.StaticFile("/docs/openapi.json", "./docs/LiveStudio.openapi.json")

	// Serve Swagger UI
	route.GET("/docs", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", []byte(`
		<!DOCTYPE html>
		<html>
		<head>
		  <title>Live Studio API Docs</title>
		  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css" />
		</head>
		<body>
		  <div id="swagger-ui"></div>
		  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
		  <script>
		    window.onload = () => {
		      SwaggerUIBundle({
		        url: '/docs/openapi.json', // sesuaikan dengan StaticFile
		        dom_id: '#swagger-ui',
		      });
		    };
		  </script>
		</body>
		</html>
		`))
	})

	// User Management Routes
	// RegisterSuperAdminRoutes(route)
	return route
}
