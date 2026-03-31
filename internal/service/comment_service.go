package service

import (
	"errors"
	"strings"

	"ticketing-system/internal/config"
	"ticketing-system/internal/model"
)

// ======================
// DTO
// ======================

type CommentResponse struct {
	ID         uint   `json:"id"`
	TicketID   uint   `json:"ticket_id"`
	UserID     uint   `json:"user_id"`
	User       string `json:"user"`
	Role       string `json:"role"`
	Content    string `json:"content"`
	IsInternal bool   `json:"is_internal"`
	CreatedAt  string `json:"created_at"`
}

func mapComment(c model.Comment) CommentResponse {
	return CommentResponse{
		ID:         c.ID,
		TicketID:   c.TicketID,
		UserID:     c.UserID,
		User:       c.User.Name,
		Role:       c.User.Role,
		Content:    c.Content,
		IsInternal: c.IsInternal,
		CreatedAt:  c.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ======================
// GET COMMENTS
// ======================

func GetComments(ticketID uint, role string) ([]CommentResponse, error) {
	var comments []model.Comment

	query := config.DB.Preload("User").Where("ticket_id = ?", ticketID)

	// user biasa tidak bisa lihat komentar internal
	if role == "user" {
		query = query.Where("is_internal = ?", false)
	}

	err := query.Order("created_at ASC").Find(&comments).Error
	if err != nil {
		return nil, err
	}

	result := make([]CommentResponse, 0, len(comments))
	for _, c := range comments {
		result = append(result, mapComment(c))
	}
	return result, nil
}

// ======================
// ADD COMMENT
// ======================

func AddComment(ticketID, userID uint, content string, isInternal bool, role string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("comment content cannot be empty")
	}
	if len(content) > 5000 {
		return errors.New("comment too long, max 5000 characters")
	}

	// hanya agent/admin yang bisa buat komentar internal
	if isInternal && role == "user" {
		return errors.New("only agent or admin can post internal comments")
	}

	// pastikan ticket ada
	var ticket model.Ticket
	if err := config.DB.First(&ticket, ticketID).Error; err != nil {
		return errors.New("ticket not found")
	}

	comment := model.Comment{
		TicketID:   ticketID,
		UserID:     userID,
		Content:    content,
		IsInternal: isInternal,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		return err
	}

	// Log history
	action := "commented"
	if isInternal {
		action = "internal_comment"
	}
	logHistory(ticketID, userID, action, "", "")

	// Jika agent/admin reply → ubah status ke "replied" jika masih "open"
	if (role == "agent" || role == "admin") && !isInternal && ticket.Status == "open" {
		config.DB.Model(&ticket).Update("status", "replied")
		logHistory(ticketID, userID, "status_changed", "open", "replied")
	}

	return nil
}

// ======================
// DELETE COMMENT
// ======================

func DeleteComment(commentID, userID uint, role string) error {
	var comment model.Comment
	if err := config.DB.First(&comment, commentID).Error; err != nil {
		return errors.New("comment not found")
	}

	// hanya admin atau pemilik komentar yang bisa hapus
	if role != "admin" && comment.UserID != userID {
		return errors.New("not authorized to delete this comment")
	}

	return config.DB.Delete(&comment).Error
}
