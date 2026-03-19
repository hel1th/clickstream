CREATE TABLE IF NOT EXISTS users (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL    DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_created_at ON users (created_at);