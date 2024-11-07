package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kataras/iris/v12"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	app := newApp()

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := InitDB(db); err != nil {
		log.Fatal(err)
	}

	log.Println("Database initialized and schemas migrated")

	err2 := app.Listen(":" + os.Getenv("PORT"))
	if err2 != nil {
		panic(err2)
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
		testEndpoints.Get("/db-connection", func(ctx iris.Context) {
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
