# Quickstart: 优化 Makefile

**Feature**: 007-optimize-makefile
**Date**: 2026-07-16

> 本文档提供重构后 Makefile 的验证场景，开发者可按顺序执行以确认所有目标正常工作。详细信息参见 [contracts/makefile-targets.md](./contracts/makefile-targets.md)。

---

## Prerequisites

- Go ≥ 1.22（`go version`）
- GNU Make（`make --version`，macOS 自带）
- （可选）gofumpt：`go install mvdan.cc/gofumpt@latest`

---

## Validation Scenarios

### 1. Help Target (FR-007, SC-001)

```bash
# Test: bare make should show help
make

# Test: explicit help target
make help
```

**Expected**:
- All 8 targets displayed grouped by category (Development / Build / Test / Deploy)
- Output appears in < 1 second

---

### 2. Format Code (FR-004, SC-002)

```bash
# Test with gofumpt (if installed)
make fmt

# Test fallback: uninstall gofumpt temporarily then run
make fmt
```

**Expected**:
- With gofumpt: prints `Formatting with gofumpt...`
- Without gofumpt: prints fallback message with install hint
- All `.go` files formatted
- Completes in < 5 seconds

---

### 3. Static Analysis (FR-005)

```bash
make lint
```

**Expected**:
- Passes with exit 0 on clean code
- Reports vet issues with exit 1 if any

---

### 4. Build (FR-002, FR-010, SC-003)

```bash
# Test: current platform build
make build

# Test: verify binary exists
ls -la bin/

# Test: cross-compile override
make build GOOS=darwin GOARCH=arm64
ls -la bin/
```

**Expected**:
- `bin/` directory auto-created if absent
- Binary `bin/greeting` (or project name) generated
- `make build GOOS=linux GOARCH=amd64` produces Linux binary
- Completes in < 30 seconds (without dependency download)

---

### 5. Linux Build (FR-011)

```bash
make build-linux
file bin/greeting
```

**Expected**:
- `file` output shows `ELF 64-bit LSB executable, x86-64, statically linked`
- Binary stripped (`-ldflags="-s -w"`)

---

### 6. Run Tests (FR-003, SC-004)

```bash
make test
```

**Expected**:
- All tests pass, exit code 0
- Verbose output shows individual test results
- Failures clearly identify which test failed

---

### 7. Clean (FR-006)

```bash
make clean
ls bin/   # should fail: no such file or directory
```

**Expected**:
- `bin/` directory removed
- Silent success on re-run (nothing to clean)

---

### 8. Local Development Server (FR-001)

```bash
# Start server (Ctrl+C to stop)
make rundev
```

**Expected**:
- Prints `Starting local dev server...`
- Server starts and listens on configured port
- After Ctrl+C, process exits cleanly

---

### 9. Deploy Guard (FR-008, FR-012, SC-006)

```bash
# Test: missing required variables
make deploy-qa
```

**Expected**:
- Error message listing missing `DEPLOY_HOST` and `DEPLOY_SUPERVISOR`
- Shows usage example
- Exit code 1, no partial execution

```bash
# Test: all variables set (requires actual server, skip in CI)
# make deploy-qa DEPLOY_HOST=10.0.0.1 DEPLOY_SUPERVISOR=myapp DEPLOY_PATH=/opt/myapp
```

---

### 10. All PHONY Targets (FR-009)

```bash
# Verify no target conflicts with filesystem
make help
make rundev  # then Ctrl+C
make fmt
make lint
make build
make build-linux
make test
make clean
```

**Expected**:
- All targets execute their intended commands
- No "make: `xxx` is up to date." messages for `.PHONY` targets
