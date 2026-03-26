package handler

import (
	"net/http"

	"ticketing-system/internal/service"

	"github.com/gin-gonic/gin"
)

func GetDashboard(c *gin.Context) {
	stats, err := service.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch dashboard",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}
