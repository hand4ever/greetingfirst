# Implementation Plan: 006 Audit User Model

**Branch**: `006-audit-user-model` | **Date**: 2026-07-15 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/006-audit-user-model/spec.md`

## Summary

将 `model.User` 结构体及 CRUD 函数与数据库实际 `users` 表建表语句完全对齐。核心变更：

1. `Name` → `Realname` + `Username`（已完成）
2. `ID`: `uint` → `int`
3. `DeletedAt`: `gorm.DeletedAt` → `*time.Time`，手动软删除
4. `Phone` 长度修正、`PasswordHash` 新增（已完成）
5. 手动软删除模式扩展到所有查询函数

## Technical Context

**Language/Version**: Go 1.22+

**Primary Dependencies**: Echo v5, GORM v2 (MySQL driver)

**Storage**: MySQL (existing `users` table with canonical schema)

**Testing**: `go test -v ./... -count=1`; handlers use `httptest.NewRequest` + `echo.New().NewContext`; models test against MySQL

**Target Platform**: Linux/macOS server

**Project Type**: Web service (REST API)

**Performance Goals**: N/A (model refactor, no new endpoints)

**Constraints**: 无新增依赖；所有变更向后不兼容（API 字段名变更）

**Scale/Scope**: ~7 files modified, 0 new files needed

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. 分层架构 | ✅ PASS | 变更严格限定在各层职责内（model/entity/handler） |
| II. 统一响应格式 | ✅ PASS | 使用 `response.Ok` / `response.NotOk`，无变更 |
| III. 可复制为模板 | ✅ PASS | 无新增依赖，DSN 注入方式不变 |
| IV. 英文代码产物 | ✅ PASS | 注释/commit message 使用英文 |
| V. 测试覆盖 | ✅ PASS | 现有测试已更新字段引用，Delete/List 测试覆盖手动软删除 |
| VI. 错误及时抛出 | ✅ PASS | 所有 error 显式检查并返回 |
| VII. 用户管理 Schema | ✅ PASS | `migrations/001_user.mysql.sql` 已更新为 canonical DDL |

**Gate Result**: ✅ ALL PASS — 无违规，无需 Complexity Tracking。

## Project Structure

### Documentation (this feature)

```text
specs/006-audit-user-model/
├── spec.md              # Feature specification
├── plan.md              # This file
├── research.md          # Phase 0: technical decisions
├── data-model.md        # Phase 1: entity & CRUD definitions
├── contracts/
│   └── api.md           # Phase 1: API endpoint contracts
└── quickstart.md        # Phase 1: validation guide
```

### Source Code (repository root)

```text
model/
├── user.go              # User struct: ID→int, DeletedAt→*time.Time, manual soft delete
└── user_test.go         # Updated field refs + soft delete tests

handler/
├── user.go              # extractUserID: uint→int; already uses Realname/Username
├── user_test.go         # Updated field refs
├── demo.go              # Test user creation: Name→Realname
└── demo_test.go         # Updated field refs

entity/user/
└── user.go              # CreateReq/UpdateReq: Name→Realname+Username

migrations/
└── 001_user.mysql.sql   # Canonical DDL matching actual schema
```

## Remaining Work (from clarifications)

Based on spec clarifications that are not yet implemented:

| # | Task | Files | Status |
|---|------|-------|--------|
| 1 | `ID`: `uint` → `int` | `model/user.go`, `handler/user.go` | **Pending** |
| 2 | `DeletedAt`: `gorm.DeletedAt` → `*time.Time` | `model/user.go` | **Pending** |
| 3 | `DeleteUser`: use `UPDATE deleted_at = NOW()` | `model/user.go` | **Pending** |
| 4 | `GetUserByID`: add `WHERE deleted_at IS NULL` | `model/user.go` | **Pending** |
| 5 | `GetUserByPhone`: add `WHERE deleted_at IS NULL` | `model/user.go` | **Pending** |
| 6 | `extractUserID`: `uint` → `int` | `handler/user.go` | **Pending** |
| 7 | Update test files for ID type change | `model/user_test.go`, `handler/user_test.go` | **Pending** |
| 8 | Name→Realname+Username (fields + API) | Multiple | ✅ Done |
| 9 | Phone length, PasswordHash | `model/user.go` | ✅ Done |
| 10 | Migration SQL | `migrations/001_user.mysql.sql` | ✅ Done |
| 11 | README + api.http | `README.md`, `api.http` | ✅ Done |

## Complexity Tracking

> No constitution violations. Section intentionally empty.
