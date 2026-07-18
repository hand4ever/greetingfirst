---
description: Generate an actionable, dependency-ordered tasks.md for the feature based on available design artifacts.
handoffs:
  - label: Analyze For Consistency
    agent: speckit.analyze
    prompt: Run a project analysis for consistency
    send: true
  - label: Implement Project
    agent: speckit.implement
    prompt: Start the implementation in phases
    send: true
---

## 用户输入(User Input)

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## 执行前检查(Pre-Execution Checks)

**Check for extension hooks (before tasks generation / 任务生成前的扩展钩子检查)**:
- Check if `.specify/extensions.yml` exists in the project root.
- If it exists, read it and look for entries under the `hooks.before_tasks` key
- If the YAML cannot be parsed or is invalid, skip hook checking silently and continue normally
- Filter out hooks where `enabled` is explicitly `false`. Treat hooks without an `enabled` field as enabled by default.
- For each remaining hook, do **not** attempt to interpret or evaluate hook `condition` expressions:
  - If the hook has no `condition` field, or it is null/empty, treat it as executable
  - If the hook defines a non-empty `condition`, skip the hook and leave condition evaluation to the HookExecutor implementation
- For each executable hook, output the following based on its `optional` flag:
  - **Optional hook** (`optional: true` / 可选钩子):
    ```
    ## 扩展钩子(Extension Hooks)

    **Optional Pre-Hook**: {extension}
    Command: `/{command}`
    Description: {description}

    Prompt: {prompt}
    To execute: `/{command}`
    ```
  - **Mandatory hook** (`optional: false` / 强制钩子):
    ```
    ## 扩展钩子(Extension Hooks)

    **Automatic Pre-Hook**: {extension}
    Executing: `/{command}`
    EXECUTE_COMMAND: {command}

    Wait for the result of the hook command before proceeding to the Outline.
    ```
    After emitting the block above you MUST actually invoke the hook and wait for it to finish before continuing. Run it the same way you would run the command yourself in this agent/session (the invocation may differ from the literal `{command}` id shown above, e.g. a skills-mode agent runs it as `/skill:speckit-...` or `$speckit-...`). Emitting the block alone does not run the hook.
- If no hooks are registered or `.specify/extensions.yml` does not exist, skip silently

## 大纲(Outline)

1. **Setup** (初始化): Run `.specify/scripts/bash/setup-tasks.sh --json` from repo root and parse FEATURE_DIR, TASKS_TEMPLATE, and AVAILABLE_DOCS list. `FEATURE_DIR` and `TASKS_TEMPLATE` must be absolute paths when provided. `AVAILABLE_DOCS` is a list of document names/relative paths available under `FEATURE_DIR` (for example `research.md` or `contracts/`). For single quotes in args like "I'm Groot", use escape syntax: e.g 'I'\''m Groot' (or double-quote if possible: "I'm Groot").

2. **Load design documents** (加载设计文档): Read from FEATURE_DIR:
   - **Required** (必填): plan.md (tech stack, libraries, structure), spec.md (user stories with priorities)
   - **Optional** (可选): data-model.md (entities), contracts/ (interface contracts), research.md (decisions), quickstart.md (test scenarios)
   - **IF EXISTS**: Load `.specify/memory/constitution.md` for project principles and governance constraints
   - Note: Not all projects have all documents. Generate tasks based on what's available.

3. **Execute task generation workflow** (执行任务生成工作流):
   - Load plan.md and extract tech stack, libraries, project structure
   - Load spec.md and extract user stories with their priorities (P1, P2, P3, etc.)
   - If data-model.md exists: Extract entities and map to user stories
   - If contracts/ exists: Map interface contracts to user stories
   - If research.md exists: Extract decisions for setup tasks
   - Generate tasks organized by user story (see Task Generation Rules below / 参见下方任务生成规则)
   - Generate dependency graph showing user story completion order
   - Create parallel execution examples per user story
   - Validate task completeness (each user story has all needed tasks, independently testable / 每用户故事拥有全部所需任务且可独立测试)

4. **Generate tasks.md** (生成 tasks.md): Read the tasks template from TASKS_TEMPLATE (from the JSON output above) and use it as structure. If TASKS_TEMPLATE is empty, fall back to `.specify/templates/tasks-template.md`. Fill with:
   - Correct feature name from plan.md (来自 plan.md 的正确功能名)
   - Phase 1: Setup tasks (project initialization / 初始化任务)
   - Phase 2: Foundational tasks (blocking prerequisites for all user stories / 所有用户故事的阻塞前置任务)
   - Phase 3+: One phase per user story (in priority order from spec.md / 按 spec.md 优先级每用户故事一个阶段)
   - Each phase includes: story goal, independent test criteria, tests (if requested), implementation tasks
   - Final Phase: Polish & cross-cutting concerns (打磨与横切关注点)
   - All tasks must follow the strict checklist format (see Task Generation Rules below / 严格遵循检查清单格式)
   - Clear file paths for each task (每个任务的清晰文件路径)
   - Dependencies section showing story completion order (展示故事完成顺序的依赖章节)
   - Parallel execution examples per story (每故事并行执行示例)
   - Implementation strategy section (MVP first, incremental delivery / 实现策略章节：先 MVP，增量交付)

## 强制后置钩子(Mandatory Post-Execution Hooks)

**You MUST complete this section before reporting completion to the user.**

Check if `.specify/extensions.yml` exists in the project root.
- If it does not exist, or no hooks are registered under `hooks.after_tasks`, skip to the Completion Report.
- If it exists, read it and look for entries under the `hooks.after_tasks` key.
- If the YAML cannot be parsed or is invalid, skip hook checking silently and continue to the Completion Report.
- Filter out hooks where `enabled` is explicitly `false`. Treat hooks without an `enabled` field as enabled by default.
- For each remaining hook, do **not** attempt to interpret or evaluate hook `condition` expressions:
  - If the hook has no `condition` field, or it is null/empty, treat it as executable
  - If the hook defines a non-empty `condition`, skip the hook and leave condition evaluation to the HookExecutor implementation
- For each executable hook, output the following based on its `optional` flag:
  - **Mandatory hook** (`optional: false` / 强制钩子) — **You MUST emit `EXECUTE_COMMAND:` for each mandatory hook**:
    ```
    ## 扩展钩子(Extension Hooks)

    **Automatic Hook**: {extension}
    Executing: `/{command}`
    EXECUTE_COMMAND: {command}
    ```
    After emitting the block above you MUST actually invoke the hook and wait for it to finish before continuing. Run it the same way you would run the command yourself in this agent/session (the invocation may differ from the literal `{command}` id shown above, e.g. a skills-mode agent runs it as `/skill:speckit-...` or `$speckit-...`). Emitting the block alone does not run the hook.
  - **Optional hook** (`optional: true` / 可选钩子):
    ```
    ## 扩展钩子(Extension Hooks)

    **Optional Hook**: {extension}
    Command: `/{command}`
    Description: {description}

    Prompt: {prompt}
    To execute: `/{command}`
    ```

## 完成报告(Completion Report)

Output path to generated tasks.md and summary:
- Total task count (任务总数)
- Task count per user story (每用户故事任务数)
- Parallel opportunities identified (识别的并行机会)
- Independent test criteria for each story (每故事的独立测试标准)
- Suggested MVP scope (typically just User Story 1 / 建议的 MVP 范围，通常仅用户故事 1)
- Format validation: Confirm ALL tasks follow the checklist format (checkbox, ID, labels, file paths / 格式校验：确认所有任务遵循检查清单格式)

Context for task generation: $ARGUMENTS

The tasks.md should be immediately executable - each task must be specific enough that an LLM can complete it without additional context.

## 任务生成规则(Task Generation Rules)

**CRITICAL** (关键): Tasks MUST be organized by user story to enable independent implementation and testing (任务必须按用户故事组织，以支持独立实现与测试).

**Tests are OPTIONAL** (测试可选): Only generate test tasks if explicitly requested in the feature specification or if user requests TDD approach.

### 检查清单格式(必填)(Checklist Format (REQUIRED))

Every task MUST strictly follow this format (每个任务必须严格遵循此格式):

```text
- [ ] [TaskID] [P?] [Story?] Description with file path
```

**Format Components** (格式组成):
1. **Checkbox** (复选框): ALWAYS start with `- [ ]` (markdown checkbox)
2. **Task ID** (任务 ID): Sequential number (T001, T002, T003...) in execution order
3. **[P] marker** ([P] 标记): Include ONLY if task is parallelizable (different files, no dependencies on incomplete tasks / 仅当任务可并行时包含：不同文件、对未完成任务无依赖)
4. **[Story] label** ([故事] 标签): REQUIRED for user story phase tasks only
   - Format: [US1], [US2], [US3], etc. (maps to user stories from spec.md / 映射到 spec.md 中的用户故事)
   - Setup phase: NO story label (初始化阶段：无故事标签)
   - Foundational phase: NO story label (基础阶段：无故事标签)
   - User Story phases: MUST have story label (用户故事阶段：必须有故事标签)
   - Polish phase: NO story label (打磨阶段：无故事标签)
5. **Description** (描述): Clear action with exact file path (含确切文件路径的清晰动作)

**Examples** (示例):
- ✅ CORRECT (正确): `- [ ] T001 Create project structure per implementation plan`
- ✅ CORRECT (正确): `- [ ] T005 [P] Implement authentication middleware in src/middleware/auth.py`
- ✅ CORRECT (正确): `- [ ] T012 [P] [US1] Create User model in src/models/user.py`
- ✅ CORRECT (正确): `- [ ] T014 [US1] Implement UserService in src/services/user_service.py`
- ❌ WRONG (错误): `- [ ] Create User model` (missing ID and Story label / 缺少 ID 与故事标签)
- ❌ WRONG (错误): `T001 [US1] Create model` (missing checkbox / 缺少复选框)
- ❌ WRONG (错误): `- [ ] [US1] Create User model` (missing Task ID / 缺少任务 ID)
- ❌ WRONG (错误): `- [ ] T001 [US1] Create model` (missing file path / 缺少文件路径)

### 任务组织(Task Organization)

1. **From User Stories (spec.md)** (来自用户故事) - PRIMARY ORGANIZATION (主要组织方式):
   - Each user story (P1, P2, P3...) gets its own phase (每个用户故事拥有独立阶段)
   - Map all related components to their story (将相关组件映射到其故事)：
     - Models needed for that story (该故事所需模型)
     - Services needed for that story (该故事所需服务)
     - Interfaces/UI needed for that story (该故事所需接口/UI)
     - If tests requested: Tests specific to that story (若请求测试：该故事专属测试)
   - Mark story dependencies (most stories should be independent / 标记故事依赖，多数故事应独立)

2. **From Contracts** (来自契约):
   - Map each interface contract → to the user story it serves
   - If tests requested: Each interface contract → contract test task [P] before implementation in that story's phase

3. **From Data Model** (来自数据模型):
   - Map each entity to the user story(ies) that need it
   - If entity serves multiple stories: Put in earliest story or Setup phase
   - Relationships → service layer tasks in appropriate story phase

4. **From Setup/Infrastructure** (来自初始化/基础设施):
   - Shared infrastructure → Setup phase (Phase 1 / 共享基础设施 → 初始化阶段)
   - Foundational/blocking tasks → Foundational phase (Phase 2 / 基础/阻塞任务 → 基础阶段)
   - Story-specific setup → within that story's phase (故事专属初始化 → 在该故事阶段内)

### 阶段结构(Phase Structure)

- **Phase 1** (阶段 1): Setup (project initialization / 初始化)
- **Phase 2** (阶段 2): Foundational (blocking prerequisites - MUST complete before user stories / 阻塞前置——必须在用户故事前完成)
- **Phase 3+** (阶段 3+): User Stories in priority order (P1, P2, P3...)
  - Within each story: Tests (if requested) → Models → Services → Endpoints → Integration
  - Each phase should be a complete, independently testable increment (每个阶段应为完整、可独立测试的成果增量)
- **Final Phase** (最终阶段): Polish & Cross-Cutting Concerns (打磨与横切关注点)

## 完成条件(Done When)

- [ ] tasks.md generated with all phases, task IDs, and file paths (tasks.md 已生成，含所有阶段、任务 ID 与文件路径)
- [ ] Extension hooks dispatched or skipped according to the rules in Mandatory Post-Execution Hooks above (扩展钩子已按上述规则派发或跳过)
- [ ] Completion reported to user with task count, story breakdown, and MVP scope (已向用户报告任务数、故事拆分与 MVP 范围)
