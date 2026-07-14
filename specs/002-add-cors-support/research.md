# Research: CORS Middleware with Echo v5

## Decision: Use Echo v5 built-in `middleware.CORS()`

**Rationale**: Echo v5 provides a production-ready CORS middleware in `github.com/labstack/echo/v5/middleware`. It handles all standard CORS behaviors: OPTIONS preflight, `Access-Control-Allow-Origin`, `Access-Control-Allow-Methods`, `Access-Control-Allow-Headers`, `Access-Control-Max-Age`, and `Access-Control-Allow-Credentials`. Using the built-in middleware means zero new dependencies, which aligns with Constitution III (Copy-Ready Template).

**Alternatives considered**:

| Alternative | Rejected Because |
|-------------|-----------------|
| Custom CORS middleware from scratch | Reinventing the wheel; Echo v5 built-in is well-tested and maintained |
| `rs/cors` third-party library | Adds a new dependency with no benefit over the built-in solution |
| `echo-cors` community middleware | Unnecessary when framework provides the feature natively |

## Echo v5 CORS Middleware Configuration

The middleware accepts `middleware.CORSConfig` with these relevant fields:

```go
middleware.CORSConfig{
    AllowOrigins:     []string{"*"},                          // FR-003: configurable origins
    AllowMethods:     []string{http.MethodGet, ...},          // FR-004: configurable methods
    AllowHeaders:     []string{"Content-Type", ...},          // FR-005: configurable headers
    AllowCredentials: false,                                  // FR-006: credentials off by default
    MaxAge:           86400,                                   // preflight cache duration
}
```

### Default Configuration Decision

| Field | Default Value | Reason |
|-------|---------------|--------|
| `AllowOrigins` | `["*"]` | Development-friendly; allows all origins per spec assumption |
| `AllowMethods` | `[GET, POST, PUT, DELETE, OPTIONS, PATCH, HEAD]` | Covers RESTful and common HTTP methods |
| `AllowHeaders` | `[Content-Type, Authorization, X-Requested-With, Accept, Origin]` | Common request headers for APIs |
| `AllowCredentials` | `false` | Per clarification Q1: no credential mode by default |
| `MaxAge` | `86400` (24 hours) | Standard preflight cache duration to reduce OPTIONS requests |

### Middleware Registration Order

Current middleware chain in `main.go`:
```
1. RequestLogger
2. Recover
3. RequestID
4. CostTime
```

CORS middleware should be inserted after `Recover` (panic recovery should come before CORS) and before `RequestID` (so CORS headers are added early and present even for preflight):

```
1. RequestLogger
2. Recover
3. CORS          ← NEW (after Recover, before business middleware)
4. RequestID
5. CostTime
```

**Rationale**: CORS evaluation is cheap and should happen early. Placing it after `Recover` ensures panic recovery covers all middleware. Placing it before `RequestID` means CORS headers are present even if downstream middleware short-circuits.

## Test Strategy

- Test file: `middle/cors_test.go`
- Use `echo.New()` + `httptest.NewRequest` to simulate cross-origin requests
- Test scenarios:
  1. OPTIONS preflight with valid Origin → 200 with correct CORS headers
  2. GET with Origin → response includes `Access-Control-Allow-Origin`
  3. Same-origin request (no Origin) → no CORS headers added
  4. OPTIONS without `Access-Control-Request-Method` → 200 without CORS headers
  5. Disallowed origin → no `Access-Control-Allow-Origin` header

## Performance Impact

CORS middleware does minimal work: checks request headers and conditionally sets response headers. Expected overhead: < 1ms (well within the 50ms budget defined in SC-002).
