# API Contracts: MySQL User CRUD

**Feature**: `004-mysql-support` | **Date**: 2026-07-15

MySQL CRUD endpoints under `/demo/usr`. SQLite existing endpoints under `/demo` remain unchanged (not documented here).

All responses use the unified `response.ErrMsg` envelope:
```json
{ "code": 0, "message": "", "data": ..., "trace_id": "...", "cost": "..." }
```
- `code == 0` → success; `code != 0` → error (`response.ErrCodeCustom = 100001`).
- All timestamps serialized as `2006-01-02 15:04:05` via `model.LocalTime`.

Base path: `/demo/usr` (list: `/demo/usrs`)

---

## 1. Create User — `POST /demo/usr`

**Request body** (`UserCreateReq`):
```json
{ "name": "张三", "phone": "13800138000", "age": 25 }
```
- `name` (required), `phone` (required, unique among non-deleted MySQL records), `age` (optional, default 0).

**Success** `200`:
```json
{
  "code": 0, "message": "",
  "data": { "id": 1, "phone": "13800138000", "name": "张三", "age": 25,
            "created_at": "2026-07-15 10:00:00", "updated_at": "2026-07-15 10:00:00" }
}
```

**Errors**:
- missing `name`/`phone` → `code 100001`, `message: "name and phone are required"`
- duplicate `phone` (active record) → `code 100001`, `message: "phone already exists"`

---

## 2. Get User — `GET /demo/usr/:id`

**Path**: `id` = user ID (uint).

**Success** `200`: `data` = user object (same shape as create).

**Error**: not found → `code 100001`, `message: "user not found"`.

---

## 3. Update User — `PUT /demo/usr/:id`

**Request body** (`UserUpdateReq`, partial):
```json
{ "name": "张三丰", "age": 30 }
```
- Only provided fields are updated; `phone` is not updatable; omitted fields keep original value.

**Success** `200`: `data` = updated user object.

**Error**: not found → `code 100001`, `message: "user not found"`.

---

## 4. Delete User — `DELETE /demo/usr/:id`

**Path**: `id` = user ID.

**Success** `200`: `data` = `""` (soft delete; record remains with `deleted_at` set).

**Error**: not found → `code 100001`, `message: "user not found"`.

---

## 5. List Users — `GET /demo/usrs`

**Success** `200`: `data` = array of user objects, ordered by `created_at DESC`, excluding soft-deleted records. Empty list `[]` when no users.

---

## Error Contract

| Condition | `code` | `message` |
|-----------|--------|-----------|
| Missing required field | 100001 | `name and phone are required` |
| Duplicate phone (active record) | 100001 | `phone already exists` |
| Record not found | 100001 | `user not found` |
| MySQL connection failure (startup) | — | panics with type+address, non-zero exit |
| SQLite connection failure (startup) | — | panics with type+address, non-zero exit |
