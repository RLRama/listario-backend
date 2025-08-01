package router

import (
	"github.com/RLRama/listario-backend/handler"
	"github.com/RLRama/listario-backend/logger"
	"github.com/RLRama/listario-backend/models"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
)

func SetupRoutes(app *iris.Application, userHandler *handler.UserHandler, taskHandler *handler.TaskHandler, verifier *jwt.Verifier) {
	verifyMiddleware := verifier.Verify(func() interface{} {
		return new(models.UserClaims)
	})
	// Misc routes
	app.Get("/health", func(ctx iris.Context) {
		logger.Info().Msg("Received /health check request")
		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString("OK")
	})

	// Public routes
	authAPI := app.Party("/auth")
	{
		authAPI.Post("/register", userHandler.Register)
		authAPI.Post("/login", userHandler.Login)
		authAPI.Post("/refresh", userHandler.RefreshToken)
	}

	// Protected routes
	userAPI := app.Party("/users")
	userAPI.Use(verifyMiddleware)
	{
		userAPI.Get("/me", userHandler.GetMyDetails)
		userAPI.Put("/me", userHandler.UpdateMyDetails)
		userAPI.Get("/logout", userHandler.Logout)
	}
	taskAPI := app.Party("/tasks")
	taskAPI.Use(verifyMiddleware)
	{
		taskAPI.Post("/", taskHandler.CreateTask)
		taskAPI.Get("/", taskHandler.GetMyTasks)
		taskAPI.Get("/{id:uint}", taskHandler.GetTask)
		taskAPI.Put("/{id:uint}", taskHandler.UpdateTask)
		taskAPI.Delete("/{id:uint}", taskHandler.DeleteTask)
	}
}
