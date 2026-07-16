# 数据模型(Data Model): 部署用户变量化与 QA 快捷目标

**功能(Feature)**: 009-deploy-usr-shortcut
**日期(Date)**: 2026-07-16

> 本次变更不涉及数据库实体，但 Makefile 中的目标(Target)和变量(Variable)构成了其"数据模型"。以下记录变更后的完整目录，变更项以 **[新增(NEW)]** 或 **[修改(MODIFIED)]** 标注。

## Makefile 目标目录(Target Catalog)

### 分类(Category): 帮助(Help)（默认）

| 目标(Target) | `.PHONY` | 依赖(Dependencies) | 规格引用(Spec Ref) | 变更(Change) |
|-------------|----------|-------------------|-------------------|-------------|
| `help` | 是 | 无 | FR-007 | 自动拾取新增的 `runqa` 目标（无代码变更） |

### 分类(Category): 部署(Deploy)

| 目标(Target) | `.PHONY` | 依赖(Dependencies) | 规格引用(Spec Ref) | 变更(Change) |
|-------------|----------|-------------------|-------------------|-------------|
| `deploy-qa` | 是 | `build-linux` | FR-002, FR-003 | **[修改(MODIFIED)]** `scp`/`ssh` 命令中 `root@` → `$(DEPLOY_USR)@`，`check_required` 用法提示更新 |
| `runqa` | 是 | `deploy-qa` | FR-004, FR-005, FR-006 | **[新增(NEW)]** `deploy-qa` 的别名，零参数即可部署到 QA |

- **`deploy-qa`**（修改(MODIFIED)）：将 `DEPLOY_HOST` 和 `DEPLOY_USR` 绑定为 `$(DEPLOY_USR)@$(DEPLOY_HOST)` 格式用于 scp/ssh；所有变量现通过 `?=` 拥有 QA 默认值
- **`runqa`**（新增(NEW)）：`deploy-qa` 的别名；所有默认值已在变量层级设置

## 变量目录(Variable Catalog)

### 部署变量（FR-001）

| 变量(Variable) | 默认值(Default) | 可覆盖(Overridable) | 使用者(Required by) | 变更(Change) |
|---------------|----------------|--------------------|-------------------|-------------|
| `DEPLOY_HOST` | `111.229.4.203` | 环境变量 / 命令行 | `deploy-qa` | **[修改(MODIFIED)]**（原为空） |
| `DEPLOY_USR` | `ubuntu` | 环境变量 / 命令行 | `deploy-qa` | **[新增(NEW)]** |
| `DEPLOY_PATH` | `/opt/project/greeting` | 环境变量 / 命令行 | `deploy-qa` | **[修改(MODIFIED)]**（原为 `/opt/src/main`） |
| `DEPLOY_SUPERVISOR` | `greeting` | 环境变量 / 命令行 | `deploy-qa` | **[修改(MODIFIED)]**（原为空） |

> 所有变量使用 `?=` 赋值：未覆盖时应用 QA 默认值。命令行参数优先级最高。

## 状态流转(State Transitions)

无状态机。Makefile 是纯函数式构建工具，每个目标(target)是独立的命令序列。

`deploy-qa` 执行流程(Execution flow)：
```
check_required → build-linux → rm old binary → scp upload → ssh restart
```

`runqa` 执行流程(Execution flow)：
```
deploy-qa（同上流程，无额外逻辑）
```
