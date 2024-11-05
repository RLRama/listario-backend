package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
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
