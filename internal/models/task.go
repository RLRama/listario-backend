package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	DueDate     time.Time
	Priority    string `gorm:"default:'medium'"` // low, medium, high
	Status      string `gorm:"default:'todo'"`   // todo, in_progress, done
	ProjectID   uint   `gorm:"not null"`
	Project     Project
	AssigneeID  uint
	Assignee    User `gorm:"foreignKey:AssigneeID"`
	CreatorID   uint `gorm:"not null"`
	Creator     User `gorm:"foreignKey:CreatorID"`
}
