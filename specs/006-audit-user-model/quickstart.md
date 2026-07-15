# Quickstart: 006 Audit User Model

## Prerequisites

- Go 1.22+
- MySQL instance with `users` table created (see `migrations/001_user.mysql.sql`)
- Config file `config.toml` with valid MySQL DSN

## Schema Setup

Apply the canonical schema (if not already present):

```bash
mysql -h <host> -u <user> -p <database> < migrations/001_user.mysql.sql
```

## Run Tests

```bash
# All tests require a MySQL connection
go test -v ./model/... -count=1
go test -v ./handler/... -count=1
```

## Build & Run

```bash
go build ./...
./greeting
```

## Validation Scenarios

### 1. Create user with realname and username

```bash
curl -s -X POST http://localhost:1323/demo/usr \
  -H 'Content-Type: application/json' \
  -d '{"phone":"13900000001","realname":"张三","username":"zs","age":25}' | jq .
```

**Expected**: `code: 0`, data includes `realname: "张三"`, `username: "zs"`, no `password_hash`.

### 2. Get user by ID

```bash
curl -s http://localhost:1323/demo/usr/1 | jq .
```

**Expected**: Returns user. No `name` field; has `realname` and `username`.

### 3. Update user (partial)

```bash
curl -s -X PUT http://localhost:1323/demo/usr/1 \
  -H 'Content-Type: application/json' \
  -d '{"realname":"张三丰"}' | jq .
```

**Expected**: `realname` updated to "张三丰", other fields unchanged.

### 4. Soft delete user

```bash
curl -s -X DELETE http://localhost:1323/demo/usr/1 | jq .
```

**Expected**: `code: 0`. Then `GET /demo/usr/1` should return "user not found".

### 5. List excludes deleted users

```bash
curl -s http://localhost:1323/demo/usrs | jq .
```

**Expected**: Array does not include the deleted user.

### 6. Duplicate phone blocked for active users

```bash
curl -s -X POST http://localhost:1323/demo/usr \
  -H 'Content-Type: application/json' \
  -d '{"phone":"13900000001","realname":"duplicate"}' | jq .
```

**Expected**: `"phone already exists"` if the previous user with this phone is still active.

## Key Changes to Verify

- [ ] Response JSON has `realname` and `username`, NOT `name`
- [ ] `password_hash` is NOT present in any API response
- [ ] Soft-deleted users are excluded from `GET /demo/usr/:id` and `GET /demo/usrs`
- [ ] Duplicate phone check only considers non-deleted users
- [ ] `id` field works correctly as `int` type (no overflow issues)
