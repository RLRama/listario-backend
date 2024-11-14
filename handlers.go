package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
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
	}

	user, err := db.GetUserByUsernameOrEmail(req.Identifier)
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Authentication error").
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

	token, err := generateJWT(user.Username)
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal server error").
			Detail(err.Error()).
			Status(iris.StatusInternalServerError))
		return
	}

	ctx.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
	})
}

// ══════════════════════════ Various handlers ══════════════════════════

func resolveErrorsDocumentation(ctx iris.Context) {
	ctx.WriteString("This is planned to redirect to a page that should document to web developers or users of the API on how to resolve the validation errors")
}
