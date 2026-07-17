-- SQLite test_user table DDL
-- User-owned schema asset: run manually against greeting.db before starting the app:
--   sqlite3 greeting.db < migrations/002_test_user.sql
-- The application MUST NOT create or migrate this table (constitution principle VII).
-- Unit tests reuse this exact script in TestMain against an in-memory database.

CREATE TABLE IF NOT EXISTS test_user (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        VARCHAR(100) NOT NULL DEFAULT '',
    phone       VARCHAR(20)  NOT NULL DEFAULT '',
    age         INTEGER      NOT NULL DEFAULT 0,
    created_at  DATETIME,
    updated_at  DATETIME,
    deleted_at  DATETIME
);
