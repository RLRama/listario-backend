package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Tasks    []Task `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Task struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	Completed   bool
	UserID      uint `gorm:"not null"`
	User        User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CategoryID  uint
	DueDate     *time.Time
	Tags        []Tag `gorm:"many2many:task_tags;"`
}

type Category struct {
	gorm.Model
	Name  string `gorm:"unique;not null"`
	Tasks []Task `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Tag struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}
