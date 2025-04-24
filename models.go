package main

import (
	"time"

	"gorm.io/gorm"
)

// ══════════════════════════ Database table models ══════════════════════════

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username" validate:"required,min=3,max=50"`
	Password string `gorm:"not null" json:"password" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789,specialchar"`
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
	GetUserByID(userID uint) (*User, error)
	UpdateUser(user *User) error
	UpdateUserPassword(userID uint, hashedPassword string) error
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

type UpdateUserRequest struct {
	Username string `json:"username" validate:"omitempty,min=3,max=50"`
	Email    string `json:"email" validate:"omitempty,email"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ,containsany=abcdefghijklmnopqrstuvwxyz,containsany=0123456789,specialchar"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaskRequest struct {
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description" validate:"omitempty"`
	Completed   bool       `json:"completed"`
	CategoryID  uint       `json:"category_id" validate:"required"`
	DueDate     *time.Time `json:"due_date" validate:"omitempty"`
	TagIDs      []uint     `json:"tag_ids" validate:"omitempty"`
}

type TaskResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	UserID      uint       `json:"user_id"`
	CategoryID  uint       `json:"category_id"`
	DueDate     *time.Time `json:"due_date"`
	Tags        []Tag      `json:"tags"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
