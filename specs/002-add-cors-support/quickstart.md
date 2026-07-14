# Quickstart: CORS 跨域支持验证

**Feature**: 002-add-cors-support

## Prerequisites

- Go ≥ 1.22 installed
- Project dependencies: `go mod tidy`
- Server running on `http://localhost:1323` (start with `go run .`)

## Validation Scenarios

### Scenario 1: OPTIONS Preflight Request

Verify that preflight requests return correct CORS headers.

```bash
curl -s -i -X OPTIONS \
  -H "Origin: https://example.com" \
  -H "Access-Control-Request-Method: POST" \
  http://localhost:1323/demo/search
```

**Expected**: HTTP 200 (or 204). Response headers include:
- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods` containing `POST`

---

### Scenario 2: Simple GET with Origin

Verify that cross-origin GET requests include `Access-Control-Allow-Origin`.

```bash
curl -s -i \
  -H "Origin: https://example.com" \
  http://localhost:1323/demo/sha256?text=hello
```

**Expected**: HTTP 200. Response headers include `Access-Control-Allow-Origin: *`.

---

### Scenario 3: Same-Origin Request (No Origin Header)

Verify that non-cross-origin requests do NOT have CORS headers added.

```bash
curl -s -i http://localhost:1323/demo/sha256?text=hello
```

**Expected**: HTTP 200. Response headers do NOT include `Access-Control-Allow-Origin`.

---

### Scenario 4: OPTIONS Without Access-Control-Request-Method

Verify that OPTIONS without the proper preflight header is handled gracefully.

```bash
curl -s -i -X OPTIONS \
  -H "Origin: https://example.com" \
  http://localhost:1323/demo/sha256
```

**Expected**: HTTP 200 (or 204). Echo CORS middleware still adds CORS headers as a safe default.

---

### Scenario 5: Non-Simple Request (POST with Content-Type)

Verify full preflight + request flow for non-simple requests.

```bash
# Preflight
curl -s -i -X OPTIONS \
  -H "Origin: https://example.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  http://localhost:1323/demo/sha256

# Actual request
curl -s -i -X POST \
  -H "Origin: https://example.com" \
  -H "Content-Type: application/json" \
  http://localhost:1323/demo/sha256
```

**Expected**: Preflight returns 200 with correct CORS headers (including `Access-Control-Allow-Headers: Content-Type`). Actual request returns normal response.

---

### Scenario 6: Disallowed Origin (when AllowOrigins is specific)

Modify `AllowOrigins` in `main.go` to `["https://myapp.com"]`, restart server, then:

```bash
curl -s -i \
  -H "Origin: https://evil.com" \
  http://localhost:1323/demo/sha256?text=hello
```

**Expected**: Response does NOT include `Access-Control-Allow-Origin` header (browser would block the response).

---

### Scenario 7: Run Unit Tests

```bash
go test -v ./middle/... -count=1 -run TestCORS
```

**Expected**: All 6 CORS test functions pass with `PASS` output.
