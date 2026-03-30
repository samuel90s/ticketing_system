package service

import (
	"errors"
	"strings"

	"ticketing-system/internal/config"
	"ticketing-system/internal/model"
)

// ======================
// DTO RESPONSE
// ======================

type TicketResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	User        string `json:"user"`
	Assignee    string `json:"assignee"`
	Status      string `json:"status"`
}

type PaginatedTickets struct {
	Data  []TicketResponse `json:"data"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

// ======================
// VALID STATUSES
// ======================

var validStatuses = map[string]bool{
	"open":        true,
	"in_progress": true,
	"closed":      true,
}

// ======================
// HELPER MAPPING
// ======================

func mapToResponse(tickets []model.Ticket) []TicketResponse {
	result := make([]TicketResponse, 0, len(tickets))

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

func sanitizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return page, limit
}

// ======================
// CREATE TICKET
// ======================

func CreateTicket(title, description string, userID uint) error {
	ticket := model.Ticket{
		Title:       title,
		Description: description,
		UserID:      userID,
		Status:      "open",
	}

	return config.DB.Create(&ticket).Error
}

// ======================
// GET USER TICKETS
// ======================

func GetTicketsByUser(userID uint, page, limit int, search string) (PaginatedTickets, error) {
	page, limit = sanitizePagination(page, limit)

	query := config.DB.Model(&model.Ticket{}).
		Preload("User").
		Preload("Assignee").
		Where("user_id = ?", userID)

	// Search aman pakai named param — tidak bisa SQL injection
	if search != "" {
		query = query.Where("title LIKE ?", "%"+strings.TrimSpace(search)+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return PaginatedTickets{}, err
	}

	var tickets []model.Ticket
	err := query.
		Offset((page - 1) * limit).
		Limit(limit).
		Order("created_at DESC").
		Find(&tickets).Error

	return PaginatedTickets{
		Data:  mapToResponse(tickets),
		Total: total,
		Page:  page,
		Limit: limit,
	}, err
}

// ======================
// GET ALL TICKETS (ADMIN)
// ======================

func GetAllTickets(page, limit int, search, status string) (PaginatedTickets, error) {
	page, limit = sanitizePagination(page, limit)

	query := config.DB.Model(&model.Ticket{}).
		Preload("User").
		Preload("Assignee")

	if search != "" {
		query = query.Where("title LIKE ?", "%"+strings.TrimSpace(search)+"%")
	}

	// Filter by status (opsional)
	if status != "" {
		if !validStatuses[strings.ToLower(status)] {
			return PaginatedTickets{}, errors.New("invalid status filter")
		}
		query = query.Where("status = ?", strings.ToLower(status))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return PaginatedTickets{}, err
	}

	var tickets []model.Ticket
	err := query.
		Offset((page - 1) * limit).
		Limit(limit).
		Order("created_at DESC").
		Find(&tickets).Error

	return PaginatedTickets{
		Data:  mapToResponse(tickets),
		Total: total,
		Page:  page,
		Limit: limit,
	}, err
}

// ======================
// UPDATE STATUS
// ======================

func UpdateTicketStatus(ticketID, userID uint, role, status string) error {
	var ticket model.Ticket

	if err := config.DB.First(&ticket, ticketID).Error; err != nil {
		return errors.New("ticket not found")
	}

	// Cek akses
	isOwner := ticket.UserID == userID
	isAssignee := ticket.AssigneeID != nil && *ticket.AssigneeID == userID
	isAdmin := role == "admin"

	if !isOwner && !isAssignee && !isAdmin {
		return errors.New("not authorized to update this ticket")
	}

	status = strings.ToLower(strings.TrimSpace(status))
	if !validStatuses[status] {
		return errors.New("invalid status, allowed: open, in_progress, closed")
	}

	return config.DB.Model(&ticket).Update("status", status).Error
}

// ======================
// ASSIGN TICKET (ADMIN)
// ======================

func AssignTicket(ticketID, assigneeID uint) error {
	var ticket model.Ticket
	if err := config.DB.First(&ticket, ticketID).Error; err != nil {
		return errors.New("ticket not found")
	}

	var user model.User
	if err := config.DB.First(&user, assigneeID).Error; err != nil {
		return errors.New("assignee user not found")
	}

	return config.DB.Model(&ticket).Update("assignee_id", assigneeID).Error
}

// ======================
// GET TICKET BY ID (ADMIN)
// ======================

func GetTicketByID(ticketID uint) (TicketResponse, error) {
	var ticket model.Ticket

	err := config.DB.
		Preload("User").
		Preload("Assignee").
		First(&ticket, ticketID).Error

	if err != nil {
		return TicketResponse{}, errors.New("ticket not found")
	}

	assigneeName := ""
	if ticket.Assignee != nil {
		assigneeName = ticket.Assignee.Name
	}

	return TicketResponse{
		ID:          ticket.ID,
		Title:       ticket.Title,
		Description: ticket.Description,
		User:        ticket.User.Name,
		Assignee:    assigneeName,
		Status:      ticket.Status,
	}, nil
}

// ======================
// EDIT TICKET (ADMIN)
// ======================

func EditTicket(ticketID uint, title, description string) error {
	var ticket model.Ticket

	if err := config.DB.First(&ticket, ticketID).Error; err != nil {
		return errors.New("ticket not found")
	}

	updates := map[string]interface{}{}
	if title != "" {
		updates["title"] = strings.TrimSpace(title)
	}
	if description != "" {
		updates["description"] = strings.TrimSpace(description)
	}

	return config.DB.Model(&ticket).Updates(updates).Error
}

// ======================
// DELETE TICKET (ADMIN)
// ======================

func DeleteTicket(ticketID uint) error {
	var ticket model.Ticket

	if err := config.DB.First(&ticket, ticketID).Error; err != nil {
		return errors.New("ticket not found")
	}

	return config.DB.Delete(&ticket).Error
}

// ======================
// GET ALL USERS (untuk assign dropdown)
// ======================

func GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := config.DB.Select("id, name, email, role").Find(&users).Error
	return users, err
}
