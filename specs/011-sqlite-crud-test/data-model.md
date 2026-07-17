# 数据模型(Data Model): 新增 SQLite 实例与 CRUD 测试接口

**功能分支(Feature Branch)**: `011-sqlite-crud-test`
**日期(Date)**: 2026-07-17
**关联规格(Spec)**: [spec.md](./spec.md) | **研究报告**: [research.md](./research.md)

本文件提取规格中的实体、字段、校验规则与状态转换（Phase 1 输出）。

---

## 1. 实体：TestUser（SQLite 测试实体）

仅用于测试的 `test_user` 表映射模型，与 MySQL `users` / `User` 实体完全隔离。字段定义见研究任务 4 决策（GORM 原生软删除）。

### 字段(Field) 定义

| 字段(Field) | 类型(Type) | GORM Tag | JSON | 说明 |
|-------------|-----------|----------|------|------|
| ID | `int` | `primaryKey` | `id` | 自增主键 |
| Name | `string` | `type:varchar(100)` | `name` | 用户姓名（创建必填） |
| Phone | `string` | `type:varchar(20);not null` | `phone` | 手机号（创建必填；活跃记录内唯一，应用层校验） |
| Age | `int` | `default:0` | `age` | 年龄（创建选填，默认 0） |
| CreatedAt | `model.LocalTime` | — | `created_at` | 创建时间，格式 `2006-01-02 15:04:05` |
| UpdatedAt | `model.LocalTime` | — | `updated_at` | 更新时间 |
| DeletedAt | `gorm.DeletedAt` | `index` | `deleted_at` | 软删除标记；非 NULL 表示已删除（GORM 自动过滤） |

> 说明：`CreatedAt` / `UpdatedAt` 复用现有 `model.LocalTime` 类型（统一时间格式，符合原则 II）。`DeletedAt` 使用 GORM 原生 `gorm.DeletedAt`，删除自动置值、查询自动排除。

### 表结构映射(Table Mapping)
- 表名：`test_user`（GORM 默认复数规则，或显式 `TableName() string { return "test_user" }` 以确保精确匹配用户已建表）。
- 该表由用户预先创建（见规格 Clarifications 与 FR-010），应用不创建、不迁移。迁移脚本见 [migrations/002_test_user.sql](#用户自管迁移脚本)。

---

## 2. 实体：SQLite 实例（全局连接）

| 项 | 定义 |
|----|------|
| 全局变量 | `model.SQLiteDB *gorm.DB` |
| 初始化 | `model.InitSQLite(dsn string) error` |
| DSN 来源 | `config.Cfg.Database.SQLite.DSN`（默认 `greeting.db`） |
| 失败策略 | 连接或 Ping 失败返回带类型与地址的错误；`main.go` panic 退出（fail-fast） |
| 与 MySQL 关系 | 与 `model.DB` 并存，各自独立连接池与生命周期 |

---

## 3. 请求实体(Request Entities)

位于 `entity/sqliteusr/`（按模块分子目录，符合原则 I）。Tag 使用 `json`（请求体）。

### TestUserCreateReq — 创建
```go
type TestUserCreateReq struct {
    Name  string `json:"name"`  // required
    Phone string `json:"phone"` // required
    Age   *int   `json:"age"`   // optional, nil means 0
}
```

### TestUserUpdateReq — 部分更新
```go
type TestUserUpdateReq struct {
    Name  *string `json:"name"`  // optional
    Phone *string `json:"phone"` // optional
    Age   *int    `json:"age"`   // optional
}
```

### 路径参数
`/sqlite/testuser/:id` 中的 `id` 通过 `c.Param("id")` 提取（参考现有 `extractUserID`）。

---

## 4. 校验规则(Validation Rules)

| 规则 | 来源 | 处理 |
|------|------|------|
| `name` 必填 | FR-006 | 为空返回 `NotOk(c, "name is required")` |
| `phone` 必填 | FR-006 | 为空返回 `NotOk(c, "phone is required")` |
| `phone` 活跃记录唯一 | Edge Case | 创建前查活跃同号；命中返回 `NotOk(c, "phone already exists")` |
| 部分更新未传字段不变 | FR-007 | 仅对非空指针字段赋值后 `Save` |
| 删除为软删除 | FR-008 | `SQLiteDB.Delete(&TestUser{}, id)`，不物理删除 |
| 列表仅返回未删除 | FR-009 | GORM 原生 `Find` 自动排除 `deleted_at IS NOT NULL` |

---

## 5. 状态转换(State Transitions)

```text
        ┌─────────────┐   Create    ┌─────────────┐
        │  (不存在)    │ ──────────▶ │  Active     │
        └─────────────┘             │ deleted_at  │
                                    │   IS NULL   │
                                    └──────┬──────┘
                                           │ Delete（软删除）
                                           ▼
                                    ┌─────────────┐
                                    │  Deleted    │
                                    │ deleted_at  │
                                    │  NOT NULL   │
                                    └─────────────┘
   Update：在 Active 状态修改 Name/Phone/Age 字段（未传保持不变）
   软删除后同 phone 可再次 Create（活跃唯一性校验不命中已删除记录）
```

---

## 6. 用户自管迁移脚本(migrations/002_test_user.sql)

描述为 `test_user` 建表的 SQL（用户手动对其 `greeting.db` 执行；测试 `TestMain` 对 `:memory:` 执行同一脚本）。字段须与上方实体一致：

```sql
CREATE TABLE IF NOT EXISTS test_user (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        VARCHAR(100) NOT NULL DEFAULT '',
    phone       VARCHAR(20)  NOT NULL DEFAULT '',
    age         INTEGER      NOT NULL DEFAULT 0,
    created_at  DATETIME,
    updated_at  DATETIME,
    deleted_at  DATETIME
);
```

> 说明：应用**不会**运行此脚本；仅作为用户自管资产与测试初始化复用（满足原则 VII / V）。用户已确认 `greeting.db` 中该表已存在。
