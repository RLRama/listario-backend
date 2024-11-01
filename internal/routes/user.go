package routes

import (
	"go/token"
	"listario-backend/internal/database"
	"listario-backend/internal/models"
	"log"

	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx iris.Context) {
    var user models.User
    if err := ctx.ReadJSON(&user); err != nil {
        ctx.StatusCode(iris.StatusBadRequest)
        ctx.JSON(iris.Map{"error": "Invalid input"})
        return
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Printf("Password hashing error: %v", err)
        ctx.StatusCode(iris.StatusInternalServerError)
        ctx.JSON(iris.Map{"error": "Server error"})
        return
    }
    user.Password = string(hashedPassword)

    // Save the user in the database
    db := database.DBConn()
    if err := db.Create(&user).Error; err != nil {
        ctx.StatusCode(iris.StatusConflict)
        ctx.JSON(iris.Map{"error": "Email already exists"})
        return
    }

    ctx.StatusCode(iris.StatusCreated)
    ctx.JSON(iris.Map{"message": "User registered successfully"})
}

func LoginUser(ctx iris.Context) {
    var credentials struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := ctx.ReadJSON(&credentials); err!= nil {
        ctx.StatusCode(iris.StatusBadRequest)
        ctx.JSON(iris.Map{"error": "Invalid input"})
        return
    }

    // Find user by email
    var user models.User
    db := database.DBConn()
    if err := db.Where("email =?", credentials.Email).First(&user).Error; err!= nil {
        ctx.StatusCode(iris.StatusNotFound)
        ctx.JSON(iris.Map{"error": "Invalid email or password"})
        return
    }

    // Check password hash
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
        ctx.StatusCode(iris.StatusUnauthorized)
        ctx.JSON(iris.Map{"error": "Invalid email or password"})
        return
    }

    token := 
}