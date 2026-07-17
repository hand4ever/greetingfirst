# 快速开始与验证指南(Quickstart): 新增 SQLite 实例与 CRUD 测试接口

**功能分支(Feature Branch)**: `011-sqlite-crud-test`
**日期(Date)**: 2026-07-17
**关联规格(Spec)**: [spec.md](./spec.md) | **契约**: [contracts/api.md](./contracts/api.md) | **数据模型**: [data-model.md](./data-model.md)

本文件提供端到端验证场景，证明特性可用。**不含实现代码**，实现细节见后续 `tasks.md` 与编码阶段。

---

## 前置条件(Prerequisites)

1. Go ≥ 1.22（当前 `go.mod` 为 1.26.3）。
2. 已存在 `greeting.db` 且其中包含 `test_user` 表（用户已手动创建；如需重建，执行 `migrations/002_test_user.sql`）。
3. MySQL 实例可用（现有 `/demo/usr` 不受影响，SC-007）。
4. 安装 sqlite 驱动依赖：`go get github.com/glebarez/sqlite`（见研究任务 1）。

---

## 场景 A：启动验证 SQLite 实例共存（SC-001 / SC-002）

```bash
# 默认使用 greeting.db；如需自定义，编辑 config.toml 的 [database.sqlite].dsn
go run . 
# 或：go build ./... && ./greeting
```

**期望(Expected)**:
- 服务正常监听 `:1323`，启动日志显示 MySQL 与 SQLite 均连接成功。
- 不修改任何 MySQL 配置即可启用 SQLite（SC-002）。
- 若 `greeting.db` 无法打开（如只读目录），服务打印明确 SQLite 错误并以非零状态码退出（SC-005、原则 VI）。

---

## 场景 B：CRUD 全链路（SC-003）

使用 `api.http` 或任意 HTTP 客户端，依次执行（契约见 [contracts/api.md](./contracts/api.md)）：

```http
### 1. 创建
POST /sqlite/testuser
Content-Type: application/json

{ "name": "张三", "phone": "13900138000", "age": 28 }

### 2. 按 ID 查询（假设返回 id=1）
GET /sqlite/testuser/1

### 3. 更新
PUT /sqlite/testuser/1
Content-Type: application/json

{ "name": "张三丰", "age": 30 }

### 4. 再查询确认更新生效
GET /sqlite/testuser/1

### 5. 删除（软删除）
DELETE /sqlite/testuser/1

### 6. 再查询应返回 user not found
GET /sqlite/testuser/1

### 7. 列表（空或部分）
GET /sqlite/testusers
```

**期望(Expected)**:
- 步骤 1 返回含 `id` 与 `created_at` 的完整用户（契约 §1）。
- 步骤 2 / 4 字段一致；步骤 4 的 `name=张三丰`、`age=30`，未传字段不变（FR-007）。
- 步骤 5 成功；步骤 6 返回 `user not found`（FR-008 软删除）。
- 步骤 7 返回未删除用户列表；若仅建/删该用户则返回 `[]`（FR-009）。

---

## 场景 C：phone 唯一性与软删复用（Edge Case）

```http
POST /sqlite/testuser  { "name":"A", "phone":"13900001234" }   # 成功
POST /sqlite/testuser  { "name":"B", "phone":"13900001234" }   # 返回 phone already exists
DELETE /sqlite/testuser/<A的id>                                # 软删除 A
POST /sqlite/testuser  { "name":"C", "phone":"13900001234" }   # 复用成功（活跃唯一校验不命中已删除）
```

**期望(Expected)**: 第二次创建报错；软删后第三次创建成功（研究任务 5）。

---

## 场景 D：双库隔离（SC-006）

- 通过 `/sqlite/testuser` 写入的数据，在 MySQL `users` 表中不可见；反之同理。可在 MySQL 侧 `SELECT * FROM users` 验证无 SQLite 写入记录。

---

## 场景 E：单元测试（原则 V / 研究任务 2、3）

```bash
go test -v ./... -count=1
```

**期望(Expected)**:
- `model` 包：新增 `TestCreateTestUser` / `TestGetTestUserByID` / `TestUpdateTestUser` / `TestDeleteTestUser` / `TestListTestUsers` 等，使用 `:memory:` SQLite + `TestMain` 执行 `migrations/002_test_user.sql`（`SetMaxOpenConns(1)`），全部通过。
- `handler` 包：新增 `TestSqliteUser_*` 用例，使用 `httptest.NewRequest` + `echo.New().NewContext` 覆盖创建/查询/更新/删除/列表与失败分支。
- 现有 MySQL 测试与 CRUD 行为不变（SC-007）。
- `go build ./...` 编译通过（提交前必须）。

---

## 验证检查清单(Checklist)

- [ ] SC-001 SQLite 3 秒内连接、与 MySQL 共存
- [ ] SC-002 仅配置 `[database.sqlite]` 即可启用
- [ ] SC-003 CRUD 全链路数据一致
- [ ] SC-004（可选）并发创建 10 条不同 phone 全部成功
- [ ] SC-005 任一 DB 失败 5 秒内退出并明确报错
- [ ] SC-006 双库数据完全隔离
- [ ] SC-007 现有 MySQL 接口行为不变
- [ ] 原则 V/VII：`go test` 使用 SQL 脚本建测试表，应用零自动建表
