# Research: CORS Middleware with Echo v5

## Decision: Use Echo v5 built-in `middleware.CORSWithConfig()`

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
var corsConfig = middleware.CORSConfig{
    AllowOrigins:     []string{"*"},
    AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch, http.MethodHead},
    AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With", "Accept", "Origin"},
    AllowCredentials: false,
    MaxAge:           86400,
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

### Config Storage Decision

CORS configuration is stored as a Go package-level variable in `main.go`. The project now has a `config/` package with TOML support (from 003-config-file), but CORS config remains as Go constants for the following reasons:

- **Compile-time safety**: Invalid CORS config would be caught at build time, not runtime
- **Simplicity**: The CORS feature spec describes configuration via "constants or config file" — constants mode satisfies the requirement
- **No additional dependency**: Does not require changes to `config.toml` or `config.Config` struct
- **Future path**: If needed, CORS config can be migrated to `config.toml` as a `[cors]` section with a straightforward code change

### Middleware Registration Order

Middleware chain in `main.go`:

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
- Use `newCORSApp()` helper to create Echo instance with CORS middleware
- Use `e.ServeHTTP(rec, req)` for integration-level middleware testing
- Test scenarios:
  1. `TestCORSPreflight` — OPTIONS with Origin + ACRM → 200/204 with CORS headers
  2. `TestCORSGETWithOrigin` — GET with Origin → response includes `Access-Control-Allow-Origin`
  3. `TestCORSSameOrigin` — request without Origin → no CORS headers added
  4. `TestCORSOptionsWithoutRequestMethod` — OPTIONS without ACRM → handled gracefully
  5. `TestCORSSpecificOriginAllowed` — matched origin → receives correct CORS header
  6. `TestCORSSpecificOriginDisallowed` — unmatched origin → no CORS headers

## Performance Impact

CORS middleware does minimal work: checks request headers and conditionally sets response headers. Expected overhead: < 1ms (well within the 50ms budget defined in SC-002).
