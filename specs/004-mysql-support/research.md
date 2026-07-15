# Research: MySQL 数据库支持

**Feature**: `004-mysql-support` | **Date**: 2026-07-15

This document resolves all design decisions for the dual-database coexistence architecture.
MySQL and SQLite run independently with separate global instances, models, and endpoints.

---

## R-01 — GORM MySQL driver & DSN format

**Decision**: Add `gorm.io/driver/mysql` as a direct dependency. Use the canonical MySQL DSN format:
`user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local`.

**Rationale**:
- `parseTime=True` required for `time.Time` columns (project uses custom `model.LocalTime` wrapping `time.Time`).
- `charset=utf8mb4` prevents Chinese mojibake (spec Edge Case).
- `gorm.io/driver/mysql` depends on `github.com/go-sql-driver/mysql`, the de-facto MySQL driver for GORM v2.

**Alternatives considered**:
- Raw `database/sql` — rejected: violates Constitution "ORM: GORM v2" constraint.
- `gorm.io/driver/sqlserver` / `postgres` — out of scope.

---

## R-02 — Dual global DB instances

**Decision**: Two `*gorm.DB` global variables:
- `model.DB` → MySQL connection
- `model.SQLiteDB` → SQLite connection

`model.InitDB(mysqlDSN, sqliteDSN string)` opens both connections sequentially. Either failure causes immediate exit.

**Rationale**: Matches spec FR-001 and FR-010. Using `model.DB` as the MySQL instance aligns with the user's instruction that MySQL is the "main" data connection. `model.SQLiteDB` naming indicates it is the secondary/companion instance.

**Alternatives considered**:
- Symmetric naming (`MysqlDB` + `SqliteDB`) — rejected: user explicitly chose Option A.
- Single `map[string]*gorm.DB` — rejected: loses type safety; `model.DB` is referenced extensively.

---

## R-03 — Independent model structs

**Decision**: Two separate GORM model structs:
- `model.User` — MySQL entity (new)
- `model.SQLiteUser` — SQLite entity (existing `User` renamed)

Each has its own CRUD functions using its respective global DB instance:
- `model.CreateUser(user *User) error` uses `model.DB`
- `model.CreateSQLiteUser(user *SQLiteUser) error` uses `model.SQLiteDB`

**Rationale**: Matches spec FR-010. Independent types prevent accidental cross-DB operations. Type safety ensures compile-time enforcement of DB boundaries.

**Alternatives considered**:
- Shared model struct + two DB instances — rejected: spec explicitly says "完全独立" (completely independent). Shared struct could lead to confusion about which DB an instance belongs to.
- Generic model with type parameter — rejected: over-engineering for a simple template project.

---

## R-04 — Cross-dialect schema compatibility

**Decision**: Ship two dialect-specific migration files:
- `migrations/001_user.mysql.sql` — `BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY`
- `migrations/001_user.sqlite.sql` — `INTEGER PRIMARY KEY AUTOINCREMENT`

**Rationale**: MySQL uses `AUTO_INCREMENT`; SQLite uses `AUTOINCREMENT`. A single shared DDL cannot satisfy both dialects (Constitution Principle VII).

**Alternatives considered**:
- One file with dialect-neutral syntax — impossible (id column grammar differs).
- GORM `AutoMigrate` — forbidden by Principle VII.

---

## R-05 — Fail-fast connection error reporting

**Decision**: `model.InitDB` opens both DBs sequentially. If MySQL fails → report MySQL error + exit. If SQLite fails → report SQLite error + exit. Call `sqlDB.Ping()` after each `gorm.Open` for immediate error detection.

**Rationale**: Satisfies Constitution Principle VI and spec FR-009. GORM's `gorm.Open` does not eagerly connect; `Ping()` ensures errors surface within SC-005's 5s budget.

---

## R-06 — MySQL route prefix `/demo/usr`

**Decision**: MySQL CRUD routes use `/demo/usr` prefix:
- `POST /demo/usr` — create
- `GET /demo/usr/:id` — get by ID
- `PUT /demo/usr/:id` — update
- `DELETE /demo/usr/:id` — delete
- `GET /demo/usrs` — list all

SQLite existing routes (`/demo/search`, `/demo/user/phone`, etc.) remain unchanged.

**Rationale**: Matches spec FR-003 and Clarifications. `/demo/usr` is under the same `/demo` group but clearly differentiated from SQLite endpoints.

---

## R-07 — Phone uniqueness scoped to non-deleted records

**Decision**: Enforce phone uniqueness only among active (non-deleted) records. Implementation uses generated columns:

- **MySQL**: `phone_active VARCHAR(32) GENERATED ALWAYS AS (IF(deleted_at IS NULL, phone, NULL)) STORED UNIQUE`
- **SQLite**: `phone_active TEXT GENERATED ALWAYS AS (CASE WHEN deleted_at IS NULL THEN phone END) STORED UNIQUE`

**Rationale**: Spec FR-004 allows soft-deleted phone reuse. A global `UNIQUE(phone)` would block re-creating users with the same phone after soft-delete. Generated column approach is dialect-portable and avoids application-level locking.

**Alternatives considered**:
- Application-level check before insert — risky under concurrency.
- Composite `UNIQUE(phone, deleted_at)` — unreliable (MySQL `NULL != NULL` in unique indexes).
- `RestoreUserByPhone` (existing code) — removed: user chose "新建并存" over "复活".

---

## R-08 — Test schema loading without AutoMigrate

**Decision**: `model.ApplySchema(db *gorm.DB) error` reads `migrations/001_user.sqlite.sql`, splits on `;`, skip comments, executes each statement. `TestMain` in both `model/user_test.go` and `handler/` test files calls `ApplySchema` instead of `DB.AutoMigrate`.

**Rationale**: Satisfies Constitution Principle V and VII. SQL file resolved via `runtime.Caller(0)` from `model/db.go` for cwd-independence.

---

## R-09 — Existing demo handler migration

**Decision**: `handler/demo.go`'s `GetUserByPhoneTest` migrates from `model.DB` + `model.User` to `model.SQLiteDB` + `model.SQLiteUser`. No other demo handler methods are affected (they don't use DB).

**Rationale**: Spec FR-012. The only demo method using `model.User` is `GetUserByPhoneTest`. Other methods (`Search`, `ErrDebug`) operate on request parameters only.

---

## R-10 — `config.toml` dual-DB structure

**Decision**: Replace single `[database]` block with two independent sub-sections:

```toml
[database.mysql]
dsn = "root:password@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"

[database.sqlite]
dsn = "greeting.db"
```

No `type` selector. Both sections are always active.

**Rationale**: Matches spec FR-002. Each DB configuration is self-contained; modifying one does not affect the other. This is simpler than a `type` selector for the coexistence model.
