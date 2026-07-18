---
description: Execute the implementation planning workflow using the plan template to generate design artifacts.
handoffs:
  - label: Create Tasks
    agent: speckit.tasks
    prompt: Break the plan into tasks
    send: true
  - label: Create Checklist
    agent: speckit.checklist
    prompt: Create a checklist for the following domain...
---

## 用户输入(User Input)

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## 执行前检查(Pre-Execution Checks)

**Check for extension hooks (before planning / 规划前的扩展钩子检查)**:
- Check if `.specify/extensions.yml` exists in the project root.
- If it exists, read it and look for entries under the `hooks.before_plan` key
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

1. **Setup** (初始化): Run `.specify/scripts/bash/setup-plan.sh --json` from repo root and parse JSON for FEATURE_SPEC, IMPL_PLAN, SPECS_DIR, BRANCH. For single quotes in args like "I'm Groot", use escape syntax: e.g 'I'\''m Groot' (or double-quote if possible: "I'm Groot").

2. **Load context** (加载上下文): Read FEATURE_SPEC and `.specify/memory/constitution.md`. Load IMPL_PLAN template (already copied).

3. **Execute plan workflow** (执行规划工作流): Follow the structure in IMPL_PLAN template to:
   - Fill Technical Context (mark unknowns as "NEEDS CLARIFICATION" / 标记未知为"需澄清")
   - Fill Constitution Check section from constitution (从宪法填充宪法检查章节)
   - Evaluate gates (ERROR if violations unjustified / 若违规未被论证则报错)
   - Phase 0: Generate research.md (resolve all NEEDS CLARIFICATION / 生成研究文档，解决所有需澄清项)
   - Phase 1: Generate data-model.md, contracts/, quickstart.md (生成数据模型、契约、快速上手指南)
   - Re-evaluate Constitution Check post-design (设计后重新评估宪法检查)

## 强制后置钩子(Mandatory Post-Execution Hooks)

**You MUST complete this section before reporting completion to the user.**

Check if `.specify/extensions.yml` exists in the project root.
- If it does not exist, or no hooks are registered under `hooks.after_plan`, skip to the Completion Report.
- If it exists, read it and look for entries under the `hooks.after_plan` key.
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

Command ends after Phase 1 design. Report branch, IMPL_PLAN path, and generated artifacts.

## 阶段(Phases)

### 阶段 0：大纲与研究(Phase 0: Outline & Research)

1. **Extract unknowns from Technical Context** (从技术上下文中提取未知项) above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents** (生成并派发研究代理):

   ```text
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** (整合结论) in `research.md` using format:
   - Decision: [what was chosen] (决策：所选方案)
   - Rationale: [why chosen] (理由：选择原因)
   - Alternatives considered: [what else evaluated] (已考虑的备选方案)

**Output** (输出): research.md with all NEEDS CLARIFICATION resolved (所有需澄清项已解决)

### 阶段 1：设计与契约(Phase 1: Design & Contracts)

**Prerequisites** (前置条件): `research.md` complete

1. **Extract entities from feature spec** (从功能规格提取实体) → `data-model.md`:
   - Entity name, fields, relationships (实体名、字段、关系)
   - Validation rules from requirements (来自需求的校验规则)
   - State transitions if applicable (状态转换，若适用)

2. **Define interface contracts** (定义接口契约, if project has external interfaces) → `/contracts/`:
   - Identify what interfaces the project exposes to users or other systems
   - Document the contract format appropriate for the project type
   - Examples: public APIs for libraries, command schemas for CLI tools, endpoints for web services, grammars for parsers, UI contracts for applications
   - Skip if project is purely internal (build scripts, one-off tools, etc.)

3. **Create quickstart validation guide** (创建快速上手验证指南) → `quickstart.md`:
   - Document runnable validation scenarios that prove the feature works end-to-end
   - Include prerequisites, setup commands, test/run commands, and expected outcomes
   - Use links or references to contracts and data model details instead of duplicating them
   - Do not include full implementation code, model/service/controller bodies, migrations, or complete test suites
   - Keep this artifact as a validation/run guide; implementation details belong in `tasks.md` and the implementation phase

**Output** (输出): data-model.md, /contracts/*, quickstart.md

## 关键规则(Key rules)

- Use absolute paths for filesystem operations; use project-relative paths for references in documentation (文件系统操作用绝对路径；文档内引用用项目相对路径)
- ERROR on gate failures or unresolved clarifications (门禁失败或存在未解决澄清项时报错)

## 完成条件(Done When)

- [ ] Plan workflow executed and design artifacts generated (规划工作流已执行且设计产物已生成)
- [ ] Extension hooks dispatched or skipped according to the rules in Mandatory Post-Execution Hooks above (扩展钩子已按上述规则派发或跳过)
- [ ] Completion reported to user with branch, plan path, and generated artifacts (已向用户报告分支、计划路径与生成产物)
