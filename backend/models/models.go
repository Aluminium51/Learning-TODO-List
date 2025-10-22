package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"unique;not null" json:"username"`
	Email        string    `gorm:"unique;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Todos        []Todo    `json:"todos"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Todo struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	UserID        uint       `gorm:"not null" json:"user_id"`
	Title         string     `gorm:"not null" json:"title"`
	Description   string     `json:"description"`
	Completed     bool       `gorm:"default:false" json:"completed"`
	DueDate       *time.Time `json:"due_date"`
	AttachmentURL string     `json:"attachment_url"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
