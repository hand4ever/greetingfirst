# 功能规格说明(Feature Specification): 新增 SQLite 实例与 CRUD 测试接口

**Feature Branch**: `011-sqlite-crud-test`

**创建日期(Created)**: 2026-07-17

**状态(Status)**: 草稿(Draft)

**输入(Input)**: 用户描述(User description): "新增sqlite实例，并写个测试接口测试sqlite的增删改查"

## Clarifications

### Session 2026-07-17

- Q: 是否需要为 SQLite 侧新增与 MySQL 对等的「用户实体」与「用户表」？ → A: 否。本次不新增产品级用户实体、不新增（与生产 `users` 对等的）用户表；直接使用已预先存在的、仅用于测试的 `test_user` 表（应用不负责创建）。
- Q: SQLite 测试表是否需要通过 SQL 迁移脚本管理、应用启动是否自动建表？ → A: 否，且不是应用创建。`test_user` 测试表已由用户预先创建并存在，应用 MUST NOT 创建或迁移表结构（完全符合宪法原则 VII：schema 用户自管理）；若表缺失则按 fail-fast 报错退出。
- Q: SQLite 测试接口是否需要与现有 MySQL 接口（`/demo/usr`）对齐（字段语义、路由结构等）？ → A: 否。无需对齐 MySQL 接口，可独立、简单地设计测试接口。
- Q: SQLite 测试接口的路由命名是否采用 `testuser`？ → A: 是。单条操作统一前缀 `/sqlite/testuser`（POST 创建、GET/PUT/DELETE `/:id`），列表接口为 `/sqlite/testusers`；据此消除 plan/tasks 中 `/sqlite/usr`(+`/sqlite/usrs`) 与 contracts 中 `/sqlite/usr/s` 的命名不一致（CHK011）。

## 用户场景与测试(User Scenarios & Testing) *(必填)*

### 用户故事(User Story) 1 - 启动服务时建立 SQLite 独立实例 (优先级(Priority): P1)

作为开发者，我希望服务启动时能够初始化一个独立的 SQLite 数据库实例，与现有 MySQL 实例共存、互不干扰，用于轻量级数据读写与测试。

**优先级理由(Why this priority)**: 独立的 SQLite 实例是整个特性的基础，没有实例就无法进行任何数据持久化与接口测试。

**独立测试(Independent Test)**: 在 `config.toml` 中配置 `[database.sqlite]` 连接信息（默认指向本地 `greeting.db` 文件），启动服务，观察启动日志确认 SQLite 成功连接，同时现有 MySQL 实例仍正常工作。

**验收场景(Acceptance Scenarios)**:

1. **假设(Given)** 配置文件中 `[database.sqlite]` 已正确配置且目标文件可写，**当(When)** 启动服务，**Then(Then)** 服务成功启动，日志显示 SQLite 连接成功，MySQL 实例同时可用。
2. **假设(Given)** SQLite 配置指向一个只读目录（文件无法创建），**When** 启动服务，**Then** 服务启动失败并打印明确的 SQLite 连接错误信息，退出（不影响错误信息的明确性）。
3. **假设(Given)** 未显式配置 `[database.sqlite]`，**When** 启动服务，**Then** 使用默认 DSN（`greeting.db`）成功建立 SQLite 实例。

---

### 用户故事(User Story) 2 - 配置文件中独立管理 SQLite 连接参数 (优先级(Priority): P1)

作为开发者，我希望在 `config.toml` 中通过 `[database.sqlite]` 独立子节配置 SQLite 连接参数，与 MySQL 配置互不依赖，可分别调整。

**优先级理由(Why this priority)**: 独立配置是数据库可替换性的关键，保证 SQLite 与 MySQL 各自拥有独立配置与生命周期。

**独立测试(Independent Test)**: 打开 `config.toml`，确认存在 `[database.sqlite]` 子节（含 `dsn` 字段），默认值指向 `greeting.db`，且与 `[database.mysql]` 完全独立。

**验收场景(Acceptance Scenarios)**:

1. **假设(Given)** 项目首次拉取，**When** 查看 `config.toml` 默认配置，**Then** 同时包含 `[database.sqlite]`（默认 `greeting.db`）与现有的 `[database.mysql]`，两者独立无耦合。
2. **假设(Given)** 开发者修改 `[database.sqlite]` 的 DSN，**When** 重启服务，**Then** 服务使用新的 SQLite 连接信息，MySQL 不受影响。

---

### 用户故事(User Story) 3 - SQLite 用户数据 CRUD 测试接口 (优先级(Priority): P2)

作为开发者/测试人员，我希望通过一组专用接口对 SQLite 中的用户数据进行增删改查，以验证 SQLite 数据持久化与读写正确性，且不影响现有 MySQL 接口。

**优先级理由(Why this priority)**: 该测试接口是本次需求的核心交付物，直接体现 "用接口测试 SQLite 增删改查" 的意图。

**独立测试(Independent Test)**: 使用 HTTP 客户端依次调用 SQLite CRUD 接口（创建→按 ID 查询→更新→再查询→删除→再查询返回不存在），验证每一步返回的数据一致性与 SQLite 持久化生效。

**验收场景(Acceptance Scenarios)**:

1. **假设(Given)** SQLite 实例正常运行且用户表已创建，**When** 调用创建接口（带 name、phone、age）创建用户，**Then** 返回成功，响应包含新用户的 ID 与创建时间。
2. **假设(Given)** 已创建用户 ID=1，**When** 调用按 ID 查询接口，**Then** 返回该用户信息（name、phone、age 与原创建一致）。
3. **假设(Given)** 已创建用户 ID=1，**When** 调用更新接口修改 name 与 age，**Then** 返回成功，再次查询确认字段已更新且未传字段保持不变。
4. **假设(Given)** 已创建用户 ID=1，**When** 调用删除接口，**Then** 返回成功，再次查询 ID=1 返回"用户不存在"错误。
5. **假设(Given)** 用户表为空，**When** 调用列表查询接口，**Then** 返回空列表。
6. **假设(Given)** 已创建多个用户，**When** 调用列表查询接口，**Then** 返回所有未删除用户的列表。

---

### 用户故事(User Story) 4 - 独立的 SQLite 测试表(test_user) (优先级(Priority): P2)

作为开发者，我希望 SQLite 侧使用一个仅用于测试的 `test_user` 表，与 MySQL 用户数据完全隔离，避免任何数据交叉；该测试表由用户预先创建并存在（遵循宪法原则 VII，应用不创建、不迁移表结构）。

**优先级理由(Why this priority)**: 测试表与 MySQL 生产数据隔离是底线，确保 SQLite 测试数据不会污染 MySQL；同时省去迁移脚本，降低测试环境搭建成本。

**独立测试(Independent Test)**: 在 SQLite 中通过测试接口写入 `test_user` 后，查询 MySQL 用户表，确认其中不出现任何 SQLite 创建的记录；反之同理。

**验收场景(Acceptance Scenarios)**:

1. **假设(Given)** 通过 SQLite 测试接口写入 `test_user` 一条记录，**When** 查询 MySQL 用户表，**Then** MySQL 中不存在该记录。
2. **假设(Given)** `test_user` 表已由用户预先创建并存在，**When** 启动服务并调用 SQLite 测试接口，**Then** 接口正常工作，应用不执行任何建表或迁移操作。

---

### 边界情况(Edge Cases)

- SQLite 实例初始化失败（文件权限、磁盘满等）时，服务应打印明确错误信息并退出；一个数据库连接失败即整体启动失败（与 MySQL 失败策略一致）。
- `test_user` 表由用户预先创建并存在，应用 MUST NOT 在启动时创建或迁移表结构；若表缺失，应用应打印明确错误并按 fail-fast 退出（遵循宪法原则 VII）。
- 创建用户时 phone 重复（在未删除记录中）应返回明确错误；软删除后同 phone 可新建并存。
- 更新接口未传字段应保持不变（部分更新语义）。
- 删除接口应使用软删除，不物理删除记录。
- 并发调用创建接口时，SQLite 实例应能正确处理（GORM 连接池/事务支持），无数据丢失。
- `:memory:` 模式仅用于单元测试，不影响生产文件型实例。

## 需求(Requirements) *(必填)*

### 功能需求(Functional Requirements)

- **FR-001**: 系统 MUST 在 `config.toml` 中新增 `[database.sqlite]` 子节（含独立 `dsn`），默认值为 `"greeting.db"`，与现有 `[database.mysql]` 完全独立。
- **FR-002**: `model` 层 MUST 新增独立的全局 SQLite 实例变量（与现有 MySQL 实例 `model.DB` 并存），并在初始化时从 `[database.sqlite]` 读取 DSN 建立连接、执行 Ping 验证。
- **FR-003**: 初始化时任一数据库（MySQL 或 SQLite）连接失败 MUST 打印具体错误原因（含数据库类型与目标地址）并以非零状态码退出（fail-fast）。
- **FR-004**: SQLite 侧 MUST 提供仅用于测试的 `TestUser` 模型（映射 `test_user` 表，字段含 ID、Name、Phone、Age、CreatedAt、UpdatedAt、DeletedAt 等），与 MySQL 用户实体完全分离，仅用于测试接口，不作为产品级用户实体。
- **FR-005**: 系统 MUST 提供一组专用于 SQLite `test_user` 表的 CRUD 测试接口：创建、按 ID 查询、更新、删除、列表查询；接口设计独立于现有 MySQL 接口，可自由定义路由与字段语义。
- **FR-006**: 创建用户接口 MUST 接受 `name`（必填）、`phone`（必填、在未删除记录中唯一）、`age`（选填）参数，返回新用户完整信息。
- **FR-007**: 更新用户接口 MUST 支持部分字段更新，未传字段保持原值不变。
- **FR-008**: 删除用户接口 MUST 使用软删除，不物理删除记录。
- **FR-009**: 查询用户列表接口 MUST 返回所有未软删除的用户（可按创建时间倒序）。
- **FR-010**: `test_user` 为测试专用表，由用户预先创建并存在；应用启动 MUST NOT 创建或迁移该表结构（遵循项目宪法原则 VII：schema 用户自管理）。若表缺失，应用 MUST 按 fail-fast 报错退出。
- **FR-011**: 所有接口 MUST 使用项目统一的 `response.Ok` / `response.NotOk` / `response.NotOkWithCode` 返回 JSON。
- **FR-012**: 现有 MySQL CRUD 接口与行为 MUST 保持不变，不受本次 SQLite 特性影响。

### 关键实体(Key Entities) *(涉及数据时填写)*

- **TestUser（SQLite 测试实体）**: 仅用于测试的 `test_user` 表映射模型，字段包括 ID、Name、Phone、Age、CreatedAt、UpdatedAt、DeletedAt（软删除）。与 MySQL 用户实体完全独立、互不干涉，仅服务于测试接口。
- **SQLite 实例**: 新增的全局 SQLite 数据库连接实例，与现有 MySQL 实例共存，各自拥有独立的连接池与生命周期。

## 成功标准(Success Criteria) *(必填)*

### 可衡量成果(Measurable Outcomes)

- **SC-001**: 启动服务时，SQLite 实例在 3 秒内完成连接，与 MySQL 实例同时可用，端口正常监听。
- **SC-002**: 开发者仅需配置 `[database.sqlite]` 子节（或接受默认值）即可启用 SQLite，无需改动任何 MySQL 配置。
- **SC-003**: SQLite CRUD 5 个接口全部调用成功，创建→查询→更新→再查询→删除→再查询（返回不存在）的完整链路数据一致。
- **SC-004**: 并发调用 SQLite 创建接口 10 次（不同 phone），10 条记录全部创建成功，无数据丢失或冲突。
- **SC-005**: 任一数据库连接失败时，服务在 5 秒内退出并输出明确错误信息（而非静默降级为单库运行）。
- **SC-006**: 通过 SQLite 接口写入的数据在 MySQL 用户表中完全不可见，反之亦然，双库数据完全隔离。
- **SC-007**: 现有 MySQL CRUD 接口行为不变，回归测试全部通过。

## 假设(Assumptions)

- 项目中已具备 GORM v2 及 `gorm.io/driver/sqlite` 驱动（若已在 005-remove-sqlite 中移除，则需重新引入该依赖）。
- SQLite 实例采用文件型数据库（默认 `greeting.db`），与现有 MySQL 实例共存，二者独立。
- SQLite CRUD 测试接口的路由与字段语义可自由设计，无需与现有 MySQL `/demo/usr` 接口对齐，只要不与现有路由冲突即可。
- `test_user` 表由用户预先创建并存在，应用不创建、不迁移表结构，完全符合宪法原则 VII（schema 用户自管理）。
- 单元测试可借助 SQLite `:memory:` 模式为模型层/Handler 层测试提供轻量数据库支撑，使 `go test` 在无 MySQL 实例时也能运行（本特性为后续完善测试基础设施铺路）。
- 默认 SQLite DSN 为 `greeting.db`，位于项目同级目录，开发者可自行修改路径。

## 依赖(Dependencies)

- 依赖 `gorm.io/driver/sqlite` 驱动包（若已在 005 中移除则需重新加入 `go.mod`）。
- 依赖已有的 `config/config.go` 配置加载模块与 `model` 数据库初始化模块。
- 依赖项目现有 `response` 统一响应封装与 `router` 路由注册机制。
