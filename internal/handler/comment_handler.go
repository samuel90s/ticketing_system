package handler

import (
	"net/http"
	"strconv"

	"ticketing-system/internal/service"

	"github.com/gin-gonic/gin"
)

type AddCommentRequest struct {
	Content    string `json:"content"     binding:"required,min=1,max=5000"`
	IsInternal bool   `json:"is_internal"`
}

// GET /api/tickets/:id/comments
func GetComments(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	role, _ := getRole(c)

	comments, err := service.GetComments(uint(ticketID), role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// POST /api/tickets/:id/comments
func AddComment(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, _ := getRole(c)

	var req AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := service.AddComment(uint(ticketID), userID, req.Content, req.IsInternal, role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "comment added successfully"})
}

// DELETE /api/tickets/:id/comments/:comment_id
func DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, _ := getRole(c)

	if err := service.DeleteComment(uint(commentID), userID, role); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
}
