package model

import "time"

type Comment struct {
	ID       uint   `gorm:"primaryKey"`
	TicketID uint   `gorm:"not null;index"`
	Ticket   Ticket `gorm:"foreignKey:TicketID"`

	UserID uint `gorm:"not null"`
	User   User `gorm:"foreignKey:UserID"`

	Content    string `gorm:"type:text;not null"`
	IsInternal bool   `gorm:"default:false"` // true = hanya agent/admin

	CreatedAt time.Time
	UpdatedAt time.Time
}
