---
description: Assess the current codebase against the feature's spec, plan, and tasks, then append any remaining unbuilt work as new tasks to tasks.md so implement can complete it.
---

## 用户输入(User Input)

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## 执行前检查(Pre-Execution Checks)

**Check for extension hooks (before convergence / 收敛前的扩展钩子检查)**:

- Check if `.specify/extensions.yml` exists in the project root.
- If it exists, read it and look for entries under the `hooks.before_converge` key
- If the YAML cannot be parsed or is invalid, skip hook checking silently and continue normally
- Filter out hooks where `enabled` is explicitly `false`. Treat hooks without an `enabled` field as enabled by default.
- For each remaining hook, do **not** attempt to interpret or evaluate hook `condition` expressions:
  - If the hook has no `condition` field, or it is null/empty, treat it as executable
  - If the hook defines a non-empty `condition`, skip the hook and leave condition evaluation to the HookExecutor implementation
- For each executable hook, output the following based on its `optional` flag:
  - **Optional hook** (`optional: true` / 可选钩子):

    ```text
    ## 扩展钩子(Extension Hooks)

    **Optional Pre-Hook**: {extension}
    Command: `/{command}`
    Description: {description}

    Prompt: {prompt}
    To execute: `/{command}`
    ```

  - **Mandatory hook** (`optional: false` / 强制钩子):

    ```text
    ## 扩展钩子(Extension Hooks)

    **Automatic Pre-Hook**: {extension}
    Executing: `/{command}`
    EXECUTE_COMMAND: {command}

    Wait for the result of the hook command before proceeding to the Goal.
    ```
    After emitting the block above you MUST actually invoke the hook and wait for it to finish before continuing. Run it the same way you would run the command yourself in this agent/session (the invocation may differ from the literal `{command}` id shown above, e.g. a skills-mode agent runs it as `/skill:speckit-...` or `$speckit-...`). Emitting the block alone does not run the hook.

- If no hooks are registered or `.specify/extensions.yml` does not exist, skip silently

## 目标(Goal)

Close the gap between what a feature's specification, plan, and tasks call for and what the
codebase currently implements. Read `spec.md`, `plan.md`, and `tasks.md` as the **sole
source of intent** (with the constitution as governing constraints), assess the current
state of the code, determine which requirements, acceptance criteria, plan decisions, and
existing tasks are unmet, incomplete, or only partially satisfied, and **append each piece
of remaining work as a new, traceable task** at the bottom of `tasks.md` so that
`/speckit.implement` can complete it. This command MUST run only after
`/speckit.implement` has run on the current `tasks.md`, and after `/speckit.tasks` has produced a complete `tasks.md`.

This is **not** a diff tool and does **not** track changes. It assesses the present state
of the code relative to the feature's artifacts — no git, no branch comparison, no history.

## 运行约束(Operating Constraints)

**APPEND-ONLY, NEVER REWRITE** (仅追加，绝不重写): The command's **only** write is appending a new
`## Phase N: Convergence` (## 阶段 N：收敛) section to `tasks.md`. It MUST NOT:

- modify `spec.md` or `plan.md` in any way;
- rewrite, renumber, reorder, or delete any existing task (including tasks from a prior
  Convergence phase);
- modify, create, or delete any application code — completing the appended tasks is the
  job of `/speckit.implement`.

When the codebase already satisfies everything, the command MUST leave `tasks.md`
**byte-for-byte unchanged** (no empty Convergence header / 不写空的收敛标题) and report a clean result.

**Constitution Authority** (宪法权威): The project constitution (`.specify/memory/constitution.md`) is
**non-negotiable**. Code that violates a MUST principle is the highest-severity finding and
produces a corresponding remediation task. If the constitution is an unfilled template,
skip constitution checks gracefully rather than failing.

## 执行步骤(Execution Steps)

### 1. 初始化收敛上下文(Initialize Convergence Context)

Run `.specify/scripts/bash/check-prerequisites.sh --json --require-tasks --include-tasks` once from repo root and parse JSON for FEATURE_DIR and AVAILABLE_DOCS. Derive absolute paths:

- SPEC = FEATURE_DIR/spec.md
- PLAN = FEATURE_DIR/plan.md
- TASKS = FEATURE_DIR/tasks.md
- CONSTITUTION = `.specify/memory/constitution.md` (if present)
If `spec.md`, `plan.md`, or `tasks.md` is missing, STOP with a clear, actionable message naming the
prerequisite command to run (`/speckit.specify` for a missing spec, `/speckit.plan` for a missing plan,
`/speckit.tasks` for missing tasks). Do not produce partial output.
For single quotes in args like "I'm Groot", use escape syntax: e.g 'I'\''m Groot' (or double-quote if possible: "I'm Groot").

### 2. 加载产物(Load Artifacts, 渐进式披露)

Load only the minimal necessary context from each artifact:

**From spec.md (来自规格文件):**
- Functional Requirements (FR-### / 功能需求)
- Success Criteria (SC-###) — include only items requiring buildable work; exclude
  post-launch outcome metrics and business KPIs
- User Stories and their Acceptance Scenarios (用户故事及其验收场景)
- Edge Cases (边界情况, if present)

**From plan.md (来自计划文件):**
- Architecture/stack choices and technical decisions (架构/技术栈选择与技术决策)
- Data Model references (数据模型引用)
- Phases and named touch-points (files/components the plan says will be created or edited / 阶段与命名接触点)
- Technical constraints (技术约束)

**From tasks.md (来自任务文件):**
- Task IDs (任务 ID, to compute the next ID and next phase number)
- Descriptions, phase grouping, and referenced file paths (描述、阶段分组、引用文件路径)

**From constitution (来自宪法, if not an unfilled template):**
- Principle names and MUST/SHOULD normative statements (原则名称与 MUST/SHOULD 规范性陈述)

### 3. 构建意图清单(Build the Intent Inventory)

Create an internal model (do not echo raw artifacts):
- **Requirements inventory** (需求清单): one stable key per FR-### / SC-### / user-story acceptance
  scenario (e.g. `US1/AC2`), plus the plan decisions and constitution principles that
  impose buildable obligations.
- **Code-scope map** (代码范围映射): from the file paths named in `plan.md` and `tasks.md`, plus a keyword
  search for the concepts each requirement describes, derive the set of source files and
  components in scope for assessment. Bound the assessment to these — do **not** infer
  scope beyond what the artifacts define.

### 4. 评估代码库并分类发现(Assess the Codebase and Classify Findings)

For each item in the intent inventory, inspect the current code in scope and produce a
`Finding` only where there is a gap. Classify every finding by **gap type** (缺口类型):

- **`missing`** (缺失): the required work is absent from the code entirely.
- **`partial`** (部分): the work exists but does not yet fully satisfy the requirement /
  acceptance criterion / plan decision.
- **`contradicts`** (冲突): the code does something that conflicts with stated intent or a
  constitution MUST principle.
- **`unrequested`** (未请求): the code contains work not called for by the spec, plan, or tasks
  (surfaced for awareness — converge does **not** delete code, it only appends a task to
  review/justify or remove it).

Each `Finding` records: a stable id, the `source-ref` it traces to, the `gap-type`, a
severity, and a short human-readable description with the evidence (the file/area observed).

**Edge cases** (边界情况):
- **Little or no code yet** (代码极少或尚无): treat the entire specified scope as `missing` remaining work
  rather than failing.
- **Nothing remains** (无剩余): produce zero findings and follow the converged branch in Step 7.

### 5. 分配严重度(Assign Severity)

- **CRITICAL** (严重): violates a constitution MUST principle, or a `missing`/`contradicts` gap
  that blocks baseline functionality of a P1 user story.
- **HIGH** (高): a `missing` or `partial` gap on a core functional requirement or acceptance
  criterion.
- **MEDIUM** (中): a `partial` gap on a secondary requirement, or an `unrequested` addition with
  unclear justification.
- **LOW** (低): minor partial gaps, polish, or low-risk `unrequested` additions.

### 6. 呈现会话内发现摘要(Present the In-Session Findings Summary)

Before appending anything, output a compact, severity-graded summary (no file writes yet):

## 收敛发现(Convergence Findings)

| ID | Gap Type | Severity | Source | Evidence | Remaining Work |
|----|----------|----------|--------|----------|----------------|
| F1 | missing  | HIGH     | FR-008 | Example: no append-only guard detected in path/to/module.py when writing tasks.md | Add append-only enforcement |

**Summary metrics** (汇总指标):
- Requirements / acceptance criteria checked (已检查的需求 / 验收标准)
- Plan decisions checked (已检查的计划决策)
- Constitution principles checked (or "skipped — template" / 已检查的宪法原则，或"跳过——模板")
- Findings by gap type (missing / partial / contradicts / unrequested / 按缺口类型)
- Findings by severity (按严重度)

### 7. 追加收敛任务(或报告已收敛)(Append Convergence Tasks or report converged)

**If there are one or more actionable findings** (`tasks_appended` outcome / 存在可操作发现):

Append to the **end** of `tasks.md`, per the append contract:
1. Scan all existing task IDs; let `M` be the maximum. Determine the next phase number `N`
   (highest existing phase + 1).
2. Write a single new section header `## Phase N: Convergence` (## 阶段 N：收敛).
3. Emit one checklist item per actionable finding, ordered CRITICAL/HIGH first, assigning
   zero-padded IDs `T{M+1:03d}, T{M+2:03d}, …`:

   ```markdown
   - [ ] T042 <imperative description> per <source-ref> (<gap-type>)
   ```

   `<source-ref>` traces the task to its origin: e.g. `FR-003`, `SC-002`,
   `US1/AC2`, `plan: storage decision`, `Constitution II`.

   `<gap-type>` is one of `missing`, `partial`, `contradicts`, `unrequested`.

   Constitution-violation tasks MUST be emitted first and described as
   `CRITICAL`.
4. Never reuse or renumber existing IDs. If a prior Convergence phase exists, add a new,
   separately-numbered one below it — do not touch the old one.

**If there are no actionable findings** (`converged` outcome / 无可操作发现):
- Do **not** modify `tasks.md` at all — no empty phase header.
- Report: **"✅ Converged — the implementation satisfies the spec, plan, and tasks."** (✅ 已收敛——实现满足规格、计划与任务)
- Include the summary counts of what was checked.

### 8. 提供下一步动作(交接)(Provide Next Actions / Handoff)

- On `tasks_appended`: state how many tasks were appended under which phase, and recommend
  running `/speckit.implement` to complete them; note that a follow-up converge
  run will find fewer or no remaining items.
- On `converged`: recommend proceeding to review / opening a PR. No further implement pass
  is needed for this feature's specified scope.

### 9. 检查扩展钩子(Check for extension hooks)

After producing the result, check if `.specify/extensions.yml` exists in the project root.

- If it exists, read it and look for entries under the `hooks.after_converge` key
- If the YAML cannot be parsed or is invalid, skip hook checking silently and continue normally
- Filter out hooks where `enabled` is explicitly `false`. Treat hooks without an `enabled` field as enabled by default.
- For each remaining hook, do **not** attempt to interpret or evaluate hook `condition` expressions:
  - If the hook has no `condition` field, or it is null/empty, treat it as executable
  - If the hook defines a non-empty `condition`, skip the hook and leave condition evaluation to the HookExecutor implementation
- Report the convergence outcome (`converged` or `tasks_appended`) in-session before listing
  any hooks, so users can decide whether to run optional follow-up commands.
- For each executable hook, output the following based on its `optional` flag:
  - **Optional hook** (`optional: true` / 可选钩子):

    ```text
    ## 扩展钩子(Extension Hooks)

    **Optional Hook**: {extension}
    Command: `/{command}`
    Description: {description}

    Prompt: {prompt}
    To execute: `/{command}`
    ```

  - **Mandatory hook** (`optional: false` / 强制钩子):

    ```text
    ## 扩展钩子(Extension Hooks)

    **Automatic Hook**: {extension}
    Executing: `/{command}`
    EXECUTE_COMMAND: {command}
    ```
    After emitting the block above you MUST actually invoke the hook and wait for it to finish before continuing. Run it the same way you would run the command yourself in this agent/session (the invocation may differ from the literal `{command}` id shown above, e.g. a skills-mode agent runs it as `/skill:speckit-...` or `$speckit-...`). Emitting the block alone does not run the hook.

- If no hooks are registered or `.specify/extensions.yml` does not exist, skip silently
