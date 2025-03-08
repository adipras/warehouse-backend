package middleware

import (
	"net/http"
	"strings"
	"warehouse-backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware melindungi endpoint dengan JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bypass middleware jika request OPTIONS (agar CORS berjalan)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := utils.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
