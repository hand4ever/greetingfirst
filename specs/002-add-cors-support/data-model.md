# Data Model: CORS 跨域支持

**Feature**: 002-add-cors-support

## Summary

This feature is a middleware-only enhancement and does not introduce any new data entities, database tables, or persistent state.

## Configuration Constants

The following configuration values are defined as Go constants/variables for the CORS middleware:

| Name | Type | Description | Default |
|------|------|-------------|---------|
| `corsConfig` | `middleware.CORSConfig` | CORS middleware configuration struct | See defaults below |

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
