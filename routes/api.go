package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/middleware"
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
	route.Use(middleware.CORS)

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
