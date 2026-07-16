# Implementation Plan: 优化 Makefile

**Branch**: `007-optimize-makefile` | **Date**: 2026-07-16 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/007-optimize-makefile/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command; its definition describes the execution workflow.

## Summary

重构项目根目录的 `Makefile`，将现有的 2 个目标（`rundev`、`buildqa`）扩展为完整的构建工具链。新增 `build`、`test`、`fmt`、`lint`、`clean`、`help` 等标准目标，将 `buildqa` 中硬编码的 IP、路径、supervisor 进程名抽取为可配置变量，并修复 `-ldflags` 悬空等多个技术问题。目标是把 Makefile 从简单的启动/部署脚本升级为覆盖开发、构建、测试、部署全流程的完整构建配置。

## Technical Context

**Language/Version**: GNU Make (macOS 自带 / Linux apt install make)

**Primary Dependencies**: Go ≥ 1.22、gofumpt（可选，不可用时回退 go fmt）、scp、ssh

**Storage**: N/A（纯构建工具配置，不涉及数据存储）

**Testing**: `go test -v ./... -count=1`（通过 `make test` 目标调用）

**Target Platform**: 开发者本地 macOS / Linux；部署目标 Linux amd64

**Project Type**: build tool configuration（Makefile）

**Performance Goals**: `make help` < 1s、`make fmt` < 5s、`make build` < 30s（不含依赖下载）

**Constraints**: GNU Make 语法、POSIX shell 兼容（macOS 内置 zsh 和 bash）、须保留 `rundev` 目标名以兼容现有开发者习惯

**Scale/Scope**: 约 8-10 个 `.PHONY` 目标，6-8 个可配置变量，单一 Makefile 文件（< 100 行）

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. 分层架构 | N/A | Makefile 是构建工具配置，不属于应用分层代码 |
| II. 统一响应格式 | N/A | Makefile 不涉及 API 响应 |
| III. 可复制为模板 | ✅ PASS | 变量驱动设计，无硬编码 IP/路径，新项目可直接复用 |
| IV. 英文代码产物 | ✅ PASS | Makefile 注释使用英文 |
| V. 测试覆盖 | ✅ PASS | `make test` 目标正确执行 `go test -v ./... -count=1` |
| VI. 错误及时抛出 | ✅ PASS | 部署目标缺配置时立即终止；构建/测试失败传播非零退出码 |
| VII. 数据库表由用户创建 | N/A | Makefile 不涉及数据库操作 |

**Gate Result**: All applicable principles pass. No violations. Proceed to Phase 0.

## Project Structure

### Documentation (this feature)

```text
specs/007-optimize-makefile/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
Makefile            # UPDATED: 重构后的 Makefile（唯一变更文件）
```

**Structure Decision**: 本次优化只涉及单个文件 `Makefile`，无新增目录或代码文件。所有目标均在根目录 Makefile 中集中定义。

## Constitution Check (Post-Design Re-evaluation)

*Re-check after Phase 1 design artifacts (data-model.md, contracts/, quickstart.md) completed.*

| Principle | Status | Verification |
|-----------|--------|-------------|
| I. 分层架构 | N/A | 不变 |
| II. 统一响应格式 | N/A | 不变 |
| III. 可复制为模板 | ✅ PASS | 8 个变量通过 `?=` 定义默认值，部署变量用空占位符 + fail-fast guard；无硬编码项目特定值 |
| IV. 英文代码产物 | ✅ PASS | contracts/ 中 exit code 注释均为英文；help 输出的目标描述可使用中文（spec 要求） |
| V. 测试覆盖 | ✅ PASS | `make test` 目标映射到 `go test -v ./... -count=1`，与 Constitution 要求一致 |
| VI. 错误及时抛出 | ✅ PASS | `deploy-qa` 使用 `check_required` 宏在第一步做必填校验，缺配置立即 exit 1；`build`/`test`/`lint` 均通过 `go` 命令自动传播非零退出码 |
| VII. 数据库表由用户创建 | N/A | 不变 |

**Post-Design Gate Result**: All applicable principles still pass. No new violations introduced by design decisions.

## Complexity Tracking

> No violations. No entries needed.
