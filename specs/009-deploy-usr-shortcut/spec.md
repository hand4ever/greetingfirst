# 功能规格说明(Feature Specification): 部署用户变量化与 QA 快捷目标

**功能分支(Feature Branch)**: `009-deploy-usr-shortcut`

**创建日期(Created)**: 2026-07-16

**状态(Status)**: 草稿(Draft)

**输入(Input)**: 用户描述(User description): "追加：1 makefile里root改为ubuntu，或者改为传参 2 将下面这个命令做个shortfor，如make runqa:make deploy-qa DEPLOY_HOST=111.229.4.203 DEPLOY_USR=ubuntu DEPLOY_SUPERVISOR=greeting DEPLOY_PATH=/opt/project/greeting"

## 用户场景与测试(User Scenarios & Testing) *(必填)*

### 用户故事(User Story) 1 - 通过变量指定 SSH 部署用户 (优先级(Priority): P1)

作为运维/后端开发者，我希望通过 Makefile 变量指定 SSH 部署用户（如 `ubuntu`），而非在 Makefile 中硬编码 `root`，以便适配不同目标服务器的用户配置。

**优先级理由(Why this priority)**: 硬编码 `root` 用户是安全隐患，且不同服务器可能使用不同用户（如 `ubuntu`、`deploy`），变量化是安全性和灵活性的基础要求。

**独立测试(Independent Test)**: 执行 `make deploy-qa`（使用默认值），验证 scp/ssh 命令使用 `ubuntu@111.229.4.203` 而非 `root@111.229.4.203`。

**验收场景(Acceptance Scenarios)**:

1. **假设(Given)** 部署变量默认值已设为 `ubuntu`, **当(When)** 开发者执行 `make deploy-qa`, **那么(Then)** scp 和 ssh 命令均使用 `ubuntu@111.229.4.203` 进行连接，部署成功后重启服务
2. **假设(Given)** 开发者通过命令行覆盖 `DEPLOY_USR=root`, **当(When)** 执行 `make deploy-qa DEPLOY_USR=root`, **那么(Then)** 使用 `root@111.229.4.203` 进行 SSH 连接
3. **假设(Given)** `DEPLOY_USR` 设置为不存在的用户, **当(When)** 开发者执行 `make deploy-qa DEPLOY_USR=nobody`, **那么(Then)** SSH 连接失败并显示认证错误，终止部署

---

### 用户故事(User Story) 2 - QA 环境一键部署快捷目标 (优先级(Priority): P2)

作为一名后端开发者，我希望通过执行 `make runqa`（`deploy-qa` 的别名）一键部署到 QA 服务器，无需每次输入完整的 `DEPLOY_HOST`、`DEPLOY_USR`、`DEPLOY_SUPERVISOR`、`DEPLOY_PATH` 等参数。

**优先级理由(Why this priority)**: QA 服务器的部署参数相对固定，每次部署都需要手动输入多个参数是一种重复性工作，快捷目标能显著提升部署效率。

**独立测试(Independent Test)**: 执行 `make runqa`，验证等价于执行 `make deploy-qa`（所有变量均已设为 QA 默认值），编译上传并重启服务。

**验收场景(Acceptance Scenarios)**:

1. **假设(Given)** Makefile 变量已设置 QA 默认值（`DEPLOY_HOST=111.229.4.203`、`DEPLOY_USR=ubuntu`、`DEPLOY_SUPERVISOR=greeting`、`DEPLOY_PATH=/opt/project/greeting`）, **当(When)** 开发者执行 `make runqa`, **那么(Then)** 使用这些默认值执行完整部署流程
2. **假设(Given)** 开发者执行 `make runqa`, **当(When)** QA 服务器不可达, **那么(Then)** 在 scp 步骤失败并显示连接错误，不执行后续重启步骤
3. **假设(Given)** 开发者想临时覆盖某个参数（如 `DEPLOY_HOST`）, **当(When)** 执行 `make runqa DEPLOY_HOST=other-host`, **那么(Then)** 使用覆盖后的值进行部署

---

### 边界情况(Edge Cases)

- `DEPLOY_USR` 变量设为空字符串时，`?=` 默认值生效（`ubuntu`）
- `runqa` 目标与 `deploy-qa` 在功能上完全等价（同为 QA 默认值）
- 当 QA 环境参数需要更新时（如更换 QA 服务器 IP），只需修改 Makefile 中对应的 `?=` 默认值

## 需求(Requirements) *(必填)*

### 功能需求(Functional Requirements)

- **FR-001**: Makefile 必须(MUST) 新增 `DEPLOY_USR` 变量，默认值为 `ubuntu`，支持通过环境变量或 `make` 变量覆盖；同时将 `DEPLOY_HOST`、`DEPLOY_PATH`、`DEPLOY_SUPERVISOR` 的默认值设为 QA 环境值（`111.229.4.203`、`/opt/project/greeting`、`greeting`）
- **FR-002**: `deploy-qa` 目标中的 `scp` 和 `ssh` 命令必须(MUST) 使用 `$(DEPLOY_USR)@$(DEPLOY_HOST)` 替代硬编码的 `root@$(DEPLOY_HOST)`
- **FR-003**: `check_required` 宏的用法提示必须(MUST) 更新为展示所有可选参数
- **FR-004**: Makefile 必须(MUST) 新增 `runqa` 目标，作为 `deploy-qa` 的别名（所有变量已设为 QA 默认值，无需额外传参）
- **FR-005**: `runqa` 目标必须(MUST) 声明为 `.PHONY`
- **FR-006**: `runqa` 目标必须(MUST) 允许开发者通过命令行覆盖任意参数
- **FR-007**: `help` 输出必须(MUST) 包含 `runqa` 目标的中英双语说明

### 关键实体(Key Entities)

无新增数据实体。仅涉及 Makefile 中的变量定义和目标扩展：

- **构建变量(Build Variable)**: `DEPLOY_USR` — 新增的 SSH 部署用户变量，默认值 `ubuntu`
- **Make 目标(Make Target)**: `runqa` — 新增的 QA 一键部署快捷目标

## 成功标准(Success Criteria) *(必填)*

### 可衡量成果(Measurable Outcomes)

- **SC-001**: 执行 `make deploy-qa`（不传任何参数）时，所有 SSH 连接均使用 `ubuntu@111.229.4.203`
- **SC-002**: 执行 `make deploy-qa DEPLOY_USR=root` 时，可覆盖用户为 `root`
- **SC-003**: 开发者执行 `make runqa` 完成部署所需输入字符数减少 80% 以上（从完整参数输入变为单一命令）
- **SC-004**: `make help` 中 `runqa` 目标正确展示其用途说明

## 假设(Assumptions)

- 基于 `007-optimize-makefile` 重构后的 Makefile 进行增量修改，当前 `deploy-qa` 目标中 `scp`/`ssh` 使用硬编码 `root@`（实际代码第 92、94 行）
- QA 服务器配置（IP `111.229.4.203`、用户 `ubuntu`、路径 `/opt/project/greeting`、supervisor `greeting`）相对稳定，作为 Makefile 变量的默认值
- 所有部署变量均通过 `?=` 设置默认值，命令行参数可覆盖；`deploy-qa` 和 `runqa` 在功能上等价
- `runqa` 目标作为 `deploy-qa` 的别名实现，无额外逻辑
