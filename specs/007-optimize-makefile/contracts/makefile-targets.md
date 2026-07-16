# Makefile Target Contracts

**Feature**: 007-optimize-makefile
**Date**: 2026-07-16
**Version**: 1.0

> 本文档定义 Makefile 各目标的"接口契约"——调用方式、前置条件、预期行为和输出格式。这些契约是开发者与 Makefile 之间的约定，也是验证实现的依据。

---

## Contract: `help` (default)

**Call**: `make` or `make help`

**Preconditions**: GNU Make ≥ 3.81

**Behavior**:
1. Print header: `Usage: make [target]`
2. Print targets grouped by category using `## <category>: <description>` comment format
3. Categories displayed in order: Development → Build → Test → Deploy

**Exit codes**:

| Code | Meaning |
|------|---------|
| 0 | Displayed help successfully |

**Example output**:
```
Usage: make [target]

Development:
  make fmt           Format Go source code
  make lint          Run static analysis (go vet)
  make rundev        Start local development server

Build:
  make build         Compile for current platform
  make build-linux   Cross-compile for Linux amd64
  make clean         Remove build artifacts

Test:
  make test          Run all unit tests

Deploy:
  make deploy-qa     Build, upload, restart QA server
```

---

## Contract: `rundev`

**Call**: `make rundev`

**Preconditions**: Go ≥ 1.22, port available (from `config.toml`)

**Behavior**:
1. Print: `Starting local dev server...`
2. Execute: `go run main.go`
3. On exit (Ctrl+C): service stops naturally

**Exit codes**:

| Code | Meaning |
|------|---------|
| 0 | Server exited cleanly |
| 1 | Go not found or port already in use |

---

## Contract: `fmt`

**Call**: `make fmt`

**Preconditions**: none (Go required, gofumpt optional)

**Behavior**:
1. Detect if gofumpt is available via `command -v gofumpt`
2. If found: run `gofumpt -l -w .` and print `Formatting with gofumpt...`
3. If not found: run `gofmt -l -w .` and print `gofumpt not found, using go fmt. Install: go install mvdan.cc/gofumpt@latest`
4. If any files were changed, list them

**Exit codes**:

| Code | Meaning |
|------|---------|
| 0 | Formatted successfully (or no changes needed) |
| 1 | Format command failed |

---

## Contract: `lint`

**Call**: `make lint`

**Preconditions**: Go ≥ 1.22

**Behavior**:
1. Execute: `go vet ./...`
2. Print any vet warnings/errors

**Exit codes**:

| Code | Meaning |
|------|---------|
| 0 | No vet issues found |
| 1 | Vet found issues |

---

## Contract: `build`

**Call**: `make build [GOOS=<os>] [GOARCH=<arch>]`

**Preconditions**: Go ≥ 1.22, `main.go` exists

**Behavior**:
1. Create `$(BIN_DIR)` if it doesn't exist
2. Print: `Building for $(GOOS)/$(GOARCH)...`
3. Execute: `go build -o $(BIN_DIR)/$(BIN_NAME) main.go`
4. Print: `Build complete: $(BIN_DIR)/$(BIN_NAME)`

**Exit codes**:

| Code | Meaning |
|------|---------|
| 0 | Build successful |
| 1 | Build failed (compile error) |

---

## Contract: `build-linux`

**Call**: `make build-linux`

**Preconditions**: Go ≥ 1.22

**Behavior**:
1. Create `$(BIN_DIR)` if it doesn't exist
2. Print: `Building for linux/amd64...`
3. Execute: `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BIN_NAME) main.go`
4. Print: `Build complete: $(BIN_DIR)/$(BIN_NAME)`

**Exit codes**: Same as `build`

---

## Contract: `test`

**Call**: `make test`

**Preconditions**: Go ≥ 1.22

**Behavior**:
1. Print: `Running tests...`
2. Execute: `go test -v ./... -count=1`
3. Pass: print `All tests passed.`
4. Fail: print failing test names and exit 1

**Exit codes**:

| Code | Meaning |
|------|---------|
| 0 | All tests passed |
| 1 | One or more tests failed |

---

## Contract: `clean`

**Call**: `make clean`

**Preconditions**: none

**Behavior**:
1. Print: `Cleaning build artifacts...`
2. Remove `$(BIN_DIR)` directory
3. Print: `Done.`

**Exit codes**:

| Code | Meaning |
|------|---------|
| 0 | Cleaned successfully |

---

## Contract: `deploy-qa`

**Call**: `make deploy-qa DEPLOY_HOST=<host> DEPLOY_SUPERVISOR=<name> [DEPLOY_PATH=<path>]`

**Preconditions**: `DEPLOY_HOST` and `DEPLOY_SUPERVISOR` must be set; ssh/scp access configured

**Behavior**:
1. Validate required variables (`DEPLOY_HOST`, `DEPLOY_SUPERVISOR`)
   - Missing → print error with usage example, exit 1
2. Invoke `build-linux` (implicit dependency)
3. Print: `Uploading to $(DEPLOY_HOST):$(DEPLOY_PATH)/ ...`
4. Execute: `scp -O $(BIN_DIR)/$(BIN_NAME) root@$(DEPLOY_HOST):$(DEPLOY_PATH)/`
5. Print: `Restarting service $(DEPLOY_SUPERVISOR) on $(DEPLOY_HOST)...`
6. Execute: `ssh root@$(DEPLOY_HOST) "supervisorctl restart $(DEPLOY_SUPERVISOR)"`
7. Print: `Deploy complete!`

**Exit codes**:

| Code | Meaning |
|------|---------|
| 0 | Deploy successful |
| 1 | Pre-check failed (missing required variables) |

**Variable validation error output**:
```
Error: DEPLOY_HOST is not set.
  Usage: make deploy-qa DEPLOY_HOST=10.0.0.1 DEPLOY_SUPERVISOR=myapp
```
