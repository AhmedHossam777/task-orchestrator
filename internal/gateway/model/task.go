package model

import "time"

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "PENDING"
	TaskStatusRunning   TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted TaskStatus = "COMPLETED"
	TaskStatusFailed    TaskStatus = "FAILED"
)

type Task struct {
	ID          string
	Type        string
	Payload     map[string]any
	Status      TaskStatus
	Result      map[string]any
	Error       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time // Pointer because it's nullable
}

// ? create task request represent what it needed to create a task
type CreateTaskRequest struct {
	Type    string
	Payload map[string]any
}
