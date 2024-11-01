package main

import (
	"listario-backend/internal/database"
	"listario-backend/internal/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading.env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := iris.New()

	db, err := database.Connect(os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	routes.Register(app, db)

	app.Listen(":" + port)
}
