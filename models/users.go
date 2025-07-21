package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null" json:"username" validate:"required,min=3,max=30"`
	Email    string `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password string `gorm:"not null" json:"-" validate:"required,password"`
	Tasks    []Task
}

type UserClaims struct {
	UserID uint `json:"user_id"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" validate:"omitempty,email"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
