package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
)

// TenantMiddleware memastikan tenant_id tersedia di context request
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, _ := c.Get("tenant_id")

		var tenant string
		if t, ok := tenantID.(string); ok && t != "" {
			tenant = t
		}
		if tenant == "" {
			tenant = c.GetHeader("X-Tenant-ID")
		}

		if tenant == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "tenant id not found. please login again.",
			})
			return
		}

		c.Set("tenant_id", tenant)

		// Buat context baru dengan tenant_id
		ctx := tenantdb.AttachTenant(c.Request.Context(), tenant)
		// store request path in global fallback and in the context
		tenantdb.SetRequestPath(c.Request.URL.Path)
		ctx = context.WithValue(ctx, "path", c.Request.URL.Path)

		// Simpan ke request context agar GORM bisa akses
		c.Request = c.Request.WithContext(ctx)

		// Simpan database yang sudah aware tenant
		c.Set("DB", database.DB.WithContext(ctx))

		c.Next()
	}
}
