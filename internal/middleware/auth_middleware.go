package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/royhairul/live-studio-api/internal/pkg/utils"
)

func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized: missing or invalid token",
			})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		tokenStr = strings.TrimSpace(tokenStr)

		claims, err := utils.VerifyTokenJWT(tokenStr, os.Getenv("JWT_SECRET"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		userRole := claims.Role
		isAuthorized := false
		for _, role := range roles {
			if userRole == role {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden: insufficient permissions"})
			return
		}

		c.Set("user_id", claims.ID)
		c.Set("name", claims.Name)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		// ✅ Tambahkan tenant_id dari JWT jika tersedia
		if claims.TenantID != "" {
			c.Set("tenant_id", claims.TenantID)
			ctx := context.WithValue(c.Request.Context(), "tenant_id", claims.TenantID)
			c.Request = c.Request.WithContext(ctx)
		}

		log.Print(claims.TenantID)

		c.Next()
	}
}
