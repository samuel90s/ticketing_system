package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware — hanya admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access only"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AgentMiddleware — admin atau agent
func AgentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || (role != "admin" && role != "agent") {
			c.JSON(http.StatusForbidden, gin.H{"error": "agent or admin access only"})
			c.Abort()
			return
		}
		c.Next()
	}
}
