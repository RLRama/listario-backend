package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDBConnection() (string, error) {

	dsn := os.Getenv("DSN")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return "", fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return "", fmt.Errorf("failed to ping database: %v", err)
	}

	var version string
	err = db.QueryRow("SELECT version();").Scan(&version)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve PostgreSQL version: %v", err)
	}

	return version, nil
}

func setupDatabase() (*GormDatabase, error) {
	dsn := os.Getenv("DSN")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	db.AutoMigrate(
		&Category{},
		&Tag{},
		&User{},
		&Task{},
	)

	if err := populateSampleData(db); err != nil {
		return nil, err
	}

	return &GormDatabase{db}, nil
}

func populateSampleData(db *gorm.DB) error {
	var count int64

	if err := db.Model(&Category{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		categories := []*Category{
			{Name: "Work"},
			{Name: "Personal"},
			{Name: "Family"},
			{Name: "Shopping"},
			{Name: "Health"},
			{Name: "Finance"},
		}

		if err := db.Create(&categories).Error; err != nil {
			return err
		}
	}

	if err := db.Model(&Tag{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		tags := []*Tag{
			{Name: "Urgent"},
			{Name: "Important"},
			{Name: "Low Priority"},
			{Name: "High Priority"},
			{Name: "Medium Priority"},
		}

		if err := db.Create(&tags).Error; err != nil {
			return err
		}
	}

	return nil
}

// ══════════════════════════ Database operations ══════════════════════════

func (db *GormDatabase) CreateUser(user *User) error {
	return db.DB.Create(user).Error
}

func (db *GormDatabase) GetUserByUsernameOrEmail(identifier string) (*User, error) {
	var user User
	if err := db.DB.Where("username = ?", identifier).Or("email = ?", identifier).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *GormDatabase) GetUserByID(userID uint) (*User, error) {
	var user User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *GormDatabase) UpdateUser(user *User) error {
	return db.DB.Save(user).Error
}

func (db *GormDatabase) UpdateUserPassword(userID uint, hashedPassword string) error {
	return db.DB.Model(&User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}
