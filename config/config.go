package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
)

var (
	TestDB *gorm.DB
	ProdDB *gorm.DB
)

func GetDSN(prefix string) string {
	host := os.Getenv(prefix + "_DB_HOST")
	user := os.Getenv(prefix + "_DB_USER")
	password := os.Getenv(prefix + "_DB_PASSWORD")
	dbName := os.Getenv(prefix + "_DB_NAME")
	port := os.Getenv(prefix + "_DB_PORT")
	sslMode := os.Getenv(prefix + "_DB_SSL_MODE")
	timeZone := os.Getenv(prefix + "_DB_TIME_ZONE")
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, port, user, password, dbName, sslMode, timeZone)
}

func connectDB(dsn) *gorm.DB {
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	return db
}
