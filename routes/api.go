package routes

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/middleware"
	"github.com/royhairul/live-studio-api/internal/pkg/response"
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

	api := route.Group("/api")

	api.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.NewBaseResponse("Everything OK!", nil))
	})

	// Setup documentation routes with security
	setupDocsRoutes(route)

	return route
}

// setupDocsRoutes configures the API documentation routes based on APP_ENV
func setupDocsRoutes(route *gin.Engine) {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	// Disable docs in production and maintenance modes
	if appEnv == "production" || appEnv == "maintenance" {
		route.GET("/docs", func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Documentation is not available in " + appEnv + " mode"})
		})
		route.GET("/docs/openapi.json", func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Documentation is not available in " + appEnv + " mode"})
		})
		return
	}

	// Create docs group
	docsGroup := route.Group("/docs")

	// In staging mode, require authentication using superadmin credentials
	if appEnv == "staging" {
		superadminName := os.Getenv("SUPERADMIN_NAME")
		superadminPass := os.Getenv("SUPERADMIN_PASS")

		if superadminName != "" && superadminPass != "" {
			docsGroup.Use(docsBasicAuth(superadminName, superadminPass))
		}
	}
	// In development mode, docs are open without authentication

	// Serve OpenAPI JSON dynamically with server URL from environment
	docsGroup.GET("/openapi.json", serveOpenAPIWithDynamicServer)

	// Serve Swagger UI
	docsGroup.GET("", serveSwaggerUI)
}

// docsBasicAuth returns a Basic Auth middleware for documentation
func docsBasicAuth(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, hasAuth := c.Request.BasicAuth()
		
		if !hasAuth {
			c.Header("WWW-Authenticate", `Basic realm="API Documentation"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Authentication required to access documentation",
			})
			return
		}
		
		// Use constant-time comparison to prevent timing attacks
		userMatch := subtle.ConstantTimeCompare([]byte(user), []byte(username)) == 1
		passMatch := subtle.ConstantTimeCompare([]byte(pass), []byte(password)) == 1
		
		if !userMatch || !passMatch {
			c.Header("WWW-Authenticate", `Basic realm="API Documentation"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid credentials",
			})
			return
		}
		
		c.Next()
	}
}

// serveSwaggerUI serves the Swagger UI HTML page
func serveSwaggerUI(c *gin.Context) {
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
	        url: '/docs/openapi.json',
	        dom_id: '#swagger-ui',
	      });
	    };
	  </script>
	</body>
	</html>
	`))
}

// serveOpenAPIWithDynamicServer reads the OpenAPI JSON and injects the server URL from environment
func serveOpenAPIWithDynamicServer(c *gin.Context) {
	// Read the OpenAPI JSON file
	data, err := os.ReadFile("./docs/LiveStudio.openapi.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read OpenAPI file"})
		return
	}

	// Parse JSON
	var openAPI map[string]interface{}
	if err := json.Unmarshal(data, &openAPI); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse OpenAPI file"})
		return
	}

	// Get environment status
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	// Get API_URL from environment variable
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		// Fallback to constructing from SERVER_HOST and SERVER_PORT
		host := os.Getenv("SERVER_HOST")
		port := os.Getenv("SERVER_PORT")
		if host == "" {
			host = "localhost"
		}
		if port == "" {
			port = "8080"
		}
		apiURL = "http://" + host + ":" + port + "/api"
	}

	// Update info with environment status in description only
	if info, ok := openAPI["info"].(map[string]interface{}); ok {
		// Keep original title, add environment to description
		originalDesc := ""
		if desc, ok := info["description"].(string); ok {
			originalDesc = desc
		}
		info["description"] = "**Environment:** `" + appEnv + "`\n\n" + originalDesc
	}

	// Add environment as a tag at the beginning
	envBadge := getEnvBadge(appEnv)
	envTag := map[string]interface{}{
		"name":        envBadge,
		"description": "Current environment: " + appEnv,
	}

	// Prepend environment tag to existing tags
	if tags, ok := openAPI["tags"].([]interface{}); ok {
		openAPI["tags"] = append([]interface{}{envTag}, tags...)
	} else {
		openAPI["tags"] = []interface{}{envTag}
	}

	// Update servers with dynamic URL and environment
	openAPI["servers"] = []map[string]interface{}{
		{
			"url":         apiURL,
			"description": appEnv + " server",
		},
	}

	// Return modified JSON
	c.JSON(http.StatusOK, openAPI)
}

// getEnvBadge returns an emoji badge for the environment
func getEnvBadge(env string) string {
	switch env {
	case "production":
		return "🟢 Production"
	case "staging":
		return "🟡 Staging"
	case "development":
		return "🔵 Development"
	case "maintenance":
		return "🔴 Maintenance"
	default:
		return "⚪ " + env
	}
}
