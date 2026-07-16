# 技术研究(Research): 部署用户变量化与 QA 快捷目标

**功能(Feature)**: 009-deploy-usr-shortcut
**日期(Date)**: 2026-07-16

> 本次变更在当前 Makefile 基础上：(1) 将所有部署变量的默认值设为 QA 环境值；(2) 新增 `DEPLOY_USR` 变量替代硬编码 `root`；(3) 新增 `runqa` 作为 `deploy-qa` 的别名。

## R1: 部署变量默认值策略

**决策(Decision)**: 所有部署变量通过 `?=` 设置 QA 环境默认值；`deploy-qa` **无需传参即可运行**；`DEPLOY_USR` 新增，默认值 `ubuntu`

**理由(Rationale)**:
- `?=` 仅在变量未定义时赋值，命令行参数 `make DEPLOY_HOST=X` 优先级最高
- QA 环境参数（`111.229.4.203`、`ubuntu`、`greeting`、`/opt/project/greeting`）是项目主要部署目标，设为默认值后 `deploy-qa` 零参数即可运行
- `DEPLOY_USR ?= ubuntu` 代替硬编码 `root`，同时提供更安全的默认值
- 所有默认值统一管理在变量定义区，易于发现和修改

**实现(Implementation)**:
```makefile
# Deploy variables (QA defaults)
DEPLOY_HOST ?= 111.229.4.203
DEPLOY_PATH ?= /opt/project/greeting
DEPLOY_USR ?= ubuntu
DEPLOY_SUPERVISOR ?= greeting
```

```makefile
deploy-qa: build-linux
	$(call check_required,DEPLOY_HOST,DEPLOY_HOST)
	$(call check_required,DEPLOY_SUPERVISOR,DEPLOY_SUPERVISOR)
	scp -O $(BIN_DIR)/$(BIN_NAME) $(DEPLOY_USR)@$(DEPLOY_HOST):$(DEPLOY_PATH)/
	ssh $(DEPLOY_USR)@$(DEPLOY_HOST) "sudo supervisorctl restart $(DEPLOY_SUPERVISOR)"
```

**备选方案(Alternatives considered)**:
- 保持变量空默认值 + `runqa` 目标特定变量覆盖：拒绝理由(rejected)，变量定义分散在多处，不够直观；且 `deploy-qa` 仍需手动传参
- 每个部署目标（QA/Prod）各设一套默认值：拒绝理由(rejected)，本项目的 QA 环境是唯一稳定部署目标，过度设计

---

## R2: `runqa` 目标实现方式

**决策(Decision)**: `runqa` 直接依赖 `deploy-qa`，无需目标特定变量或额外逻辑

**理由(Rationale)**:
- 所有部署变量已通过 `?=` 设为 QA 默认值，`deploy-qa` 零参数即可部署到 QA
- `runqa` 只需作为语义清晰的别名，降低用户记忆负担
- 无需目标特定变量、`$(eval ...)` 或子 make 调用，实现极简
- 命令行覆盖依然生效：`make runqa DEPLOY_HOST=other` 等同于 `make deploy-qa DEPLOY_HOST=other`

**实现(Implementation)**:
```makefile
.PHONY: runqa
runqa: deploy-qa ## deploy: QA 一键部署 (Quick deploy to QA server)
```

**备选方案(Alternatives considered)**:
- 使用目标特定变量预设 QA 参数：拒绝理由(rejected)，变量默认值已包含 QA 配置，冗余且容易造成默认值不一致
- 创建独立 `deploy-runqa` 目标：拒绝理由(rejected)，代码重复，背离 DRY 原则
- 使用 `export VAR=... && $(MAKE) deploy-qa` 子 make：拒绝理由(rejected)，额外进程开销

---

## R3: 变量覆盖优先级验证

**决策(Decision)**: 依赖 GNU Make 标准变量优先级机制，所有变量均可通过命令行覆盖

**理由(Rationale)**:
GNU Make 变量优先级（从高到低）：
1. 命令行参数 `make VAR=val`（最高）
2. 全局 `?=` / `=` / `:=` 赋值
3. 环境变量（最低）

行为验证：
- `make deploy-qa` → 使用 `?=` 默认值（QA 环境）
- `make deploy-qa DEPLOY_HOST=other` → 命令行覆盖 HOST，其余保持默认
- `DEPLOY_HOST=other make deploy-qa` → 环境变量传入，`?=` 已定义则不覆盖；命令行传参优先级更高

---

## R4: Help 文本格式

**决策(Decision)**: `runqa` 目标使用与现有目标一致的 `## deploy:` 注释格式，中英双语描述

**理由(Rationale)**:
- 沿用 007/R6 的自文档化模式：`## <category>: <描述>`
- 部署类目标归属 `deploy` 分类，`help` 目标自动展示
- 中英双语格式与现有 `deploy-qa` 保持一致

**实现(Implementation)**:
```makefile
runqa: deploy-qa ## deploy: QA 一键部署 (Quick deploy to QA server)
```

---

## R5: `check_required` 用法提示更新

**决策(Decision)**: 更新用法提示，列出所有可用参数（含 `DEPLOY_USR`）

**理由(Rationale)**:
- 现有提示 `DEPLOY_HOST=<host> DEPLOY_SUPERVISOR=<name> [DEPLOY_PATH=<path>]` 未包含 `DEPLOY_USR`
- 所有变量已有默认值，`check_required` 仅在用户显式设为空时触发，提示应展示完整参数列表

**实现(Implementation)**:
```makefile
echo "  Usage: make deploy-qa [DEPLOY_HOST=<host>] [DEPLOY_USR=<user>] [DEPLOY_PATH=<path>] [DEPLOY_SUPERVISOR=<name>]";
```

---

## 总结(Summary)

所有技术决策已收敛，无需要澄清(NEEDS CLARIFICATION)残留：
- 所有部署变量通过 `?=` 设为 QA 默认值，`deploy-qa` 零参数即可运行
- `runqa` 作为 `deploy-qa` 的极简别名
- 变量优先级依赖 GNU Make 标准机制，命令行可覆盖任意参数
- 所有变更仅涉及 `Makefile` 一个文件
