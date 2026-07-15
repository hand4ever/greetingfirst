-- SQLite users table DDL
-- Phone uniqueness: partial unique index scoped to non-deleted rows.
-- Table name: sl_users (sqlite_ prefix is reserved by SQLite internally)

CREATE TABLE IF NOT EXISTS sl_users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    phone VARCHAR(32) NOT NULL,
    name VARCHAR(64) NOT NULL,
    age INT DEFAULT 0,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_phone_active ON sl_users (phone) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_deleted_at ON sl_users (deleted_at);
