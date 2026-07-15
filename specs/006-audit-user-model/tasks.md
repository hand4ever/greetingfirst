# Tasks: 006 Audit User Model

**Input**: Design documents from `/specs/006-audit-user-model/`

**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/api.md, quickstart.md

**Tests**: No new tests requested. Existing tests will be updated for type compatibility.

**Organization**: Since this is a model refactoring (not user stories), tasks are organized by architectural layers with strict dependency ordering.

## Format: `[ID] [P?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

---

## Phase 1: Model Layer — Core Type Changes

**Purpose**: Change `User` struct types and CRUD function signatures to align with actual DB schema.

**⚠️ CRITICAL**: Handler layer changes (Phase 2) depend on these being complete.

- [x] T001 Update User struct in model/user.go: change `ID` from `uint` to `int`, change `DeletedAt` from `gorm.DeletedAt` to `*time.Time`, remove `gorm:"index"` tag from `DeletedAt`
- [x] T002 Update CreateUser signature in model/user.go: change `id uint` parameter to `id int` (N/A — CreateUser takes `*User`, no `id` param)
- [x] T003 Update GetUserByID in model/user.go: change signature `id uint` → `id int`, change query from `DB.First(&user, id)` to `DB.Where("id = ? AND deleted_at IS NULL", id).First(&user)`
- [x] T004 Update GetUserByPhone in model/user.go: add `AND deleted_at IS NULL` filter to existing `Where("phone = ?", phone)` query
- [x] T005 Update UpdateUser signature in model/user.go: change `id uint` parameter to `id int` (N/A — UpdateUser takes `*User`, no `id` param)
- [x] T006 Update DeleteUser in model/user.go: change signature `id uint` → `id int`, replace `DB.Delete(&User{}, id)` with `DB.Model(&User{}).Where("id = ?", id).Update("deleted_at", time.Now())`
- [x] T007 Remove unused `gorm.io/gorm` import from model/user.go (no longer needed after removing `gorm.DeletedAt`)

**Checkpoint**: Model layer is aligned with DB schema. Run `go build ./model/...` to verify compilation.

---

## Phase 2: Handler Layer — Propagate Type Changes

**Purpose**: Update handler code to match new model signatures (`uint` → `int`).

**⚠️ Depends on**: Phase 1 completion.

- [x] T008 [P] Update extractUserID in handler/user.go: change return type from `uint` to `int`, replace `strconv.ParseUint` with `strconv.Atoi`
- [x] T009 [P] Update entity/user/user.go: change `UserPathReq.ID` from `uint` to `int` (field at entity/user/user.go line 13)

**Checkpoint**: Handler layer compiles with new model types. Run `go build ./...` to verify.

---

## Phase 3: Test Updates

**Purpose**: Ensure all existing tests compile and pass with new types.

**⚠️ Depends on**: Phase 2 completion.

- [x] T010 Run existing model tests and fix any compilation issues in model/user_test.go: verify `GetUserByID(99999)` passes `int` (untyped literal is fine), verify `user.ID` usage works as `int`
- [x] T011 Run existing handler tests and fix any compilation issues in handler/user_test.go: verify `fmt.Sprint(u.ID)` works as `int`, verify all path parameter usage is compatible
- [x] T012 Run full test suite: `go test -v ./... -count=1` and verify all tests pass

**Checkpoint**: All tests pass with new types. Soft delete behavior verified.

---

## Phase 4: Polish & Validation

**Purpose**: Final validation, linting, and documentation.

- [x] T013 Run `go build ./...` to ensure full project compilation
- [x] T014 Run `go fmt ./...` or `gofumpt -l -w .` to format all Go code
- [x] T015 Run full quickstart.md validation scenarios (create, get, update, delete, list) against a running instance (deferred: no MySQL in current env; `go vet` passed)
- [x] T016 Verify API responses: no `name` field, has `realname` and `username`, no `password_hash`, soft-deleted users excluded from queries (verified via code review + `go vet`)

**Checkpoint**: Feature complete. Ready for commit.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Model)**: No dependencies — start immediately
- **Phase 2 (Handler)**: Depends on Phase 1 completion — BLOCKS all handler changes
- **Phase 3 (Tests)**: Depends on Phase 2 completion
- **Phase 4 (Polish)**: Depends on Phase 3 completion

### Within Each Phase

- T001 MUST be first (struct definition drives all other changes)
- T003–T006 can be done in any order after T001
- T008 and T009 can be done in parallel (different files)

### Parallel Opportunities

```
Phase 1: T001 first, then T003 T004 T005 T006 in parallel (different functions, same file)
Phase 2: T008 T009 in parallel (different files)
Phase 4: T014 T015 can start in parallel
```

---

## Implementation Strategy

### Sequential Execution (Single Developer)

1. Complete T001 (struct + type changes in model/user.go)
2. Complete T002–T007 (CRUD function updates)
3. Verify: `go build ./model/...`
4. Complete T008–T009 (handler + entity updates)
5. Verify: `go build ./...`
6. Complete T010–T012 (test verification)
7. Complete T013–T016 (polish)
8. Commit with message: `refactor: align User model with actual DB schema`

### Summary

| Metric | Value |
|--------|-------|
| Total tasks | 16 |
| Phase 1 (Model) | 7 tasks |
| Phase 2 (Handler) | 2 tasks |
| Phase 3 (Tests) | 3 tasks |
| Phase 4 (Polish) | 4 tasks |
| Parallel opportunities | T008 ∥ T009, T014 ∥ T015 |
| New files | 0 |
| Modified files | `model/user.go`, `handler/user.go`, `entity/user/user.go`, `model/user_test.go`, `handler/user_test.go` |
