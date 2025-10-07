package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	helpers "github.com/royhairul/live-studio-api/internal/pkg/utils"
)

func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized", "error": c.GetHeader("Authorization")})
			return
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		// Verifikasi JWT
		claims, err := helpers.VerifyTokenJWT(tokenStr, os.Getenv("JWT_SECRET"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		// Cek apakah role cocok
		userRole := claims.Role
		isAuthorized := false
		for _, r := range roles {
			if r == userRole {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden: insufficient permissions"})
			return
		}

		// Simpan ke context (jika ingin digunakan di handler berikutnya)
		c.Set("name", claims.Name)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		if claims.Role == "superadmin" {
			c.Set("superadmin_id", claims.ID)
		} else {
			c.Set("user_id", claims.ID)
		}

		// Lanjut ke handler berikutnya
		c.Next()
	}
}
