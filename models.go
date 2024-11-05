package main

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string
	Password string
}

type Task struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	Completed   bool
	UserID      uint `gorm:"not null"`
	CategoryID  uint
}

type Category struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}
