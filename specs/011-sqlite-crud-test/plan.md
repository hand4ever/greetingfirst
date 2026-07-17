# 实施计划(Implementation Plan): 新增 SQLite 实例与 CRUD 测试接口

**分支(Branch)**: `011-sqlite-crud-test` | **日期(Date)**: 2026-07-17 | **规格(Spec)**: [spec.md](./spec.md)

**输入(Input)**: Feature specification from `/specs/011-sqlite-crud-test/spec.md`

**说明(Note)**: This template is filled in by the `/speckit.plan` command; its definition describes the execution workflow.

## 概述(Summary)

本特性在现有 MySQL 实例（`model.DB`）之外新增一个**独立共存**的 SQLite 实例（`model.SQLiteDB`），并暴露一组专用于 `test_user` 表的 CRUD 测试接口（`/sqlite/testuser`）。`test_user` 表由用户预先创建（应用不创建、不迁移，符合宪法原则 VII）。技术决策（驱动选型、`:memory:` 测试、软删除、phone 唯一性等）已在 [research.md](./research.md) 中确定，无未决项。现有 MySQL 接口与行为保持不变（FR-012）。

## 技术上下文(Technical Context)

**语言/版本(Language/Version)**: Go 1.26.3（go.mod），遵循 ≥1.22 约束

**主要依赖(Primary Dependencies)**:
- Echo v5（`github.com/labstack/echo/v5`）
- GORM v2（`gorm.io/gorm`）
- **新增**：`github.com/glebarez/sqlite`（纯 Go SQLite 驱动，无 CGO；研究任务 1）
- 现有：`gorm.io/driver/mysql`、`github.com/BurntSushi/toml`

**存储(Storage)**: SQLite 文件型数据库（默认 `greeting.db`），与 MySQL 并存；表 `test_user` 用户预建

**测试框架(Testing)**: Go 原生 `testing` + `net/http/httptest` + GORM `:memory:`（测试用 `SetMaxOpenConns(1)`，研究任务 2）

**目标平台(Target Platform)**: macOS / Linux 服务端（纯 Go 驱动，无平台工具链依赖）

**项目类型(Project Type)**: web-service（Echo v5 HTTP 服务）

**性能目标(Performance Goals)**: SC-001 启动 3 秒内建立 SQLite 连接；SC-004 并发创建 10 条不同 phone 全部成功

**约束(Constraints)**:
- 应用 MUST NOT 创建/迁移 `test_user`（原则 VII / FR-010）
- 任一数据库连接失败 MUST fail-fast 退出（原则 VI / FR-003）
- 统一响应 `response.Ok/NotOk`（原则 II）
- 代码注释英文、commit message 英文（原则 IV）

**规模/范围(Scale/Scope)**: 单库实例 + 5 个测试接口 + 1 个测试实体；不影响现有 MySQL 模块

## 宪法检查(Constitution Check)

*门禁(GATE): Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则 | 检查项 | 结论 |
|------|--------|------|
| I. 分层架构 | router/handler/entity/model/response 分层，新模块在 `router/router.go` 注册 | ✅ PASS |
| II. 统一响应 | 全部接口使用 `response.Ok/NotOk/NotOkWithCode` | ✅ PASS |
| III. 可复制模板 | `InitSQLite` 仅连接不建表；纯 Go 驱动无 CGO 复制负担；DSN 来自配置 | ✅ PASS |
| IV. 英文产物 | 代码注释/commit 英文（实现阶段遵守） | ✅ PASS |
| V. 测试覆盖 | 每个 handler/model 方法单测；TestMain 用 SQL 脚本（非 AutoMigrate）建 `:memory:` 测试表 | ✅ PASS |
| VI. 错误及时抛出 | DB 连接/Ping 失败返回明确错误，`main.go` panic 退出 | ✅ PASS |
| VII. 用户自管 Schema | 应用零 AutoMigrate/建表；`migrations/002_test_user.sql` 用户自管；测试用同脚本 | ✅ PASS |

**Phase 1 复查结论**: 设计（[data-model.md](./data-model.md)、[contracts/api.md](./contracts/api.md)）未引入任何宪法违规。应用不建表、测试不 AutoMigrate，与原则 V/VII 完全一致，无需复杂度追踪例外。

## 项目结构(Project Structure)

### 文档(Documentation)（本特性）

```text
specs/011-sqlite-crud-test/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   └── api.md
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### 源代码(Source Code)（仓库根目录）

```text
config/
└── config.go            # 新增 Database.SQLite SQLiteConfig{DSN}；defaultConfig 增加默认 dsn="greeting.db"

model/
├── db.go                # 新增 SQLiteDB *gorm.DB 与 InitSQLite(dsn) error；复用 InitDB 模式
└── testuser.go          # 新增 TestUser 模型 + CRUD 方法（使用 SQLiteDB）

entity/
└── sqliteusr/
    └── sqliteusr.go     # TestUserCreateReq / TestUserUpdateReq

handler/
├── sqliteuser.go        # _SqliteUser handler：Create/Get/Update/Delete/List
└── sqliteuser_test.go   # handler 单测（httptest）

router/
└── sqlite.go            # sqlite(e) 注册 /sqlite/testuser 路由组

migrations/
└── 002_test_user.sql    # 用户自管 test_user 建表脚本（应用不执行）

main.go                  # 调用 model.InitSQLite，失败 panic

config.toml             # 新增 [database.sqlite] dsn
api.http                # 新增 /sqlite/testuser* REST Client 用例
README.md               # 更新「API 列表」「更新日志」
go.mod / go.sum         # 新增 github.com/glebarez/sqlite
```

**结构决策(Structure Decision)**: 沿用现有分层（router→handler→entity→model→response）。新增 `model.SQLiteDB` 与 `model.DB` 并存；新增 `handler/sqliteuser.go`、`entity/sqliteusr/`、`router/sqlite.go` 三个文件，严格遵循目录规范。迁移脚本置于 `migrations/`（用户自管资产，应用不执行）。

## 复杂度追踪(Complexity Tracking)

> 宪法检查全部 PASS，无违规，无需填写本表。

| 违规(Violation) | 必要性(Why Needed) | 拒绝更简单方案的原因(Simpler Alternative Rejected Because) |
|-----------|------------|-------------------------------------|
| （无） | — | — |

## 实现任务预览(Implementation Tasks Preview)

> 以下为 Phase 2 (`/speckit.tasks`) 将细化的任务概览，便于评审计划完整性：

1. **配置层**：`config.go` 增加 `SQLiteConfig`；`config.toml` 增加 `[database.sqlite]`（默认 `greeting.db`）。
2. **模型层**：`model/db.go` 增加 `SQLiteDB` 与 `InitSQLite`；新增 `model/testuser.go`（TestUser + CRUD 方法，含 phone 活跃唯一校验）。
3. **请求实体**：`entity/sqliteusr/sqliteusr.go` 定义 Create/Update 请求体。
4. **处理器**：`handler/sqliteuser.go` 实现 5 个接口，统一 `response` 封装。
5. **路由**：`router/sqlite.go` 注册 `/sqlite/testuser` 组；`router/router.go` 调用。
6. **启动**：`main.go` 调用 `InitSQLite`，失败 panic。
7. **Schema 资产**：新增 `migrations/002_test_user.sql`。
8. **依赖**：`go get github.com/glebarez/sqlite`。
9. **测试**：`model/testuser_test.go`（TestMain + `:memory:` + SQL 脚本）、`handler/sqliteuser_test.go`。
10. **文档/用例**：`api.http` 用例、`README.md` API 列表与更新日志。
11. **校验**：`go build ./...` 与 `go test -v ./... -count=1` 通过；`gofmt`。
