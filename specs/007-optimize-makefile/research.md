# Research: 优化 Makefile

**Feature**: 007-optimize-makefile
**Date**: 2026-07-16

## R1: Go 项目 Makefile 目标命名与最佳实践

**Decision**: 采用 Go 社区标准目标命名，按功能类别分组

**Rationale**:
- Go 社区普遍使用 `build`、`test`、`fmt`、`lint`、`clean` 作为标准目标名
- 分组展示（开发/构建/测试/部署）提升可读性，符合 User Story 4 的新人上手需求
- 保留现有 `rundev` 目标名，避免破坏开发者的肌肉记忆

**Standard targets mapped to spec requirements**:

| Target | Spec Reference | Purpose |
|--------|---------------|---------|
| `help` (default) | FR-007 | Show all available commands |
| `rundev` | FR-001 | Start local dev server |
| `build` | FR-002 | Compile for current platform |
| `test` | FR-003 | Run all unit tests |
| `fmt` | FR-004 | Format code (gofumpt with go fmt fallback) |
| `lint` | FR-005 | Static analysis via go vet |
| `clean` | FR-006 | Remove build artifacts |
| `deploy-qa` | FR-008 | Build + upload + restart on QA server |

**Alternatives considered**:
- `run` instead of `rundev`: rejected, existing developers already use `rundev`
- `format` instead of `fmt`: rejected, `fmt` is the Go ecosystem convention
- All-in-one `check` target: could be added later as a CI convenience target combining `fmt` + `lint` + `test`

---

## R2: 变量配置与覆盖模式

**Decision**: 使用 `?=` 赋值操作符设置默认值，允许通过环境变量或命令行覆盖

**Rationale**:
- `?=` 仅在变量未定义时赋值，优先级：命令行 `make VAR=val` > 环境变量 `VAR=val` > `?=` 默认值
- 部署相关变量使用空占位符作为默认值，配合必填检查（参见 R5）
- 构建相关变量（GOOS、GOARCH）通过 `$(shell go env ...)` 获取宿主平台默认值

**Variable definitions pattern**:

```makefile
# Build variables (FR-002)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BIN_DIR ?= bin
BIN_NAME ?= $(shell basename $(CURDIR))

# Deploy variables (FR-008)
DEPLOY_HOST ?=
DEPLOY_PATH ?= /opt/src/main
DEPLOY_SUPERVISOR ?=
```

**Alternatives considered**:
- `=` (recursive expansion): rejected, may cause unexpected re-evaluation with `$(shell ...)`
- `:=` (immediate expansion): rejected for overridable variables, prevents env var override
- Config file (`.env` / `makefile.inc`): considered but rejected as unnecessary complexity for < 10 variables

---

## R3: 交叉编译模式

**Decision**: `build` 目标默认编译当前平台，`build-linux` 作为单独的便利目标用于部署前的 Linux 编译；`deploy-qa` 依赖 `build-linux`

**Rationale**:
- `make build` 输出当前平台二进制（开发者日常使用）
- `make build-linux` 始终输出 Linux amd64 二进制（部署前使用）
- `deploy-qa` 作为部署的编排目标，先 `build-linux` 再 scp + restart
- 符合 spec FR-002（支持交叉编译）和 FR-011（部署使用 CGO_ENABLED=0 + -ldflags="-s -w"）

**Cross-compilation command**:
```makefile
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BIN_NAME) main.go
```

**Alternatives considered**:
- 单一 `build` 目标接收参数：虽然更简洁，但增加了使用复杂度（`make build GOOS=linux`），不如显式目标直观
- `build-all` 目标编译所有平台：项目当前仅需 Linux 部署，不需要多平台矩阵

---

## R4: gofumpt 工具检测与回退

**Decision**: 使用 `command -v gofumpt` 检测可用性，不可用时回退到 `go fmt`

**Rationale**:
- `command -v` 是 POSIX 兼容的二进制检测方式，macOS 和 Linux 均可工作
- 项目 Constitution 推荐 `gofumpt`（见 GO_STYLE.mdc），但不强制要求
- Spec FR-004 明确要求回退机制
- 检测失败时应给出安装指引（`go install mvdan.cc/gofumpt@latest`）

**Implementation pattern**:
```makefile
FMT_CMD := $(shell command -v gofumpt 2>/dev/null)
ifeq ($(FMT_CMD),)
  FMT_CMD := gofmt
  FMT_NOTE := "(gofumpt not found, using go fmt. Install: go install mvdan.cc/gofumpt@latest)"
endif
```

**Alternatives considered**:
- 将 gofumpt 设为必需依赖：rejected，增加环境搭建负担，违背 Principle III（可复制性）
- `go run mvdan.cc/gofumpt@latest`：可避免手动安装，但每次执行都需要下载，速度慢且不离线友好

---

## R5: 部署目标的 Fail-Fast 设计

**Decision**: 部署目标在执行任何步骤前检查必填变量，缺失则打印错误信息并 exit 1

**Rationale**:
- 符合 Constitution Principle VI（错误及时抛出）
- 半完成部署（编译了但未上传或未重启）比不部署更危险，难以排查
- 错误信息应包含变量名、用途和示例，帮助开发者快速修复

**Guard pattern**:
```makefile
define check_required
  @if [ -z "$($1)" ]; then \
    echo "Error: $(2) is not set."; \
    echo "  Usage: make deploy-qa DEPLOY_HOST=10.0.0.1 DEPLOY_SUPERVISOR=myapp"; \
    exit 1; \
  fi
endef
```

**Required variables for deploy-qa**:
- `DEPLOY_HOST`: 目标服务器 IP 或域名（无默认值，必填）
- `DEPLOY_SUPERVISOR`: supervisor 进程名（无默认值，必填）
- `DEPLOY_PATH`: 部署路径（默认 `/opt/src/main`，可覆盖）

**Alternatives considered**:
- 使用默认值静默执行：rejected，IP 无法有合理默认值
- 交互式提示输入：rejected，破坏自动化流程（CI/CD 不可用）
- 锁文件机制防止并行部署冲突：rejected in clarification session，部署为低频操作，团队内部沟通即可避免，不引入不必要的复杂度

---

## R6: Help 目标设计

**Decision**: `help` 设为 `.DEFAULT_GOAL`，使用 `@grep -E` 提取注释并按类别分组展示；帮助文本使用中英双语格式（`中文 (English)`）

**Rationale**:
- 设置 `.DEFAULT_GOAL = help` 使无参数 `make` 等同于 `make help`
- 使用连贯的注释格式 `## category: 中文描述 (English description)` 实现自文档化
- 不需要手动维护目标列表，添加新目标时只需在注释中遵循格式
- 中英双语格式：帮助文本属于文档范畴（Constitution IV 允许中文），同时保留英文关键词方便开发者对照

**Self-documenting pattern with bilingual comments**:
```makefile
.DEFAULT_GOAL := help

help: ## help: 显示帮助 (Show available commands)
	@echo "Usage: make [target]"
	@echo ""
	@echo "开发 (Development):"
	@grep -E '^[a-zA-Z_-]+:.*?## dev:.*$$' $(firstword $(MAKEFILE_LIST)) | \
		awk -F':.*?## dev: ' '{printf "  make %-15s %s\n", $$1, $$2}' | sort
	@echo ""
	@echo "构建 (Build):"
	@grep -E '^[a-zA-Z_-]+:.*?## build:.*$$' $(firstword $(MAKEFILE_LIST)) | \
		awk -F':.*?## build: ' '{printf "  make %-15s %s\n", $$1, $$2}' | sort
	# ... other categories
```

**Alternatives considered**:
- 手写 help 文本：简单直接，但需双维护（写目标 + 写 help），容易不同步
- 外部工具（如 `makehelp`）：引入额外依赖，违背简化原则
- 纯中文 / 纯英文 help：rejected in clarification session，中英双语兼顾可读性和国际化

---

## Summary

所有技术决策已收敛，无 NEEDS CLARIFICATION 残留：
- 8 个标准目标覆盖开发/构建/测试/部署全流程
- 变量通过 `?=` + 空占位符 + fail-fast guard 实现灵活配置
- gofumpt 通过 `command -v` 检测，自动回退 go fmt
- help 通过 grep 注释实现自文档化，无需双维护
- 所有变更仅涉及 `Makefile` 一个文件
