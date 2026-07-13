# Research: SHA256 Demo API

**Feature**: SHA256 Demo API | **Date**: 2026-07-13

## Decision Log

### R1: SHA256 Implementation

**Decision**: Use Go standard library `crypto/sha256`.

**Rationale**:
- No external dependency needed; aligns with Constitution III (可复制为模板)
- `sha256.Sum256()` is a one-liner computation, no configuration required
- Performance: in-memory CPU-bound, well within 1s target for inputs up to 1MB

**Alternatives considered**:
- Custom SHA256 implementation: unnecessary complexity, no benefit
- Third-party library: violates Constitution III (minimal dependencies)

### R2: Query Parameter Binding

**Decision**: Use Echo's built-in query parameter binding via struct tag `query:"text"`.

**Rationale**:
- Echo v5 natively supports `c.Bind()` with `query` tags for GET requests
- Follows existing project pattern (entity struct with query/param/json tags)
- No manual `c.QueryParam()` parsing needed

**Alternatives considered**:
- Manual `c.QueryParam("text")`: works but less idiomatic, violates entity layer pattern
- Path parameter: spec requires query string, not path

### R3: Error Handling for Missing Parameter

**Decision**: Return `response.NotOk` with HTTP 400 and descriptive message when `text` is missing.

**Rationale**:
- FR-003 explicitly requires error response with readable message
- Aligns with Constitution II (Unified Response Format)
- Echo's `c.Bind()` returns error on missing required fields, which we wrap

### R4: No Database Interaction

**Decision**: This is a pure computation endpoint with no state.

**Rationale**:
- Feature spec does not require persistence
- Simplifies implementation (no entity/model mapping)
- No GORM AutoMigrate needed for this feature
- Handler test does not need database setup

## Standard Library Usage Confirmed

| Package | Purpose | Import Path |
|---------|---------|-------------|
| `crypto/sha256` | Hash computation | `crypto/sha256` |
| `fmt` | Hex encoding of hash bytes | `fmt` |
