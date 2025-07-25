package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title     string `gorm:"not null" json:"title"`
	Content   string `json:"content"`
	Completed bool   `gorm:"default:false" json:"completed"`
	UserID    uint   `gorm:"not null" json:"user_id"`
}

type CreateTaskRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=100"`
	Content string `json:"content"`
}

type UpdateTaskRequest struct {
	Title     string `json:"title" validate:"omitempty,min=1,max=100"`
	Content   string `json:"content"`
	Completed *bool  `json:"completed"`
}

type TaskResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
