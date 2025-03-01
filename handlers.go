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

func updateUserPassword(ctx iris.Context, db Database) {
	var req UpdatePasswordRequest
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
			Detail("Current password is incorrect: "+err.Error()).
			Status(iris.StatusInternalServerError))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		ctx.StopWithProblem(iris.StatusUnauthorized, iris.NewProblem().
			Title("Authentication error").
			Detail("Invalid current password: "+err.Error()).
			Status(iris.StatusUnauthorized))
		return
	}

	hashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal server error").
			Detail(err.Error()).
			Status(iris.StatusInternalServerError))
		return
	}

	if err := db.UpdateUserPassword(userID, hashedPassword); err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal server error").
			Detail(err.Error()).
			Status(iris.StatusInternalServerError))
		return
	}

	ctx.JSON(iris.Map{
		"message": "Password updated successfully",
	})
}

func logoutUser(ctx iris.Context) {
	ctx.RemoveCookie("token")
	ctx.JSON(iris.Map{
		"message": "Logged out successfully",
	})
}

func refreshToken(ctx iris.Context, db Database) {
	claims := ctx.Values().Get("claims").(jwt.MapClaims)
	userID := uint(claims["sub"].(float64))

	user, err := db.GetUserByID(userID)
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Invalid user").
			Detail(err.Error()).
			Status(iris.StatusInternalServerError))
		return
	}

	token, err := generateJWTToken(user)
	if err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Token generation error").
			Detail(err.Error()).
			Status(iris.StatusInternalServerError))
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
		"message": "Token refreshed successfully",
		"token":   token,
	})
}

func getUserDetails(ctx iris.Context, db Database) {
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

	response := UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	ctx.JSON(response)
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
