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

	return &GormDatabase{db}, nil
}

func populateSampleData(db *gorm.DB) error {

}

// Database operations

func (db *GormDatabase) CreateUser(user *User) error {
	return db.DB.Create(user).Error
}
