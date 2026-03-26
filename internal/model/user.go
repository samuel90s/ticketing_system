package model

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `json:"-"`
	Role      string `gorm:"type:varchar(20);default:'user'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
