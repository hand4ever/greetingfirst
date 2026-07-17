# 任务(Tasks): 新增 SQLite 实例与 CRUD 测试接口

**输入(Input)**: Design documents from `/specs/011-sqlite-crud-test/`

**前置条件(Prerequisites)**: plan.md（必需）、spec.md（用户故事必需）、research.md、data-model.md、contracts/

**测试(Tests)**: 本特性将测试作为交付物的一部分（规格假设与 quickstart 场景 E 明确要求 `go test` 覆盖；宪法原则 V 要求每个 handler/model 方法单测）。故包含 model 与 handler 单测任务。

**组织(Organization)**: 任务按用户故事分组，使每个故事可独立实现与测试。

## 格式(Format): `[ID] [P?] [Story] 描述(Description)`

- **[P]**: 可并行执行（不同文件，无依赖）
- **[Story]**: 所属用户故事（US1–US4）
- 描述含具体文件路径

---

## 阶段(Phase) 1: 初始化(Setup)（共享基础设施）

**目的(Purpose)**: 引入依赖与用户自管的 `test_user` schema 资产（应用不执行该脚本，见原则 VII / FR-010）。

- [X] T001 [P] 在 go.mod 引入纯 Go SQLite 驱动：`go get github.com/glebarez/sqlite`，确认 go.sum 更新（research 任务 1）
- [X] T002 [P] 新增用户自管迁移脚本 `migrations/002_test_user.sql`，描述 `test_user` 表结构（id/name/phone/age/created_at/updated_at/deleted_at），`CREATE TABLE IF NOT EXISTS`；应用不执行，仅供用户手动建表与测试 `TestMain` 复用（research 任务 3、data-model §6）

**检查点(Checkpoint)**: 依赖就绪、schema 脚本可用（暂不建表）。

---

## 阶段(Phase) 2: 用户故事(User Story) 2 - 配置文件独立管理 SQLite 参数 (优先级(Priority): P1)

**目标(Goal)**: 在 `config.toml` 通过 `[database.sqlite]` 独立子节配置 SQLite 连接，与 `[database.mysql]` 解耦（FR-001、US2）。

**独立测试(Independent Test)**: 查看 `config.toml` 含 `[database.sqlite]`（默认 `greeting.db`），与 `[database.mysql]` 独立；修改 DSN 重启不影响 MySQL。

- [X] T003 [US2] 在 `config/config.go` 新增 `SQLiteConfig{ DSN string }` 结构体（toml tag `sqlite`）与 `DatabaseConfig.SQLite SQLiteConfig` 字段；在 `defaultConfig()` 中设置 `SQLite: SQLiteConfig{DSN: "greeting.db"}`（research 任务 6）
- [X] T004 [P] [US2] 在 `config.toml` 新增 `[database.sqlite]` 子节，`dsn = "greeting.db"`，与现有 `[database.mysql]` 并列

**检查点(Checkpoint)**: 配置层支持独立 SQLite 连接参数，默认指向 `greeting.db`。

---

## 阶段(Phase) 3: 用户故事(User Story) 1 - 启动建立 SQLite 独立实例 (优先级(Priority): P1) 🎯 MVP

**目标(Goal)**: 服务启动时初始化独立 SQLite 实例（`model.SQLiteDB`），与 `model.DB`（MySQL）共存；连接/Ping 失败或 `test_user` 表缺失即 fail-fast 退出（FR-002、FR-003、US1、原则 VI）。

**独立测试(Independent Test)**: `go run .` 启动日志显示 SQLite 连接成功且与 MySQL 共存；将 DSN 指向只读/不存在路径时服务明确报错退出；表缺失时 fail-fast（quickstart 场景 A）。

- [X] T005 [US1] 在 `model/db.go` 新增全局 `SQLiteDB *gorm.DB` 与 `InitSQLite(dsn string) error`：用 `glebarez/sqlite` 打开连接、`Ping()` 验证，失败返回带类型与地址的 `fmt.Errorf`（复用 `InitDB` 模式，research 任务 6）
- [X] T006 [US1] 在 `model/db.go` 的 `InitSQLite` 中，连接成功后执行 `SELECT 1 FROM test_user LIMIT 1` 校验表是否存在；表缺失返回明确 fail-fast 错误（FR-010、原则 VI/VII：仅校验不建表）
- [X] T007 [US1] 在 `main.go` 于 `model.InitDB(...)` 之后调用 `model.InitSQLite(config.Cfg.Database.SQLite.DSN)`；失败 `panic` 退出（与 MySQL 失败策略一致）

**检查点(Checkpoint)**: SQLite 实例随服务启动建立并与 MySQL 共存；连接或表缺失均 fail-fast。MVP 止于此可验证「实例可用」。

---

## 阶段(Phase) 4: 用户故事(User Story) 4 - 独立的 SQLite 测试表模型(TestUser) (优先级(Priority): P2)

**目标(Goal)**: 提供仅用于测试的 `TestUser` 模型与 CRUD 方法，映射预存的 `test_user` 表，与 MySQL `User` 完全隔离（FR-004、US4）。

**独立测试(Independent Test)**: `model` 包单测对 `:memory:` 实例执行 `002_test_user.sql` 后，可创建/查询/更新/软删/列表 `TestUser`；软删后同 phone 可复用（quickstart 场景 E）。

- [X] T008 [US4] 在 `model/testuser.go` 定义 `TestUser` 结构体（字段 ID int `primaryKey`、Name string、Phone string、Age int、CreatedAt/UpdatedAt `model.LocalTime`、DeletedAt `gorm.DeletedAt`），实现 `TableName() string { return "test_user" }`（data-model §1、research 任务 4）
- [X] T009 [US4] 在 `model/testuser.go` 实现基于 `SQLiteDB` 的 CRUD 方法：`CreateTestUser`、`GetTestUserByID`（GORM 自动排除软删）、`UpdateTestUser`、`DeleteTestUser`（软删）、`ListTestUsers`（`Find` 自动排除软删）；新增 `phoneActiveExists(phone) (bool, error)` 用于创建前活跃同号校验（research 任务 5）
- [X] T010 [P] [US4] 在 `model/testuser_test.go` 编写 `TestMain`：用 `glebarez/sqlite` 开 `:memory:` 并 `SetMaxOpenConns(1)`（research 任务 2），执行 `migrations/002_test_user.sql` 建表；用例覆盖 Create/Get/Update/Delete/List/phone 活跃唯一（含软删复用）；定义 `logOK` 辅助函数（原则 V、quickstart 场景 E）

**检查点(Checkpoint)**: `TestUser` 模型与单测就绪，可独立验证 SQLite 数据读写与软删语义。

---

## 阶段(Phase) 5: 用户故事(User Story) 3 - SQLite CRUD 测试接口 (优先级(Priority): P2)

**目标(Goal)**: 暴露专用于 `test_user` 的 5 个 CRUD 测试接口（`/sqlite/testuser`），独立于 MySQL `/demo/usr`（FR-005、FR-006、FR-007、FR-008、FR-009、FR-011、US3）。

**独立测试(Independent Test)**: 依次调用创建→查询→更新→再查询→删除→再查询（返回 not found），全链路数据一致（SC-003、quickstart 场景 B/C）。

### 用户故事(User Story) 3 的测试

- [X] T014 [P] [US3] 在 `handler/sqliteuser_test.go` 用 `httptest.NewRequest` + `echo.New().NewContext` 覆盖 5 个接口的成败分支（含 name/phone 必填、phone 重复、id 非法、not found、软删后再查）；定义/复用 `logOK`（原则 V、quickstart 场景 E）

### 用户故事(User Story) 3 的实现

- [X] T011 [P] [US3] 在 `entity/sqliteusr/sqliteusr.go` 定义 `TestUserCreateReq{ Name string \`json:"name"\`; Phone string \`json:"phone"\`; Age *int \`json:"age"\` }` 与 `TestUserUpdateReq{ Name/Phone *string; Age *int }`（data-model §3）
- [X] T012 [US3] 在 `handler/sqliteuser.go` 实现 `_SqliteUser` 处理器（包级变量 `var SqliteUser = &_SqliteUser{}`），方法 Create/Get/Update/Delete/List：绑定请求体、字段必填校验、创建前调 `phoneActiveExists` 校验、统一 `response.Ok/NotOk` 返回（原则 II）；复用 `extractUserID` 思路提取 `:id`（contracts/api.md）
- [X] T013 [US3] 在 `router/sqlite.go` 新增 `sqlite(e *echo.Echo)` 注册 `e.Group("/sqlite/testuser")` 的 POST ``、GET `/:id`、PUT `/:id`、DELETE `/:id`、GET `testusers`（`/sqlite/testusers`）；在 `router/router.go` 的 `Router(e)` 中调用 `sqlite(e)`（原则 I）

**检查点(Checkpoint)**: 5 个 SQLite CRUD 接口完整可用，与 MySQL 接口隔离、行为独立。

---

## 阶段(Phase) 6: 收尾与横切关注点(Polish & Cross-Cutting Concerns)

**目的(Purpose)**: 文档、用例与整体校验（FR-012 不受影响）。

- [X] T015 [P] 在 `api.http` 新增 `/sqlite/testuser` 的 REST Client 用例：创建、按 ID 查询、更新、删除、列表（quickstart 场景 B/C）
- [X] T016 [P] 在 `README.md` 的「API 列表」与「更新日志」中记录新增 SQLite CRUD 接口（开发流程规范）
- [X] T017 执行整体校验：`go build ./...` 编译通过；`go test -v ./... -count=1` 全部通过；`gofmt`/`gofumpt` 格式化（提交前必须）；确认现有 MySQL 接口与 `go test` 行为不变（SC-007）

---

## 依赖与执行顺序(Dependencies & Execution Order)

### 阶段依赖(Phase Dependencies)

- **初始化(Setup)（阶段 1）**: 无依赖 — 可立即开始
- **US2 配置（阶段 2）**: 依赖阶段 1 驱动可用（T001）；T004 独立
- **US1 实例（阶段 3）**: 依赖 T003（配置结构）+ T001（驱动）
- **US4 模型（阶段 4）**: 依赖阶段 3 `SQLiteDB` 与 T002 迁移脚本
- **US3 接口（阶段 5）**: 依赖阶段 4 模型方法 + 阶段 2 配置
- **收尾（阶段 6）**: 依赖阶段 3–5 全部完成

### 用户故事依赖(User Story Dependencies)

- **US2（P1）**: 阶段 1 后开始；不与其它故事冲突
- **US1（P1）**: 依赖 US2 的配置结构 + 驱动；MVP 核心
- **US4（P2）**: 依赖 US1 的 `SQLiteDB`
- **US3（P2）**: 依赖 US4 的模型方法 + US2 配置；UI/接口层

### 用户故事内部(Within Each User Story)

- 模型先于方法（T008 → T009）
- 模型/实体先于处理器（T009/T011 → T012）
- 处理器先于路由（T012 → T013）
- 测试（T010、T014）可在对应实现文件创建后编写并先于运行失败
- 当前故事完成后再进入下一优先级

### 并行机会(Parallel Opportunities)

- 阶段 1：T001 与 T002 可并行（[P]）
- 阶段 2：T004（config.toml）与 T003（config.go）可并行（[P]）
- 阶段 4：T010（测试文件）与 T008/T009 可并行编写（[P]，不同文件）
- 阶段 5：T011（实体）可并行；T014（handler 测试）可并行（[P]，不同文件）
- 阶段 6：T015 与 T016 可并行（[P]）

---

## 并行示例(Parallel Example): 用户故事(User Story) 3

```bash
# 同时创建实体与（稍后）测试文件：
Task: "在 entity/sqliteusr/sqliteusr.go 定义请求体"            # T011 [P]
Task: "在 handler/sqliteuser_test.go 编写 5 接口单测"          # T014 [P]

# 实现处理器后再注册路由（串行）：
Task: "在 handler/sqliteuser.go 实现 _SqliteUser"              # T012
Task: "在 router/sqlite.go 注册 /sqlite/testuser 路由"              # T013
```

---

## 实施策略(Implementation Strategy)

### MVP 优先(MVP First)（US2 + US1）

1. 完成阶段 1：初始化（驱动 + 迁移脚本）
2. 完成阶段 2：US2 配置参数
3. 完成阶段 3：US1 SQLite 实例（[T005][T006][T007]）
4. **停止并验证**: `go run .` 确认 SQLite 与 MySQL 共存、表缺失 fail-fast
5. 如需最小可用，可至此演示；否则继续

### 增量交付(Incremental Delivery)

1. 初始化 + US2 + US1 → SQLite 实例可用（MVP）
2. 添加 US4 → TestUser 模型与单测（数据层就绪）
3. 添加 US3 → 5 个 CRUD 测试接口（核心交付）
4. 收尾 → api.http / README / 整体校验
5. 每阶段独立测试，不影响已有 MySQL 模块（SC-007）

### 并行团队策略(Parallel Team Strategy)

- 单人顺序执行即可（特性规模小）
- 若多人：A 负责阶段 1–3（基础），B 并行负责阶段 4–5（模型+接口），收尾共同校验

---

## 备注(Notes)

- [P] 任务 = 不同文件，无依赖，可并行
- [Story] 标签将任务映射到用户故事（US1–US4），便于追踪
- 每个用户故事应可独立完成与测试
- 应用**绝不**建表/迁移（原则 VII / FR-010）；`migrations/002_test_user.sql` 仅供用户手动执行与测试 `TestMain` 复用
- 所有代码注释英文、commit message 英文（原则 IV）
- 每任务或逻辑组完成后提交；可在任意检查点暂停验证
