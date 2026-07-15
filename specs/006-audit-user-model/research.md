# Research: 006 Audit User Model

## Decision 1: Remove GORM Soft Delete, Use Manual Soft Delete

**Decision**: Replace `gorm.DeletedAt` with `*time.Time` and implement manual soft delete logic.

**Rationale**:
- The actual DB column `deleted_at` is a plain `datetime` without GORM's soft-delete semantics.
- Manual soft delete gives explicit control over query filtering and delete behavior.
- `gorm.DeletedAt` introduces implicit `WHERE deleted_at IS NULL` filtering on all queries, which may surprise developers unfamiliar with GORM conventions.

**Alternatives Considered**:
- Keep `gorm.DeletedAt`: rejected because it adds magic behavior not aligned with DB schema expectations.
- Hard delete (`DB.Unscoped().Delete`): rejected per spec clarification — manual soft delete preferred.

**Implementation Pattern**:
```go
// Delete: update deleted_at instead of physical delete
func DeleteUser(id int) error {
    now := time.Now()
    return DB.Model(&User{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// Query: explicitly filter deleted records
func GetUserByID(id int) (*User, error) {
    var user User
    err := DB.Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
    // ...
}
```

## Decision 2: ID Type — `uint` → `int`

**Decision**: Change `User.ID` from `uint` to `int` and propagate through all CRUD signatures.

**Rationale**:
- DB column `id` is `int auto_increment`. Matching DB type avoids implicit conversions and edge cases with negative values from raw queries.
- `int` is the default GORM mapping for `int` columns. Using `uint` introduces unnecessary mismatch.
- Affected signatures: `CreateUser`, `GetUserByID`, `GetUserByPhone`, `UpdateUser`, `DeleteUser`, `extractUserID`.

**Alternatives Considered**:
- Keep `uint`: rejected per spec clarification — full alignment with DB schema.

**Impact**:
- `extractUserID` in `handler/user.go`: `strconv.ParseUint` → `strconv.Atoi`.
- Test files: all `uint` ID references → `int`.

## Decision 3: List Uses Manual WHERE Clause

**Decision**: The `List` handler already uses `DB.Where("deleted_at IS NULL")` — this pattern will be extended to `GetUserByID`, `GetUserByPhone`, and any future query functions.

**Rationale**:
- Without `gorm.DeletedAt`, GORM no longer auto-filters deleted records.
- Consistent manual `WHERE deleted_at IS NULL` pattern across all query functions.
- `List` is the reference implementation.

## Decision 4: No New Dependencies

**Decision**: All changes are type/field renames and GORM usage pattern changes. No new Go modules required.

**Rationale**:
- `*time.Time` is from Go standard library.
- Existing GORM v2 API supports all needed operations.
- Aligned with Constitution Principle III (Copy-Ready Template).
