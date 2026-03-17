package model

type Ticket struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	UserID      uint   `gorm:"not null"` // pemilik ticket
	AssigneeID  *uint  // user yang ditugaskan (nullable)
	Status      string `gorm:"default:open"` // open / closed
}
