# Tasks: 全局配置文件

**Input**: Design documents from `/specs/003-config-file/`

**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, quickstart.md ✅

**Tests**: Tests are REQUIRED per constitution Principle V — handler and config package tests.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2)

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add TOML dependency and create config package skeleton.

- [x] T001 Add `github.com/BurntSushi/toml` dependency in `go.mod` via `go get`
- [x] T002 [P] Create `config/config.go` with empty package and placeholder types

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core config infrastructure that BOTH user stories depend on.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T003 Define `Config`, `AppConfig`, `ServerConfig`, `DatabaseConfig`, `ChangelogConfig` structs in `config/config.go`
- [x] T004 Implement `defaultConfig()` returning safe built-in defaults
- [x] T005 Implement `InitConfig(configPath string) error` — reads TOML, falls back to defaults on failure
- [x] T006 Declare `var Cfg = defaultConfig()` as package-level singleton
- [x] T007 Create `config.toml` at project root with app/server/database/changelog sections

**Checkpoint**: Config package is self-contained — structs defined, defaults ready, TOML loading works.

---

## Phase 3: User Story 1 - 通过配置文件管理版本和更新日志 (Priority: P1) 🎯 MVP

**Goal**: Version, changelog, and setting endpoints read data from `config.Cfg` instead of hardcoded values.

**Independent Test**: Edit `config.toml` version → restart → `GET /common/version` returns new value.

### Implementation for User Story 1

- [x] T008 [US1] Call `config.InitConfig("config.toml")` in `main.go` before any other init
- [x] T009 [P] [US1] Update `handler/common.go` Version() to read from `config.Cfg.App`
- [x] T010 [P] [US1] Update `handler/common.go` Changelog() to read from `config.Cfg.Changelog`
- [x] T011 [P] [US1] Update `handler/common.go` Setting() to read from `config.Cfg.App/Server/Database`
- [x] T012 [US1] Add REST Client test cases for `/common/version`, `/common/changelog`, `/common/setting` in `api.http`
- [x] T013 [US1] Update `README.md` with config file documentation, directory structure, and changelog

### Tests for User Story 1 ⚠️

- [ ] T014 [P] [US1] Unit test for `Config` struct TOML deserialization in `config/config_test.go`
- [ ] T015 [P] [US1] Unit test for `InitConfig` with valid TOML file in `config/config_test.go`
- [ ] T016 [US1] Handler test for `GET /common/version` returning config-driven values in `handler/common_test.go`

**Checkpoint**: Config-driven endpoints work; `config.toml` edits reflected after restart.

---

## Phase 4: User Story 2 - 配置文件缺失或格式错误时的降级处理 (Priority: P2)

**Goal**: Service starts normally with built-in defaults when `config.toml` is missing or invalid.

**Independent Test**: Delete `config.toml` → start service → endpoints return defaults.

### Implementation for User Story 2

- [x] T017 [US2] Handle `os.IsNotExist` in `InitConfig` — log warning, return nil (no error)
- [x] T018 [US2] Handle `toml.Unmarshal` errors in `InitConfig` — log warning, return nil (no error)
- [x] T019 [US2] Handle empty config file — TOML unmarshal produces zero values → defaults apply

### Tests for User Story 2 ⚠️

- [ ] T020 [P] [US2] Unit test for `InitConfig` with non-existent file in `config/config_test.go`
- [ ] T021 [P] [US2] Unit test for `InitConfig` with invalid TOML in `config/config_test.go`
- [ ] T022 [P] [US2] Unit test for `InitConfig` with empty file in `config/config_test.go`
- [ ] T023 [US2] Unit test for `defaultConfig` returning expected values in `config/config_test.go`

**Checkpoint**: Degraded mode works; service never panics on config errors.

---

## Phase 5: Polish & Cross-Cutting Concerns

**Purpose**: Validation and cleanup.

- [ ] T024 Run `quickstart.md` validation scenarios (all 5 scenarios)
- [ ] T025 Run `go test -v ./... -count=1` — all tests pass including new config tests
- [ ] T026 [P] Run `go build ./...` — zero compilation errors
- [ ] T027 [P] Run `read_lints` on `config/` — zero lint warnings

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — completed.
- **Foundational (Phase 2)**: Depends on Setup — completed.
- **User Story 1 (Phase 3)**: Depends on Foundational — implementation completed, tests pending.
- **User Story 2 (Phase 4)**: Depends on Foundational — implementation completed, tests pending.
- **Polish (Phase 5)**: Depends on all tests passing.

### User Story Dependencies

- **User Story 1 (P1)**: Independent — no dependency on US2.
- **User Story 2 (P2)**: Independent — no dependency on US1.

### Parallel Opportunities

- T014, T015 can run in parallel (different test functions in same file is fine)
- T020, T021, T022 can run in parallel
- T023 is independent of other test tasks
- T026, T027 can run in parallel

---

## Implementation Strategy

### Current Status

Implementation is complete. **Remaining work: write unit tests** (T014–T016, T020–T023).

### Remaining: Test-First Approach

1. Write `config/config_test.go` covering:
   - `defaultConfig()` returns correct defaults
   - `InitConfig()` loads valid TOML
   - `InitConfig()` handles missing file (no error, defaults preserved)
   - `InitConfig()` handles invalid TOML (no error, defaults preserved)
   - `InitConfig()` handles empty file (no error, defaults preserved)
2. Write `handler/common_test.go` covering:
   - `GET /common/version` returns config-driven values
   - `GET /common/changelog` returns config-driven entries
   - `GET /common/setting` returns config-driven settings
3. Run full test suite and polish checks

### MVP Scope

Implementation is complete. Tests are the only deliverable remaining.

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- `config/config_test.go` does NOT need `TestMain` (no database dependency)
- `handler/common_test.go` follows existing handler test patterns from project
- `quickstart.md` provides step-by-step validation for manual testing
