package service

import (
	"errors"
	"strings"

	"ticketing-system/internal/config"
	"ticketing-system/internal/model"
)

//
// ======================
// DTO RESPONSE
// ======================
//

type TicketResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	User        string `json:"user"`
	Assignee    string `json:"assignee"`
	Status      string `json:"status"`
}

//
// ======================
// HELPER MAPPING
// ======================
//

func mapToResponse(tickets []model.Ticket) []TicketResponse {
	var result []TicketResponse

	for _, t := range tickets {
		assigneeName := ""
		if t.Assignee != nil {
			assigneeName = t.Assignee.Name
		}

		result = append(result, TicketResponse{
			ID:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			User:        t.User.Name,
			Assignee:    assigneeName,
			Status:      t.Status,
		})
	}

	return result
}

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
// GET USER TICKETS (JOIN + DTO)
// ======================
func GetTicketsByUser(userID uint, page, limit int, search string) ([]TicketResponse, error) {
	var tickets []model.Ticket

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	query := config.DB.
		Preload("User").
		Preload("Assignee").
		Where("user_id = ?", userID)

	if search != "" {
		search = strings.ToLower(search)
		query = query.Where("LOWER(title) LIKE ?", "%"+search+"%")
	}

	err := query.
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&tickets).Error

	return mapToResponse(tickets), err
}

// ======================
// GET ALL TICKETS (ADMIN JOIN + DTO)
// ======================
func GetAllTickets(page, limit int, search string) ([]TicketResponse, error) {
	var tickets []model.Ticket

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	query := config.DB.
		Preload("User").
		Preload("Assignee").
		Model(&model.Ticket{})

	if search != "" {
		search = strings.ToLower(search)
		query = query.Where("LOWER(title) LIKE ?", "%"+search+"%")
	}

	err := query.
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&tickets).Error

	return mapToResponse(tickets), err
}

// ======================
// UPDATE STATUS
// ======================
func UpdateTicketStatus(ticketID uint, userID uint, role string, status string) error {
	var ticket model.Ticket

	if err := config.DB.First(&ticket, ticketID).Error; err != nil {
		return errors.New("ticket not found")
	}

	if ticket.UserID != userID &&
		(ticket.AssigneeID == nil || *ticket.AssigneeID != userID) &&
		role != "admin" {
		return errors.New("not authorized to update ticket")
	}

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

	if err := config.DB.First(&ticket, ticketID).Error; err != nil {
		return errors.New("ticket not found")
	}

	var user model.User
	if err := config.DB.First(&user, assigneeID).Error; err != nil {
		return errors.New("assignee user not found")
	}

	ticket.AssigneeID = &assigneeID

	return config.DB.Save(&ticket).Error
}
