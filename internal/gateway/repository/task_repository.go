package repository

import (
	"fmt"
	"sync"
	
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/model"
)

type TaskRepository interface {
	Create(task *model.Task) error
	GetById(id string) (*model.Task, error)
	List() ([]*model.Task, error)
	Delete(id string) error
}

var (
	ErrTaskNotFound      = fmt.Errorf("task not found")
	ErrTaskAlreadyExists = fmt.Errorf("task already exists")
)

type inMemoryTaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*model.Task
}

func NewInMemoryTaskRepository() *inMemoryTaskRepository {
	return &inMemoryTaskRepository{
		tasks: make(map[string]*model.Task),
	}
}

func (r *inMemoryTaskRepository) Create(task *model.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, exits := r.tasks[task.ID]
	if exits {
		return ErrTaskAlreadyExists
	}
	r.tasks[task.ID] = task
	return nil
}

func (r *inMemoryTaskRepository) GetById(id string) (*model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	task, exists := r.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}
	
	return task, nil
}

func (r *inMemoryTaskRepository) List() ([]*model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tasks := make([]*model.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

func (r *inMemoryTaskRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if _, exists := r.tasks[id]; !exists {
		return ErrTaskNotFound
	}
	
	delete(r.tasks, id)
	return nil
}
