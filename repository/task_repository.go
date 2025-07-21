package repository

import (
	"errors"

	"github.com/RLRama/listario-backend/models"
	"gorm.io/gorm"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type TaskRepository interface {
	Create(task *models.Task) error
	FindByID(id uint) (*models.Task, error)
	FindByUser(userID uint) ([]models.Task, error)
	Update(task *models.Task) error
	Delete(id uint) error
}

type gormTaskRepository struct {
	db *gorm.DB
}

func NewGormTaskRepository(db *gorm.DB) TaskRepository {
	return &gormTaskRepository{db: db}
}

func (r *gormTaskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *gormTaskRepository) FindByID(id uint) (*models.Task, error) {
	var task models.Task
	result := r.db.First(&task, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrTaskNotFound
	}
	return &task, result.Error
}

func (r *gormTaskRepository) FindByUser(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	result := r.db.Where("user_id = ?", userID).Find(&tasks)
	return tasks, result.Error
}

func (r *gormTaskRepository) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

func (r *gormTaskRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Task{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}
