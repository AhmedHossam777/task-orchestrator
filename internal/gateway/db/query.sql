-- name: CreateTask :exec
INSERT INTO tasks (id, type, payload, status, created_at, updated_at)
VALUES (sqlc.arg(id),
        sqlc.arg(type),
        sqlc.arg(payload),
        sqlc.arg(status),
        sqlc.arg(created_at),
        sqlc.arg(updated_at));


-- name: GetTaskByID :one
SELECT id, type, payload, status, result, error,
       created_at, updated_at, completed_at
FROM tasks
WHERE id = sqlc.arg(id);

-- name: ListTasks :many
SELECT id, type, payload, status, result, error,
       created_at, updated_at, completed_at
FROM tasks
ORDER BY created_at DESC;

-- name: DeleteTask :execresult
DELETE FROM tasks
WHERE id = sqlc.arg(id);