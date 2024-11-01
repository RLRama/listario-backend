// routes/health.go
package routes

import (
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

func Register(app *iris.Application, db *gorm.DB) {
	app.Get("/health", func(ctx iris.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{"status": "failed", "error": "could not get database handle"})
			return
		}

		// Ping the database
		if err := sqlDB.Ping(); err != nil {
			ctx.StatusCode(500)
			ctx.JSON(iris.Map{"status": "failed", "error": "could not connect to database"})
		} else {
			ctx.JSON(iris.Map{"status": "healthy"})
		}
	})
}
