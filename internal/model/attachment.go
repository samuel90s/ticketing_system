package model

import "time"

type Attachment struct {
	ID       uint   `gorm:"primaryKey"`
	TicketID uint   `gorm:"not null;index"`
	Ticket   Ticket `gorm:"foreignKey:TicketID"`

	UserID uint `gorm:"not null"`
	User   User `gorm:"foreignKey:UserID"`

	FileName    string `gorm:"not null"`
	FilePath    string `gorm:"not null"`
	FileSize    int64
	ContentType string `gorm:"type:varchar(100)"`

	CreatedAt time.Time
}
