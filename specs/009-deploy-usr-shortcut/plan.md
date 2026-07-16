# 实施计划(Implementation Plan): 部署用户变量化与 QA 快捷目标

**分支(Branch)**: `009-deploy-usr-shortcut` | **日期(Date)**: 2026-07-16 | **规格(Spec)**: [spec.md](./spec.md)

**输入(Input)**: Feature specification from `/specs/009-deploy-usr-shortcut/spec.md`

**说明(Note)**: This template is filled in by the `/speckit.plan` command; its definition describes the execution workflow.

## 概述(Summary)

在现有 Makefile 基础上进行两项增量修改：(1) 新增 `DEPLOY_USR` 变量替代硬编码的 `root` 用户，并将所有部署变量默认值设为 QA 环境值；(2) 新增 `runqa` 快捷目标作为 `deploy-qa` 的别名，开发者只需执行 `make runqa`（或 `make deploy-qa`）即可一键部署。

技术方案：所有部署变量通过 `?=` 设置 QA 默认值，命令行参数或环境变量可覆盖；`runqa` 直接依赖 `deploy-qa`，无需目标特定变量。

## 技术上下文(Technical Context)

**语言/版本(Language/Version)**: Go 1.22+ (服务端); GNU Make 3.81+ (构建工具)

**主要依赖(Primary Dependencies)**: `scp`、`ssh`、`supervisorctl`（均在目标服务器端，非本地依赖）

**存储(Storage)**: N/A（不涉及数据库或持久化存储）

**测试框架(Testing)**: 手动验证（Makefile 修改不涉及 Go 代码变动，无需 Go 单元测试更新）

**目标平台(Target Platform)**: Linux amd64 服务器（部署目标）；macOS / Linux（开发者本地执行 Make）

**项目类型(Project Type)**: 构建系统（Makefile）配置变更

**性能目标(Performance Goals)**: N/A（部署为低频操作，无性能约束）

**约束(Constraints)**:
- 参数可覆盖：所有默认值允许命令行覆盖（`make deploy-qa DEPLOY_HOST=other`）
- 仅修改 `Makefile` 单文件，不引入新依赖或脚本
- `deploy-qa` 与 `runqa` 功能等价，均基于 QA 默认值运行

**规模/范围(Scale/Scope)**: 单文件修改（`Makefile`），修改 ~4 行变量默认值 + 新增 `DEPLOY_USR` 变量 + 新增 ~3 行 `runqa` 目标

## 宪法检查(Constitution Check)

*门禁(GATE): Must pass before Phase 0 research. Re-check after Phase 1 design.*

| 原则 | 状态 | 说明 |
|------|------|------|
| I. 分层架构 | ✅ N/A | 仅修改 Makefile，不涉及 Go 分层代码 |
| II. 统一响应格式 | ✅ N/A | 不涉及 API 变更 |
| III. 可复制为模板 | ✅ PASS | Makefile 变量默认值可通过 `?=` 轻松修改；复制到新项目时只需调整默认值 |
| IV. 英文代码产物 | ✅ PASS | Makefile 变量名和目标名使用英文；help 文本沿用现有中英双语格式 `## deploy:`；commit message 使用英文 |
| V. 测试覆盖 | ✅ PASS | Makefile 变更不涉及 Go 代码，无需新增单元测试；quickstart.md 提供手动验证场景 |
| VI. 错误及时抛出 | ✅ PASS | `deploy-qa` 通过 `check_required` 宏保持 fail-fast：变量缺失时打印错误并 exit 1；SSH 连接失败时 `scp`/`ssh` 命令自然返回非零退出码 |
| VII. 数据库表由用户创建 | ✅ N/A | 不涉及数据库 |

**门禁结论**: 全部通过，无违规项。

## 项目结构(Project Structure)

### 文档(Documentation)（本特性）

```text
specs/009-deploy-usr-shortcut/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   └── makefile-targets.md
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### 源代码(Source Code)（仓库根目录）

```text
Makefile    # 唯一变更文件：修改变量默认值为 QA 环境值 + 新增 DEPLOY_USR 变量 + 新增 runqa 目标（deploy-qa 别名）+ 修改 SSH 用户引用
```

**结构决策(Structure Decision)**: 单文件变更，不新增源文件或目录。Makefile 是构建系统的唯一入口，所有构建和部署逻辑集中于此。

## 复杂度追踪(Complexity Tracking)

> 无违规项，无需填写。
