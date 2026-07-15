# Quickstart: MySQL 数据库支持

**Feature**: `004-mysql-support` | **Date**: 2026-07-15

End-to-end validation guide. Implementation details belong in `tasks.md`.

## Prerequisites

- Go ≥ 1.22
- MySQL 5.7+/8.0 (for MySQL endpoints)
- No external DB needed for SQLite endpoints or test suite

## A. Run the test suite (no external DB required)

```bash
go test -v ./... -count=1
```

**Expected**: all model + handler tests pass. `model/user_test.go` and `handler/demo_test.go` `TestMain` build tables via the canonical schema SQL (`model.ApplySchema`), **not** `AutoMigrate`. Both `User` (MySQL model) and `SQLiteUser` (SQLite model) CRUD tests pass with `:memory:` SQLite.

## B. Validate MySQL + SQLite coexistence (requires local MySQL)

1. Start MySQL and create the database:
   ```sql
   CREATE DATABASE demo CHARACTER SET utf8mb4;
   ```

2. Apply the MySQL schema:
   ```bash
   mysql -uroot -p demo < migrations/001_user.mysql.sql
   ```

3. Ensure `config.toml` has both database sections (no `type` field):
   ```toml
   [database.mysql]
   dsn = "root:password@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"

   [database.sqlite]
   dsn = "greeting.db"
   ```

4. Start the server:
   ```bash
   go run .
   ```
   **Expected**: log shows `[model] connected to mysql at 127.0.0.1:3306/demo` AND `[model] connected to sqlite at greeting.db`, server listens on `:1323` within ~3s (SC-001).

5. Exercise MySQL CRUD via `api.http` or curl:
   ```bash
   # Create
   curl -X POST http://localhost:1323/demo/usr -H "Content-Type: application/json" -d '{"name":"张三","phone":"13800138000","age":25}'
   # Get (use the returned id)
   curl http://localhost:1323/demo/usr/1
   # Update
   curl -X PUT http://localhost:1323/demo/usr/1 -H "Content-Type: application/json" -d '{"name":"张三丰","age":30}'
   # List
   curl http://localhost:1323/demo/usrs
   # Delete
   curl -X DELETE http://localhost:1323/demo/usr/1
   ```

6. Verify SQLite endpoints still work independently:
   ```bash
   curl http://localhost:1323/demo/search?tag=go
   curl http://localhost:1323/demo/user/phone
   curl http://localhost:1323/demo/sha256?text=hello
   ```

## C. Validate fail-fast on connection failure (SC-005)

1. Set `[database.mysql] dsn` to an unreachable host or wrong port.
2. `go run .` → exits within ~5s with error mentioning MySQL + address. Even if SQLite connects fine, MySQL failure causes exit (both DBs must connect).

## D. Validate missing-table pause-and-continue (FR-011 / SC-006)

1. Point `[database.mysql]` at a reachable `demo` DB without the `users` table.
2. `go run .` → service prints reminder:
   ```
   [WARN] users table not found in mysql database.
          Please create it manually:
            mysql -u root -p demo < migrations/001_user.mysql.sql
          The service will continue automatically once the table is created.
   ```
   and then **pauses** (does NOT exit, does NOT auto-create).
3. In another shell:
   ```bash
   mysql -uroot -p demo < migrations/001_user.mysql.sql
   ```
4. Service detects the table, logs continuation, and starts listening. CRUD now works.

## Schema assets

- `migrations/001_user.mysql.sql` — MySQL DDL (run manually, Principle VII)
- `migrations/001_user.sqlite.sql` — SQLite DDL (used by tests, Principle VII)
