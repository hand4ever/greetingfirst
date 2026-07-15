# Tasks: MySQL 数据库支持

**Input**: Design documents from `/specs/004-mysql-support/`

**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/

**Tests**: Required per Constitution Principle V — every handler and model method must have unit tests.

**Organization**: Tasks are grouped by user story. US1 + US2 (P1) form the foundational phase; US3 (P2) and US4 (P2) build on top.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

All paths relative to repository root `greeting/` (Go Echo v5 project with layered architecture):
- `config/` — configuration structs + TOML file
- `model/` — GORM models + DB init logic
- `handler/` — HTTP request handlers
- `entity/` — request/response entity structs
- `router/` — route registration
- `migrations/` — SQL schema files (user-owned, Principle VII)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization — dependencies, directories, and SQL schema assets

- [ ] T001 Install `gorm.io/driver/mysql` dependency via `go get gorm.io/driver/mysql`
- [ ] T002 [P] Create `migrations/001_user.mysql.sql` with MySQL DDL (BIGINT UNSIGNED AUTO_INCREMENT, phone_active generated column, utf8mb4)
- [ ] T003 [P] Create `migrations/001_user.sqlite.sql` with SQLite DDL (INTEGER PRIMARY KEY AUTOINCREMENT, phone_active generated column)
- [ ] T004 [P] Create `entity/user/` directory for MySQL request entities

---

## Phase 2: Foundational — US1 + US2 + US4 (Blocking Prerequisites)

**Purpose**: Dual-database coexistence infrastructure. Implements US1 (dual DB connection), US2 (config management), and US4 (user-owned schema + pause-and-continue). ALL user stories depend on this phase.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Configuration (US2)

- [ ] T005 Update `config/config.go` — replace `DatabaseConfig{Type, DSN}` with `MySQLConfig{DSN}` and `SQLiteConfig{DSN}` sub-structs; update `Config` struct and `defaultConfig()` with MySQL default DSN (`root:password@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local`) and SQLite default DSN (`greeting.db`)
- [ ] T006 [P] Update `config.toml` — replace `[database] type/dsn` with `[database.mysql]` and `[database.sqlite]` independent sections; remove `type` field

### Database Initialization (US1)

- [ ] T007 Rewrite `model/db.go` — declare `var DB *gorm.DB` (MySQL) and `var SQLiteDB *gorm.DB` (SQLite); implement `InitDB(mysqlDSN, sqliteDSN string) error` that opens both with `sqlDB.Ping()` fail-fast, wrapping errors with db type + address; implement `ApplySchema(db *gorm.DB) error` that reads `migrations/001_user.sqlite.sql`, splits on `;`, skips comments/blank lines, execs each; implement `EnsureUserTable(db *gorm.DB, dialect string, maxWait time.Duration) error` that logs reminder, blocks/polls with `schemaPollInterval` (default 3s) until table exists, returns nil on success; MUST NOT auto-create tables

### Model Layer (US1)

- [ ] T008 Rewrite `model/user.go` — rename existing `User` struct to `SQLiteUser`; add new `User` struct for MySQL (identical field structure with `gorm` tags); add `model.SQLiteDB`-backed CRUD for `SQLiteUser` (CreateSQLiteUser, GetSQLiteUserByID, GetSQLiteUserByPhone, UpdateSQLiteUser, DeleteSQLiteUser); add `model.DB`-backed CRUD for `User` (CreateUser, GetUserByID, GetUserByPhone, UpdateUser, DeleteUser); remove `RestoreUserByPhone`

### Main Entry (US1)

- [ ] T009 Update `main.go` — call `model.InitDB(config.Cfg.Database.MySQL.DSN, config.Cfg.Database.SQLite.DSN)` instead of current hardcoded init; call `model.EnsureUserTable(model.DB, "mysql", 0)` after init for MySQL pause-and-continue

### SQLite Handler Migration (US1, FR-012)

- [ ] T010 Update `handler/demo.go` — migrate `GetUserByPhoneTest` to use `model.SQLiteDB` + `model.SQLiteUser` instead of `model.DB` + `model.User`
- [ ] T011 Update `handler/demo_test.go` — update `TestMain` to use `model.ApplySchema` instead of `DB.AutoMigrate` (Constitution Principle V + VII); update existing test to use `SQLiteDB` + `SQLiteUser`

### Model Tests (US1)

- [ ] T012 Rewrite `model/user_test.go` — update `TestMain` to init `:memory:` SQLite DBs and call `ApplySchema` instead of `AutoMigrate`; add tests for `SQLiteUser` CRUD (CreateSQLiteUser, GetSQLiteUserByID, GetSQLiteUserByPhone, UpdateSQLiteUser, DeleteSQLiteUser); add tests for `User` CRUD (CreateUser, GetUserByID, GetUserByPhone, UpdateUser, DeleteUser); include `logOK` helper; phone-uniqueness test for active records
- [ ] T013 [P] Test `model/db.go` — create `model/db_test.go` with tests for `ApplySchema` (verify table creation from SQL file) and `EnsureUserTable` (verify pause-and-continue polling behavior with mock/short timeout); include `logOK` helper

**Checkpoint**: Foundation ready — dual DBs connect, config manages both, SQLite endpoints work, all model tests pass

---

## Phase 3: User Story 3 — MySQL Demo 模块完整 CRUD 接口 (Priority: P2) 🎯

**Goal**: Expose 5 MySQL CRUD endpoints under `/demo/usr` prefix (POST/GET/PUT/DELETE /demo/usr, GET /demo/usrs). SQLite `/demo` endpoints continue working independently.

**Independent Test**: Use HTTP client to call Create → Get → Update → Get → Delete → Get(404) on `/demo/usr`, verifying data consistency at each step. SQLite `/demo/user/phone` continues to work independently.

### Request Entities

- [ ] T014 [P] [US3] Create `entity/user/user.go` — `UserCreateReq{Name string, Phone string, Age *int}` with json tags, `UserPathReq{ID uint}` with param tag, `UserUpdateReq{Name string, Age *int}` with json tags

### Handler Implementation

- [ ] T015 [US3] Create `handler/user.go` — declare `var User = &_User{}`; implement `_User` with 5 methods:
  - `Create(*echo.Context) error` — bind `UserCreateReq`, validate name/phone non-empty, call `model.CreateUser`, handle duplicate phone `response.NotOk("phone already exists")`, return `response.Ok`
  - `Get(*echo.Context) error` — bind `UserPathReq`, call `model.GetUserByID`, handle `ErrRecordNotFound` → `response.NotOk("user not found")`, return `response.Ok`
  - `Update(*echo.Context) error` — bind `UserPathReq` + `UserUpdateReq`, fetch user, update only provided fields (zero-value int → skip with pointer), call `model.UpdateUser`, handle not-found
  - `Delete(*echo.Context) error` — bind `UserPathReq`, call `model.DeleteUser`, handle not-found, return `response.Ok(c, "")`
  - `List(*echo.Context) error` — query `model.DB.Where("deleted_at IS NULL").Order("created_at DESC").Find(&users)`, return `response.Ok`

### Route Registration

- [ ] T016 [US3] Update `router/demo.go` — add `/demo/usr` route group: `POST /demo/usr` → `handler.User.Create`, `GET /demo/usr/:id` → `handler.User.Get`, `PUT /demo/usr/:id` → `handler.User.Update`, `DELETE /demo/usr/:id` → `handler.User.Delete`, `GET /demo/usrs` → `handler.User.List`

### Handler Tests

- [ ] T017 [US3] Create `handler/user_test.go` — `TestMain` with `:memory:` SQLite + `ApplySchema`; logOK helper; tests for all 5 endpoints:
  - `TestCreate_Success` / `TestCreate_MissingName` / `TestCreate_DuplicatePhone`
  - `TestGet_Success` / `TestGet_NotFound`
  - `TestUpdate_Success` / `TestUpdate_NotFound` / `TestUpdate_PartialFields`
  - `TestDelete_Success` / `TestDelete_NotFound`
  - `TestList_Empty` / `TestList_WithUsers`
  - `TestCreate_Concurrent` — launch 10 goroutines with different phones, verify all 10 succeed without data loss (SC-004)

### Documentation & Integration

- [ ] T018 [P] [US3] Update `api.http` — add REST Client test cases for all 5 MySQL CRUD endpoints (POST /demo/usr, GET /demo/usr/:id, PUT /demo/usr/:id, DELETE /demo/usr/:id, GET /demo/usrs)
- [ ] T019 [P] [US3] Update `README.md` — add MySQL CRUD endpoints to API list section; add changelog entry for MySQL support

**Checkpoint**: MySQL CRUD fully functional, all handler tests pass, SQLite endpoints still work independently

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Final validation, formatting, and build verification

- [ ] T020 Run `go fmt ./...` (or `gofumpt`) to format all Go code
- [ ] T021 Run `go build ./...` to verify compilation
- [ ] T022 Run `go test -v ./... -count=1` to verify all tests pass
- [ ] T023 Run quickstart.md validation: start server with `config.toml`, exercise MySQL CRUD via curl, verify SQLite `/demo/user/phone` still works

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Phase 1 (T002, T003 SQL files needed by T007 ApplySchema) — **BLOCKS** all user stories
- **US3 (Phase 3)**: Depends on Phase 2 completion (needs `model.User`, `model.DB`, config)
- **Polish (Phase 4)**: Depends on Phase 3 completion

### User Story Dependencies

- **US1 (P1) — Dual DB Connection**: Implemented in Phase 2 (T007, T008, T009); no upstream user story dependencies; blocks US3
- **US2 (P1) — Config Management**: Implemented in Phase 2 (T005, T006); no upstream dependencies; blocks US3
- **US3 (P2) — MySQL CRUD**: Depends on Phase 2 (US1+US2); independently testable via `/demo/usr`
- **US4 (P2) — Manual Schema**: Implemented in Phase 2 (T002, T003 SQL files; T007 EnsureUserTable); verified via T013 model tests

### Within Each Phase

- Phase 2: T005 (config.go) → T006 (config.toml) can parallel; T007 (db.go) → T008 (user.go) → T010 (handler/demo.go) are sequential
- Phase 3: T014 (entities) → T015 (handler) → T016 (router) → T017 (tests); T018, T019 can parallel after T016

### Parallel Opportunities

- Phase 1: T002, T003, T004 can all run in parallel
- Phase 2: T005 + T006 (config) can run in parallel with T012 (model tests) since they touch different files; T013 can parallel after T007
- Phase 3: T014 (entities) can start in parallel with Phase 2 tail tasks; T018 + T019 parallel after T016
- Phase 4: T020, T021, T022 run sequentially; T023 runs after T021

---

## Parallel Example: Phase 2 Foundational

```bash
# Parallel batch 1 — config + tests can start simultaneously:
Task: "T005 Update config/config.go"
Task: "T006 Update config.toml"
Task: "T012 Rewrite model/user_test.go" (needs T007 db.go, but can develop test structure)

# Parallel batch 2 — after db.go ready:
Task: "T008 Rewrite model/user.go"
Task: "T013 Test model/db.go"
```

## Parallel Example: Phase 3 US3

```bash
# After T015 handler done, test + docs can run in parallel:
Task: "T017 Create handler/user_test.go"
Task: "T016 Update router/demo.go"
Task: "T018 Update api.http"
Task: "T019 Update README.md"
```

---

## Implementation Strategy

### MVP First (Phase 1 + Phase 2 Only)

1. Complete Phase 1: Setup (deps, SQL files, directory)
2. Complete Phase 2: Foundational (dual DB connect, config, model migration, tests)
3. **STOP and VALIDATE**: Run `go test -v ./... -count=1` — all model tests pass; `go build ./...` compiles; `go run .` connects both DBs
4. This is the MVP: dual-database infrastructure is ready

### Incremental Delivery

1. Complete Setup + Foundational → Dual DB foundation ready ✅
2. Add US3 (MySQL CRUD) → Test independently → Full feature complete ✅
3. Polish → Format, build, run full test suite ✅

### Single Developer Strategy

1. Phase 1 (T001→T004): sequential but fast
2. Phase 2 (T005→T013): T005→T006 config; T007 db.go; T008 user.go; T009 main.go; T010 handler/demo.go; T011 demo_test.go; T012 user_test.go; T013 db_test.go
3. Phase 3 (T014→T019): T014 entities; T015 handler; T016 router; T017 handler tests; T018+T019 docs in parallel
4. Phase 4 (T020→T023): fmt → build → test → quickstart validate

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Tests MUST use `ApplySchema` not `AutoMigrate` (Principle VII)
- Tests MUST include `logOK` helper per Constitution Principle V
- Commit after each phase checkpoint
- `go test -v ./... -count=1` must pass before considering any phase complete
- US3 Update handler: use `*int` for `Age` in `UserUpdateReq` (or use a dedicated update struct with pointer) to distinguish "not provided" (nil) from "set to 0"
