---

description: "任务列表模板(Task list template for feature implementation)"
---

# 任务(Tasks): [功能名称]

**输入(Input)**: Design documents from `/specs/[###-feature-name]/`

**前置条件(Prerequisites)**: plan.md（必需）、spec.md（用户故事必需）、research.md、data-model.md、contracts/

**测试(Tests)**: The examples below include test tasks. Tests are OPTIONAL - only include them if explicitly requested in the feature specification.

**组织(Organization)**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## 格式(Format): `[ID] [P?] [Story] 描述(Description)`

- **[P]**: 可并行执行 (不同文件，无依赖)
- **[Story]**: 属于哪个用户故事（如 US1、US2、US3）
- 描述中包含具体文件路径

## 路径约定(Path Conventions)

- **单体项目**: `src/`、`tests/` 在仓库根目录
- **Web 应用**: `backend/src/`、`frontend/src/`
- **移动端**: `api/src/`、`ios/src/` 或 `android/src/`
- 以下路径按单体项目展示 - 根据 plan.md 结构调整

<!--
  ============================================================================
  IMPORTANT: The tasks below are SAMPLE TASKS for illustration purposes only.

  The /speckit.tasks command MUST replace these with actual tasks based on:
  - User stories from spec.md (with their priorities P1, P2, P3...)
  - Feature requirements from plan.md
  - Entities from data-model.md
  - Endpoints from contracts/

  Tasks MUST be organized by user story so each story can be:
  - Implemented independently
  - Tested independently
  - Delivered as an MVP increment

  DO NOT keep these sample tasks in the generated tasks.md file.
  ============================================================================
-->

## 阶段(Phase) 1: 初始化(Setup)（共享基础设施）

**目的(Purpose)**: 项目初始化和基本结构搭建

- [ ] T001 按实施计划创建项目结构
- [ ] T002 用 [framework] 依赖初始化 [language] 项目
- [ ] T003 [P] 配置 linting 和格式化工具

---

## 阶段(Phase) 2: 基础设施(Foundational)（阻塞性前置）

**目的(Purpose)**: 任何用户故事开始前必须完成的核心基础设施

**⚠️ 关键**: 此阶段完成前，不得开始任何用户故事

基础任务示例（根据项目调整）：

- [ ] T004 设置数据库 schema 和迁移框架
- [ ] T005 [P] 实现认证/授权框架
- [ ] T006 [P] 设置 API 路由和中间件结构
- [ ] T007 创建所有故事依赖的基础模型/实体
- [ ] T008 配置错误处理和日志基础设施
- [ ] T009 设置环境配置管理

**检查点(Checkpoint)**: 基础设施就绪 - 可开始并行实现用户故事

---

## 阶段(Phase) 3: 用户故事(User Story) 1 - [标题] (优先级(Priority): P1) 🎯 MVP

**目标(Goal)**: [简要描述此故事交付什么]

**独立测试(Independent Test)**: [如何独立验证此故事是否可用]

### 用户故事(User Story) 1 的测试（可选 - 仅当明确要求时）⚠️

> **注意：先写测试，确保实现前测试失败**

- [ ] T010 [P] [US1] 在 tests/contract/test_[name].py 中写 [endpoint] 的契约测试
- [ ] T011 [P] [US1] 在 tests/integration/test_[name].py 中写 [用户旅程] 的集成测试

### 用户故事(User Story) 1 的实现

- [ ] T012 [P] [US1] 在 src/models/[entity1].py 中创建 [Entity1] 模型
- [ ] T013 [P] [US1] 在 src/models/[entity2].py 中创建 [Entity2] 模型
- [ ] T014 [US1] 在 src/services/[service].py 中实现 [Service]（依赖 T012、T013）
- [ ] T015 [US1] 在 src/[location]/[file].py 中实现 [endpoint/feature]
- [ ] T016 [US1] 添加验证和错误处理
- [ ] T017 [US1] 为用户故事 1 添加日志记录

**检查点(Checkpoint)**: 此时用户故事 1 应完整可用并可独立测试

---

## 阶段(Phase) 4: 用户故事(User Story) 2 - [标题] (优先级(Priority): P2)

**目标(Goal)**: [简要描述此故事交付什么]

**独立测试(Independent Test)**: [如何独立验证此故事是否可用]

### 用户故事(User Story) 2 的测试（可选 - 仅当明确要求时）⚠️

- [ ] T018 [P] [US2] 在 tests/contract/test_[name].py 中写 [endpoint] 的契约测试
- [ ] T019 [P] [US2] 在 tests/integration/test_[name].py 中写 [用户旅程] 的集成测试

### 用户故事(User Story) 2 的实现

- [ ] T020 [P] [US2] 在 src/models/[entity].py 中创建 [Entity] 模型
- [ ] T021 [US2] 在 src/services/[service].py 中实现 [Service]
- [ ] T022 [US2] 在 src/[location]/[file].py 中实现 [endpoint/feature]
- [ ] T023 [US2] 与用户故事 1 组件集成（如需要）

**检查点(Checkpoint)**: 此时用户故事 1 和 2 均应独立可用

---

## 阶段(Phase) 5: 用户故事(User Story) 3 - [标题] (优先级(Priority): P3)

**目标(Goal)**: [简要描述此故事交付什么]

**独立测试(Independent Test)**: [如何独立验证此故事是否可用]

### 用户故事(User Story) 3 的测试（可选 - 仅当明确要求时）⚠️

- [ ] T024 [P] [US3] 在 tests/contract/test_[name].py 中写 [endpoint] 的契约测试
- [ ] T025 [P] [US3] 在 tests/integration/test_[name].py 中写 [用户旅程] 的集成测试

### 用户故事(User Story) 3 的实现

- [ ] T026 [P] [US3] 在 src/models/[entity].py 中创建 [Entity] 模型
- [ ] T027 [US3] 在 src/services/[service].py 中实现 [Service]
- [ ] T028 [US3] 在 src/[location]/[file].py 中实现 [endpoint/feature]

**检查点(Checkpoint)**: 所有用户故事应均可独立使用

---

[根据需要添加更多用户故事阶段，遵循相同模式]

---

## 阶段(Phase) N: 收尾与横切关注点(Polish & Cross-Cutting Concerns)

**目的(Purpose)**: 影响多个用户故事的改进

- [ ] TXXX [P] 更新 docs/ 中的文档
- [ ] TXXX 代码清理和重构
- [ ] TXXX 跨所有故事进行性能优化
- [ ] TXXX [P] 补充单元测试（如需要）于 tests/unit/
- [ ] TXXX 安全加固
- [ ] TXXX 执行 quickstart.md 验证

---

## 依赖与执行顺序(Dependencies & Execution Order)

### 阶段依赖(Phase Dependencies)

- **初始化(Setup)（阶段 1）**: 无依赖 - 可立即开始
- **基础设施(Foundational)（阶段 2）**: 依赖初始化完成 - 阻塞所有用户故事
- **用户故事（阶段 3+）**: 全部依赖基础设施阶段完成
  - 用户故事随后可并行进行（如有人员）
  - 或按优先级顺序串行（P1 → P2 → P3）
- **收尾（最终阶段）**: 依赖所有目标用户故事完成

### 用户故事依赖(User Story Dependencies)

- **用户故事 1（P1）**: 基础设施（阶段 2）完成后即可开始 - 不依赖其他故事
- **用户故事 2（P2）**: 基础设施（阶段 2）完成后即可开始 - 可与 US1 集成但应可独立测试
- **用户故事 3（P3）**: 基础设施（阶段 2）完成后即可开始 - 可与 US1/US2 集成但应可独立测试

### 用户故事内部(Within Each User Story)

- 测试（如包含）必须(MUST)先写并在实现前失败
- 模型先于服务
- 服务先于端点
- 核心实现先于集成
- 当前故事完成后再进入下一优先级

### 并行机会(Parallel Opportunities)

- 标记 [P] 的初始化任务可并行执行
- 标记 [P] 的基础设施任务可并行执行（在阶段 2 内）
- 基础设施阶段完成后，所有用户故事可并行开始（如团队容量允许）
- 标记 [P] 的故事内测试可并行执行
- 标记 [P] 的故事内模型可并行执行
- 不同用户故事可由不同成员并行开发

---

## 并行示例(Parallel Example): 用户故事(User Story) 1

```bash
# 同时启动用户故事 1 的所有测试（如需要）：
Task: "在 tests/contract/test_[name].py 中写 [endpoint] 的契约测试"
Task: "在 tests/integration/test_[name].py 中写 [用户旅程] 的集成测试"

# 同时启动用户故事 1 的所有模型：
Task: "在 src/models/[entity1].py 中创建 [Entity1] 模型"
Task: "在 src/models/[entity2].py 中创建 [Entity2] 模型"
```

---

## 实施策略(Implementation Strategy)

### MVP 优先(MVP First)（仅用户故事 1）

1. 完成阶段 1：初始化(Setup)
2. 完成阶段 2：基础设施(Foundational)（关键 - 阻塞所有故事）
3. 完成阶段 3：用户故事 1
4. **停止并验证**: 独立测试用户故事 1
5. 如就绪则部署/演示

### 增量交付(Incremental Delivery)

1. 完成初始化 + 基础设施 → 基础就绪
2. 添加用户故事 1 → 独立测试 → 部署/演示（MVP！）
3. 添加用户故事 2 → 独立测试 → 部署/演示
4. 添加用户故事 3 → 独立测试 → 部署/演示
5. 每个故事在不断增加价值的同时不影响已有故事

### 并行团队策略(Parallel Team Strategy)

多名开发者时：

1. 团队共同完成初始化 + 基础设施
2. 基础设施完成后：
   - 开发者 A：用户故事 1
   - 开发者 B：用户故事 2
   - 开发者 C：用户故事 3
3. 各故事独立完成并集成

---

## 备注(Notes)

- [P] 任务 = 不同文件，无依赖
- [Story] 标签将任务映射到特定用户故事，便于追踪
- 每个用户故事应可独立完成和测试
- 实现前验证测试失败
- 每个任务或逻辑组完成后提交
- 可在任何检查点暂停以独立验证故事
- 避免：模糊任务、同文件冲突、破坏故事独立性的跨故事依赖
