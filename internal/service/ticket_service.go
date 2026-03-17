package service

import (
	"errors"
	"strings"

	"ticketing-system/internal/config"
	"ticketing-system/internal/model"
)

// ======================
// CREATE TICKET
// ======================
func CreateTicket(title, description string, userID uint) error {
	if title == "" || description == "" {
		return errors.New("title and description are required")
	}

	ticket := model.Ticket{
		Title:       title,
		Description: description,
		UserID:      userID,
		Status:      "open",
	}

	return config.DB.Create(&ticket).Error
}

// ======================
// GET USER TICKETS (Pagination + Search)
// ======================
func GetTicketsByUser(userID uint, page, limit int, search string) ([]model.Ticket, error) {
	var tickets []model.Ticket

	// safety pagination
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	query := config.DB.Where("user_id = ?", userID)

	if search != "" {
		search = strings.ToLower(search)
		query = query.Where("LOWER(title) LIKE ?", "%"+search+"%")
	}

	err := query.
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&tickets).Error

	return tickets, err
}

// ======================
// GET ALL TICKETS (ADMIN)
// ======================
func GetAllTickets(page, limit int, search string) ([]model.Ticket, error) {
	var tickets []model.Ticket

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	query := config.DB.Model(&model.Ticket{})

	if search != "" {
		search = strings.ToLower(search)
		query = query.Where("LOWER(title) LIKE ?", "%"+search+"%")
	}

	err := query.
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&tickets).Error

	return tickets, err
}

// ======================
// UPDATE STATUS
// ======================
func UpdateTicketStatus(ticketID uint, userID uint, role string, status string) error {
	var ticket model.Ticket

	if err := config.DB.First(&ticket, ticketID).Error; err != nil {
		return errors.New("ticket not found")
	}

	// authorization check
	if ticket.UserID != userID &&
		(ticket.AssigneeID == nil || *ticket.AssigneeID != userID) &&
		role != "admin" {
		return errors.New("not authorized to update ticket")
	}

	// validasi status
	status = strings.ToLower(status)
	if status != "open" && status != "closed" {
		return errors.New("invalid status")
	}

	ticket.Status = status

	return config.DB.Save(&ticket).Error
}

// ======================
// ASSIGN TICKET (ADMIN)
// ======================
func AssignTicket(ticketID uint, assigneeID uint) error {
	var ticket model.Ticket

	// cek ticket
	if err := config.DB.First(&ticket, ticketID).Error; err != nil {
		return errors.New("ticket not found")
	}

	// cek user tujuan ada
	var user model.User
	if err := config.DB.First(&user, assigneeID).Error; err != nil {
		return errors.New("assignee user not found")
	}

	ticket.AssigneeID = &assigneeID

	return config.DB.Save(&ticket).Error
}
