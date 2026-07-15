# Feature Specification: 移除 SQLite 相关代码

**Feature Branch**: `005-remove-sqlite`

**Created**: 2026-07-15

**Status**: Draft

**Input**: User description: "本次新需求为，去除代码里所有sqlite相关的代码"

## Clarifications

### Session 2026-07-15

- Q: 是否完全移除所有 SQLite 代码？ → A: 是，在宪法范围内移除所有 SQLite 相关代码；架构层面保留多库并存灵活性，但当前仅保留 MySQL 实现。
- Q: InitDB 启动时是否需要检查表存在？ → A: 否，InitDB 仅初始化数据库实例（打开连接 + Ping），不做任何表存在性检查；`EnsureUserTable()` 相关逻辑一并移除。
- Q: 测试环境策略——没有 MySQL 实例时测试如何运行？ → A: 采用 Option C，仅保留编译验证（`go build ./...`），暂不要求模型层/Handler 层测试通过。测试代码中的 SQLite 引用将被移除，具体测试数据库方案后续单独处理。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 移除生产环境 SQLite 依赖 (Priority: P1)

作为开发者，我希望项目的生产代码不再依赖 SQLite，只保留 MySQL 作为唯一数据库，从而简化数据库配置、减少驱动依赖、降低维护复杂度。

**Why this priority**: 这是本次需求的核心目标，直接消除双数据库架构带来的代码重复和维护负担。

**Independent Test**: 启动项目只需配置 MySQL DSN，无需 SQLite 配置即可正常运行；`config.toml` 中不再包含 `[database.sqlite]` 段。

**Acceptance Scenarios**:

1. **Given** 项目启动时只配置了 MySQL DSN，**When** 执行 `go build ./...`，**Then** 编译成功，不依赖 `gorm.io/driver/sqlite` 包。
2. **Given** 配置文件中已移除 SQLite 相关配置，**When** 启动服务，**Then** 服务正常启动，仅连接 MySQL 数据库，不做任何表存在性检查。
3. **Given** 已移除所有 SQLite 模型和 CRUD 函数，**When** 调用现有的 MySQL CRUD 接口（`/demo/usr`），**Then** 功能正常，数据读写均通过 MySQL 完成。
4. **Given** 原有的 SQLite 相关 handler 端点（`GET /demo/user/phone`）已迁移至 MySQL，**When** 调用该端点，**Then** 数据在 MySQL 的 `users` 表中完成读写。

---

### User Story 2 - 清理 SQLite 遗留文件 (Priority: P2)

作为开发者，我希望所有 SQLite 相关的迁移文件、Schema 文件以及文档中的 SQLite 引用被清除，确保项目文件干净、无歧义。

**Why this priority**: 残留文件会导致混淆，新人看到 SQLite 迁移脚本和 MySQL 迁移脚本并存时会困惑应该用哪个。

**Independent Test**: 在项目目录中搜索 "sqlite"、"SQLite"、"sl_users" 关键字，仅 `.specify/` 和 `specs/` 历史文档中存在匹配，其他源码和配置文件中无任何匹配。

**Acceptance Scenarios**:

1. **Given** `migrations/` 目录中存在 `001_user.sqlite.sql`，**When** 执行清理，**Then** 该文件被删除。
2. **Given** `model/schema.sql` 内嵌了 SQLite 建表 DDL，**When** 执行清理，**Then** 该文件或其中 SQLite 相关内容被移除。
3. **Given** `README.md` 中技术栈描述包含 "SQLite"，**When** 执行清理，**Then** 更新为仅描述 MySQL。
4. **Given** `api.http` 和 `README.md` 中存在 SQLite 相关接口文档，**When** 执行清理，**Then** 更新为描述 MySQL 对应接口。

---

### User Story 3 - 清理测试中的 SQLite 依赖 (Priority: P3)

作为开发者，我希望测试文件不再导入 `gorm.io/driver/sqlite`，但考虑到当前 MySQL 测试基础设施未就绪，暂不要求测试通过，仅确保编译不受阻。

**Why this priority**: 测试需要 MySQL 实例，当前阶段仅清理编译依赖即可，完整测试方案后续单独处理。优先级放低，不影响核心功能。

**Independent Test**: 测试文件中不存在 `gorm.io/driver/sqlite` 导入，`go build ./...` 编译通过。

**Acceptance Scenarios**:

1. **Given** `model/user_test.go` 中存在 SQLite 相关导入和 `SQLiteUser` 测试，**When** 清理后，**Then** SQLite 导入被移除，`SQLiteUser` 相关测试用例被移除或注释。
2. **Given** `model/db_test.go` 中使用 `gorm.io/driver/sqlite`，**When** 清理后，**Then** SQLite 导入被移除，依赖 `ApplySchema`/`EnsureUserTable` 的测试被移除。
3. **Given** handler 测试文件中使用 `gorm.io/driver/sqlite`，**When** 清理后，**Then** SQLite 导入被移除。

---

### User Story 4 - 更新项目宪法 (Priority: P3)

作为项目维护者，我希望项目宪法（constitution.md）中的技术栈描述移除 SQLite 引用，体现当前 MySQL 实现与多库并存架构灵活性。

**Why this priority**: 宪法是项目的架构准则，需与实际情况保持一致，但不影响功能实现，优先级较低。

**Independent Test**: 查看 `.specify/memory/constitution.md` 的技术栈约束表格，数据库一项为 "MySQL / 可替换"，且不再提及 SQLite。

**Acceptance Scenarios**:

1. **Given** 宪法技术栈表格中 `数据库` 行描述为 "MySQL / 可替换"，**When** 检查约束列，**Then** 注明「架构支持多库并存」。
2. **Given** 宪法原则中可能存在 SQLite 相关的描述，**When** 更新后，**Then** 相关内容被泛化，不再硬编码 SQLite。

---

### Edge Cases

- 若 `config.toml` 中仍保留 `[database.sqlite]` 段但程序已不支持，启动时应如何处理？→ **假定**：启动时若检测到无效配置项，打印警告日志但不阻止启动（不影响核心功能）。
- 现有的 `greeting.db` SQLite 数据库文件如何处理？→ **假定**：程序不再读取该文件，用户可自行删除；程序不负责清理文件系统中的遗留文件。
- 测试迁移：若没有可用的 MySQL 实例（如 CI 环境），测试如何运行？→ **已澄清（Option C）**：仅保留编译验证（`go build ./...`），暂不要求模型层/Handler 层测试通过。测试代码中移除 SQLite 引用，具体测试数据库方案后续单独处理。

---

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 项目 MUST 移除 `gorm.io/driver/sqlite` 依赖，`go.mod` 中不再包含该包。
- **FR-002**: `model.InitDB()` MUST 仅接受 MySQL DSN 参数，仅执行数据库连接与 Ping 验证，不执行任何表存在性检查。
- **FR-002a**: `main.go` 中 MUST 移除所有 `EnsureUserTable()` 调用；`model/db.go` 中 `EnsureUserTable()` 函数 MUST 被移除。
- **FR-003**: 全局变量 `model.SQLiteDB` MUST 被移除，项目中仅保留 `model.DB` 一个数据库实例。
- **FR-004**: `model.SQLiteUser` 结构体及其所有 CRUD 函数（`CreateSQLiteUser`、`GetSQLiteUserByID`、`GetSQLiteUserByPhone`、`UpdateSQLiteUser`、`DeleteSQLiteUser`）MUST 被移除。
- **FR-005**: `config.toml` 中 `[database.sqlite]` 配置段 MUST 被移除；`config/config.go` 中 `SQLiteConfig` 结构体 MUST 被移除。
- **FR-006**: `handler/demo.go` 中 `GetUserByPhoneTest` 方法 MUST 改为使用 MySQL（`model.DB` + `model.User`）。
- **FR-007**: `handler/common.go` 中 `Setting()` 返回的 `sqlite_dsn` 字段 MUST 被移除。
- **FR-008**: `migrations/001_user.sqlite.sql` 文件 MUST 被删除。
- **FR-009**: `model/schema.sql` 文件 MUST 被移除，或移除其中 SQLite 特定内容；`model/db.go` 中 `//go:embed schema.sql` 相关代码 MUST 被移除。
- **FR-010**: `model/db.go` 中 `ApplySchema()` 函数 MUST 被移除；`//go:embed schema.sql` 相关代码 MUST 被移除。
- **FR-011**: 所有测试文件（`model/*_test.go`、`handler/*_test.go`）MUST 不再导入 `gorm.io/driver/sqlite`；测试通过与否不在当前需求范围内（后续单独处理 MySQL 测试方案）。
- **FR-012**: `model/user_test.go` 中所有 `SQLiteUser` 专属测试用例 MUST 被移除。
- **FR-013**: `model/db_test.go` 中依赖 `ApplySchema` 或 `EnsureUserTable` 的测试 MUST 被移除（因这些函数将被移除）。
- **FR-014**: `README.md` 中技术栈、配置说明、API 列表等涉及 SQLite 的描述 MUST 更新为仅描述 MySQL。
- **FR-015**: 项目宪法 `.specify/memory/constitution.md` 技术栈约束表格中数据库一项 MUST 更新为 MySQL。

### Key Entities

- **数据库实例**: 从「MySQL + SQLite 双实例」简化为「MySQL 单实例」，全局变量从 `{DB, SQLiteDB}` 简化为 `{DB}`。
- **用户实体**: 从 `User`（MySQL）和 `SQLiteUser`（SQLite）两个实体合并为一个 `User` 实体，映射到 `users` 表。

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: `go.mod` 中 `gorm.io/driver/sqlite` 依赖已移除，`require` 块中不再出现该包名。
- **SC-002**: 项目生产代码（`*.go` 文件，不含 `_test.go`）中搜索 "sqlite"（不区分大小写）返回零匹配。
- **SC-003**: 项目配置与迁移文件中搜索 "sqlite"（不区分大小写）返回零匹配。
- **SC-004**: `go build ./...` 编译通过，无 SQLite 驱动相关编译错误。
- **SC-005**: 测试文件中不导入 `gorm.io/driver/sqlite`；`go build ./...` 编译通过（测试通过不在当前需求范围内）。
- **SC-006**: 原有的 MySQL CRUD 接口和迁移后的 `/demo/user/phone` 接口行为不变，数据正确读写 MySQL `users` 表。
- **SC-007**: `go vet ./...` 无警告。

---

## Assumptions

- 项目已具备 MySQL 数据库支持（基于 004-mysql-support 的成果），此次仅在已有 MySQL 基础上移除 SQLite。
- 架构层面保留多库并存的灵活性（符合宪法 v1.3.1），但当前实现仅保留 MySQL，移除所有 SQLite 具体代码。
- `handler/demo.go` 中的 `GetUserByPhoneTest` 功能仍需保留，仅将其数据存储从 SQLite 迁移至 MySQL。
- `InitDB()` 仅负责打开数据库连接并验证可达性，不检查任何表是否存在；表管理由用户通过迁移脚本负责（符合宪法 Principle VII）。
- 测试策略：采用 Option C——仅保留编译验证（`go build ./...`），测试通过不在当前需求范围内。测试代码中的 SQLite 引用仍需移除，MySQL 测试方案后续单独处理。
- 配置文件中移除 `[database.sqlite]` 段后，若用户配置文件仍保留该段，程序以兼容方式忽略（打警告日志）。
- 文件系统中的 `greeting.db` 遗留文件不由此需求处理。
