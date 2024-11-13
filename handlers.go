package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
)

// ----------------- User handlers -----------------

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

func resolveErrorsDocumentation(ctx iris.Context) {
	ctx.WriteString("This is planned to redirect to a page that should document to web developers or users of the API on how to resolve the validation errors")
}
