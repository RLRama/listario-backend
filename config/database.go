package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kataras/iris/v12"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
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

func GetDBNameFromEnv(prefix string) string {
	dbName := os.Getenv(prefix + "_DB_NAME")
	if dbName == "" {
		panic(fmt.Sprintf("Environment variable %s_DB_NAME is not set", prefix))
	}
	return dbName
}

func connectDB(dsn, name string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Str("db", name).Msg("Failed to connect to database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Str("db", name).Msg("Failed to get underlying SQL DB")
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal().Err(err).Str("db", name).Msg("failed to ping database")
	}

	log.Info().Str("db", name).Msg("database connection successful")
	return db
}

func testDBConnections(ctx iris.Context) {
	var result1, result2 int

	if err := TestDB.Raw("SELECT 1").Scan(&result1).Error; err != nil {
		log.Error().Err(err).Msg("TestDB connection failed")
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.
	}
}