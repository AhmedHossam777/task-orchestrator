package handler

import (
	"errors"
	"net/http"
	"time"
	
	"github.com/gin-gonic/gin"
	
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/model"
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/service"
	"github.com/AhmedHossam777/task-orchestrator/pkg/apperror"
	"github.com/AhmedHossam777/task-orchestrator/pkg/response"
)

type CreateTaskRequest struct {
	Type    string         `json:"type" binding:"required,min=1,max=100"`
	Payload map[string]any `json:"payload"`
}

type TaskResponse struct {
	ID          string         `json:"id"`
	Type        string         `json:"type"`
	Payload     map[string]any `json:"payload,omitempty"`
	Status      string         `json:"status"`
	Result      map[string]any `json:"result,omitempty"`
	Error       string         `json:"error,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
}

type TaskHandler struct {
	svc service.TaskService
}

func NewTaskHandler(svc service.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	
	task, err := h.svc.CreateTask(
		model.CreateTaskRequest{
			Type:    req.Type,
			Payload: req.Payload,
		},
	)
	if err != nil {
		handleError(c, err)
		return
	}
	
	response.Created(c, toTaskResponse(task))
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	
	task, err := h.svc.GetTask(id)
	if err != nil {
		handleError(c, err)
		return
	}
	
	response.OK(c, toTaskResponse(task))
}

func (h *TaskHandler) List(c *gin.Context) {
	tasks, err := h.svc.ListTasks()
	if err != nil {
		handleError(c, err)
		return
	}
	
	taskResponses := make([]TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		taskResponses = append(taskResponses, toTaskResponse(t))
	}
	
	response.OK(c, taskResponses)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	
	if err := h.svc.DeleteTask(id); err != nil {
		handleError(c, err)
		return
	}
	
	response.OK(c, gin.H{"message": "task deleted successfully"})
}

func handleError(c *gin.Context, err error) {
	if appErr, ok := errors.AsType[*apperror.Apperror](err); ok {
		response.Fail(c, appErr.Status, appErr.Code, appErr.Message)
		return
	}
	
	response.Fail(
		c, http.StatusInternalServerError,
		"INTERNAL_ERROR", "an unexpected error occurred",
	)
}

func toTaskResponse(t *model.Task) TaskResponse {
	return TaskResponse{
		ID:          t.ID,
		Type:        t.Type,
		Payload:     t.Payload,
		Status:      string(t.Status),
		Result:      t.Result,
		Error:       t.Error,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		CompletedAt: t.CompletedAt,
	}
}
