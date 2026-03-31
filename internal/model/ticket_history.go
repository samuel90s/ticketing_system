package model

import "time"

type TicketHistory struct {
	ID       uint   `gorm:"primaryKey"`
	TicketID uint   `gorm:"not null;index"`
	Ticket   Ticket `gorm:"foreignKey:TicketID"`

	UserID uint `gorm:"not null"`
	User   User `gorm:"foreignKey:UserID"`

	// Jenis perubahan: status_changed | assigned | commented
	Action   string `gorm:"type:varchar(50);not null"`
	OldValue string `gorm:"type:varchar(100)"`
	NewValue string `gorm:"type:varchar(100)"`

	CreatedAt time.Time
}
