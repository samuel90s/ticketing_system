package service

import (
	"errors"
	"strings"
	"time"

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
	UserID      uint   `json:"user_id"`
	Assignee    string `json:"assignee"`
	AssigneeID  *uint  `json:"assignee_id"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type PaginatedTickets struct {
	Data  []TicketResponse `json:"data"`
	Total int64            `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

type TicketHistoryResponse struct {
	ID        uint   `json:"id"`
	TicketID  uint   `json:"ticket_id"`
	UserID    uint   `json:"user_id"`
	User      string `json:"user"`
	Action    string `json:"action"`
	OldValue  string `json:"old_value"`
	NewValue  string `json:"new_value"`
	CreatedAt string `json:"created_at"`
}

// ======================
// VALID STATUSES
// ======================

// open → replied → on_progress → done
var validStatuses = map[string]bool{
	"open":        true,
	"replied":     true,
	"on_progress": true,
	"done":        true,
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
			UserID:      t.UserID,
			Assignee:    assigneeName,
			AssigneeID:  t.AssigneeID,
			Status:      t.Status,
			CreatedAt:   t.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   t.UpdatedAt.Format("2006-01-02 15:04:05"),
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
// LOG HISTORY
// ======================

func logHistory(ticketID, userID uint, action, oldVal, newVal string) {
	h := model.TicketHistory{
		TicketID:  ticketID,
		UserID:    userID,
		Action:    action,
		OldValue:  oldVal,
		NewValue:  newVal,
		CreatedAt: time.Now(),
	}
	config.DB.Create(&h)
}

// ======================
// CREATE TICKET
// ======================

func CreateTicket(title, description string, userID uint) (uint, error) {
	ticket := model.Ticket{
		Title:       title,
		Description: description,
		UserID:      userID,
		Status:      "open",
	}

	if err := config.DB.Create(&ticket).Error; err != nil {
		return 0, err
	}

	logHistory(ticket.ID, userID, "created", "", "open")
	return ticket.ID, nil
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
// GET ALL TICKETS (ADMIN/AGENT)
// ======================

func GetAllTickets(page, limit int, search, status string) (PaginatedTickets, error) {
	page, limit = sanitizePagination(page, limit)

	query := config.DB.Model(&model.Ticket{}).
		Preload("User").
		Preload("Assignee")

	if search != "" {
		query = query.Where("title LIKE ?", "%"+strings.TrimSpace(search)+"%")
	}

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
// GET TICKET BY ID (ADMIN/AGENT)
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
		UserID:      ticket.UserID,
		Assignee:    assigneeName,
		AssigneeID:  ticket.AssigneeID,
		Status:      ticket.Status,
		CreatedAt:   ticket.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   ticket.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// ======================
// GET TICKET BY ID (USER - only own ticket)
// ======================

func GetUserTicketByID(ticketID, userID uint) (TicketResponse, error) {
	var ticket model.Ticket
	err := config.DB.
		Preload("User").
		Preload("Assignee").
		Where("id = ? AND user_id = ?", ticketID, userID).
		First(&ticket).Error

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
		UserID:      ticket.UserID,
		Assignee:    assigneeName,
		AssigneeID:  ticket.AssigneeID,
		Status:      ticket.Status,
		CreatedAt:   ticket.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   ticket.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// ======================
// UPDATE STATUS
// ======================

func UpdateTicketStatus(ticketID, userID uint, role, status string) error {
	var ticket model.Ticket
	if err := config.DB.First(&ticket, ticketID).Error; err != nil {
		return errors.New("ticket not found")
	}

	isOwner := ticket.UserID == userID
	isAssignee := ticket.AssigneeID != nil && *ticket.AssigneeID == userID
	isAdmin := role == "admin"
	isAgent := role == "agent"

	if !isOwner && !isAssignee && !isAdmin && !isAgent {
		return errors.New("not authorized to update this ticket")
	}

	status = strings.ToLower(strings.TrimSpace(status))
	if !validStatuses[status] {
		return errors.New("invalid status, allowed: open, replied, on_progress, done")
	}

	oldStatus := ticket.Status
	if err := config.DB.Model(&ticket).Update("status", status).Error; err != nil {
		return err
	}

	logHistory(ticket.ID, userID, "status_changed", oldStatus, status)
	return nil
}

// ======================
// ASSIGN TICKET (ADMIN)
// ======================

func AssignTicket(ticketID, assigneeID, adminID uint) error {
	var ticket model.Ticket
	if err := config.DB.First(&ticket, ticketID).Error; err != nil {
		return errors.New("ticket not found")
	}

	var user model.User
	if err := config.DB.First(&user, assigneeID).Error; err != nil {
		return errors.New("assignee user not found")
	}

	if err := config.DB.Model(&ticket).Update("assignee_id", assigneeID).Error; err != nil {
		return err
	}

	logHistory(ticket.ID, adminID, "assigned", "", user.Name)
	return nil
}

// ======================
// EDIT TICKET (ADMIN/AGENT)
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
// GET TICKET HISTORY
// ======================

func GetTicketHistory(ticketID uint) ([]TicketHistoryResponse, error) {
	var histories []model.TicketHistory
	err := config.DB.
		Preload("User").
		Where("ticket_id = ?", ticketID).
		Order("created_at ASC").
		Find(&histories).Error

	if err != nil {
		return nil, err
	}

	result := make([]TicketHistoryResponse, 0, len(histories))
	for _, h := range histories {
		result = append(result, TicketHistoryResponse{
			ID:        h.ID,
			TicketID:  h.TicketID,
			UserID:    h.UserID,
			User:      h.User.Name,
			Action:    h.Action,
			OldValue:  h.OldValue,
			NewValue:  h.NewValue,
			CreatedAt: h.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return result, nil
}

// ======================
// GET ALL USERS (untuk assign dropdown)
// ======================

func GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := config.DB.Select("id, name, email, role").Find(&users).Error
	return users, err
}
