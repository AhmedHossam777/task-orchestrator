package model

import "time"

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "Pending"
	TaskStatusRunning   TaskStatus = "Running"
	TaskStatusCompleted TaskStatus = "Completed"
	TaskStatusFailed    TaskStatus = "Failed"
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
