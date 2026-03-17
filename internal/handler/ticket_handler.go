package handler

import (
	"net/http"
	"strconv"

	"ticketing-system/internal/service"

	"github.com/gin-gonic/gin"
)

//
// ======================
// REQUEST STRUCTS
// ======================
//

type CreateTicketRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type AssignTicketRequest struct {
	AssigneeID uint `json:"assignee_id" binding:"required"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"` // open / closed
}

// ======================
// CREATE TICKET
// ======================
func CreateTicket(c *gin.Context) {
	var req CreateTicketRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	err := service.CreateTicket(req.Title, req.Description, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket created"})
}

// ======================
// GET USER TICKETS
// ======================
func GetTickets(c *gin.Context) {
	userID, _ := c.Get("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.DefaultQuery("search", "")

	tickets, err := service.GetTicketsByUser(userID.(uint), page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tickets"})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// ======================
// ASSIGN TICKET (ADMIN)
// ======================
func AssignTicket(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	var req AssignTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = service.AssignTicket(uint(ticketID), req.AssigneeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket assigned"})
}

// ======================
// UPDATE STATUS
// ======================
func UpdateTicketStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	ticketID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = service.UpdateTicketStatus(
		uint(ticketID),
		userID.(uint),
		role.(string),
		req.Status,
	)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}
