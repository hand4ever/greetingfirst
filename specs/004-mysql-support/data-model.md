# Data Model: MySQL 数据库支持

**Feature**: `004-mysql-support` | **Date**: 2026-07-15

Two independent database entities: `User` (MySQL) and `SQLiteUser` (SQLite, renamed from existing `User`).

## Entities

### User (MySQL — `model.User`)

| Field | Go Type | MySQL Column | Constraint | Notes |
|-------|---------|-------------|------------|-------|
| `ID` | `uint` | `BIGINT UNSIGNED AUTO_INCREMENT` | PRIMARY KEY | Auto-generated |
| `Phone` | `string` | `VARCHAR(32)` | `NOT NULL` | Unique among active records via generated column `phone_active` (see R-07); reusable after soft-delete |
| `Name` | `string` | `VARCHAR(64)` | `NOT NULL` | — |
| `Age` | `int` | `INT` | `DEFAULT 0` | Optional |
| `CreatedAt` | `LocalTime` | `DATETIME` | — | Auto-set by GORM |
| `UpdatedAt` | `LocalTime` | `DATETIME` | — | Auto-updated by GORM |
| `DeletedAt` | `gorm.DeletedAt` | `DATETIME NULL` | INDEX | Soft delete marker |
| *(implicit)* | — | `VARCHAR(32)` | `UNIQUE` | Generated column `phone_active = IF(deleted_at IS NULL, phone, NULL) STORED` — not mapped to Go struct |

Table name: GORM default `users` (pluralized).

### SQLiteUser (SQLite — `model.SQLiteUser`)

| Field | Go Type | SQLite Column | Constraint | Notes |
|-------|---------|-------------|------------|-------|
| `ID` | `uint` | `INTEGER PRIMARY KEY AUTOINCREMENT` | PRIMARY KEY | Auto-generated |
| `Phone` | `string` | `VARCHAR(32)` | `NOT NULL` | Unique among active records via generated column; reusable after soft-delete |
| `Name` | `string` | `VARCHAR(64)` | `NOT NULL` | — |
| `Age` | `int` | `INT` | `DEFAULT 0` | Optional |
| `CreatedAt` | `LocalTime` | `DATETIME` | — | Auto-set by GORM |
| `UpdatedAt` | `LocalTime` | `DATETIME` | — | Auto-updated by GORM |
| `DeletedAt` | `gorm.DeletedAt` | `DATETIME NULL` | INDEX | Soft delete marker |
| *(implicit)* | — | `TEXT` | `UNIQUE` | Generated column `phone_active = CASE WHEN deleted_at IS NULL THEN phone END STORED` — not mapped to Go struct |

Table name: GORM default `sqlite_users` (pluralized).

**Key difference from existing**: `SQLiteUser` is a rename of the existing `model.User`. All existing fields and behaviors are preserved. The `RestoreUserByPhone` function is removed (soft-deleted phone reuse is now handled by the generated column at the DB level).

## Request / Response Entities

### UserCreateReq — `POST /demo/usr` (MySQL only)
```go
// entity/user/user.go
type UserCreateReq struct {
    Name  string `json:"name"`  // required
    Phone string `json:"phone"` // required, unique among non-deleted MySQL records
    Age   int    `json:"age"`   // optional, default 0
}
```

### UserPathReq — path param `:id`
```go
// entity/user/user.go
type UserPathReq struct {
    ID uint `param:"id" json:"-"`
}
```

### UserUpdateReq — `PUT /demo/usr/:id` (MySQL only)
```go
// entity/user/user.go
type UserUpdateReq struct {
    Name string `json:"name"`  // optional, partial update
    Age  int    `json:"age"`   // optional, partial update
    // Phone is NOT updatable
}
```

### Response Data

- MySQL endpoints: `data` = `model.User` (serialized via `response.Ok`)
- SQLite endpoints: `data` = `model.SQLiteUser` (existing behavior unchanged)

## State Transitions

### MySQL (User)

| Operation | Endpoint | Behavior |
|-----------|----------|----------|
| Create | `POST /demo/usr` | INSERT into `users`; `DeletedAt = NULL`. If soft-deleted record with same phone exists, new record is allowed (old record's `phone_active = NULL`, new record's `phone_active = phone`). |
| Get | `GET /demo/usr/:id` | SELECT where `id = :id AND deleted_at IS NULL`. Returns 404 if not found. |
| List | `GET /demo/usrs` | SELECT where `deleted_at IS NULL`, ordered by `created_at DESC`. Returns `[]` if empty. |
| Update | `PUT /demo/usr/:id` | UPDATE only provided fields; `UpdatedAt` refreshed. `Phone` not updatable. Returns 404 if not found. |
| Delete | `DELETE /demo/usr/:id` | Soft delete: SET `deleted_at = NOW()`. `phone_active` becomes NULL, releasing phone for reuse. Returns 404 if not found. |

### SQLite (SQLiteUser)

Existing behavior preserved. `GetUserByPhoneTest` continues to use `model.SQLiteDB` + `model.SQLiteUser` with phone `"13636311005"`.

## Validation Rules

| Rule | Applies To | Error Response |
|------|-----------|----------------|
| `Name` and `Phone` must be non-empty on create | MySQL `POST /demo/usr` | `response.NotOk("name and phone are required")` |
| `Phone` must be unique among active records | MySQL `POST /demo/usr` | `response.NotOk("phone already exists")` |
| Record must exist (not soft-deleted) | MySQL GET/PUT/DELETE | `response.NotOk("user not found")` |
