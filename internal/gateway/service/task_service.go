package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"
	
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/model"
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/repository"
	"github.com/AhmedHossam777/task-orchestrator/pkg/apperror"
)

type TaskService interface {
	CreateTask(task model.CreateTaskRequest) (*model.Task, error)
	GetTask(id string) (*model.Task, error)
	ListTasks() ([]*model.Task, error)
	DeleteTask(id string) error
}

type taskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) TaskService {
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
		return nil, apperror.Internal(
			"CREATE_FAILED",
			"failed to create task",
		)
	}
	
	return newTask, nil
}

func (s *taskService) GetTask(id string) (*model.Task, error) {
	task, err := s.repo.GetById(id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return nil, apperror.NotFound("TASK_NOT_FOUND", "task not found")
		}
		return nil, apperror.NotFound("GET_FAILED", "failed to retrieve task")
	}
	return task, nil
}

func (s *taskService) ListTasks() ([]*model.Task, error) {
	tasks, err := s.repo.List()
	if err != nil {
		return nil, apperror.Internal(
			"LIST_FAILED",
			"failed to list tasks",
		)
	}
	return tasks, nil
}

func (s *taskService) DeleteTask(id string) error {
	task, err := s.repo.GetById(id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			return apperror.NotFound(
				"TASK_NOT_FOUND",
				"no task found with the given ID",
			)
		}
		return apperror.Internal(
			"DELETE_FAILED",
			"failed to retrieve task for deletion",
		)
	}
	
	// only PENDING tasks can be deleted
	if task.Status != model.TaskStatusPending {
		return apperror.Conflict(
			"TASK_NOT_DELETABLE",
			"only tasks with PENDING status can be deleted",
		)
	}
	
	if err := s.repo.Delete(task.ID); err != nil {
		return apperror.Internal(
			"DELETE_FAILED",
			"failed to delete task",
		)
	}
	
	return nil
}

func generateID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return fmt.Sprintf("task_%x", b)
}
