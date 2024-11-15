package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
)

// ══════════════════════════ User handlers ══════════════════════════

func postUser(ctx iris.Context, db Database) {
	var user User
	if err := ctx.ReadJSON(&user); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors := wrapValidationErrors(errs)

			ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
				Title("Invalid input").
				Detail("One or more fields failed validation").
				Type("/user/validation-errors").
				Key("errors", validationErrors))
			return
		}
	}

	if !validatePassword(user.Password) {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid password").
			Detail("Password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one number, and one special character"))
		return
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal server error").
			Detail(err.Error()).
			Status(iris.StatusInternalServerError))
		return
	}

	user.Password = hashedPassword

	if err := db.CreateUser(&user); err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal server error").
			Detail(err.Error()).
			Status(iris.StatusInternalServerError))
		return
	}

	ctx.JSON(iris.Map{"message": "User registered successfully"})
}

func loginUser(ctx iris.Context, db Database) {
	var req LoginRequest
	if err := ctx.ReadJSON(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors := wrapValidationErrors(errs)

			ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
				Title("Invalid input").
				Detail("One or more fields failed validation").
				Type("/user/validation-errors").
				Key("errors", validationErrors))
			return
		}
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid input").
			Detail(err.Error()).
			Status(iris.StatusBadRequest))
		return
	}

	user, err := db.GetUserByUsernameOrEmail(req.Identifier)
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Invalid login credentials").
			Detail("Invalid username or email or password: "+err.Error()).
			Status(iris.StatusUnauthorized))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Authentication error").
			Detail("Invalid username or email or password: "+err.Error()).
			Status(iris.StatusUnauthorized))
		return
	}

	token, err := generateJWTToken(user)
	if err != nil {
		ctx.StopWithProblem(http.StatusInternalServerError, iris.NewProblem().
			Title("Internal server error").
			Detail(err.Error()).
			Status(http.StatusInternalServerError))
		return
	}

	ctx.SetCookie(&iris.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	ctx.JSON(iris.Map{
		"message": "Login successful",
		"token":   token,
	})
}

func updateUser(ctx iris.Context, db Database) {
	var req UpdateUserRequest
	if err := ctx.ReadJSON(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors := wrapValidationErrors(errs)

			ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
				Title("Invalid input").
				Detail("One or more fields failed validation").
				Type("/user/validation-errors").
				Key("errors", validationErrors))
			return
		}
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid input").
			Detail(err.Error()).
			Status(iris.StatusBadRequest))
		return
	}

	claims := ctx.Values().Get("claims").(jwt.MapClaims)
	userID := uint(claims["sub"].(float64))

	user, err := db.GetUserByID(userID)
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal server error").
			Detail(err.Error()).
			Status(iris.StatusInternalServerError))
		return
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := db.UpdateUser(user); err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal server error").
			Detail(err.Error()).
			Status(iris.StatusInternalServerError))
		return
	}

	ctx.JSON(iris.Map{
		"message": "User updated successfully",
	})
}

// ══════════════════════════ Various handlers ══════════════════════════

func resolveErrorsDocumentation(ctx iris.Context) {
	ctx.WriteString("This is planned to redirect to a page that should document to web developers or users of the API on how to resolve the validation errors")
}

func protectedHandler(ctx iris.Context) {
	token, err := ctx.Request().Cookie("token")
	if err != nil || token.Value == "" {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail("No token found in cookie: "+err.Error()).
			Status(iris.StatusUnauthorized))
		return
	}

	signingKey := []byte(os.Getenv("JWT_SECRET"))
	tokenClaims, err := jwt.Parse(token.Value, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil || !tokenClaims.Valid {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Unauthorized").
			Detail("Invalid token: "+err.Error()).
			Status(iris.StatusUnauthorized))
		return
	}

	ctx.JSON(iris.Map{
		"message": "This is a protected endpoint",
	})
}
