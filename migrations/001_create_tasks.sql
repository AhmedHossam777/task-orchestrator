CREATE TABLE IF NOT EXISTS tasks
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

        CONSTRAINT valid_status CHECK (status IN
                                       ('PENDING', 'IN_PROGRESS', 'COMPLETED',
                                        'FAILED'))
);

CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks (status);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks (created_at DESC);