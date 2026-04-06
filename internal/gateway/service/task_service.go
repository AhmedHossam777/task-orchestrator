package service

import (
	"crypto/rand"
	"fmt"
	"time"
	
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/model"
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/repository"
)

type TaskService interface {
	CreateTask(task model.CreateTaskRequest) (*model.Task, error)
	GetTask(id string) (*model.Task, error)
	ListTasks() ([]*model.Task, error)
	DeleteTask(id string) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{
		repo,
	}
}

func (s *taskService) CreateTask(task model.CreateTaskRequest) (
	*model.Task, error,
) {
	newTask := &model.Task{
		ID:        generateID(),
		Type:      task.Type,
		Payload:   task.Payload,
		Status:    model.TaskStatusPending,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	
	err := s.repo.Create(newTask)
	if err != nil {
		return nil, fmt.Errorf("creating task: %w", err)
	}
	
	return newTask, nil
}

func (s *taskService) GetTask(id string) (*model.Task, error) {
	task, err := s.repo.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("getting task %s: %w", id, err)
	}
	return task, nil
}

func (s *taskService) ListTasks() ([]*model.Task, error) {
	tasks, err := s.repo.List()
	if err != nil {
		return nil, fmt.Errorf("listing tasks: %w", err)
	}
	return tasks, nil
}

func (s *taskService) DeleteTask(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("Deleting Task %s: %w", id, err)
	}
	return nil
}

func generateID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return fmt.Sprintf("task_%x", b)
}
