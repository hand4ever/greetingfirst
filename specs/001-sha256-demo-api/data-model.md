# Data Model: SHA256 Demo API

**Feature**: SHA256 Demo API | **Date**: 2026-07-13

## Overview

This feature has no persistent data model. It is a stateless computation endpoint that takes a query parameter and returns a computed result.

## Request Entity

### Sha256Request

Query parameter binding struct for the `/demo/sha256` endpoint.

| Field | Type | Tag | Required | Description |
|-------|------|-----|----------|-------------|
| Text | string | `query:"text"` | Yes | Input text to compute SHA256 hash for |

**Validation rules**:
- `text` parameter is required; missing it returns HTTP 400

**Go struct**:
```go
type Sha256Request struct {
    Text string `query:"text"`
}
```

## Response Entity

### Sha256Response

Data payload embedded in the unified `response.ErrMsg.Data` field.

| Field | Type | JSON Key | Description |
|-------|------|----------|-------------|
| Input | string | `input` | Original input text (echoed back for verification) |
| Hash  | string | `hash`  | SHA256 hash in lowercase hexadecimal format |

**Go struct**:
```go
type Sha256Response struct {
    Input string `json:"input"`
    Hash  string `json:"hash"`
}
```

## State Transitions

N/A — this endpoint is stateless. No state machine or lifecycle applies.

## Relationships

No relationships to other entities. This is a standalone computation endpoint.
