# Tasks: SHA256 Demo API

**Input**: Design documents from `/specs/001-sha256-demo-api/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Included per Constitution V — every handler must have unit tests.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1)
- Include exact file paths in descriptions

## Path Conventions

- Go project with Echo v5 layered architecture
- `entity/<module>/` → request structs
- `handler/` → request handler + test files
- `router/` → route registration
- `api.http` → REST Client test cases
- `README.md` → API documentation

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Ensure module directory structure exists

- [x] T001 Create entity/demo/ directory in entity/demo/

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Nothing needed — project is fully initialized (Go + Echo v5 + GORM already configured). No database migration required (stateless endpoint).

**⚠️ SKIP**: This feature is stateless and requires no foundational infrastructure beyond what already exists.

**Checkpoint**: Foundation ready — proceed directly to User Story 1.

---

## Phase 3: User Story 1 - Compute SHA256 Hash (Priority: P1) 🎯 MVP

**Goal**: 调用方通过 GET `/demo/sha256?text=<input>` 获取 SHA256 哈希值，响应包含 `input` 和 `hash` 两个字段。

**Independent Test**: `curl "http://localhost:1323/demo/sha256?text=hello"` 返回 `{ "input": "hello", "hash": "2cf24dba..." }`

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T002 [P] [US1] Create Sha256Request struct with `query:"text"` tag in entity/demo/demo.go
- [x] T003 [P] [US1] Write handler unit tests with logOK helper covering: basic hash, Chinese text, empty string, missing param in handler/demo_sha256_test.go
- [x] T004 [US1] Run `go test -v ./handler/ -run TestSha256 -count=1` — confirm tests FAIL (red phase)

### Implementation for User Story 1

- [x] T005 [US1] Implement Sha256Response struct in entity/demo/demo.go (add Input and Hash fields with json tags)
- [x] T006 [US1] Implement _Sha256 handler with Compute method using crypto/sha256 in handler/demo_sha256.go
- [x] T007 [US1] Run `go test -v ./handler/ -run TestSha256 -count=1` — confirm tests PASS (green phase)
- [x] T008 [US1] Register GET /demo/sha256 route in router/router.go (demo group, bind to handler.Sha256.Compute)
- [x] T009 [US1] Add REST Client test cases for SHA256 demo in api.http
- [x] T010 [US1] Update README.md: add /demo/sha256 to API list and changelog entries

**Checkpoint**: At this point, User Story 1 should be fully functional — `go run main.go` + `curl localhost:1323/demo/sha256?text=hello` returns correct SHA256 hash.

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and cleanup

- [x] T011 Run `go build ./...` to verify compilation
- [x] T012 Run quickstart.md all 6 validation scenarios
- [x] T013 Run full test suite: `go test -v ./... -count=1`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: SKIP — project already initialized
- **User Story 1 (Phase 3)**: Depends on Phase 1 (entity/demo/ directory) completion
- **Polish (Phase 4)**: Depends on Phase 3 completion

### User Story Dependencies

- **User Story 1 (P1)**: No dependencies on other stories — only story in this feature

### Within User Story 1

- Tests (T002, T003, T004) MUST be written and FAIL before implementation
- Entity structs (T002, T005) before handler (T006)
- Handler (T006) before route registration (T008)
- Route (T008) before REST Client test (T009) and docs (T010)
- All implementation complete before Polish phase

### Parallel Opportunities

- T002 (entity struct) and T003 (handler test) can be done in parallel (different files)
- T009 (api.http) and T010 (README.md) can be done in parallel

---

## Parallel Example: User Story 1

```bash
# Step 1: Create entity struct + write test (parallel)
Task: "Create Sha256Request struct in entity/demo/demo.go"
Task: "Write handler unit tests in handler/demo_sha256_test.go"

# Step 2: Confirm tests fail (serial)
Task: "Run tests — confirm FAIL"

# Step 3: Implement handler + response entity (serial, depends on Step 1)
Task: "Implement handler in handler/demo_sha256.go"

# Step 4: Register route (serial)

# Step 5: Update docs (parallel)
Task: "Add REST Client test in api.http"
Task: "Update README.md"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (create entity/demo/)
2. SKIP Phase 2: Foundational (not needed)
3. Complete Phase 3: User Story 1 (TDD: red → green → register → docs)
4. **STOP and VALIDATE**: Run quickstart.md scenarios
5. Complete Phase 4: Polish
6. Feature ready for commit

### Incremental Delivery

Since this feature has only one user story, MVP = full feature:
1. Create entity struct + write tests → FAIL
2. Implement handler → tests PASS
3. Register route + update docs → complete
4. Polish → done

---

## Notes

- [P] tasks = different files, no dependencies
- [US1] label maps task to the single user story
- No database migrations needed (stateless)
- No external dependencies needed (crypto/sha256 is stdlib)
- Follow Constitution: English comments, unified response format, logOK test helper
- Commit after each task or logical group
- TDD: tests fail first, then implement, then pass
