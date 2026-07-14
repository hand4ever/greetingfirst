# Data Model: CORS 跨域支持

**Feature**: 002-add-cors-support

## Summary

This feature is a middleware-only enhancement and does not introduce any new data entities, database tables, or persistent state.

## Configuration

CORS configuration is defined as a package-level Go variable in `main.go`:

```go
var corsConfig = middleware.CORSConfig{...}
```

### CORSConfig Fields

| Field | Type | Default | Notes |
|-------|------|---------|-------|
| `AllowOrigins` | `[]string` | `["*"]` | Allowed origin patterns; `"*"` means any origin |
| `AllowMethods` | `[]string` | `[GET, POST, PUT, DELETE, OPTIONS, PATCH, HEAD]` | Allowed HTTP methods |
| `AllowHeaders` | `[]string` | `[Content-Type, Authorization, X-Requested-With, Accept, Origin]` | Allowed request headers |
| `AllowCredentials` | `bool` | `false` | Whether to allow credentials (cookies, auth headers) |
| `MaxAge` | `int` | `86400` | Preflight cache duration in seconds (24h) |

## State Transitions

No state involved — CORS is configured at startup and remains static for the lifetime of the server.

## Relationship to config.toml

The project has a `config.toml` file and `config/` package (from 003-config-file) handling application metadata, server port, database DSN, and changelog entries. CORS configuration is intentionally kept separate as Go constants in `main.go` rather than in `config.toml`, for compile-time safety and simplicity. This can be migrated to `config.toml` in the future if dynamic CORS configuration is needed.
