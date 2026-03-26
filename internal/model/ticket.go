package model

import "time"

type Ticket struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`

	UserID uint `gorm:"not null"` // owner
	User   User `gorm:"foreignKey:UserID"`

	AssigneeID *uint
	Assignee   *User `gorm:"foreignKey:AssigneeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Status string `gorm:"type:varchar(20);default:'open';index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
