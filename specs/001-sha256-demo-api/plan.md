# Implementation Plan: SHA256 Demo API

**Branch**: `001-sha256-demo-api` | **Date**: 2026-07-13 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-sha256-demo-api/spec.md`

## Summary

Add a single GET endpoint `/demo/sha256?text=<input>` that computes SHA256 hash of the query parameter and returns `{input, hash}` in the unified response format. Uses only Go standard library (`crypto/sha256`), requires no database, and follows the existing layered project architecture.

## Technical Context

**Language/Version**: Go ≥ 1.22

**Primary Dependencies**: Echo v5 (web framework)

**Storage**: N/A (stateless computation, no persistence needed)

**Testing**: go test -v ./... -count=1

**Target Platform**: Linux/macOS server

**Project Type**: web-service (REST API)

**Performance Goals**: <1s response per request (spec SC-001)

**Constraints**: Must follow project layered architecture (router → handler → entity → response), no new external dependencies

**Scale/Scope**: Single demo endpoint, single handler file

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. 分层架构 | ✅ PASS | Follows existing layers: entity/demo/, handler/, router/router.go, response/ |
| II. 统一响应格式 | ✅ PASS | Uses response.Ok / response.NotOk for all returns |
| III. 可复制为模板 | ✅ PASS | Uses stdlib only (`crypto/sha256`), zero new dependencies |
| IV. 英文代码产物 | ✅ PASS | Code comments and commit messages will be English |
| V. 测试覆盖 | ✅ PASS | Will include handler test with httptest + logOK |

**Result**: All gates pass. No violations to justify.

## Project Structure

### Documentation (this feature)

```text
specs/001-sha256-demo-api/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (repository root)

```text
entity/
└── demo/
    └── demo.go             # Sha256Request struct with `query:"text"` tag

handler/
└── demo_sha256.go          # _Sha256 handler, Compute method
└── demo_sha256_test.go     # Unit tests with logOK

router/
└── router.go               # Add demo route group registration

api.http                     # Add REST Client test case
README.md                    # Add API listing and changelog
```

**Structure Decision**: Follows the project's existing layered architecture. `entity/demo/` is a new module subdirectory for demo-related request structs. No model/ changes needed since this is a stateless computation endpoint.

## Complexity Tracking

> No violations detected. This section is intentionally empty.
