# Feature Specification: 优化 Makefile

**Feature Branch**: `007-optimize-makefile`

**Created**: 2026-07-16

**Status**: Draft

**Input**: User description: "本地运行用make rundev，这个在makefile里，本地任务目标就是优化makefile文件"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 开发者本地快速启动服务 (Priority: P1)

作为一名开发者，我希望在本地执行一个简单的 `make` 命令就能启动开发服务，不需要记忆复杂的 `go run` 参数，且能获得清晰的启动反馈。

**Why this priority**: 这是开发者最频繁使用的日常操作，影响开发效率的核心入口。

**Independent Test**: 在项目根目录执行 `make rundev`，服务能成功启动并监听端口，终端显示清晰的启动信息。

**Acceptance Scenarios**:

1. **Given** 项目代码已就绪、Go 环境已配置, **When** 开发者执行 `make rundev`, **Then** 服务启动并在终端显示监听地址和端口
2. **Given** 服务已在运行, **When** 开发者再次执行 `make rundev`, **Then** 终端提示端口已被占用，给出明确错误信息
3. **Given** Go 未安装或版本不兼容, **When** 开发者执行 `make rundev`, **Then** 终端给出明确的 Go 环境缺失提示

---

### User Story 2 - 开发者本地编译、格式化与测试 (Priority: P1)

作为一名开发者，我希望通过 `make` 命令一键完成代码格式化检查、编译验证和单元测试，确保提交前的代码质量。

**Why this priority**: 本地开发闭环的关键环节，与启动服务同等重要，确保代码质量门禁。

**Independent Test**: 执行 `make build` 完成编译、`make test` 运行全部测试、`make fmt` 检查代码格式。

**Acceptance Scenarios**:

1. **Given** 项目代码已修改, **When** 开发者执行 `make fmt`, **Then** 自动格式化所有 Go 代码（通过 `gofumpt`），并报告有无格式问题
2. **Given** 项目代码已修改, **When** 开发者执行 `make build`, **Then** 编译整个项目，编译通过无错误，编译失败则显示具体错误信息
3. **Given** 项目代码已修改, **When** 开发者执行 `make test`, **Then** 运行所有单元测试并显示测试结果摘要（通过/失败数量）
4. **Given** 开发者不确定 Makefile 有哪些可用命令, **When** 执行 `make help` 或仅执行 `make`, **Then** 显示所有可用目标的列表及简要说明，说明使用中英双语格式（如 `编译项目 (Build binary)`）

---

### User Story 3 - 运维人员部署到测试服务器 (Priority: P2)

作为一名运维或后端开发者，我希望通过 Makefile 中的部署目标一键完成编译、上传和远程重启，且 IP 地址、路径等可变参数能够方便地配置，而非硬编码在 Makefile 中。

**Why this priority**: 部署是上线前的必要步骤，但频率低于本地开发，配置灵活性能避免多环境维护 Makefile 副本。

**Independent Test**: 设置部署相关变量后执行 `make deploy-qa`，编译产出上传到指定服务器并重启服务。

**Acceptance Scenarios**:

1. **Given** 部署目标服务器 IP 和路径已通过变量配置, **When** 执行 `make deploy-qa`, **Then** 编译 Linux 二进制文件、上传到目标服务器、重启服务，全程展示进度提示
2. **Given** 部署变量未配置（如 IP 为空）, **When** 执行 `make deploy-qa`, **Then** 明确提示缺少必要配置并给出示例格式，不执行部署
3. **Given** 服务器连接失败, **When** 执行 `make deploy-qa`, **Then** 在失败步骤终止并显示连接错误信息，不继续执行后续步骤

---

### User Story 4 - 项目新人快速上手 (Priority: P3)

作为一名新加入项目的开发者，我希望通过执行 `make help` 或仅 `make` 就能了解项目的所有可用操作，无需翻阅文档或阅读 Makefile 源码。

**Why this priority**: 降低新人上手门槛，但属于一次性需求，优先级低于日常高频操作。

**Independent Test**: 执行 `make`（不带参数）或 `make help`，终端打印出所有目标及其用途说明。

**Acceptance Scenarios**:

1. **Given** 项目刚克隆到本地, **When** 开发者执行 `make` 或 `make help`, **Then** 展示所有可用目标列表，每个目标附有中英双语的用途说明（如 `build        编译项目 (Build binary)`）
2. **Given** Makefile 已包含所有目标, **When** 开发者执行 `make help`, **Then** 帮助信息按类别分组（如：开发、构建、测试、部署）

---

### Edge Cases

- 当项目不在 Git 仓库根目录执行 make 时，路径相关目标如何处理？
- 当 `go.mod` 中定义的模块名与实际不一致时，`make build` 是否能给出清晰错误？
- 当目标目录（如 `bin/`）不存在时，构建目标是否能自动创建？
- 当 `gofumpt` 未安装时，`make fmt` 是否给出安装指引？
- 当多开发者并行执行 `make deploy-qa` 时，是否可能出现冲突？如何处理？→ **已澄清：接受风险，不处理。部署为低频操作，团队内部沟通即可避免并行冲突。**

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Makefile MUST 提供 `rundev` 目标，使用 `go run main.go` 启动本地开发服务，并在启动前后打印明确的状态提示
- **FR-002**: Makefile MUST 提供 `build` 目标，支持交叉编译（通过变量指定 GOOS/GOARCH），默认编译当前平台
- **FR-003**: Makefile MUST 提供 `test` 目标，执行 `go test -v ./... -count=1` 并展示测试结果
- **FR-004**: Makefile MUST 提供 `fmt` 目标，使用 `gofumpt` 格式化代码，若 `gofumpt` 不可用则回退到 `go fmt`
- **FR-005**: Makefile MUST 提供 `lint` 目标，使用 `go vet ./...` 进行静态分析检查
- **FR-006**: Makefile MUST 提供 `clean` 目标，清理编译产物（`bin/` 目录）
- **FR-007**: Makefile MUST 提供 `help` 目标作为默认目标（`.DEFAULT_GOAL`），展示所有可用目标及其用途说明，按类别分组；帮助文本使用中英双语格式（如 `build        编译项目 (Build binary)`）
- **FR-008**: 部署目标（`deploy-qa`）MUST 从变量读取服务器 IP、端口、部署路径、supervisor 进程名，而非硬编码；变量缺失时 MUST 明确提示并终止
- **FR-009**: 所有对外目标 MUST 声明 `.PHONY`，避免与同名文件冲突
- **FR-010**: 构建目标 MUST 自动创建 `bin/` 输出目录（若不存在）
- **FR-011**: 部署目标的编译步骤 MUST 使用 `CGO_ENABLED=0` 和 `-ldflags="-s -w"` 生成精简静态链接二进制，减小文件体积
- **FR-012**: Makefile 中的可变参数（部署 IP、路径等）MUST 支持通过环境变量或 make 变量覆盖，并提供合理默认值（若无法确定目标值，使用占位符并提示配置）

### Key Entities

本次优化不涉及数据库实体，仅涉及 Makefile 文件本身。

- **Make Target**: Makefile 中的一个构建目标，包含名称、依赖、执行命令、用途描述
- **Build Variable**: Makefile 中可覆盖的变量，用于控制构建行为（如 GOOS、部署服务器 IP 等）

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 开发者执行 `make help` 能在 1 秒内看到完整的命令列表和说明
- **SC-002**: `make fmt` 在 5 秒内完成全项目代码格式化检查
- **SC-003**: `make build` 在 30 秒内完成全项目编译（不含依赖下载）
- **SC-004**: `make test` 执行全部测试用例，失败时清晰标识失败用例名称
- **SC-005**: 新开发者无需阅读 Makefile 源码，仅通过 `make help` 即可了解所有可用操作
- **SC-006**: 部署目标在缺失必要配置时立即终止（不执行部分步骤），避免产生半完成部署状态

## Clarifications

### Session 2026-07-16

- **Q: 并行部署冲突如何处理？** → **A: 接受风险，不处理。** 部署为低频操作，团队内部沟通即可避免并行冲突，无需引入锁机制增加复杂度。
- **Q: 帮助文本（help targets）用中文还是英文？** → **A: 中英双语。** 帮助文本属于文档范畴，可使用中文；同时保留英文关键词方便开发者对照。格式：`build        编译项目 (Build binary)`。

## Assumptions

- 开发者本地已安装 Go ≥ 1.22，符合项目 Constitution 中的技术栈约束
- `gofumpt` 虽非必需依赖，但 Makefile 会优先尝试使用，不可用则回退 `go fmt`
- 部署目标服务器为 Linux amd64 架构
- Makefile 不负责 Go 依赖管理（`go mod download` 由开发者手动或 CI 执行），但可在 `build` 目标中可选地自动下载缺失依赖
- 项目 Makefile 使用 GNU Make 语法，开发者环境需安装 `make` 工具（macOS 自带、Linux 需安装）
