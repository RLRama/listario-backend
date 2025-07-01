package main

import (
	"os"

	"github.com/RLRama/listario-backend/db"
	"github.com/RLRama/listario-backend/handler"
	"github.com/RLRama/listario-backend/logger"
	"github.com/RLRama/listario-backend/middleware"
	"github.com/RLRama/listario-backend/repository"
	"github.com/RLRama/listario-backend/router"
	"github.com/RLRama/listario-backend/service"
	"github.com/RLRama/listario-backend/utils"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kataras/iris/v12"
)

func main() {
	logger.SetupLogger()

	logger.Info().Msgf("Starting Listario backend on port %s...", os.Getenv("PORT"))

	signer, verifier, err := utils.SetupJWT()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to set up JWT signer and verifier")
	}

	database, err := db.InitDB(db.GetDSN("PROD"))
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize database")
	}

	app := iris.Default()

	userRepository := repository.NewGormUserRepository(database)
	userService := service.NewUserService(userRepository, signer)
	userHandler := handler.NewUserHandler(userService, verifier)

	app.Validator = utils.NewCustomValidator()
	app.Use(middleware.RequestLogger())

	router.SetupRoutes(app, userHandler, verifier)

	if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start the server")
	}
}
