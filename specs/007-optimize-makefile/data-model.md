# Data Model: 优化 Makefile

**Feature**: 007-optimize-makefile
**Date**: 2026-07-16

> 本次优化不涉及数据库实体（如 spec.md 所述），但 Makefile 中的 Targets 和 Variables 构成了其自身的"数据模型"，以下结构化描述确保实现一致性。

## Makefile Target Catalog

### Category: Help (default)

| Target | `.PHONY` | Dependencies | Spec Ref |
|--------|----------|-------------|----------|
| `help` | yes | none | FR-007 |

- **Description**: Display all available make targets grouped by category
- **Behavior**: Parses `Makefile` comments to extract and display target descriptions
- **Default**: Set as `.DEFAULT_GOAL`, so bare `make` shows help

### Category: Development

| Target | `.PHONY` | Dependencies | Spec Ref |
|--------|----------|-------------|----------|
| `rundev` | yes | none | FR-001 |
| `fmt` | yes | none | FR-004 |
| `lint` | yes | none | FR-005 |

- **`rundev`**: Runs `go run main.go`, prints start/stop status messages
- **`fmt`**: Formats all Go source files using `gofumpt` (fallback: `go fmt`), prints detection result
- **`lint`**: Runs `go vet ./...` for static analysis

### Category: Build

| Target | `.PHONY` | Dependencies | Spec Ref |
|--------|----------|-------------|----------|
| `build` | yes | none | FR-002 |
| `build-linux` | yes | none | FR-002, FR-011 |
| `clean` | yes | none | FR-006, FR-010 |

- **`build`**: Compiles binary for current platform (`$(GOOS)`/`$(GOARCH)`), outputs to `$(BIN_DIR)/$(BIN_NAME)`
- **`build-linux`**: Cross-compiles for Linux amd64 with `CGO_ENABLED=0 -ldflags="-s -w"`, outputs to `$(BIN_DIR)/$(BIN_NAME)`
- **`clean`**: Removes `$(BIN_DIR)` directory

### Category: Test

| Target | `.PHONY` | Dependencies | Spec Ref |
|--------|----------|-------------|----------|
| `test` | yes | none | FR-003 |

- **`test`**: Runs `go test -v ./... -count=1`, exit code reflects test results

### Category: Deploy

| Target | `.PHONY` | Dependencies | Spec Ref |
|--------|----------|-------------|----------|
| `deploy-qa` | yes | `build-linux` | FR-008, FR-009 |

- **`deploy-qa`**: Pre-check required variables → build-linux → scp upload → ssh restart
- **Guard**: Fails immediately if `DEPLOY_HOST` or `DEPLOY_SUPERVISOR` is empty

## Variable Catalog

### Build Variables (FR-002, FR-010)

| Variable | Default | Overridable | Description |
|----------|---------|-------------|-------------|
| `GOOS` | `$(shell go env GOOS)` | env / cli | Target OS for `build` |
| `GOARCH` | `$(shell go env GOARCH)` | env / cli | Target architecture for `build` |
| `BIN_DIR` | `bin` | env / cli | Binary output directory |
| `BIN_NAME` | `$(shell basename $(CURDIR))` | env / cli | Binary file name |

### Deploy Variables (FR-008, FR-012)

| Variable | Default | Overridable | Required by |
|----------|---------|-------------|-------------|
| `DEPLOY_HOST` | _(empty)_ | env / cli | `deploy-qa` |
| `DEPLOY_PATH` | `/opt/src/main` | env / cli | `deploy-qa` |
| `DEPLOY_SUPERVISOR` | _(empty)_ | env / cli | `deploy-qa` |

### Format Variable (FR-004)

| Variable | Default | Overridable | Description |
|----------|---------|-------------|-------------|
| `FMT_CMD` | `gofumpt` or `gofmt` (auto-detect) | no | Formatter command |

## State Transitions

无状态机。Makefile 是纯函数式构建工具，每个 target 是独立的命令序列，不维护持久化状态。

唯一的状态感知来自文件系统（`$(BIN_DIR)` 是否存在）和工具链（`gofumpt` 是否安装），均在执行时动态检测。
