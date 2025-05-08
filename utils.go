package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ══════════════════════════ Common utilities ══════════════════════════

func wrapValidationErrors(errs validator.ValidationErrors) []validationError {
	validationErrors := make([]validationError, 0, len(errs))
	for _, validationErr := range errs {
		validationErrors = append(validationErrors, validationError{
			ActualTag: validationErr.ActualTag(),
			Namespace: validationErr.Namespace(),
			Kind:      validationErr.Kind().String(),
			Type:      validationErr.Type().String(),
			Value:     fmt.Sprintf("%v", validationErr.Value()),
			Param:     validationErr.Param(),
		})
	}

	return validationErrors
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func generateJWTToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func registerCustomValidators(v *validator.Validate) {
	v.RegisterValidation("specialchar", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		specialChars := "!@#$%^&*()-_=+[]{}|;:',.<>?/"
		for _, ch := range value {
			if strings.ContainsRune(specialChars, ch) {
				return true
			}
		}
		return false
	})
}

func configureCORS() iris.Handler {
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	var allowedOrigins []string
	if allowedOriginsStr != "" {
		allowedOrigins = strings.Split(allowedOriginsStr, ",")
	} else {
		allowedOrigins = []string{"*"}
	}

	return cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})
}

// ══════════════════════════ Task utilities ══════════════════════════

func validateCategory(db *GormDatabase, categoryID uint) (*Category, error) {
	var category Category
	if err := db.First(&category, categoryID).Error; err != nil {
		return nil, fmt.Errorf("Category does not exist: %w", err)
	}
	return &category, nil
}

func validateTags(db *GormDatabase, tagIDs []uint) ([]Tag, error) {
	var tags []Tag
	if len(tagIDs) > 0 {
		if err := db.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
			return nil, fmt.Errorf("tags do not exist: %w", err)
		}
		if len(tags) != len(tagIDs) {
			return nil, fmt.Errorf("some tag IDs are invalid")
		}
	}
	return tags, nil
}

func mapTaskToResponse(task *Task) TaskResponse {
	tagNames := make([]string, len(task.Tags))
	for i, tag := range task.Tags {
		tagNames[i] = tag.Name
	}

	return TaskResponse{
		ID:           task.ID,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		UserID:       task.UserID,
		CategoryName: task.Category.Name,
		DueDate:      task.DueDate,
		TagNames:     tagNames,
		CreatedAt:    task.CreatedAt,
		UpdatedAt:    task.UpdatedAt,
	}
}

func checkTaskOwnership(db *GormDatabase, taskID, userID uint) (*Task, error) {
	var task Task
	if err := db.DB.Preload("Tags").Preload("Category").Where("id = ? AND user_id = ?", taskID, userID).First(&task).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("task does not exist or does not belong to the user")
		}
		return nil, fmt.Errorf("failed to check task ownership: %w", err)
	}
	return &task, nil
}
