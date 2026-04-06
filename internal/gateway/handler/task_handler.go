package handler

import (
	"net/http"
	"time"
	
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/model"
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/service"
	"github.com/AhmedHossam777/task-orchestrator/pkg/response"
	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	Type    string         `json:"type" binding:"required"`
	Payload map[string]any `json:"payload"`
}

type TaskHandler struct {
	taskService service.TaskService
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

func newTaskHandler(svc service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: svc,
	}
}

func (h *TaskHandler) CreateTaskHandler(c *gin.Context) {
	var req CreateTaskRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}
	task, err := h.taskService.CreateTask(
		model.CreateTaskRequest{
			Type:    req.Type,
			Payload: req.Payload,
		},
	)
	
	if err != nil {
		response.Fail(
			c, http.StatusInternalServerError, "CREATE_FAILED", err.Error(),
		)
		return
	}
	
	response.Created(c, toTaskResponse(task))
}

func (h *TaskHandler) GetAllTasks() {
	
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
