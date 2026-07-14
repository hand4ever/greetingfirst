# Implementation Plan: 添加 CORS 跨域支持

**Branch**: `002-add-cors-support` | **Date**: 2026-07-14 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/002-add-cors-support/spec.md`

## Summary

Add CORS (Cross-Origin Resource Sharing) middleware to enable frontend applications from different origins to call the API. Use Echo v5 built-in CORS middleware with configurable parameters (allowed origins, methods, headers). Register it in the global middleware chain at the appropriate position.

## Technical Context

**Language/Version**: Go ≥ 1.22

**Primary Dependencies**: Echo v5 (`github.com/labstack/echo/v5`), built-in `middleware.CORS()` — no new third-party dependencies required

**Storage**: N/A (middleware feature, no data persistence)

**Testing**: Standard Go testing (`testing` package) + `httptest.NewRequest` + `echo.New().NewContext`

**Target Platform**: Linux/macOS server (same as existing project)

**Project Type**: web-service middleware

**Performance Goals**: OPTIONS preflight response within 50ms overhead (per SC-002)

**Constraints**: Must use Echo v5 built-in CORS middleware (no new dependencies), must be configurable via constants (per Constitution III)

**Scale/Scope**: Global middleware applied to all routes; development default allows all origins (`*`)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. 分层架构 | ✅ PASS | CORS as global middleware in `main.go` middleware chain, no cross-layer violation |
| II. 统一响应格式 | ✅ PASS | CORS middleware adds response headers only, does not modify response body |
| III. 可复制为模板 | ✅ PASS | Uses Echo v5 built-in middleware (`middleware.CORS()`), zero new dependencies; configurable via Go constants |
| IV. 英文代码产物 | ✅ PASS | Comments and commit message in English |
| V. 测试覆盖 | ✅ PASS | Will add CORS middleware integration test in `middle/` |

**Gate Result: ALL PASS — no violations to justify.**

## Project Structure

### Documentation (this feature)

```text
specs/002-add-cors-support/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (minimal — no data entities)
├── quickstart.md        # Phase 1 output
├── contracts/           # Not applicable (no new API endpoints)
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (repository root)

```text
# This feature touches only:
main.go                  # Register CORS middleware in middleware chain
middle/cors_test.go      # CORS integration tests (new)
```

**Structure Decision**: No new directories required. CORS middleware is registered directly in `main.go` using Echo v5's built-in `middleware.CORS()` with a configurable CORS config variable. Integration tests live in `middle/` alongside existing middleware tests.

## Complexity Tracking

> No violations — all constitution checks passed.
