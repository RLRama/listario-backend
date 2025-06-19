package main

import (
	"os"

	"github.com/RLRama/listario-backend/db"
	"github.com/RLRama/listario-backend/logger"
	"github.com/kataras/iris/v12"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger.SetupLogger()

	logger.Info().Msgf("Starting Listario Backend on port %s...", os.Getenv("PORT"))

	database, err := db.InitDB(db.GetDSN("DB"))
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize database")
	}

	app := iris.New()
}
