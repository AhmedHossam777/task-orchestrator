package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/db/generated"
	"github.com/AhmedHossam777/task-orchestrator/internal/gateway/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	pool    *pgxpool.Pool
	queries *generated.Queries
}

var (
	ErrTaskNotFound      = fmt.Errorf("task not found")
	ErrTaskAlreadyExists = fmt.Errorf("task already exists")
	ErrTaskNotDeletable  = fmt.Errorf("error while deleting the task")
)

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		pool:    pool,
		queries: generated.New(pool),
	}
}

func (r *TaskRepository) Create(task *model.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	payloadJSON, err := json.Marshal(task.Payload)
	if err != nil {
		return fmt.Errorf("marshaling payload: %w", err)
	}
	
	err = r.queries.CreateTask(
		ctx, generated.CreateTaskParams{
			ID:        task.ID,
			Type:      task.Type,
			Payload:   payloadJSON,
			Status:    string(task.Status),
			CreatedAt: toPgTimestamptz(task.CreatedAt),
			UpdatedAt: toPgTimestamptz(task.UpdatedAt),
		},
	)
	if err != nil {
		return fmt.Errorf("inserting task: %w", err)
	}
	return nil
}

func (r *TaskRepository) GetById(id string) (*model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	dbTask, err := r.queries.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("querying task: %w", err)
	}
	return toDomainTask(dbTask), nil
}

func (r *TaskRepository) List() ([]*model.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	dbTasks, err := r.queries.ListTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying tasks: %w", err)
	}
	
	tasks := make([]*model.Task, 0, len(dbTasks))
	for _, dbTask := range dbTasks {
		tasks = append(tasks, toDomainTask(dbTask))
	}
	
	return tasks, nil
}

func (r *TaskRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	commandTag, err := r.queries.DeleteTask(ctx, id)
	
	if err != nil {
		return fmt.Errorf("deleting task: %w", err)
	}
	
	if commandTag.RowsAffected() == 0 {
		return ErrTaskNotFound
	}
	
	return nil
}

func toDomainTask(dbTask generated.Task) *model.Task {
	task := &model.Task{
		ID:        dbTask.ID,
		Type:      dbTask.Type,
		Status:    model.TaskStatus(dbTask.Status),
		CreatedAt: dbTask.CreatedAt.Time,
		UpdatedAt: dbTask.UpdatedAt.Time,
	}
	
	if dbTask.Payload != nil {
		_ = json.Unmarshal(dbTask.Payload, &task.Payload)
	}
	
	if dbTask.Result != nil {
		_ = json.Unmarshal(dbTask.Result, &task.Result)
	}
	
	if dbTask.Error.Valid {
		task.Error = dbTask.Error.String
	}
	
	if dbTask.CompletedAt.Valid {
		t := dbTask.CompletedAt.Time
		task.CompletedAt = &t
	}
	
	return task
}

func toPgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}
