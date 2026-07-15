# API Contracts: User CRUD

## POST /demo/usr — Create User

**Request**:
```json
{
  "phone": "13800138000",
  "realname": "张三",
  "username": "zhangsan",
  "age": 25
}
```

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `phone` | string | **Yes** | Unique among non-deleted users |
| `realname` | string | No | Max 100 chars |
| `username` | string | No | Max 20 chars |
| `age` | int | No | Pointer in Go, omit if unset |

**Response** (200):
```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "id": 1,
    "phone": "13800138000",
    "realname": "张三",
    "username": "zhangsan",
    "age": 25,
    "created_at": "2026-07-15 10:30:00",
    "updated_at": "2026-07-15 10:30:00"
  }
}
```

**Errors**:
- `phone is required` — phone missing
- `phone already exists` — duplicate phone among active users

---

## GET /demo/usr/:id — Get User

**Response** (200): Same as Create response data shape.

**Errors**:
- `invalid path parameter` — non-numeric id
- `user not found` — id doesn't exist or user soft-deleted

---

## PUT /demo/usr/:id — Update User

**Request** (partial update):
```json
{
  "realname": "张三丰",
  "age": 30
}
```

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `realname` | string | No | Pointer, nil = skip |
| `username` | string | No | Pointer, nil = skip |
| `age` | int | No | Pointer, nil = skip |

**Response** (200): Updated user object.

---

## DELETE /demo/usr/:id — Delete User (Soft)

Soft-deletes by setting `deleted_at = NOW()`. Deleted users are excluded from all queries.

**Response** (200):
```json
{
  "code": 0,
  "message": "ok",
  "data": ""
}
```

---

## GET /demo/usrs — List Users

Returns all non-deleted users ordered by `created_at DESC`.

**Response** (200):
```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": 1,
      "phone": "13800138000",
      "realname": "张三",
      "username": "zhangsan",
      "age": 25,
      "created_at": "2026-07-15 10:30:00",
      "updated_at": "2026-07-15 10:30:00"
    }
  ]
}
```

## Breaking Changes

| Field | Before | After | Impact |
|-------|--------|-------|--------|
| Request: `name` | string, required | — | **Removed** |
| Request: `realname` | — | string, optional | **New** |
| Request: `username` | — | string, optional | **New** |
| Response: `name` | present | — | **Removed** |
| Response: `realname` | — | present | **New** |
| Response: `username` | — | present | **New** |
| Response: `id` type | uint | int | No visible change in JSON |
