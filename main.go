package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kataras/iris/v12"
)

func main() {
	app := newApp()

	err := app.Listen(":" + os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
}

func newApp() *iris.Application {
	app := iris.Default()

	// test endpoints
	testEndpoints := app.Party("/test")
	{
		testEndpoints.Get("/hello", func(ctx iris.Context) {
			ctx.JSON(iris.Map{
				"message": "Hello, Iris!",
			})
		})
		testEndpoints.Get("/test-db", func(ctx iris.Context) {
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
	}

	return app
}
