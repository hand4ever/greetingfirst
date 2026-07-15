# Plan: 移除 SQLite 相关代码

**Feature**: 005-remove-sqlite
**Created**: 2026-07-15

## Tech Stack

| Category | Selection | Notes |
|----------|-----------|-------|
| Language | Go 1.26.3 | |
| Web Framework | Echo v5.2.1 | |
| ORM | GORM v1.31.2 | MySQL only |
| Database | MySQL | Architecture supports multi-DB pluggability |
| Config | TOML | `config.toml` |

## Constitution Check

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Layered Architecture | ✅ PASS | No structural changes |
| II. Unified Response Format | ✅ PASS | No response format changes |
| III. Copy-Ready Template | ✅ PASS | Simpler DB setup (single DB) |
| IV. English-Only Code Artifacts | ✅ PASS | Comments and commits in English |
| V. Test Coverage | ⚠️ PARTIAL | Tests cleaned of SQLite, pass not required (Option C) |
| VI. Fail Fast | ✅ PASS | InitDB still uses Ping fail-fast |
| VII. User-Owned Schema | ✅ PASS | EnsureUserTable removed; user runs migrations manually |

## File Changes

### Removed Files
- `migrations/001_user.sqlite.sql`
- `model/schema.sql`
- `model/db_test.go` (all tests depend on ApplySchema/EnsureUserTable)

### Modified Files

| File | Change Summary |
|------|----------------|
| `go.mod` | Remove `gorm.io/driver/sqlite` + `mattn/go-sqlite3` |
| `model/db.go` | Remove SQLite init, `ApplySchema`, `EnsureUserTable` |
| `model/user.go` | Remove `SQLiteUser` struct and all `*SQLiteUser` CRUD functions |
| `model/user_test.go` | Remove SQLite imports, remove `SQLiteUser` tests, keep `User` tests |
| `config/config.go` | Remove `SQLiteConfig` struct, update defaults |
| `config.toml` | Remove `[database.sqlite]` section |
| `main.go` | Remove `SQLiteDB` init, remove `EnsureUserTable` call |
| `handler/demo.go` | Switch from `SQLiteDB`/`SQLiteUser` to `DB`/`User` |
| `handler/demo_test.go` | Remove SQLite imports, remove TestMain DB init |
| `handler/user_test.go` | No changes needed (already uses `DB`/`User` only) |
| `handler/common.go` | Remove `sqlite_dsn` setting item |
| `README.md` | Remove SQLite references |
| `api.http` | Update comments |

## Test Strategy (Option C)

- Test files cleaned of SQLite imports
- `go build ./...` must pass
- Test pass (go test) NOT required in this scope
