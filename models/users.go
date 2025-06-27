package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"not null" json:"username" validate:"required,min=3,max=30"`
	Email    string `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Password string `gorm:"not null" json:"-" validate:"required,password"`
}

type UserClaims struct {
	UserID uint `json:"user_id"`
}
