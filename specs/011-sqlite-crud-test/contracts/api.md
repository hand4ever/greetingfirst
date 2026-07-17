# 接口契约(Contracts): 新增 SQLite 实例与 CRUD 测试接口

**功能分支(Feature Branch)**: `011-sqlite-crud-test`
**日期(Date)**: 2026-07-17
**关联规格(Spec)**: [spec.md](./spec.md) | **数据模型**: [data-model.md](../data-model.md)

本项目为 Web 服务（Echo v5），对外暴露 HTTP JSON 接口。契约覆盖新增的 SQLite CRUD 测试接口。所有响应统一使用 `response.ErrMsg` 结构（见原则 II）：

```json
{ "code": 0, "message": "", "data": ..., "trace_id": "...", "cost": "..." }
```

- 成功：`code == 0`（`ErrCodeOk`），`data` 为业务数据。
- 失败：`code != 0`（`ErrCodeCustom` = 100001），`message` 为明确错误信息。

路由前缀 `/sqlite/testuser` 独立于现有 MySQL 的 `/demo/usr`（FR-005）。

---

## 1. 创建用户 — `POST /sqlite/testuser`

**请求体(Request Body)**：
```json
{
  "name": "张三",        // required
  "phone": "13900138000", // required
  "age": 28              // optional, 缺省 0
}
```

**成功响应(200, code=0)**：
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 1,
    "name": "张三",
    "phone": "13900138000",
    "age": 28,
    "created_at": "2026-07-17 10:00:00",
    "updated_at": "2026-07-17 10:00:00",
    "deleted_at": null
  }
}
```

**失败响应**：
- `name` 为空 → `{"code":100001,"message":"name is required"}`
- `phone` 为空 → `{"code":100001,"message":"phone is required"}`
- 活跃同号已存在 → `{"code":100001,"message":"phone already exists"}`

---

## 2. 按 ID 查询 — `GET /sqlite/testuser/:id`

**路径参数(Path)**：`id` (int)

**成功响应(200, code=0)**：`data` 为单个 TestUser（字段同上）。

**失败响应**：
- `id` 非整数 → `{"code":100001,"message":"invalid path parameter"}`
- 不存在或已软删除 → `{"code":100001,"message":"user not found"}`

---

## 3. 更新用户 — `PUT /sqlite/testuser/:id`

**路径参数(Path)**：`id` (int)
**请求体(Request Body)**（仅传需要修改的字段）：
```json
{ "name": "张三丰", "age": 30 }
```

**成功响应(200, code=0)**：`data` 为更新后的 TestUser，未传字段保持原值。

**失败响应**：
- `id` 非整数 → `{"code":100001,"message":"invalid path parameter"}`
- 不存在或已软删除 → `{"code":100001,"message":"user not found"}`
- 更新失败 → `{"code":100001,"message":"update user failed: <err>"}`

---

## 4. 删除用户（软删除）— `DELETE /sqlite/testuser/:id`

**路径参数(Path)**：`id` (int)

**成功响应(200, code=0)**：
```json
{ "code": 0, "message": "", "data": "" }
```

**失败响应**：
- `id` 非整数 → `{"code":100001,"message":"invalid path parameter"}`
- 删除失败 → `{"code":100001,"message":"delete user failed: <err>"}`

> 注：删除为软删除，记录 `deleted_at` 被置值，后续查询/列表不再返回。

---

## 5. 用户列表 — `GET /sqlite/testusers`

**查询参数(Query)**：无（可按需扩展分页/排序，本期默认按 `created_at DESC`）。

**成功响应(200, code=0)**：`data` 为未软删除的 TestUser 数组（可能为空数组 `[]`）。

**失败响应**：
- 查询失败 → `{"code":100001,"message":"query users failed: <err>"}`

---

## 6. 错误码约定(Error Codes)

| code | 常量 | 含义 |
|------|------|------|
| 0 | `ErrCodeOk` | 成功 |
| 100001 | `ErrCodeCustom` | 业务/参数/资源错误（本特性统一使用） |

---

## 7. 与现有接口的隔离性

| 接口组 | 前缀 | 数据库 |
|--------|------|--------|
| MySQL CRUD | `/demo/usr` | `model.DB`（MySQL `users`） |
| SQLite CRUD（本特性） | `/sqlite/testuser` | `model.SQLiteDB`（SQLite `test_user`） |

两者路由不冲突、数据完全隔离（SC-006）。
