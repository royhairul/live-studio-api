package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
)

func CORS(ctx *gin.Context) {
	allowed := os.Getenv("ALLOWED_ORIGINS")
	if strings.TrimSpace(allowed) == "" {
		allowed = "http://localhost:5173"
	}

	allowedOrigins := map[string]struct{}{}
	for _, o := range strings.Split(allowed, ",") {
		origin := strings.TrimSpace(o)
		if origin != "" {
			allowedOrigins[origin] = struct{}{}
		}
	}

	origin := ctx.GetHeader("Origin")
	if origin != "" {
		if _, ok := allowedOrigins[origin]; ok {
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Vary", "Origin")
			ctx.Header("Access-Control-Allow-Credentials", "true")
		}
	}

	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT,PATCH, DELETE, OPTIONS")
	ctx.Header("Access-Control-Max-Age", "86400")

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
	}

	// store current request path globally as a fallback for tenantdb whitelist
	tenantdb.SetRequestPath(ctx.Request.URL.Path)

	ctx.Next()
}
