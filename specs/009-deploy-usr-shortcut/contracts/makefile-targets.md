# Makefile 目标契约(Contracts)

**功能(Feature)**: 009-deploy-usr-shortcut
**日期(Date)**: 2026-07-16
**版本(Version)**: 2.0

> 本文档基于 007-optimize-makefile 的契约文档，记录本次变更涉及的目标契约。仅列出变更的目标；未变更目标（`help`、`rundev`、`fmt`、`lint`、`build`、`build-linux`、`test`、`clean`）的契约参见 [specs/007-optimize-makefile/contracts/makefile-targets.md](../../007-optimize-makefile/contracts/makefile-targets.md)。

---

## 契约(Contract): `deploy-qa`（修改(MODIFIED)）

**调用(Call)**: `make deploy-qa [DEPLOY_HOST=<host>] [DEPLOY_USR=<user>] [DEPLOY_PATH=<path>] [DEPLOY_SUPERVISOR=<name>]`

**默认值(Defaults)**:
| 变量(Variable) | 默认值(Default) |
|---------------|----------------|
| `DEPLOY_HOST` | `111.229.4.203` |
| `DEPLOY_USR` | `ubuntu` |
| `DEPLOY_PATH` | `/opt/project/greeting` |
| `DEPLOY_SUPERVISOR` | `greeting` |

**前置条件(Preconditions)**: 已配置 `$(DEPLOY_USR)@$(DEPLOY_HOST)` 的 ssh/scp 访问权限；所有变量均有默认值，零参数调用即可运行

**行为(Behavior)**:
1. 校验必需变量（`DEPLOY_HOST`、`DEPLOY_SUPERVISOR`）
   - 仅当用户显式清空时触发；正常使用下默认值已填充，不会触发
2. 调用 `build-linux`（隐式依赖）
3. 输出：`Removing old binary on $(DEPLOY_HOST)...`
4. 执行：`ssh $(DEPLOY_USR)@$(DEPLOY_HOST) "rm -f $(DEPLOY_PATH)/$(BIN_NAME)"`
5. 输出：`Uploading to $(DEPLOY_HOST):$(DEPLOY_PATH)/ ...`
6. 执行：`scp -O $(BIN_DIR)/$(BIN_NAME) $(DEPLOY_USR)@$(DEPLOY_HOST):$(DEPLOY_PATH)/`
7. 输出：`Restarting service $(DEPLOY_SUPERVISOR) on $(DEPLOY_HOST)...`
8. 执行：`ssh $(DEPLOY_USR)@$(DEPLOY_HOST) "sudo supervisorctl restart $(DEPLOY_SUPERVISOR)"`
9. 输出：`Deploy complete!`

**变更说明(Changes from 007)**:

| 组件(Component) | 旧值(Before)（007） | 新值(After)（009） |
|-----------------|-------------------|-------------------|
| `DEPLOY_HOST` 默认值 | _(空)_ | `111.229.4.203` |
| `DEPLOY_USR` 变量 | _(不存在)_ | `ubuntu` |
| `DEPLOY_PATH` 默认值 | `/opt/src/main` | `/opt/project/greeting` |
| `DEPLOY_SUPERVISOR` 默认值 | _(空)_ | `greeting` |
| scp 命令 | `root@$(DEPLOY_HOST)` | `$(DEPLOY_USR)@$(DEPLOY_HOST)` |
| ssh 命令 | `root@$(DEPLOY_HOST)` | `$(DEPLOY_USR)@$(DEPLOY_HOST)` |
| 用法提示 | `DEPLOY_HOST=<host> DEPLOY_SUPERVISOR=<name> [DEPLOY_PATH=<path>]` | `[DEPLOY_HOST=<host>] [DEPLOY_USR=<user>] [DEPLOY_PATH=<path>] [DEPLOY_SUPERVISOR=<name>]` |

**并行执行(Parallel Execution)**: 不提供锁机制（同 007）。

**退出码(Exit codes)**:

| 退出码(Code) | 含义(Meaning) |
|-------------|---------------|
| 0 | 部署成功 |
| 1 | 前置检查失败或 scp/ssh 失败 |

---

## 契约(Contract): `runqa`（新增(NEW)）

**调用(Call)**: `make runqa [DEPLOY_HOST=<host>] [DEPLOY_USR=<user>] [DEPLOY_PATH=<path>] [DEPLOY_SUPERVISOR=<name>]`

**前置条件(Preconditions)**: 同 `deploy-qa`（所有变量均有 QA 默认值）

**行为(Behavior)**:
1. 直接调用 `deploy-qa`（无额外逻辑）
2. 所有变量默认值和覆盖行为继承自 `deploy-qa`

**等效于(Equivalent to)**（无覆盖时）:
```bash
make deploy-qa
```

**覆盖示例(Override example)**:
```bash
# 仅覆盖 host，其余保持 QA 默认值
make runqa DEPLOY_HOST=other-qa-server

# 完全自定义部署
make runqa DEPLOY_HOST=prod.example.com DEPLOY_USR=root DEPLOY_PATH=/app DEPLOY_SUPERVISOR=app
```

**退出码(Exit codes)**: 同 `deploy-qa`（继承）。

---

## 变更总结(Summary of Changes)

| 契约(Contract) | 状态(Status) |
|---------------|-------------|
| `deploy-qa` | 修改(MODIFIED) — 变量默认值更新 + 2 处命令（`root@` → `$(DEPLOY_USR)@`）+ 用法提示更新 |
| `runqa` | 新增(NEW) — `deploy-qa` 的别名 |
| 所有其他 | 未变更(UNCHANGED) — 参见 007 契约 |
