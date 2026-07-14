# Tasks: 添加 CORS 跨域支持

**Input**: Design documents from `/specs/002-add-cors-support/`

**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, quickstart.md

**Tests**: Included per project Constitution V (Test Coverage requirement).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

## Path Conventions

- Project root: `/Users/bigbao/iproject/greeting/`
- Go source files at repository root, `middle/`, `handler/`, etc.
- Test files co-located with source files (e.g., `middle/cors_test.go`)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: No additional setup needed — project already initialized with Go + Echo v5. CORS uses built-in `middleware.CORS()` with zero new dependencies.

> Skipped: Project structure, dependencies, and tooling already in place.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Define CORS configuration and register middleware in the global chain. This MUST complete before any user story testing can begin.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T001 Define `corsConfig` CORS configuration variable with default values (AllowOrigins: `["*"]`, AllowMethods: `[GET,POST,PUT,DELETE,OPTIONS,PATCH,HEAD]`, AllowHeaders: `[Content-Type,Authorization,X-Requested-With,Accept,Origin]`, AllowCredentials: `false`, MaxAge: `86400`) in `main.go`
- [x] T002 Register CORS middleware using `e.Use(middleware.CORSWithConfig(corsConfig))` in the global middleware chain in `main.go` — insert after `Recover`, before `RequestID`

**Checkpoint**: CORS middleware is active with default config. All routes now respond with CORS headers for cross-origin requests.

---

## Phase 3: User Story 1 - 前端跨域调用 API (Priority: P1) 🎯 MVP

**Goal**: Frontend applications running on different origins can successfully call API endpoints via XMLHttpRequest/Fetch API. OPTIONS preflight requests return correct CORS headers, and actual requests include `Access-Control-Allow-Origin`.

**Independent Test**: Start server, send curl requests simulating browser cross-origin behavior (OPTIONS preflight + GET with Origin), verify correct CORS response headers.

### Tests for User Story 1

> **NOTE: Write these tests, ensure they PASS with the implemented CORS middleware**

- [x] T003 [P] [US1] Create test file with `logOK` helper and `TestCORSPreflight` test (OPTIONS with valid Origin and `Access-Control-Request-Method`) in `middle/cors_test.go`
- [x] T004 [P] [US1] Add `TestCORSGETWithOrigin` test (GET with Origin header returns `Access-Control-Allow-Origin`) in `middle/cors_test.go`
- [x] T005 [P] [US1] Add `TestCORSSameOrigin` test (request without Origin header does NOT add CORS headers) in `middle/cors_test.go`
- [x] T006 [US1] Add `TestCORSOptionsWithoutRequestMethod` test (OPTIONS without `Access-Control-Request-Method` header) in `middle/cors_test.go`

### Implementation for User Story 1

> Implementation already completed in Phase 2 (T001, T002). CORS middleware with default config satisfies all US1 acceptance scenarios.

- [x] T007 [US1] Run `go test -v ./middle/... -count=1` and verify all CORS tests pass

**Checkpoint**: User Story 1 fully functional. Cross-origin requests from any origin succeed with correct CORS headers. Same-origin requests remain unaffected.

---

## Phase 4: User Story 2 - 自定义允许的来源域名 (Priority: P2)

**Goal**: Developers can configure specific allowed origin domains instead of wildcard `*`, improving production security. Only matched origins receive CORS headers.

**Independent Test**: Configure specific allowed origin, verify matched origin gets `Access-Control-Allow-Origin`, unmatched origin does not.

### Tests for User Story 2

- [x] T008 [US2] Add `TestCORSSpecificOriginAllowed` and `TestCORSSpecificOriginDisallowed` tests (configure specific origin, verify allow/deny behavior) in `middle/cors_test.go`

### Implementation for User Story 2

> The config is already configurable via `corsConfig` variable (T001). US2 adds documentation and validation of the configuration pattern.

- [x] T009 [US2] Document CORS configuration customization (modify `AllowOrigins`, `AllowMethods`, `AllowHeaders`, `AllowCredentials`) in `README.md` under a new "CORS 配置" section

**Checkpoint**: User Story 2 complete. Users can modify `corsConfig` to restrict origins, and documentation explains how.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Final validation and quality assurance.

- [x] T010 Run full validation: `go fmt ./... && go build ./... && go test -v ./... -count=1`
- [x] T011 Execute all 6 quickstart.md validation scenarios manually (or via `api_test.sh`) and confirm expected results

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Skipped — project already initialized
- **Foundational (Phase 2)**: No dependencies — start immediately. BLOCKS all user stories.
- **User Story 1 (Phase 3)**: Depends on Foundational (Phase 2) completion
- **User Story 2 (Phase 4)**: Depends on Foundational (Phase 2) completion. Independently testable from US1.
- **Polish (Phase 5)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) — No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) — Uses same config but independently testable

### Within Each User Story

- Tests should be written and verified against the already-implemented CORS middleware (Phase 2)
- Test verification before moving to next story
- Story complete before moving to next priority

### Parallel Opportunities

- T003, T004, T005 can be created in parallel (different test functions, same file but non-overlapping)
- Once Foundational phase completes, US1 and US2 can proceed in parallel (if team capacity allows)

---

## Parallel Example: User Story 1

```bash
# Launch all independent test tasks for User Story 1 together:
Task: "Create test file with logOK helper and TestCORSPreflight in middle/cors_test.go"
Task: "Add TestCORSGETWithOrigin in middle/cors_test.go"
Task: "Add TestCORSSameOrigin in middle/cors_test.go"

# Then add the dependent test:
Task: "Add TestCORSOptionsWithoutRequestMethod in middle/cors_test.go"

# Finally verify:
Task: "Run go test -v ./middle/... -count=1"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 2: Foundational (T001, T002) — CORS middleware active
2. Complete Phase 3: User Story 1 (T003-T007) — Tests validate core CORS behavior
3. **STOP and VALIDATE**: All cross-origin requests work, tests pass
4. Deploy/demo — MVP is ready

### Incremental Delivery

1. Complete Foundational → CORS active with wildcard origin
2. Add User Story 1 → Tests pass → Deploy/Demo (MVP!)
3. Add User Story 2 → Specific origin tests pass → Documentation complete
4. Each story adds value without breaking previous stories

### Single Developer Strategy

Since this feature involves only 2 files (`main.go`, `middle/cors_test.go`) and 11 tasks, a single developer can complete all tasks sequentially:

1. T001 → T002 (Phase 2: ~5 min)
2. T003 → T004 → T005 → T006 → T007 (Phase 3: ~15 min)
3. T008 → T009 (Phase 4: ~10 min)
4. T010 → T011 (Phase 5: ~5 min)

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- CORS config is defined as a package-level variable in `main.go` — no separate config file needed per Constitution III (Copy-Ready Template)
