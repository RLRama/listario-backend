package service

import (
	"errors"

	"github.com/RLRama/listario-backend/models"
	"github.com/RLRama/listario-backend/repository"
)

var (
	ErrTaskAccessDenied = errors.New("access to the requested task is denied")
)

type TaskService interface {
	CreateTask(userID uint, title, content string) (*models.Task, error)
	GetTask(taskID, userID uint) (*models.Task, error)
	GetTasksByUser(userID uint) ([]models.Task, error)
	UpdateTask(taskID, userID uint, title, content string, completed *bool) (*models.Task, error)
	DeleteTask(taskID, userID uint) error
}

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{
		taskRepo: repo,
	}
}

func (s *taskService) CreateTask(userID uint, title, content string) (*models.Task, error) {
	task := &models.Task{
		Title:   title,
		Content: content,
		UserID:  userID,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) GetTask(taskID, userID uint) (*models.Task, error) {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, err
	}

	if task.UserID != userID {
		return nil, ErrTaskAccessDenied
	}
	return task, nil
}

func (s *taskService) GetTasksByUser(userID uint) ([]models.Task, error) {
	return s.taskRepo.FindByUser(userID)
}

func (s *taskService) UpdateTask(taskID, userID uint, title, content string, completed *bool) (*models.Task, error) {
	task, err := s.GetTask(taskID, userID)
	if err != nil {
		return nil, err
	}

	if title != "" {
		task.Title = title
	}
	if content != "" || (content == "" && task.Content != "") {
		task.Content = content
	}
	if completed != nil {
		task.Completed = *completed
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) DeleteTask(taskID, userID uint) error {
	task, err := s.GetTask(taskID, userID)
	if err != nil {
		return err
	}
	return s.taskRepo.Delete(task.ID)
}
