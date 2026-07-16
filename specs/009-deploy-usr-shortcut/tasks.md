# 任务(Tasks): 部署用户变量化与 QA 快捷目标

**输入(Input)**: Design documents from `/specs/009-deploy-usr-shortcut/`

**前置条件(Prerequisites)**: plan.md（必需）、spec.md（用户故事必需）、research.md、data-model.md、contracts/

**测试(Tests)**: 本特性不涉及 Go 代码变更，无需新增单元测试。验证通过 quickstart.md 手动执行。

**组织(Organization)**: 所有变更均在 `Makefile` 单一文件中，任务按用户故事分组。

## 格式(Format): `[ID] [P?] [Story] 描述(Description)`

- **[P]**: 可并行执行（不同文件，无依赖）
- **[Story]**: 属于哪个用户故事（US1、US2）
- 描述中包含具体文件路径

## 路径约定(Path Conventions)

- Makefile: 仓库根目录 `Makefile`

---

## 阶段(Phase) 1: 初始化(Setup)

**目的(Purpose)**: 无需初始化，本特性仅修改已有 `Makefile` 文件，不引入新依赖或文件结构。

**说明(Note)**: 跳过此阶段 — 无项目初始化需求。

---

## 阶段(Phase) 2: 基础设施(Foundational)

**目的(Purpose)**: 无需基础设施变更 — Makefile 已存在，本特性为增量修改。

**说明(Note)**: 跳过此阶段 — 无阻塞性前置依赖。

---

## 阶段(Phase) 3: 用户故事(User Story) 1 - 通过变量指定 SSH 部署用户 (优先级(Priority): P1) 🎯 MVP

**目标(Goal)**: 将所有部署变量设置为 QA 默认值，新增 `DEPLOY_USR` 变量替代硬编码 `root`，使 `make deploy-qa` 零参数即可运行

**独立测试(Independent Test)**: 执行 `make deploy-qa`（不传任何参数），验证 scp/ssh 使用 `ubuntu@111.229.4.203`，更新后的 `check_required` 提示包含所有参数

### 用户故事(User Story) 1 的实现

- [X] T001 [US1] 更新部署变量默认值为 QA 环境值：`DEPLOY_HOST ?= 111.229.4.203`、`DEPLOY_PATH ?= /opt/project/greeting`、`DEPLOY_USR ?= ubuntu`（新增）、`DEPLOY_SUPERVISOR ?= greeting` 在 `Makefile` 第 7-10 行
- [X] T002 [US1] 替换 `deploy-qa` 目标中硬编码的 `root@$(DEPLOY_HOST)` 为 `$(DEPLOY_USR)@$(DEPLOY_HOST)`：修改 `Makefile` 第 92 行 scp 命令和第 94 行 ssh 命令
- [X] T003 [US1] 更新 `check_required` 宏的用法提示为 `make deploy-qa [DEPLOY_HOST=<host>] [DEPLOY_USR=<user>] [DEPLOY_PATH=<path>] [DEPLOY_SUPERVISOR=<name>]` 在 `Makefile` 第 18 行

**检查点(Checkpoint)**: `make deploy-qa` 可零参数运行，SSH 使用 `ubuntu@111.229.4.203`；`make deploy-qa DEPLOY_USR=root` 可覆盖用户

---

## 阶段(Phase) 4: 用户故事(User Story) 2 - QA 环境一键部署快捷目标 (优先级(Priority): P2)

**目标(Goal)**: 新增 `runqa` 目标作为 `deploy-qa` 的别名，让开发者通过更简短命令完成 QA 部署

**依赖(Dependency)**: 依赖 US1 完成（变量默认值已设）

**独立测试(Independent Test)**: 执行 `make runqa`，验证等价于 `make deploy-qa`；`make help` 显示 `runqa` 说明

### 用户故事(User Story) 2 的实现

- [X] T004 [US2] 新增 `runqa` 目标：在 `Makefile` 中 `.PHONY` 声明附近添加 `.PHONY: runqa`，在 `deploy-qa` 目标附近添加 `runqa: deploy-qa ## deploy: QA 一键部署 (Quick deploy to QA server)` 规则

**检查点(Checkpoint)**: `make runqa` 一键部署到 QA；`make runqa DEPLOY_HOST=X` 可覆盖参数；`make help` 显示 runqa 说明

---

## 阶段(Phase) 5: 收尾与验证(Polish & Verification)

**目的(Purpose)**: 验证所有变更行为正确，确保质量

- [X] T005 执行 quickstart.md 验证场景：确认 `make deploy-qa` 零参数运行、`make runqa` 等价性、`check_required` 提示更新、`make help` 显示 runqa
- [X] T006 运行 `go build ./...` 确保 Go 代码编译正常（Makefile 变更不影响 Go 编译，但作为提交前检查）

---

## 依赖与执行顺序(Dependencies & Execution Order)

### 阶段依赖(Phase Dependencies)

- **阶段 1（初始化）**: 跳过
- **阶段 2（基础设施）**: 跳过
- **阶段 3（US1）**: 无前置依赖 — 可直接开始
- **阶段 4（US2）**: 依赖 US1 完成（`runqa` 依赖 `deploy-qa` 拥有 QA 默认值）
- **阶段 5（验证）**: 依赖 US1 + US2 完成

### 用户故事依赖(User Story Dependencies)

- **US1（P1）**: 无依赖 — 可立即实现
- **US2（P2）**: 依赖 US1（变量默认值必须先完成，`runqa` 才能零参数运行）

### 用户故事内部(Within Each User Story)

- US1: T001（变量默认值）→ T002（scp/ssh 命令）→ T003（用法提示）。T002 和 T003 可并行。
- US2: T004 单任务，无内部依赖。

### 并行机会(Parallel Opportunities)

- US1 内部：T002 和 T003 可并行执行（修改 Makefile 不同行且无逻辑依赖）
- 所有变更在同一文件 `Makefile`，串行提交可减少冲突

---

## 实施策略(Implementation Strategy)

### MVP 优先(MVP First)（仅用户故事 1）

1. 完成阶段 3：US1（T001 → T002 + T003）
2. **停止并验证**: `make deploy-qa` 零参数运行，使用 `ubuntu@111.229.4.203`
3. 此时 MVP 已完成 — `deploy-qa` 开箱即用

### 增量交付(Incremental Delivery)

1. US1 完成 → `make deploy-qa` 零参数即可部署到 QA（MVP！）
2. US2 完成 → `make runqa` 更简短，提升效率
3. 阶段 5 验证 → 质量确认

---

## 备注(Notes)

- 所有变更仅在 `Makefile` 一个文件中
- 提交前确保 `go build ./...` 通过（虽然 Makefile 变更不影响 Go 编译，但作为惯例检查）
- Commit message 格式: `feat: add DEPLOY_USR variable and runqa shortcut target`
- 无 Go 代码变更，无需更新单元测试
