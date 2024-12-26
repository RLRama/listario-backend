package main

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kataras/iris/v12"
)

func main() {

	db, err := setupDatabase()
	if err != nil {
		panic(err)
	}

	app := newApp(db)

	log.Println("Database initialized and schemas migrated")

	err2 := app.Listen(":" + os.Getenv("PORT"))
	if err2 != nil {
		panic(err2)
	}
}

func newApp(db Database) *iris.Application {
	app := iris.Default()
	v := validator.New(validator.WithRequiredStructEnabled())
	registerCustomValidators(v)
	app.Validator = v

	// ══════════════════════════ Test routes ══════════════════════════
	testRouter := app.Party("/test")
	{
		testRouter.Get("/hello", func(ctx iris.Context) {
			ctx.JSON(iris.Map{
				"message": "Hello, Iris!",
			})
		})
		testRouter.Get("/db-connection", func(ctx iris.Context) {
			version, err := TestDBConnection()
			if err != nil {
				ctx.JSON(iris.Map{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}

			ctx.JSON(iris.Map{
				"status":  "success",
				"message": "Database connection successful",
				"version": version,
			})
		})
		authorizedTestRouter := testRouter.Party("/protected")
		authorizedTestRouter.Use(AuthenticationMiddleware)
		authorizedTestRouter.Get("/", func(ctx iris.Context) {
			protectedHandler(ctx)
		})
	}

	// ══════════════════════════ User routes ══════════════════════════
	userRouter := app.Party("/user")
	{
		userRouter.Get("/validation-errors", resolveErrorsDocumentation)
		userRouter.Post("/register", func(ctx iris.Context) {
			postUser(ctx, db)
		})
		userRouter.Post("/login", func(ctx iris.Context) {
			loginUser(ctx, db)
		})
		authorizedUserRouter := userRouter.Party("/protected")
		authorizedUserRouter.Use(AuthenticationMiddleware)
		authorizedUserRouter.Put("/update", func(ctx iris.Context) {
			updateUser(ctx, db)
		})
		authorizedUserRouter.Put("/update-password", func(ctx iris.Context) {
			updateUserPassword(ctx, db)
		})
	}

	return app
}
