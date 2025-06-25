package main

import (
	"os"

	"github.com/RLRama/listario-backend/logger"
	"github.com/RLRama/listario-backend/middleware"
	"github.com/kataras/iris/v12"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	logger.SetupLogger()

	logger.Info().Msgf("Starting Listario backend on port %s...", os.Getenv("PORT"))

	/*
		database, err := db.InitDB(db.GetDSN("DB"))
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to initialize database")
		}
	*/
	// will uncomment once first (user) service which consumes the database is implemented

	app := iris.Default()

	// Placeholder for upcoming service and handler initialization

	app.Use(middleware.RequestLogger())

	// Placeholder for upcoming routes

	app.Get("/health", func(ctx iris.Context) {
		logger.Info().Msg("Received /health check request")
		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString("OK")
	})

	err := app.Listen(":" + os.Getenv("PORT"))
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start the server")
	}
}
