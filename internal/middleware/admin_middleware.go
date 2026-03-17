package middleware

import (
	"net/http"

	"ticketing-system/internal/service"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		user, err := service.GetUserByID(userID.(uint))
		if err != nil || user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access only"})
			c.Abort()
			return
		}

		c.Next()
	}
}
