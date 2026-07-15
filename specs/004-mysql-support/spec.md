# Feature Specification: MySQL 数据库支持

**Feature Branch**: `004-mysql-support`
**Created**: 2026-07-15
**Status**: Draft
**Input**: User description: "增加对mysql的支持，数据的主要联动是通过mysql来完成，建立全局mysql实例，在配置文件里增加对mysql的配置，默认连接本地localhost，demo数据库。并在demo里写个测试增删改查的相关接口。"

## Clarifications

### Session 2026-07-15

- Q: 配置文件如何同时保留 SQLite 与 MySQL 连接信息并允许切换？ → A: （已废弃，见下方重新澄清）原来采用 Option B `type` 选择器，后被推翻。
- Q: GET /demo/users 列表接口是否需要分页？ → A: 不分页，返回全量（保持当前 spec）。
- Q: 软删除后通过相同 phone 创建新用户的行为？ → A: 允许新建并存——phone 在未删除记录中唯一，软删除后可复用（不复活旧记录，新建独立记录）。
- Q: MySQL 与 SQLite 是二选一切换还是共存独立？ → A: **共存独立**。MySQL 和 SQLite 同时连接，各自拥有独立的 model、独立的 DB 实例、独立的接口，互不干涉。
- Q: MySQL CRUD 接口路由前缀？ → A: `/demo/usr` 前缀（如 `POST /demo/usr`、`GET /demo/usr/:id`），SQLite 现有 demo 接口保持不变。
- Q: Model 结构体如何分离？ → A: SQLite 使用 `model.SQLiteUser`（原 `User` 重命名），MySQL 使用 `model.User`。各自独立。
- Q: 全局 DB 实例命名？ → A: `model.DB` 指向 MySQL，`model.SQLiteDB` 指向 SQLite。
- Q: 配置文件是否需要 `type` 选择器？ → A: 不需要。`[database.mysql]` 和 `[database.sqlite]` 各自独立配置，启动时同时连接两个数据库。
- Q: 现有 `handler/demo.go` 使用哪个 DB？ → A: 迁移到 `model.SQLiteDB` + `model.SQLiteUser`，保持原有 SQLite 行为不变。MySQL 使用全新 `handler/user.go`。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - MySQL 与 SQLite 同时连接启动服务 (Priority: P1)

作为开发者，我希望启动服务时系统自动同时连接 MySQL 和 SQLite 两个数据库，各自独立运行，互不干扰。

**Why this priority**: 双数据库共存连接是整个 MySQL 支持的基础，没有连接无法进行任何后续操作。

**Independent Test**: 配置 `config.toml` 中 `[database.mysql]` 和 `[database.sqlite]` 两个独立子节，启动服务，观察日志确认同时成功连接 MySQL 和 SQLite，服务正常运行。

**Acceptance Scenarios**:

1. **Given** 配置文件中 `[database.mysql]` 和 `[database.sqlite]` 均已正确配置，MySQL 服务已启动且 demo 数据库已创建，**When** 启动服务，**Then** 服务成功启动，日志显示 MySQL 和 SQLite 均连接成功，服务端口正常监听。

2. **Given** MySQL 连接信息错误或服务未启动，**When** 启动服务，**Then** 服务启动失败，打印明确的 MySQL 连接错误信息并退出（SQLite 连接成功也不影响 MySQL 连接失败时的退出）。

3. **Given** SQLite 连接失败（如文件权限问题），**When** 启动服务，**Then** 服务启动失败，打印明确的 SQLite 连接错误信息并退出。

---

### User Story 2 - 配置文件中独立管理 MySQL 与 SQLite 连接参数 (Priority: P1)

作为开发者，我希望在 `config.toml` 中分别独立配置 MySQL 和 SQLite 的连接参数，两个数据库各自拥有独立的配置子节，互不依赖。

**Why this priority**: 配置管理是数据库可替换性的关键，没有配置就无法灵活管理两个数据库。

**Independent Test**: 打开 `config.toml`，确认 `[database.mysql]` 和 `[database.sqlite]` 两个独立子节，各含 `dsn` 字段，MySQL 默认指向本地 `demo` 数据库，SQLite 默认 `greeting.db`。

**Acceptance Scenarios**:

1. **Given** 项目首次拉取，**When** 查看 `config.toml` 默认配置，**Then** 同时包含 `[database.mysql]`（dsn 默认指向本地 localhost 的 `demo` 数据库）与 `[database.sqlite]`（dsn 默认 `greeting.db`）两个独立子节，无 `type` 选择器。

2. **Given** 开发者修改 `config.toml` 中任一数据库的连接信息，**When** 重启服务，**Then** 服务使用新的连接信息连接对应数据库，另一个数据库不受影响。

---

### User Story 3 - MySQL Demo 模块完整 CRUD 接口 (Priority: P2)

作为前端开发者，我希望通过 `/demo/usr` 前缀的接口对 MySQL 数据库中的用户数据进行增删改查操作，验证 MySQL 数据持久化能力。SQLite 侧现有 `/demo` 接口保持不变。

**Why this priority**: CRUD 接口是数据联动的具体体现，也是测试 MySQL 读写正确性的方式。

**Independent Test**: 使用 HTTP 客户端依次调用 `/demo/usr` 系列接口（创建→查询→更新→查询→删除），验证每一步返回的数据一致性。

**Acceptance Scenarios**:

1. **Given** MySQL 服务正常运行且 `users` 表已创建，**When** 调用 `POST /demo/usr` 带 name="张三"、phone="13800138000"、age=25 创建用户，**Then** 返回成功，响应中包含新用户的 ID 和创建时间。

2. **Given** 已创建用户 ID=1，**When** 调用 `GET /demo/usr/1` 查询用户，**Then** 返回用户信息，name="张三"、phone="13800138000"、age=25。

3. **Given** 已创建用户 ID=1，**When** 调用 `PUT /demo/usr/1` 带 name="张三丰"、age=30 更新用户，**Then** 返回成功，再次查询确认 name 和 age 已更新，phone 不变。

4. **Given** 已创建用户 ID=1，**When** 调用 `DELETE /demo/usr/1` 删除用户，**Then** 返回成功，再次查询用户 ID=1 返回"用户不存在"错误。

5. **Given** MySQL 中 `users` 表为空，**When** 调用 `GET /demo/usrs` 查询用户列表，**Then** 返回空列表。

6. **Given** 已创建多个用户，**When** 调用 `GET /demo/usrs` 查询用户列表，**Then** 返回所有用户信息列表。

---

### User Story 4 - 用户手动建表 (Priority: P2)

作为开发者，我通过项目提供的 SQL 迁移脚本手动创建 MySQL 和 SQLite 各自的数据库表结构，启动服务时应用直接连接已存在的表，不自动建表。

**Why this priority**: 表结构属于数据资产，由用户显式管理可避免 schema 漂移与环境不一致（遵循 Constitution Principle VII）。

**Independent Test**: 使用提供的 SQL 脚本在 demo 数据库中创建 MySQL `users` 表，启动服务，调用 `/demo/usr` CRUD 接口，确认接口正常返回。

**Acceptance Scenarios**:

1. **Given** demo 数据库中已通过 SQL 脚本创建 users 表（字段含 id、phone、name、age、created_at、updated_at、deleted_at），**When** 启动服务，**Then** MySQL 连接正常，`/demo/usr` 接口正常工作，不执行任何建表或迁移操作。

2. **Given** demo 数据库中**尚未**创建 users 表，**When** 启动服务，**Then** 服务打印明确提醒（指明使用 `migrations/001_user.mysql.sql` 建表）并暂停等待；用户在另一端执行建表脚本后，服务自动检测到表并继续启动（不退出、不自动建表）。

---

### Edge Cases

- MySQL 连接失败（密码错误、服务未启动、网络不通）时，服务应明确报错退出，不能静默降级为仅 SQLite 运行。
- SQLite 连接失败（文件权限问题等）时，服务应明确报错退出；一个数据库连接失败即整体启动失败。
- MySQL `users` 表缺失时，应用打印提醒并**暂停等待**，用户通过 `migrations/001_user.mysql.sql` 建表后服务自动继续，不自动建表、不退出。
- SQLite 表结构由 SQL 脚本管理，测试环境通过 `ApplySchema` 加载（不与 MySQL 冲突）。
- DSN 中字符集设置不正确时可能导致中文乱码，应在默认配置中使用 `utf8mb4`。
- 并发请求时，两个全局 DB 实例各自应能正确处理（GORM 连接池默认支持）。
- MySQL `users` 表 phone 字段在未删除记录中保持唯一约束，同 phone 创建未删除用户应返回明确错误信息；软删除后同 phone 可新建并存。
- 软删除用户后，通过 phone 创建新用户应允许（与原记录不冲突）。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: `model.InitDB` MUST 同时连接 MySQL 和 SQLite 两个数据库，分别从 `[database.mysql]` 和 `[database.sqlite]` 读取 DSN。MySQL 实例存为 `model.DB`，SQLite 实例存为 `model.SQLiteDB`。
- **FR-002**: 系统 MUST 在 `config.toml` 中包含 `[database.mysql]` 与 `[database.sqlite]` 两个独立子节（各含独立 `dsn`），无 `type` 选择器。启动时同时连接两个库，任一连不上即报错退出。MySQL 默认 `dsn` 为 `"root:password@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"`，SQLite 默认 `dsn` 为 `"greeting.db"`。
- **FR-003**: MySQL 侧 MUST 提供用户 CRUD 接口：创建用户 `POST /demo/usr`、查询单个用户 `GET /demo/usr/:id`、更新用户 `PUT /demo/usr/:id`、删除用户 `DELETE /demo/usr/:id`、查询用户列表 `GET /demo/usrs`。SQLite 侧现有 `/demo` 接口保持不变。
- **FR-004**: 创建用户接口 MUST 接受 `name`（必填）、`phone`（必填、在未删除记录中唯一）、`age`（选填，使用 `*int` 指针类型，`nil` 表示未传、`0` 表示显式传 0）参数，返回新用户的完整信息。Phone 仅与当前未软删除的记录冲突；若已有同名 phone 的软删除记录，允许新建并存。
- **FR-005**: 查询用户列表接口 MUST 返回所有未软删除的用户，按创建时间倒序排列。
- **FR-006**: 更新用户接口 MUST 支持部分字段更新，未传字段保持原值不变。
- **FR-007**: 删除用户接口 MUST 使用软删除（GORM DeletedAt），不物理删除记录。
- **FR-008**: MySQL 数据库表结构 MUST 由用户通过 SQL 迁移脚本手动创建；应用启动 MUST NOT 自动建表或执行 AutoMigrate。SQLite 表结构同样通过独立 SQL 脚本管理（测试环境通过 `ApplySchema` 加载）。
- **FR-009**: 任一数据库连接失败时 MUST 打印具体错误原因（包含数据库类型、目标地址），并以非零状态码退出。
- **FR-010**: `model.DB` 指向 MySQL 实例，`model.SQLiteDB` 指向 SQLite 实例，两者同时运行、互不干扰。MySQL 使用 `model.User`，SQLite 使用 `model.SQLiteUser`（现有 `User` 重命名）。
- **FR-011**: 当 MySQL 数据库中**不存在**所需的 `users` 表时，应用 MUST 打印明确提醒（指明应使用 `migrations/001_user.mysql.sql` 脚本手动建表），随后**暂停等待**（无限期，`maxWait=0`），以 3 秒间隔周期性检测，直到用户创建该表后**自动继续**；MUST NOT 自动建表、MUST NOT 直接退出。
- **FR-012**: 现有 `handler/demo.go` 中的接口 MUST 迁移至使用 `model.SQLiteDB` + `model.SQLiteUser`，保持原有 SQLite 行为不变。

### Key Entities

- **User（MySQL 实体，`model.User`）**: 新增 MySQL 用户实体，字段包括 ID、Phone（在未删除记录中唯一，软删除后可复用）、Name、Age、CreatedAt、UpdatedAt、DeletedAt（软删除）。
- **SQLiteUser（SQLite 实体，`model.SQLiteUser`）**: 现有 `model.User` 重命名为此，保持原有 SQLite 用户数据结构不变。与 MySQL `User` 完全独立。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 启动服务时，MySQL 和 SQLite 两个数据库各自在 3 秒内完成连接（总计不超过 3 秒），连接完成后端口开始监听。
- **SC-002**: 开发者仅需配置 `[database.mysql]` 和 `[database.sqlite]` 两个独立子节，无需 `type` 选择器，两个数据库同时连接、共存运行。
- **SC-003**: MySQL CRUD 5 个接口全部调用成功，创建→查询→更新→再查询→删除→再查询（返回不存在）的完整链路数据一致。
- **SC-004**: 并发调用创建接口 10 次（不同 phone），10 条记录全部创建成功，无数据丢失或冲突。
- **SC-005**: 任一数据库连接失败时，服务启动时间不超过 5 秒即退出并输出错误信息（而非无限等待或静默降级为单库运行）。
- **SC-006**: MySQL `users` 表缺失时，服务打印提醒并暂停，用户建表后自动继续（不退出、不自动建表）。连接失败（非表缺失）仍按 SC-005 退出。

## Assumptions

- 开发者本地已安装 MySQL 服务（默认 127.0.0.1:3306），且已创建 `demo` 数据库。
- MySQL 用户 `root` 对 `demo` 数据库有完整权限（CREATE TABLE、CRUD 等）。
- 项目中已有 GORM v2 依赖，无需额外引入新 ORM 框架。
- MySQL 和 SQLite 各自拥有独立的 model 结构体：`model.User`（MySQL）和 `model.SQLiteUser`（SQLite，由现有 `User` 重命名）。
- MySQL CRUD 接口路由使用 `/demo/usr` 前缀，SQLite 现有 `/demo` 路由保持不变。
- 配置文件的默认 MySQL DSN 使用 `root:password`，开发者在实际使用时需自行修改密码。
- SQLite 数据库文件默认路径为 `greeting.db`，与项目同级目录。

## Dependencies

- 需要 `gorm.io/driver/mysql` 驱动包（新增 Go 依赖）。
- 依赖本地 MySQL 5.7+ 或 8.0+ 服务。
- 依赖已有的 `config/config.go` 配置加载模块和 `model/db.go` 数据库初始化模块。
