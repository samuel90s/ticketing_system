package handler

import (
	"strconv"
	"ticketing-system/internal/service"

	"github.com/gin-gonic/gin"
)

func GetAllTickets(c *gin.Context) {
	// ambil query param
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.DefaultQuery("search", "")

	tickets, err := service.GetAllTickets(page, limit, search)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed"})
		return
	}

	c.JSON(200, tickets)
}
