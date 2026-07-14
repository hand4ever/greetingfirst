# Tasks: 添加 CORS 跨域支持

**Input**: Design documents from `specs/002-add-cors-support/`

**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, quickstart.md

**Tests**: Tests are included — Constitution V requires test coverage for all features.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing.

**Status Note**: CORS middleware implementation is already complete in `main.go`. Tasks below focus on verification, remaining gap closure, and documentation.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)
- Include exact file paths in descriptions

---

## Phase 1: Setup

**Purpose**: Verify project is ready for CORS feature validation

- [X] T001 Verify Echo v5 CORS middleware import exists in `main.go` (should already have `github.com/labstack/echo/v5/middleware`)
- [X] T002 Verify `corsConfig` variable is defined with correct defaults in `main.go` (AllowOrigins, AllowMethods, AllowHeaders, AllowCredentials, MaxAge)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Confirm middleware registration order and basic structure before user story verification

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T003 Verify CORS middleware is registered in correct position in `main.go` (after Recover, before RequestID)
- [X] T004 Verify `middle/cors_test.go` exists with `logOK` helper and `newCORSApp` test helper function

**Checkpoint**: Foundation verified — user story verification can now begin

---

## Phase 3: User Story 1 - 前端跨域调用 API (Priority: P1) 🎯 MVP

**Goal**: Frontend applications from different origins can call the API via cross-origin requests

**Independent Test**: Run `go test -v ./middle/... -count=1 -run "TestCORSPreflight|TestCORSGETWithOrigin|TestCORSSameOrigin|TestCORSOptionsWithoutRequestMethod"` and validate all pass

### Tests for User Story 1

> **NOTE**: These tests already exist in `middle/cors_test.go`. Verify they pass.

- [X] T005 [P] [US1] Verify `TestCORSPreflight` test in `middle/cors_test.go` — OPTIONS + Origin + ACRM → 200/204 with CORS headers
- [X] T006 [P] [US1] Verify `TestCORSGETWithOrigin` test in `middle/cors_test.go` — GET + Origin → response includes Access-Control-Allow-Origin
- [X] T007 [P] [US1] Verify `TestCORSSameOrigin` test in `middle/cors_test.go` — no Origin → no CORS headers added
- [X] T008 [P] [US1] Verify `TestCORSOptionsWithoutRequestMethod` test in `middle/cors_test.go` — OPTIONS + Origin without ACRM → handled gracefully

### Implementation for User Story 1

- [X] T009 [US1] Confirm `corsConfig` in `main.go` has `AllowOrigins: []string{"*"}` for US1 default behavior (allow all origins)
- [X] T010 [US1] Run `go test -v ./middle/... -count=1 -run TestCORS` and verify all 6 CORS tests pass
- [X] T011 [US1] Run quickstart Scenarios 1-5 from `specs/002-add-cors-support/quickstart.md` against running server to validate end-to-end

**Checkpoint**: User Story 1 verified — all cross-origin requests work correctly

---

## Phase 4: User Story 2 - 自定义允许的来源域名 (Priority: P2)

**Goal**: Operators can restrict CORS to specific origin domains for production security

**Independent Test**: Run `go test -v ./middle/... -count=1 -run "TestCORSSpecificOriginAllowed|TestCORSSpecificOriginDisallowed"` and validate both pass

### Tests for User Story 2

> **NOTE**: These tests already exist in `middle/cors_test.go`. Verify they pass.

- [X] T012 [P] [US2] Verify `TestCORSSpecificOriginAllowed` test in `middle/cors_test.go` — matched origin receives correct CORS header
- [X] T013 [P] [US2] Verify `TestCORSSpecificOriginDisallowed` test in `middle/cors_test.go` — unmatched origin receives no CORS headers

### Implementation for User Story 2

- [X] T014 [US2] Confirm `corsConfig.AllowOrigins` in `main.go` is modifiable — user can change from `["*"]` to `["https://specific-domain.com"]` without touching middleware logic
- [X] T015 [US2] Run quickstart Scenario 6 from `specs/002-add-cors-support/quickstart.md` (modify AllowOrigins, restart, verify disallowed origin gets no CORS header)

**Checkpoint**: User Story 2 verified — specific origin restriction works correctly

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Documentation and final integration validation

- [X] T016 [P] Update `README.md` API 列表: add CORS middleware entry noting default config (allow all origins, no credentials)
- [X] T017 [P] Update `api.http` or `api_test.sh`: add curl examples for CORS verification scenarios from quickstart.md
- [X] T018 Run `go build ./...` to confirm project compiles cleanly
- [X] T019 Run `go test -v ./... -count=1` to confirm all tests pass (including non-CORS tests)
- [X] T020 Run `go fmt ./...` to ensure code formatting compliance

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion — BLOCKS all user stories
- **US1 (Phase 3)**: Depends on Foundational phase completion
- **US2 (Phase 4)**: Depends on Foundational phase completion; independent of US1
- **Polish (Phase 5)**: Depends on all user stories being verified

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational — No dependencies on US2
- **User Story 2 (P2)**: Can start after Foundational — Independent of US1

### Within Each User Story

- Test verification tasks [P] can run in parallel
- Implementation verification follows test verification

### Parallel Opportunities

- All Setup tasks (T001, T002) can run in parallel
- All US1 test verification tasks (T005-T008) can run in parallel
- All US2 test verification tasks (T012, T013) can run in parallel
- Polish documentation tasks (T016, T017) can run in parallel
- After Foundational phase, US1 and US2 can run in parallel (independently testable)

---

## Parallel Example: User Story 1

```bash
# All US1 test verifications can run together:
Task: "Verify TestCORSPreflight in middle/cors_test.go"
Task: "Verify TestCORSGETWithOrigin in middle/cors_test.go"
Task: "Verify TestCORSSameOrigin in middle/cors_test.go"
Task: "Verify TestCORSOptionsWithoutRequestMethod in middle/cors_test.go"
```

## Parallel Example: User Story 2

```bash
# All US2 test verifications can run together:
Task: "Verify TestCORSSpecificOriginAllowed in middle/cors_test.go"
Task: "Verify TestCORSSpecificOriginDisallowed in middle/cors_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T002)
2. Complete Phase 2: Foundational (T003-T004)
3. Complete Phase 3: User Story 1 (T005-T011)
4. **STOP and VALIDATE**: Confirm all US1 tests pass + quickstart scenarios work
5. CORS is now usable with default `*` origin policy

### Incremental Delivery

1. Setup + Foundational → Foundation verified
2. Add User Story 1 → Test independently → MVP (allow all origins)
3. Add User Story 2 → Test independently → Production-ready (specific origins)
4. Polish → Documentation + full test run

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 verification
   - Developer B: User Story 2 verification
3. Both stories verified independently, then merge

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story is independently testable
- CORS implementation code already exists in `main.go` and `middle/cors_test.go` — tasks verify correctness rather than build from scratch
- Commit after each phase or logical group
- Stop at any checkpoint to validate story independently
- All test verifications should be confirmed with `go test -v` output
