# Data Model: 部署用户变量化与 QA 快捷目标

**Feature**: 009-deploy-usr-shortcut
**Date**: 2026-07-16

> 本次变更不涉及数据库实体，但 Makefile 中的 Targets 和 Variables 构成了其"数据模型"。以下记录变更后的完整目录，变更项以 **[NEW]** 或 **[MODIFIED]** 标注。

## Makefile Target Catalog

### Category: Help (default)

| Target | `.PHONY` | Dependencies | Spec Ref | Change |
|--------|----------|-------------|----------|--------|
| `help` | yes | none | FR-007 | 自动拾取新增的 `runqa` 目标（无代码变更） |

### Category: Deploy

| Target | `.PHONY` | Dependencies | Spec Ref | Change |
|--------|----------|-------------|----------|--------|
| `deploy-qa` | yes | `build-linux` | FR-002, FR-003 | **[MODIFIED]** `scp`/`ssh` 命令中 `root@` → `$(DEPLOY_USR)@`，`check_required` 用法提示更新 |
| `runqa` | yes | `deploy-qa` | FR-004, FR-005, FR-006 | **[NEW]** `deploy-qa` 的别名，零参数即可部署到 QA |

- **`deploy-qa`** (MODIFIED): Bind `DEPLOY_HOST` and `DEPLOY_USR` → `$(DEPLOY_USR)@$(DEPLOY_HOST)` in scp/ssh; all variables now have QA defaults via `?=`
- **`runqa`** (NEW): Alias for `deploy-qa`; all defaults already set at variable level

## Variable Catalog

### Deploy Variables (FR-001)

| Variable | Default | Overridable | Required by | Change |
|----------|---------|-------------|-------------|--------|
| `DEPLOY_HOST` | `111.229.4.203` | env / cli | `deploy-qa` | **[MODIFIED]** (was empty) |
| `DEPLOY_USR` | `ubuntu` | env / cli | `deploy-qa` | **[NEW]** |
| `DEPLOY_PATH` | `/opt/project/greeting` | env / cli | `deploy-qa` | **[MODIFIED]** (was `/opt/src/main`) |
| `DEPLOY_SUPERVISOR` | `greeting` | env / cli | `deploy-qa` | **[MODIFIED]** (was empty) |

> All variables use `?=` assignment: QA defaults apply when not overridden. Command-line arguments take highest priority.

## State Transitions

无状态机。Makefile 是纯函数式构建工具，每个 target 是独立的命令序列。

`deploy-qa` 执行流程：
```
check_required → build-linux → scp upload → ssh restart
```

`runqa` 执行流程：
```
deploy-qa (same flow as above, no extra logic)
```
