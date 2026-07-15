# Tasks: 移除 SQLite 相关代码

**Feature**: 005-remove-sqlite
**Created**: 2026-07-15

## Phase 1: Remove Files & Dependencies

- [X] T001 Remove `gorm.io/driver/sqlite` and `mattn/go-sqlite3` from `go.mod`, run `go mod tidy`
- [X] T002 [P] Delete `migrations/001_user.sqlite.sql`
- [X] T003 [P] Delete `model/schema.sql`
- [X] T004 [P] Delete `model/db_test.go` (all tests depend on removed ApplySchema/EnsureUserTable)

## Phase 2: Core Code - model/db.go and model/user.go

- [X] T005 Rewrite `model/db.go`: remove SQLite init, ApplySchema, EnsureUserTable; InitDB accepts only MySQL DSN
- [X] T006 Rewrite `model/user.go`: remove SQLiteUser struct and all *SQLiteUser CRUD functions

## Phase 3: Config and Main

- [X] T007 Update `config/config.go`: remove SQLiteConfig struct, update defaults and changelog
- [X] T008 [P] Update `config.toml`: remove [database.sqlite] section
- [X] T009 Update `main.go`: remove SQLiteDB init, remove EnsureUserTable call

## Phase 4: Handler Migration

- [X] T010 Update `handler/demo.go`: switch GetUserByPhoneTest from SQLiteDB/SQLiteUser to DB/User
- [X] T011 [P] Update `handler/common.go`: remove sqlite_dsn from Setting()

## Phase 5: Test Files Cleanup

- [X] T012 Update `model/user_test.go`: remove SQLite imports, remove SQLiteUser tests, add DB nil skip
- [X] T013 Update `handler/demo_test.go`: remove SQLite imports, remove TestMain DB init, remove SQLiteUser-dependent tests

## Phase 6: Documentation

- [X] T014 Update `README.md`: remove SQLite references from tech stack, config, API list, changelog
- [X] T015 [P] Update `api.http`: update comment for /demo/user/phone to reference MySQL

## Phase 7: Validation

- [X] T016 Run `go mod tidy` and `go build ./...` — PASS
- [X] T017 Run `go vet ./...` — PASS
- [X] T018 Search for remaining "sqlite" references in source code (*.go, *.toml, *.sql) — 0 matches
