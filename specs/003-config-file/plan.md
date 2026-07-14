# Implementation Plan: 全局配置文件

**Branch**: `003-config-file` | **Date**: 2026-07-14 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/003-config-file/spec.md`

## Summary

Introduce a global TOML configuration file (`config.toml`) loaded at startup. Version, changelog, and setting data are read from config instead of hardcoded. Falls back to built-in defaults when the file is missing or invalid.

Technical approach: Use `github.com/BurntSushi/toml` to parse TOML into a `Config` struct, store as package-level singleton `config.Cfg`, and expose via `InitConfig()` called from `main.go` before any other initialization.

## Technical Context

**Language/Version**: Go 1.26.3

**Primary Dependencies**: `github.com/BurntSushi/toml` v1.6.0 (sole new dependency)

**Storage**: SQLite (unchanged, DSN configurable via `config.toml` but not yet wired from config)

**Testing**: `go test -v ./... -count=1` (standard Go testing, httptest + echo)

**Target Platform**: Linux/macOS server (Echo v5 web service)

**Project Type**: Web service (Echo v5 REST API)

**Performance Goals**: Config parse once at startup (<10ms), no runtime overhead

**Constraints**: Zero panic on config errors; degraded mode with defaults; single new dependency

**Scale/Scope**: ~5 config sections, ~10 fields, stateless in-memory config

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. 分层架构 | ✅ PASS | `config/` is a new cross-cutting infra package (like `model/`); does not break existing layers |
| II. 统一响应格式 | ✅ PASS | Handlers use `response.Ok(c, data)`, no direct `c.JSON()` |
| III. 可复制为模板 | ✅ PASS | Config externalized from code; single lightweight TOML dependency; safe defaults for copy-paste |
| IV. 英文代码产物 | ✅ PASS | All code comments in English; commit messages follow convention |
| V. 测试覆盖 | ⚠️ PENDING | `config/` package needs unit tests (part of implementation phase) |

**Gate**: PASS — no violations. Test coverage to be addressed in implementation.

## Project Structure

### Documentation (this feature)

```text
specs/003-config-file/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
└── tasks.md             # Phase 2 output (not created by /speckit.plan)
```

### Source Code (repository root)

```text
config/
└── config.go            # Config struct, Cfg singleton, InitConfig(), defaultConfig()
config.toml              # Default config file (TOML format, at project root)
handler/
└── common.go            # Updated: reads from config.Cfg instead of hardcoded values
main.go                  # Updated: calls config.InitConfig("config.toml") on startup
entity/common/
└── common.go            # VersionResponse, ChangelogEntry, SettingItem (unchanged)
```

## Complexity Tracking

> No violations to justify.
