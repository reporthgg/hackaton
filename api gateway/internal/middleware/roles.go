package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PoliceOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("role_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Role not found in context"})
			c.Abort()
			return
		}

		if roleID != 2 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Доступ запрещен. Требуются права полиции"})
			c.Abort()
			return
		}

		c.Next()
	}
}
