# Implementation Plan: 添加 CORS 跨域支持

**Branch**: `002-add-cors-support` | **Date**: 2026-07-14 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/002-add-cors-support/spec.md`

## Summary

Add CORS (Cross-Origin Resource Sharing) middleware to enable frontend applications from different origins to call the API. Use Echo v5 built-in CORS middleware with configurable parameters (allowed origins, methods, headers). Register it in the global middleware chain at the appropriate position.

**Current Status**: Implementation complete. CORS middleware is registered in `main.go` using Echo v5's `middleware.CORSWithConfig()` with a package-level config variable `corsConfig`. Integration tests in `middle/cors_test.go` cover all key scenarios.

## Technical Context

**Language/Version**: Go ≥ 1.22

**Primary Dependencies**: Echo v5 (`github.com/labstack/echo/v5`), built-in `middleware.CORS()` — no new third-party dependencies required

**Storage**: N/A (middleware feature, no data persistence)

**Testing**: Standard Go testing (`testing` package) + `httptest.NewRequest` + `echo.New()` with `e.ServeHTTP()` for integration tests

**Target Platform**: Linux/macOS server (same as existing project)

**Project Type**: web-service middleware

**Performance Goals**: OPTIONS preflight response within 50ms overhead (per SC-002)

**Constraints**: Uses Echo v5 built-in CORS middleware (no new dependencies). CORS config is defined as a Go variable in `main.go`, configurable via constants (per Constitution III). The project now has a `config/` package (TOML-based, from 003-config-file), but CORS config remains as Go constants for simplicity and compile-time safety.

**Scale/Scope**: Global middleware applied to all routes; development default allows all origins (`*`)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. 分层架构 | ✅ PASS | CORS as global middleware in `main.go` middleware chain, no cross-layer violation |
| II. 统一响应格式 | ✅ PASS | CORS middleware adds response headers only, does not modify response body |
| III. 可复制为模板 | ✅ PASS | Uses Echo v5 built-in middleware (`middleware.CORS()`), zero new dependencies; configurable via Go constants |
| IV. 英文代码产物 | ✅ PASS | Comments and commit message in English |
| V. 测试覆盖 | ✅ PASS | Integration tests in `middle/cors_test.go` cover preflight, GET with Origin, same-origin, specific origins allowed/disallowed |
| VI. 错误及时抛出 | ✅ PASS | CORS is configured via Go constants at compile time; no runtime config loading that could trigger silent degradation |

**Gate Result: ALL PASS — no violations to justify.**

## Project Structure

### Documentation (this feature)

```text
specs/002-add-cors-support/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (minimal — no data entities)
├── quickstart.md        # Phase 1 output
├── tasks.md             # Phase 2 output (/speckit.tasks)
└── contracts/           # Not applicable (no new API endpoints)
```

### Source Code (repository root)

```text
# Files modified/created by this feature:
main.go                  # CORS config variable + middleware registration
middle/cors_test.go      # CORS integration tests (6 test functions)
```

**Structure Decision**: No new directories required. CORS middleware is registered directly in `main.go` using Echo v5's built-in `middleware.CORSWithConfig()` with a package-level `corsConfig` variable. Integration tests live in `middle/` alongside existing middleware tests. No changes to the `config/` package or `config.toml` — CORS config remains as Go constants per the feature's simplicity requirements.

## Complexity Tracking

> No violations — all constitution checks passed.
