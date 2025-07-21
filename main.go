package main

import (
	"os"

	"github.com/RLRama/listario-backend/db"
	_ "github.com/RLRama/listario-backend/docs"
	"github.com/RLRama/listario-backend/handler"
	"github.com/RLRama/listario-backend/logger"
	"github.com/RLRama/listario-backend/middleware"
	"github.com/RLRama/listario-backend/repository"
	"github.com/RLRama/listario-backend/router"
	"github.com/RLRama/listario-backend/service"
	"github.com/RLRama/listario-backend/utils"
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kataras/iris/v12"
)

// @title           Listario API
// @version         1.0
// @description     This is the API for the Listario application.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  rl.ramiro11@gmail.com

// @license.name  MIT License
// @license.url   https://opensource.org/license/mit/

// @host
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and a JWT.
func main() {
	logger.SetupLogger()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info().Msgf("Starting Listario backend on port %s...", port)

	signer, verifier, refreshTokenMaxAge, err := utils.SetupJWT()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to set up JWT signer and verifier")
	}

	database, err := db.InitDB(db.GetDSN("PROD"))
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize database")
	}

	app := iris.Default()

	swaggerURL := "/swagger/doc.json"
	config := &swagger.Config{
		URL:         swaggerURL,
		DeepLinking: true,
	}
	app.Get("/swagger/{any:path}", swagger.CustomWrapHandler(config, swaggerFiles.Handler))

	userRepository := repository.NewGormUserRepository(database)
	taskRepository := repository.NewGormTaskRepository(database)

	userService := service.NewUserService(userRepository, signer, refreshTokenMaxAge)
	taskService := service.NewTaskService(taskRepository)

	userHandler := handler.NewUserHandler(userService, verifier)
	taskHandler := handler.NewTaskHandler(taskService)

	app.Validator = utils.NewCustomValidator()
	app.Use(middleware.RequestLogger())

	router.SetupRoutes(app, userHandler, taskHandler, verifier)

	if err := app.Listen(":" + port); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start the server")
	}
}
