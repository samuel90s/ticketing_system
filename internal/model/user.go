package model

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Email     string `gorm:"unique"`
	Password  string `json:"-"`
	Role      string `gorm:"default:user"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
