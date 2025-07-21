package handler

import (
	"errors"

	"github.com/RLRama/listario-backend/logger"
	"github.com/RLRama/listario-backend/models"
	"github.com/RLRama/listario-backend/repository"
	"github.com/RLRama/listario-backend/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(ts service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: ts}
}

func toTaskResponse(task models.Task) models.TaskResponse {
	return models.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		Content:   task.Content,
		Completed: task.Completed,
		UserID:    task.UserID,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}

// CreateTask
// @Summary      Create a new task
// @Description  Creates a new task for the authenticated user.
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload body models.CreateTaskRequest true "Task Creation Payload"
// @Success      201 {object} models.TaskResponse
// @Failure      400 {object} object{error=string} "Invalid request format or validation failed"
// @Failure      401 {object} object{error=string} "Unauthorized"
// @Failure      500 {object} object{error=string} "Failed to create task"
// @Router       /tasks [post]
func (h *TaskHandler) CreateTask(ctx iris.Context) {
	claims := jwt.Get(ctx).(*models.UserClaims)
	userID := claims.UserID

	var req models.CreateTaskRequest
	if err := ctx.ReadJSON(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to read or validate create task request")
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid request format or validation failed", "details": err.Error()})
		return
	}

	task, err := h.taskService.CreateTask(userID, req.Title, req.Content)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create task")
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to create task"})
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(toTaskResponse(*task))
}

// GetMyTasks
// @Summary      Get all tasks for the current user
// @Description  Retrieves a list of all tasks belonging to the authenticated user.
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.TaskResponse
// @Failure      401 {object} object{error=string} "Unauthorized"
// @Failure      500 {object} object{error=string} "Could not retrieve tasks"
// @Router       /tasks [get]
func (h *TaskHandler) GetMyTasks(ctx iris.Context) {
	claims := jwt.Get(ctx).(*models.UserClaims)
	userID := claims.UserID

	tasks, err := h.taskService.GetTasksByUser(userID)
	if err != nil {
		logger.Error().Err(err).Uint("userID", userID).Msg("Failed to get tasks for user")
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Could not retrieve tasks"})
		return
	}

	response := make([]models.TaskResponse, len(tasks))
	for i, task := range tasks {
		response[i] = toTaskResponse(task)
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(response)
}

// GetTask
// @Summary      Get a single task by ID
// @Description  Retrieves details for a specific task if it belongs to the authenticated user.
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Task ID"
// @Success      200 {object} models.TaskResponse
// @Failure      401 {object} object{error=string} "Unauthorized or access denied"
// @Failure      404 {object} object{error=string} "Task not found"
// @Router       /tasks/{id} [get]
func (h *TaskHandler) GetTask(ctx iris.Context) {
	claims := jwt.Get(ctx).(*models.UserClaims)
	userID := claims.UserID

	taskID, err := ctx.Params().GetUint("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid task ID"})
		return
	}

	task, err := h.taskService.GetTask(taskID, userID)
	if err != nil {
		if errors.Is(err, service.ErrTaskAccessDenied) {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{"error": err.Error()})
		} else if errors.Is(err, repository.ErrTaskNotFound) {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{"error": err.Error()})
		} else {
			logger.Error().Err(err).Uint("taskID", taskID).Msg("Failed to get task")
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Could not retrieve task"})
		}
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(toTaskResponse(*task))
}

// UpdateTask
// @Summary      Update a task
// @Description  Updates a specific task's details if it belongs to the authenticated user.
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path  int                      true  "Task ID"
// @Param        payload body  models.UpdateTaskRequest true  "Task Update Payload"
// @Success      200 {object} models.TaskResponse
// @Failure      400 {object} object{error=string} "Invalid request format or ID"
// @Failure      401 {object} object{error=string} "Unauthorized or access denied"
// @Failure      404 {object} object{error=string} "Task not found"
// @Failure      500 {object} object{error=string} "Could not update task"
// @Router       /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(ctx iris.Context) {
	claims := jwt.Get(ctx).(*models.UserClaims)
	userID := claims.UserID

	taskID, err := ctx.Params().GetUint("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid task ID"})
		return
	}

	var req models.UpdateTaskRequest
	if err := ctx.ReadJSON(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to read or validate update task request")
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid request format or validation failed", "details": err.Error()})
		return
	}

	task, err := h.taskService.UpdateTask(taskID, userID, req.Title, req.Content, req.Completed)
	if err != nil {
		if errors.Is(err, service.ErrTaskAccessDenied) {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{"error": err.Error()})
		} else if errors.Is(err, repository.ErrTaskNotFound) {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{"error": err.Error()})
		} else {
			logger.Error().Err(err).Uint("taskID", taskID).Msg("Failed to update task")
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Could not update task"})
		}
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(toTaskResponse(*task))
}

// DeleteTask
// @Summary      Delete a task
// @Description  Deletes a specific task if it belongs to the authenticated user.
// @Tags         Tasks
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "Task ID"
// @Success      204 "No Content"
// @Failure      401 {object} object{error=string} "Unauthorized or access denied"
// @Failure      404 {object} object{error=string} "Task not found"
// @Failure      500 {object} object{error=string} "Could not delete task"
// @Router       /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(ctx iris.Context) {
	claims := jwt.Get(ctx).(*models.UserClaims)
	userID := claims.UserID

	taskID, err := ctx.Params().GetUint("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid task ID"})
		return
	}

	err = h.taskService.DeleteTask(taskID, userID)
	if err != nil {
		if errors.Is(err, service.ErrTaskAccessDenied) {
			ctx.StatusCode(iris.StatusForbidden)
			ctx.JSON(iris.Map{"error": err.Error()})
		} else if errors.Is(err, repository.ErrTaskNotFound) {
			ctx.StatusCode(iris.StatusNotFound)
			ctx.JSON(iris.Map{"error": err.Error()})
		} else {
			logger.Error().Err(err).Uint("taskID", taskID).Msg("Failed to delete task")
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(iris.Map{"error": "Could not delete task"})
		}
		return
	}

	ctx.StatusCode(iris.StatusNoContent)
}
