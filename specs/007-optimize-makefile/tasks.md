# Tasks: 优化 Makefile

**Input**: Design documents from `specs/007-optimize-makefile/`

**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/, quickstart.md

**Tests**: Not requested in feature specification. Tasks focus on implementation only.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

**Note**: All tasks modify the single file `Makefile` at repository root. Tasks within each phase should be executed sequentially to avoid merge conflicts. User story phases are additive — each phase adds its targets on top of the previous phase.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Review & Preparation)

**Purpose**: Understand the current Makefile state and prepare for refactoring

- [X] T001 Review current `Makefile` content — identify hardcoded values in `buildqa` target (IP `111.333.222.444`, path `/opt/src/main`, supervisor `xxxx`), broken `-ldflags` trailing backslash, and missing `.PHONY` declarations
- [X] T002 Review `specs/007-optimize-makefile/spec.md` user stories and `specs/007-optimize-makefile/contracts/makefile-targets.md` for target behavior contracts

---

## Phase 2: Foundational (Variable Definitions & Defaults)

**Purpose**: Define shared variables and default goal that ALL user story targets will reference

**⚠️ CRITICAL**: No user story target work should begin until these variables are defined

- [X] T003 Define build variables at top of `Makefile`: `GOOS ?= $(shell go env GOOS)`, `GOARCH ?= $(shell go env GOARCH)`, `BIN_DIR ?= bin`, `BIN_NAME ?= $(shell basename $(CURDIR))` (per research.md R2)
- [X] T004 Set `.DEFAULT_GOAL := help` at top of `Makefile` so bare `make` shows help (per FR-007, research.md R6)
- [X] T005 Define deploy variables at top of `Makefile`: `DEPLOY_HOST ?=`, `DEPLOY_PATH ?= /opt/src/main`, `DEPLOY_SUPERVISOR ?=` (per FR-008, research.md R2)

**Checkpoint**: Variable scaffolding ready — targets can now reference `$(GOOS)`, `$(BIN_DIR)`, etc. Verify with `make -p | grep -E '^(GOOS|BIN_DIR|DEPLOY_HOST)'`

---

## Phase 3: User Story 1 - 开发者本地快速启动服务 (Priority: P1) 🎯 MVP

**Goal**: `make rundev` starts the local dev server with clear status feedback

**Independent Test**: Execute `make rundev`, verify server starts and displays listening address. Ctrl+C to stop. See `specs/007-optimize-makefile/quickstart.md` scenario 8.

### Implementation for User Story 1

- [X] T006 [US1] Rewrite `rundev` target in `Makefile`: add `.PHONY: rundev`, print `Starting local dev server...` before `go run main.go`, use `@go run main.go` to suppress command echo (per FR-001, contracts/`rundev` contract)
- [X] T007 [US1] Validate `rundev` target: execute `make rundev`, confirm server starts and shows listening address; force-quit then re-run to verify port-in-use error is visible

**Checkpoint**: `make rundev` works with status messages. User Story 1 complete and independently testable.

---

## Phase 4: User Story 2 - 本地编译、格式化与测试 (Priority: P1)

**Goal**: `make build`, `make test`, `make fmt`, `make lint`, `make clean`, `make help` all functional

**Independent Test**: Execute each target and verify expected output per `specs/007-optimize-makefile/contracts/makefile-targets.md` and `specs/007-optimize-makefile/quickstart.md` scenarios 1-7.

### Implementation for User Story 2

- [X] T008 [US2] Add `fmt` target in `Makefile`: detect gofumpt via `command -v gofumpt`, use `gofumpt -l -w .` if found, fallback to `gofmt -l -w .` with install hint message, declare `.PHONY: fmt` (per FR-004, research.md R4)
- [X] T009 [US2] Add `build` target in `Makefile`: create `$(BIN_DIR)` if not exists via `mkdir -p`, print build info, execute `go build -o $(BIN_DIR)/$(BIN_NAME) main.go`, declare `.PHONY: build` (per FR-002, FR-010, contracts/`build` contract)
- [X] T010 [US2] Add `test` target in `Makefile`: print `Running tests...`, execute `go test -v ./... -count=1`, print `All tests passed.` on success, declare `.PHONY: test` (per FR-003, Constitution Principle V)
- [X] T011 [US2] Add `lint` target in `Makefile`: execute `go vet ./...`, declare `.PHONY: lint` (per FR-005)
- [X] T012 [US2] Add `clean` target in `Makefile`: print `Cleaning build artifacts...`, remove `$(BIN_DIR)` via `rm -rf`, print `Done.`, declare `.PHONY: clean` (per FR-006)
- [X] T013 [US2] Add `help` target in `Makefile`: print `Usage: make [target]` header, use `@grep -E` pattern to extract targets with `## category: description` comment format, group by category (per FR-007, research.md R6)
- [X] T014 [US2] Validate US2 targets: run `make fmt`, `make build`, `make test`, `make lint`, `make clean`, `make help` and verify each produces expected output per contracts

**Checkpoint**: All development/build/test targets functional. Developer can run full local workflow with `make` commands.

---

## Phase 5: User Story 3 - 部署到测试服务器 (Priority: P2)

**Goal**: `make deploy-qa` performs build → upload → restart with variable-driven configuration and fail-fast guards

**Independent Test**: Set `DEPLOY_HOST`, `DEPLOY_SUPERVISOR` and execute `make deploy-qa`. Verify binary upload and supervisor restart. Test missing-variable guard with bare `make deploy-qa`. See `specs/007-optimize-makefile/quickstart.md` scenario 9.

### Implementation for User Story 3

- [X] T015 [US3] Add `build-linux` target in `Makefile`: cross-compile with `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BIN_NAME) main.go`, create `$(BIN_DIR)` if needed, declare `.PHONY: build-linux` (per FR-011, research.md R3)
- [X] T016 [US3] Add `define check_required` macro in `Makefile`: accepts variable name and description, checks if empty via `[ -z "$($1)" ]`, prints error with usage example and `exit 1` (per FR-008, FR-012, research.md R5)
- [X] T017 [US3] Add `deploy-qa` target in `Makefile`: call `check_required` for `DEPLOY_HOST` and `DEPLOY_SUPERVISOR`, depend on `build-linux`, print progress at each step (upload, restart), use `scp -O $(BIN_DIR)/$(BIN_NAME) root@$(DEPLOY_HOST):$(DEPLOY_PATH)/` and `ssh root@$(DEPLOY_HOST) "supervisorctl restart $(DEPLOY_SUPERVISOR)"`, declare `.PHONY: deploy-qa` (per FR-008, FR-009, contracts/`deploy-qa` contract)
- [X] T018 [US3] Remove old `buildqa` target from `Makefile` since replaced by `deploy-qa` + `build-linux` (breaking change — document in commit message)
- [X] T019 [US3] Validate deploy guards: run `make deploy-qa` without variables, confirm error message and exit 1 with usage example

**Checkpoint**: Deployment target works with variables; fails fast on missing config. User Story 3 complete.

---

## Phase 6: User Story 4 - 新人快速上手 (Priority: P3)

**Goal**: `make help` output is well-organized with categorized self-documenting comments so new developers can understand all operations at a glance

**Independent Test**: Execute `make` (bare) or `make help` and verify output groups targets by category with descriptions. See `specs/007-optimize-makefile/quickstart.md` scenario 1.

### Implementation for User Story 4

- [X] T020 [US4] Add self-documenting comments to all targets in `Makefile`: format as `## category: description` above each target (e.g., `## dev: Start local development server` above `rundev`). Categories: `help:`, `dev:`, `build:`, `test:`, `deploy:` (per research.md R6)
- [X] T021 [US4] Enhance `help` target in `Makefile`: add category headers (Development, Build, Test, Deploy), match both `## dev:` and `## build:` and other category patterns in grep extraction, sort targets within each category (per contracts/`help` contract, FR-007)
- [X] T022 [US4] Validate help output: run `make` and `make help`; confirm all 9 targets displayed, grouped by category, with Chinese descriptions readable by new developers

**Checkpoint**: `make help` self-documents the entire Makefile. New developers can understand all operations without reading source.

---

## Phase 7: Polish & Final Validation

**Purpose**: Cross-cutting improvements and comprehensive validation

- [X] T023 Ensure all `.PHONY` declarations are grouped or listed consistently in `Makefile` — verify no target missing `.PHONY` (per FR-009). Tip: collect all `.PHONY` targets in a single declaration: `.PHONY: help rundev fmt build build-linux test lint clean deploy-qa`
- [X] T024 Run full quickstart validation: execute `specs/007-optimize-makefile/quickstart.md` scenarios 1-10 in order, verify all expected outcomes match (SC-001 through SC-006)
- [X] T025 Final review: confirm `Makefile` is under 100 lines, all comments are in English (Constitution Principle IV), no hardcoded IP/credentials (Constitution Principle III)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup (T001-T002 context) — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational (Phase 2) — needs variables defined
- **User Story 2 (Phase 4)**: Depends on Foundational (Phase 2) — can proceed in parallel with US1
- **User Story 3 (Phase 5)**: Depends on US2 (needs `build` target pattern established, `$(BIN_DIR)` tested). Also depends on Foundational (deploy variables)
- **User Story 4 (Phase 6)**: Depends on ALL prior phases since it adds comments to every target — MUST run last
- **Polish (Phase 7)**: Depends on all user stories complete

### User Story Dependencies

- **US1 (P1)**: Independent — only needs Phase 2 variables. No other story dependency.
- **US2 (P1)**: Independent — only needs Phase 2 variables. No other story dependency.
- **US3 (P2)**: Depends on US2 — needs `$(BIN_DIR)` creation pattern and build commands validated.
- **US4 (P3)**: Depends on ALL stories — self-documenting comments must reference completed targets.

### Within Each Phase

Since all tasks modify the single `Makefile`, execute tasks sequentially within each phase to avoid merge conflicts.

### Parallel Opportunities

> All tasks touch `Makefile`, so true parallelism is limited. Opportunistic parallelism:

- **US1 + US2 can be planned in parallel**: developer can draft the `rundev` rewrite and the build/test/fmt blocks independently, then merge into Makefile sequentially
- **Phase 7 tasks T023-T025**: review tasks that don't modify code can run concurrently

---

## Parallel Example: User Stories 1 & 2

```bash
# These can be drafted in separate branches or sections, then merged:
Task T006 [US1]: "Rewrite rundev target in Makefile"
Task T008-T012 [US2]: "Add fmt, build, test, lint, clean, help targets in Makefile"

# After both drafted, merge and execute T007 + T014 validation
```

---

## Implementation Strategy

### MVP First (User Stories 1 + 2 Only)

1. Complete Phase 1: Setup (T001-T002)
2. Complete Phase 2: Foundational (T003-T005)
3. Complete Phase 3: User Story 1 — `make rundev` works (T006-T007)
4. Complete Phase 4: User Story 2 — `make build/test/fmt/lint/clean/help` works (T008-T014)
5. **STOP and VALIDATE**: All P1 targets functional — developer can already use 90% of the Makefile
6. Deploy/demo: `make rundev`, `make build`, `make test`

### Incremental Delivery

1. Setup + Foundational → Variables and default goal ready
2. Add US1 (rundev) → Test independently → Developers can start the server
3. Add US2 (build/test/fmt/lint/clean/help) → Test independently → Full local workflow (MVP!)
4. Add US3 (deploy-qa) → Test independently → QA deployment unblocked
5. Add US4 (help enhancement) → Test independently → New developer onboarding improved
6. Polish → Final validation → Ready to merge

### Target Implementation Order (Single Developer, Sequential)

```
T001 → T002 → T003 → T004 → T005   (Setup + Foundational ~30 min)
    ↓
T006 → T007                         (US1 ~15 min)
    ↓
T008 → T009 → T010 → T011 → T012 → T013 → T014   (US2 ~45 min)
    ↓
T015 → T016 → T017 → T018 → T019   (US3 ~30 min)
    ↓
T020 → T021 → T022                  (US4 ~20 min)
    ↓
T023 → T024 → T025                  (Polish ~15 min)
```

**Estimated total**: ~2.5 hours

---

## Notes

- All tasks modify `Makefile` at repository root — the only file changed in this feature
- Commit after each phase (or logical task group) with English commit messages per Constitution Principle IV
- Existing `buildqa` target is removed in T018 — document in commit message as breaking change (`buildqa` replaced by `deploy-qa` + `build-linux`)
- `gofumpt` is optional — fallback to `go fmt` ensures zero-install usability per Constitution Principle III
- Validation tasks (T007, T014, T019, T022, T024) reference `specs/007-optimize-makefile/quickstart.md` scenarios
