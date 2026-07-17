# 研究报告(Research): 新增 SQLite 实例与 CRUD 测试接口

**功能分支(Feature Branch)**: `011-sqlite-crud-test`
**日期(Date)**: 2026-07-17
**关联规格(Spec)**: [spec.md](./spec.md)

本文件记录 Phase 0 中针对技术上下文（Technical Context）不确定项的研究结论。规格说明经澄清后已无 `NEEDS CLARIFICATION` 标记，以下研究针对"如何最优实现"做出技术决策。

---

## 研究任务 1：SQLite GORM 驱动选型

**任务(Task)**: 为项目选择兼容 GORM v2 的 SQLite 驱动，约束条件见宪法原则 III（可复制为模板，避免增加复制负担、优先轻量依赖）。

### Decision
采用 **`github.com/glebarez/sqlite`**（纯 Go 实现，零 CGO 依赖）作为 GORM 的 SQLite 驱动。

### Rationale
- 项目作为"可复制到新项目的模板"（原则 III），依赖应优先避免平台/工具链耦合。**标准驱动 `gorm.io/driver/sqlite` 底层依赖 `github.com/mattn/go-sqlite3`，需要 CGO 与本地 C 编译器（gcc）**，在交叉编译、Alpine 镜像、macOS 无 Xcode Command Line Tools 等场景下会中断 `go build`，增加复制负担。
- `glebarez/sqlite` 基于 `modernc.org/sqlite` 纯 Go 实现，无 CGO，开箱即用，更契合模板化复制诉求。
- 二者 API 完全一致（`gorm.Open(...)`），不引入额外心智成本。

### Alternatives considered
- `gorm.io/driver/sqlite`（CGO / mattn）：功能等价但引入 CGO 工具链依赖 → 否决（违背原则 III 轻量/可复制）。
- 直接使用 `database/sql` + 原生 `sqlite3` 驱动：放弃 GORM 统一 ORM 能力，与现有 `model.DB` 风格不一致 → 否决。

---

## 研究任务 2：`:memory:` 测试数据库的连接池陷阱

**任务(Task)**: 解决 SQLite `:memory:` 在 GORM 连接池下"第二次查询数据丢失"的问题，以满足单元测试需求（原则 V 允许测试用 `:memory:`）。

### Decision
测试环境下打开 `:memory:` SQLite 后，**调用 `sqlDB.SetMaxOpenConns(1)`** 将连接池限制为单连接；或等价地采用共享缓存 DSN `file::memory:?cache=shared&_fk=1`。本特性统一采用 `SetMaxOpenConns(1)` 方案。

### Rationale
- SQLite 的 `:memory:` 数据库作用域绑定到单个数据库连接。GORM 默认使用 `database/sql` 连接池，当第二个查询从池中取用新连接时，会得到一张**全新的空内存库**，导致先前写入不可见、测试随机失败。
- `SetMaxOpenConns(1)` 强制所有操作复用同一连接，内存库生命周期与测试一致，简单可靠。
- 此限制**仅用于 `:memory:` 测试实例**，生产文件型 `greeting.db` 仍使用默认连接池（并发读写由 SQLite/WAL 处理），互不影响。

### Alternatives considered
- 默认连接池（不处理）：写入后读取随机失败 → 否决。
- 每次操作重建 `:memory:`：无法跨操作保留数据 → 否决。

---

## 研究任务 3：test_user 表的 Schema 管理与宪法对齐

**任务(Task)**: 在宪法原则 V（测试用 SQL schema 脚本）与原则 VII（用户自管 schema、应用禁止自动建表）约束下，确定 `test_user` 表的创建与测试策略。用户已澄清：`test_user` 表已预先存在于 `greeting.db`，应用不得创建或迁移。

### Decision
1. 应用代码**绝不**在 `InitSQLite` 或任何初始化路径中执行 `AutoMigrate` / `CREATE TABLE`（满足 FR-010 与原则 VII）。
2. 新增一份用户自管的迁移脚本 **`migrations/002_test_user.sql`**，描述 `test_user` 表结构，作为项目可复制资产（满足原则 VII"schema 以独立 `.sql` 文件管理"）。用户可手动对其 `greeting.db` 执行以创建/重建该表。
3. 单元测试的 `TestMain` **对 `:memory:` 实例执行同一份 `002_test_user.sql`** 来初始化测试表（满足原则 V"TestMain 运行 SQL schema 脚本、禁止 AutoMigrate"）。

### Rationale
- 完全契合宪法：**应用零自动建表**（原则 VII），**测试零 AutoMigrate**（原则 V），且 schema 以 `.sql` 文件沉淀为资产。
- 用户已手动建好的 `greeting.db` 表结构与迁移脚本保持一致，避免漂移。
- 测试无需依赖外部文件型 `greeting.db`，`go test` 在无 MySQL/无预置 SQLite 文件时也能独立运行。

### Alternatives considered
- 应用启动时 `AutoMigrate(test_user)`：直接违反原则 VII 与 FR-010 → 否决。
- 测试中用 `AutoMigrate` 建表：违反原则 V → 否决。
- 测试直接复用用户 `greeting.db` 文件：测试与用户数据耦合，且 `go test` 强依赖该文件存在 → 否决。

---

## 研究任务 4：软删除（Soft Delete）实现方式

**任务(Task)**: 确定 `test_user` 软删除的实现方式（FR-008：删除接口使用软删除，不物理删除）。

### Decision
`TestUser` 模型使用 GORM 原生软删除类型 **`gorm.DeletedAt`**（字段 `DeletedAt gorm.DeletedAt`）。删除调用 `SQLiteDB.Delete(&TestUser{}, id)` 由 GORM 自动置 `deleted_at`；所有 `First` / `Find` 查询自动排除已软删除记录。

### Rationale
- GORM 原生 `gorm.DeletedAt` 内置软删除语义：`Delete` 自动转为 `UPDATE ... SET deleted_at`，`First/Find` 自动追加 `deleted_at IS NULL` 过滤，减少手写 SQL 与出错面。
- 规格明确"测试接口可自由设计、无需对齐 MySQL"，因此无需沿用 MySQL `User` 的 `*time.Time` + 手动 `deleted_at IS NULL` 写法，采用更简洁的原生方式即可。
- 配合研究任务 3 的 `002_test_user.sql`，表需含可空的 `deleted_at` 列（DATETIME NULL）。

### Alternatives considered
- 沿用 MySQL 的 `*time.Time` + 手动 `deleted_at IS NULL`：可行但更冗长，且非必要（无需与 MySQL 对齐）→ 不采用。
- 物理删除：违反 FR-008 → 否决。

---

## 研究任务 5：phone 唯一性（未删除记录内唯一）实现方式

**任务(Task)**: 满足 FR-006 / Edge Case"创建用户时 phone 重复（在未删除记录中）应返回明确错误；软删除后同 phone 可新建"——但 SQLite 普通 `UNIQUE` 约束会在软删除后阻止复用。

### Decision
**应用层校验**：创建前用 `SQLiteDB.Where("phone = ? AND deleted_at IS NULL", phone).First(...)` 查询是否已存在活跃同号记录；若存在则返回 `response.NotOk(c, "phone already exists")`。**不依赖数据库 UNIQUE 约束**（避免软删除后复用被阻）。

### Rationale
- 软删除 + 普通 UNIQUE 约束互斥（软删后同号仍被约束占用，无法复用）。应用层"活跃记录唯一"校验天然支持"软删后可复用"。
- 校验逻辑简单、可测试，且返回信息明确（符合原则 VI 错误明确化）。

### Alternatives considered
- 数据库 `UNIQUE(phone)` + 软删除：复用被阻 → 否决。
- SQLite 部分唯一索引 `CREATE UNIQUE INDEX ... WHERE deleted_at IS NULL`：可行但增加 schema 复杂度，且与 `002_test_user.sql` 简单化目标不符 → 不采用。

---

## 研究任务 6：全局实例与初始化策略（双库共存、fail-fast）

**任务(Task)**: 在现有 `model.DB`（MySQL）之外新增独立 SQLite 实例，并保证任一连接失败即整体 fail-fast（FR-002、FR-003、原则 VI）。

### Decision
- 新增全局变量 `model.SQLiteDB *gorm.DB`，与 `model.DB` 并存。
- 新增 `model.InitSQLite(dsn string) error`：打开连接 → `Ping()` 验证 → 赋值 `SQLiteDB`；失败返回带类型与地址的 `fmt.Errorf`。
- `main.go` 在 `InitDB`（MySQL）之后调用 `InitSQLite`；任一失败 `panic` 退出（与现有 MySQL 失败策略一致，满足原则 VI）。

### Rationale
- 完全复用现有 `InitDB` 的模式（打开 + Ping + 全局赋值），降低新增复杂度，符合模板一致性。
- 双库各自独立生命周期，互不耦合。

### Alternatives considered
- 单一 `DB` 多数据源切换：违背"各自独立连接池与生命周期"的隔离诉求 → 否决。

---

## 研究结论汇总(Decisions Summary)

| # | 决策项 | 结论 |
|---|--------|------|
| 1 | SQLite 驱动 | `github.com/glebarez/sqlite`（纯 Go，无 CGO） |
| 2 | `:memory:` 测试 | `SetMaxOpenConns(1)` 单连接 |
| 3 | test_user schema | 应用不建表；`migrations/002_test_user.sql` 用户自管；TestMain 对 `:memory:` 执行该 SQL |
| 4 | 软删除 | GORM 原生 `gorm.DeletedAt` |
| 5 | phone 唯一 | 应用层"活跃记录唯一"校验，无 DB UNIQUE 约束 |
| 6 | 双库实例 | `model.SQLiteDB` + `model.InitSQLite`，任一失败 fail-fast |

所有决策均可在不违反宪法原则的前提下落地，无残留未决项。
