package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	DueDate     time.Time
	Status      string `gorm:"default:'active'"` // active, completed, archived
	Tasks       []Task `gorm:"foreignKey:ProjectID"`
	Users       []User `gorm:"many2many:user_projects;"`
	CreatorID   uint   `gorm:"not null"`
	Creator     User   `gorm:"foreignKey:CreatorID"`
}
