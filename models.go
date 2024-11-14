package main

import (
	"time"

	"gorm.io/gorm"
)

// ══════════════════════════ Database table models ══════════════════════════

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username" validate:"required,min=3,max=50"`
	Password string `gorm:"not null" json:"password" validate:"required,min=8"`
	Email    string `gorm:"unique;not null" json:"email" validate:"required,email"`
	Tasks    []Task `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Task struct {
	gorm.Model
	Title       string `gorm:"not null" validate:"required"`
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
	Name  string `gorm:"unique;not null" validate:"required"`
	Tasks []Task `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Tag struct {
	gorm.Model
	Name string `gorm:"unique;not null" validate:"required"`
}

// ══════════════════════════ Application-specific models ══════════════════════════

type GormDatabase struct {
	*gorm.DB
}

type Database interface {
	CreateUser(user *User) error
	GetUserByUsernameOrEmail(username string) (*User, error)
}

type validationError struct {
	ActualTag string `json:"tag"`
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Param     string `json:"param"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}
