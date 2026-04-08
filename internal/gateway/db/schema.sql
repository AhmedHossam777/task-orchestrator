-- This file is read by sqlc to understand table structure.
-- Keep it in sync with your migrations.

CREATE TABLE tasks
(
    id           VARCHAR(36) PRIMARY KEY,
    type         VARCHAR(100) NOT NULL,
    payload      JSONB,
    status       VARCHAR(20)  NOT NULL DEFAULT 'PENDING',
    result       JSONB,
    error        TEXT,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);