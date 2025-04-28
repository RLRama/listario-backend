package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

// ══════════════════════════ Task handlers ══════════════════════════

func createTask(ctx iris.Context, db Database) {
	var req TaskRequest
	if err := ctx.ReadJSON(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors := wrapValidationErrors(errs)

			ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
				Title("Invalid input").
				Detail("One or more fields failed validation").
				Type("/task/validation-errors").
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

	// Verify existence of category
	if _, err := validateCategory(db.(*GormDatabase), req.CategoryID); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid category").
			Detail("Category not found: "+err.Error()).
			Status(iris.StatusBadRequest))
		return
	}

	// Fetch tags, if provided
	tags, err := validateTags(db.(*GormDatabase), req.TagIDs)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid tags").
			Detail("Tags not found: "+err.Error()).
			Status(iris.StatusBadRequest))
		return
	}

	task := Task{
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
		UserID:      userID,
		CategoryID:  req.CategoryID,
		DueDate:     req.DueDate,
		Tags:        tags,
	}

	if err := db.(*GormDatabase).DB.Create(&task).Error; err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal server error").
			Detail(err.Error()).
			Status(iris.StatusInternalServerError))
		return
	}

	ctx.JSON(iris.Map{
		"message": "Task created successfully",
		"task":    mapTaskToResponse(&task),
	})
}

func listTasks(ctx iris.Context, db Database) {
	claims := ctx.Values().Get("claims").(jwt.MapClaims)
	userID := uint(claims["sub"].(float64))

	query := db.(*GormDatabase).DB.Model(&Task{}).Where("user_id = ?", userID).Preload("Tags")

	// Filter by category and tags
	if categoryID := ctx.URLParam("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if tagIDs := ctx.URLParam("tag_ids"); tagIDs != "" {
		query = query.Joins("JOIN task_tags ON task_tags.task_id = tasks.id").
			Where("task_tags.tag_id IN (?)", strings.Split(tagIDs, ","))
	}

	var tasks []Task
	if err := query.Find(&tasks).Error; err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal server error").
			Detail(err.Error()).
			Status(iris.StatusInternalServerError))
		return
	}

	response := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		response[i] = mapTaskToResponse(&task)
	}

	ctx.JSON(response)
}

func getTask(ctx iris.Context, db Database) {
	taskID := ctx.Params().GetUintDefault("id", 0)
	claims := ctx.Values().Get("claims").(jwt.MapClaims)
	userID := uint(claims["sub"].(float64))

	task, err := checkTaskOwnership(db.(*GormDatabase), taskID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.StopWithProblem(iris.StatusNotFound, iris.NewProblem().
				Title("Task not found").
				Detail(err.Error()).
				Status(iris.StatusNotFound))
			return
		}
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal server error").
			Detail(err.Error()).
			Status(iris.StatusInternalServerError))
		return
	}

	ctx.JSON(mapTaskToResponse(task))
}

func updateTask(ctx iris.Context, db Database) {
	taskID := ctx.Params().GetUintDefault("id", 0)
	claims := ctx.Values().Get("claims").(jwt.MapClaims)
	userID := uint(claims["sub"].(float64))

	var req TaskRequest
	if err := ctx.ReadJSON(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			validationErrors := wrapValidationErrors(errs)

			ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
				Title("Invalid input").
				Detail("One or more fields failed validation").
				Type("/task/validation-errors").
				Key("errors", validationErrors))
			return
		}
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid input").
			Detail(err.Error()).
			Status(iris.StatusBadRequest))
		return
	}

	task, err := checkTaskOwnership(db.(*GormDatabase), taskID, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.StopWithProblem(iris.StatusNotFound, iris.NewProblem().
				Title("Task not found").
				Detail(err.Error()).
				Status(iris.StatusNotFound))
			return
		}
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal server error").
			Detail(err.Error()).
			Status(iris.StatusInternalServerError))
		return
	}

	if _, err := validateCategory(db.(*GormDatabase), req.CategoryID); err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid category").
			Detail("Category not found: "+err.Error()).
			Status(iris.StatusBadRequest))
		return
	}

	tags, err := validateTags(db.(*GormDatabase), req.TagIDs)
	if err != nil {
		ctx.StopWithProblem(iris.StatusBadRequest, iris.NewProblem().
			Title("Invalid tags").
			Detail("Tags not found: "+err.Error()).
			Status(iris.StatusBadRequest))
		return
	}

	task.Title = req.Title
	task.Description = req.Description
	task.Completed = req.Completed
	task.CategoryID = req.CategoryID
	task.DueDate = req.DueDate
	task.Tags = tags

	if err := db.(*GormDatabase).DB.Save(&task).Error; err != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal server error").
			Detail(err.Error()).
			Status(iris.StatusInternalServerError))
		return
	}

	ctx.JSON(iris.Map{
		"message": "Task updated successfully",
		"task":    mapTaskToResponse(task),
	})
}

func deleteTask(ctx iris.Context, db Database) {
	taskID := ctx.Params().GetUintDefault("id", 0)
	claims := ctx.Values().Get("claims").(jwt.MapClaims)
	userID := uint(claims["sub"].(float64))

	result := db.(*GormDatabase).DB.Where("id = ? AND user_id = ?", taskID, userID).Delete(&Task{})
	if result.Error != nil {
		ctx.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
			Title("Internal server error").
			Detail(result.Error.Error()).
			Status(iris.StatusInternalServerError))
		return
	}
	if result.RowsAffected == 0 {
		ctx.StopWithProblem(iris.StatusNotFound, iris.NewProblem().
			Title("Task not found").
			Detail("Task with ID "+fmt.Sprint(taskID)+" not found, or you don't have access").
			Status(iris.StatusNotFound))
		return
	}

	ctx.JSON(iris.Map{
		"message": "Task deleted successfully",
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
