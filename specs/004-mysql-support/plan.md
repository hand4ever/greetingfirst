# Implementation Plan: MySQL 数据库支持

**Branch**: `004-mysql-support` | **Date**: 2026-07-15 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/004-mysql-support/spec.md`

## Summary

Add MySQL support alongside existing SQLite — **both databases coexist and run independently**.
Each database has its own global GORM instance, its own model struct, and its own set of
API endpoints. No `type` selector; both databases connect simultaneously at startup.

- `model.DB` → MySQL (`model.User`)
- `model.SQLiteDB` → SQLite (`model.SQLiteUser`, renamed from existing `User`)
- MySQL CRUD: `POST/GET/PUT/DELETE /demo/usr`, `GET /demo/usrs`
- SQLite: existing `/demo` endpoints unchanged (migrated to use `SQLiteDB` + `SQLiteUser`)
- Schema managed by user-run SQL scripts (Principle VII, no AutoMigrate)

## Technical Context

**Language/Version**: Go 1.26.3 (≥ 1.22 per constitution)

**Primary Dependencies**: `gorm.io/gorm` (existing), `gorm.io/driver/sqlite` (existing), `gorm.io/driver/mysql` (NEW), `github.com/BurntSushi/toml` (existing), `github.com/labstack/echo/v5` (existing)

**Storage**: MySQL 5.7+/8.0 AND SQLite — two independent `*gorm.DB` instances (`model.DB` for MySQL, `model.SQLiteDB` for SQLite)

**Testing**: `go test` — `httptest.NewRequest` + `echo.New().NewContext` for handlers; `TestMain` with `:memory:` SQLite + schema SQL

**Target Platform**: Linux/macOS server (Echo HTTP service)

**Project Type**: web-service (Echo v5, layered architecture)

**Performance Goals**: connect both DBs < 3s (SC-001); startup failure exits < 5s (SC-005)

**Constraints**: file ≤ 500 lines; function ≤ 80 lines; indentation ≤ 3 levels (constitution code-quality)

**Scale/Scope**: minimal template project; single `users` table per DB, 5 MySQL CRUD endpoints + existing SQLite endpoints

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Layered Architecture | ✅ PASS | MySQL handlers in `handler/user.go`, SQLite handlers in `handler/demo.go`; entities in `entity/demo/` and `entity/user/`; routes in `router/`; models in `model/`. |
| II. Unified Response | ✅ PASS | All CRUD uses `response.Ok` / `response.NotOk`. |
| III. Copy-Ready Template | ✅ PASS | `InitDB` connects both DBs without AutoMigrate; DSN injected from config; each DB fully independent. |
| IV. English Code Artifacts | ✅ PASS | Comments/commit messages in English (enforced during impl). |
| V. Test Coverage | ✅ PASS (fixes existing violation) | `TestMain` in `model/user_test.go` & `handler/demo_test.go` move from `AutoMigrate` → `model.ApplySchema(...)`. New MySQL handler tests use `:memory:` SQLite for unit testing. |
| VI. Fail Fast | ✅ PASS | Either DB connection failure → fail fast (type+address error, non-zero exit). MySQL table-missing uses pause-and-continue per amended Principle VI exception (FR-011). |
| VII. User-Owned Schema | ✅ PASS | New `migrations/001_user.mysql.sql` and `migrations/001_user.sqlite.sql`; startup MUST NOT create tables. |

**Post-design re-check**: All gates PASS. No violations. No Complexity Tracking entries required.

- I: `handler/user.go` + `entity/user/` + `router/demo.go` follow layered architecture.
- II: All CRUD endpoints use `response.Ok`/`response.NotOk`.
- III: Two independent DBs, no AutoMigrate, DSN from config.
- IV: All plan comments in English.
- V: TestMain uses `ApplySchema`; new handler tests use `:memory:` SQLite.
- VI: Connection failure → fail fast; table-missing → pause-and-continue per exception.
- VII: `migrations/001_user.{mysql,sqlite}.sql`; no auto-create.

## Project Structure

### Documentation (this feature)

```text
specs/004-mysql-support/
├── plan.md              # This file
├── research.md          # Phase 0: research decisions
├── data-model.md        # Phase 1: User + SQLiteUser entities + request/response entities
├── quickstart.md        # Phase 1: validation guide
├── contracts/
│   └── api.md           # Phase 1: 5 MySQL CRUD endpoint contracts
└── tasks.md             # Phase 2 (NOT created here)
```

### Source Code (repository root) — affected files

```text
config/
├── config.go            # UPDATE: DatabaseConfig → MySQLConfig + SQLiteConfig sub-structs; remove Type/Dsn
└── config.toml          # UPDATE: [database.mysql] and [database.sqlite] sub-sections

model/
├── db.go                # REWRITE: InitDB connects both MySQL AND SQLite; DB for MySQL, SQLiteDB for SQLite; ApplySchema() helper
├── user.go              # REWRITE: rename existing User → SQLiteUser; add new User for MySQL; both have independent CRUD funcs
└── user_test.go         # UPDATE: TestMain → ApplySchema instead of AutoMigrate; test both User & SQLiteUser CRUD

handler/
├── demo.go              # UPDATE: migrate to model.SQLiteDB + model.SQLiteUser
├── user.go              # NEW: _User handler with 5 MySQL CRUD methods
└── demo_test.go         # UPDATE: TestMain → ApplySchema; ADD MySQL handler tests

entity/demo/
└── demo.go              # UNCHANGED (existing entities)

entity/user/
└── user.go              # NEW: UserCreateReq, UserPathReq, UserUpdateReq

router/
├── demo.go              # UPDATE: existing /demo routes use SQLite; ADD /demo/usr routes for MySQL

migrations/              # NEW user-owned schema (Principle VII)
├── 001_user.mysql.sql
└── 001_user.sqlite.sql

main.go                  # UPDATE: connect both DBs; EnsureUserTable for MySQL

api.http                 # ADD: 5 MySQL CRUD REST Client calls
README.md               # UPDATE: API list + changelog
```

**Structure Decision**: Single web-service project; changes follow the existing layered layout.
MySQL CRUD handlers go in dedicated `handler/user.go` to keep files under the 500-line limit.
MySQL request entities go in new `entity/user/` directory.
Schema lives in top-level `migrations/` directory.

## Migration / Schema Design

- **No AutoMigrate** (Principle VII). Tables are created manually by the user.
- `migrations/001_user.mysql.sql`: `id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY`.
- `migrations/001_user.sqlite.sql`: `id INTEGER PRIMARY KEY AUTOINCREMENT`.
- Other columns identical: `phone VARCHAR(32) NOT NULL`, `name VARCHAR(64) NOT NULL`, `age INT DEFAULT 0`, `created_at DATETIME`, `updated_at DATETIME`, `deleted_at DATETIME NULL` (indexed).
- Phone uniqueness scoped to non-deleted records via generated column `phone_active` (see research.md).
- `model.ApplySchema(db)` reads the sqlite variant and executes each statement; used only by tests.

## Implementation Notes (for tasks.md)

1. `config/config.go`:
   - Replace `DatabaseConfig{Type, DSN}` with `MySQLConfig{DSN}` and `SQLiteConfig{DSN}` sub-structs.
   - Update `defaultConfig()`: MySQL DSN defaults to `"root:password@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"`, SQLite DSN to `"greeting.db"`.
   - No `type` field.

2. `model/db.go`:
   - `InitDB(mysqlDSN, sqliteDSN string)` → opens MySQL → `model.DB = mysqlDB`, opens SQLite → `model.SQLiteDB = sqliteDB`.
   - Call `sqlDB.Ping()` on both for fail-fast; wrap errors with db type + address.
   - `ApplySchema(db *gorm.DB) error` → read `migrations/001_user.sqlite.sql`, split on `;`, skip comments, exec each. Used only by tests.
   - `EnsureUserTable(db *gorm.DB, dialect string, maxWait time.Duration) error` → if `db.Migrator().HasTable(&User{})` true, return nil. Otherwise log reminder, block/poll until table exists. `maxWait <= 0` = indefinite (used by `main.go`). MUST NOT create table.

3. `model/user.go`:
   - Rename existing `User` → `SQLiteUser`.
   - Add new `User` struct for MySQL (identical field structure, independent type).
   - Each model has its own CRUD functions using its respective DB instance (`DB` for `User`, `SQLiteDB` for `SQLiteUser`).
   - Remove `RestoreUserByPhone` (no longer needed; soft-deleted phone reuse handled by generated column).

4. `main.go`:
   - Call `model.InitDB(config.Cfg.Database.MySQL.DSN, config.Cfg.Database.SQLite.DSN)`.
   - Call `model.EnsureUserTable(model.DB, "mysql", 0)` for pause-and-continue on MySQL.

5. `handler/demo.go`:
   - Update `GetUserByPhoneTest` to use `model.SQLiteDB` + `model.SQLiteUser`.

6. `handler/user.go` (NEW):
   - Package-level `var User = &_User{}`.
   - 5 methods: `Create`, `Get`, `Update`, `Delete`, `List` — operate on `model.DB` + `model.User`.
   - Validate required fields; map GORM errors to `response.NotOk`.

7. `entity/user/user.go` (NEW):
   - `UserCreateReq`, `UserPathReq`, `UserUpdateReq`.

8. `router/demo.go`:
   - Keep existing `/demo` routes (SQLite).
   - Add `/demo/usr` routes → `handler.User` methods.

9. Tests:
   - Update `model/user_test.go` and `handler/demo_test.go` `TestMain` to use `ApplySchema`.
   - Add new `handler/user_test.go` for MySQL handler tests.
   - Add `model/user_test.go` tests for both `User` and `SQLiteUser` CRUD.

10. `api.http` + `README.md` updated per dev-flow constitution.

## Done When (Phase 1 exit criteria)

- [x] `research.md` resolves all unknowns (R-01..R-10)
- [x] `data-model.md` defines both `User` and `SQLiteUser` + request/response entities
- [x] `contracts/api.md` documents 5 MySQL CRUD endpoints under `/demo/usr`
- [x] `quickstart.md` gives runnable validation scenarios
- [x] Constitution Check gates PASS (pre + post design)
