package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email         string `gorm:"uniqueIndex;not null"`
	Password      string `gorm:"not null"`
	FirstName     string
	LastName      string
	Projects      []Project `gorm:"many2many:user_projects;"`
	AssignedTasks []Task    `gorm:"foreignKey:AssigneeID"`
	CreatedTasks  []Task    `gorm:"foreignKey:CreatorID"`
}
