package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/database"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
	"github.com/royhairul/live-studio-api/internal/pkg/utils"
)

// Mapping role name → role_id
var RoleMap = map[string]uint{
	"superadmin": 1,
	"admin":      2,
	"host":       3,
}

func RequireRoles(roleNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized: missing or invalid token",
			})
			return
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		// Verifikasi JWT
		claims, err := utils.VerifyTokenJWT(tokenStr, os.Getenv("JWT_SECRET"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		userRoleID := claims.Role

		// Konversi parameter roleName → roleID
		allowedRoles := make([]uint, 0)
		for _, rn := range roleNames {
			if id, ok := RoleMap[rn]; ok {
				allowedRoles = append(allowedRoles, id)
			}
		}

		// Validasi
		authorized := false
		for _, roleID := range allowedRoles {
			if userRoleID == roleID {
				authorized = true
				break
			}
		}

		if !authorized {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Forbidden: insufficient role permissions",
			})
			return
		}

		// Set context
		c.Set("claims", claims)
		c.Set("user_id", claims.ID)
		c.Set("name", claims.Name)
		c.Set("role_id", claims.Role)

		if claims.TenantID != "" {
			// store tenant in gin context map
			c.Set("tenant_id", claims.TenantID)

			// attach tenant to the request context using tenantdb helper
			ctx := tenantdb.AttachTenant(c.Request.Context(), claims.TenantID)

			// store request path for tenantdb whitelist checks (fallback)
			tenantdb.SetRequestPath(c.Request.URL.Path)
			ctx = context.WithValue(ctx, "path", c.Request.URL.Path)

			// update request context so GORM callbacks can extract tenant and path
			c.Request = c.Request.WithContext(ctx)

			// store DB with context so handlers can use tenant-aware DB
			c.Set("DB", database.DB.WithContext(ctx))
		}

		c.Next()
	}
}
