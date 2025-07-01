package db

import (
	"fmt"
	"os"

	"github.com/RLRama/listario-backend/logger"
	"github.com/RLRama/listario-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
		host, user, password, dbName, port, sslMode, timeZone)
}

func InitDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to the database")
		return nil, err
	}

	err = db.AutoMigrate(
		&models.User{},
	)

	if err != nil {
		logger.Error().Err(err).Msg("Failed to auto-migrate database models")
		return nil, err
	}

	logger.Info().Msg("Database connection established and migrations run")
	return db, nil
}
