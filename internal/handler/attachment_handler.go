package handler

import (
	"net/http"
	"strconv"

	"ticketing-system/internal/service"

	"github.com/gin-gonic/gin"
)

// POST /api/tickets/:id/attachments
func UploadAttachment(c *gin.Context) {
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

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	result, err := service.UploadAttachment(uint(ticketID), userID, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GET /api/tickets/:id/attachments
func GetAttachments(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	attachments, err := service.GetAttachments(uint(ticketID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get attachments"})
		return
	}

	c.JSON(http.StatusOK, attachments)
}

// GET /api/attachments/:id/download
// Route ini didaftarkan TANPA auth middleware supaya next/image & browser bisa load langsung
func DownloadAttachment(c *gin.Context) {
	attachmentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attachment id"})
		return
	}

	filePath, fileName, err := service.GetAttachmentPath(uint(attachmentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Header("Cache-Control", "public, max-age=86400")
	c.FileAttachment(filePath, fileName)
}

// DELETE /api/attachments/:id
func DeleteAttachment(c *gin.Context) {
	attachmentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid attachment id"})
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, _ := getRole(c)

	if err := service.DeleteAttachment(uint(attachmentID), userID, role); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attachment deleted successfully"})
}
