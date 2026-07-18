---
description: Perform a non-destructive cross-artifact consistency and quality analysis across spec.md, plan.md, and tasks.md after task generation.
---

## 用户输入(User Input)

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## 执行前检查(Pre-Execution Checks)

**Check for extension hooks (before analysis / 分析前的扩展钩子检查)**:
- Check if `.specify/extensions.yml` exists in the project root.
- If it exists, read it and look for entries under the `hooks.before_analyze` key
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

    Wait for the result of the hook command before proceeding to the Goal.
    ```
    After emitting the block above you MUST actually invoke the hook and wait for it to finish before continuing. Run it the same way you would run the command yourself in this agent/session (the invocation may differ from the literal `{command}` id shown above, e.g. a skills-mode agent runs it as `/skill:speckit-...` or `$speckit-...`). Emitting the block alone does not run the hook.
- If no hooks are registered or `.specify/extensions.yml` does not exist, skip silently

## 目标(Goal)

Identify inconsistencies, duplications, ambiguities, and underspecified items across the three core artifacts (`spec.md`, `plan.md`, `tasks.md`) before implementation. This command MUST run only after `/speckit.tasks` has successfully produced a complete `tasks.md`.

## 运行约束(Operating Constraints)

**STRICTLY READ-ONLY** (严格只读): Do **not** modify any files. Output a structured analysis report. Offer an optional remediation plan (user must explicitly approve before any follow-up editing commands would be invoked manually).

**Constitution Authority** (宪法权威): The project constitution (`.specify/memory/constitution.md`) is **non-negotiable** within this analysis scope. Constitution conflicts are automatically CRITICAL and require adjustment of the spec, plan, or tasks—not dilution, reinterpretation, or silent ignoring of the principle. If a principle itself needs to change, that must occur in a separate, explicit constitution update outside `/speckit.analyze`.

## 执行步骤(Execution Steps)

### 1. 初始化分析上下文(Initialize Analysis Context)

Run `.specify/scripts/bash/check-prerequisites.sh --json --require-tasks --include-tasks` once from repo root and parse JSON for FEATURE_DIR and AVAILABLE_DOCS. Derive absolute paths:

- SPEC = FEATURE_DIR/spec.md
- PLAN = FEATURE_DIR/plan.md
- TASKS = FEATURE_DIR/tasks.md

Abort with an error message if any required file is missing (instruct the user to run missing prerequisite command).
For single quotes in args like "I'm Groot", use escape syntax: e.g 'I'\''m Groot' (or double-quote if possible: "I'm Groot").

### 2. 加载产物(Load Artifacts, 渐进式披露)

Load only the minimal necessary context from each artifact:

**From spec.md (来自规格文件):**
- Overview/Context (概览/上下文)
- Functional Requirements (功能需求)
- Success Criteria (success criteria — measurable outcomes — e.g., performance, security, availability, user success, business impact / 可衡量成果，如性能、安全、可用性、用户成功、业务影响)
- User Stories (用户故事)
- Edge Cases (边界情况, if present)

**From plan.md (来自计划文件):**
- Architecture/stack choices (架构/技术栈选择)
- Data Model references (数据模型引用)
- Phases (阶段)
- Technical constraints (技术约束)

**From tasks.md (来自任务文件):**
- Task IDs (任务 ID)
- Descriptions (描述)
- Phase grouping (阶段分组)
- Parallel markers [P] (并行标记)
- Referenced file paths (引用文件路径)

**From constitution (来自宪法):**
- Load `.specify/memory/constitution.md` for principle validation (原则校验)

### 3. 构建语义模型(Build Semantic Models)

Create internal representations (do not include raw artifacts in output):

- **Requirements inventory** (需求清单): For each Functional Requirement (FR-###) and Success Criterion (SC-###), record a stable key. Use the explicit FR-/SC- identifier as the primary key when present, and optionally also derive an imperative-phrase slug for readability (e.g., "User can upload file" → `user-can-upload-file`). Include only Success Criteria items that require buildable work (e.g., load-testing infrastructure, security audit tooling), and exclude post-launch outcome metrics and business KPIs (e.g., "Reduce support tickets by 50%").
- **User story/action inventory** (用户故事/动作清单): Discrete user actions with acceptance criteria
- **Task coverage mapping** (任务覆盖映射): Map each task to one or more requirements or stories (inference by keyword / explicit reference patterns like IDs or key phrases)
- **Constitution rule set** (宪法规则集): Extract principle names and MUST/SHOULD normative statements

### 4. 检测遍历(Detection Passes, 高效令牌分析)

Focus on high-signal findings. Limit to 50 findings total; aggregate remainder in overflow summary.

#### A. 重复检测(Duplication Detection)
- Identify near-duplicate requirements
- Mark lower-quality phrasing for consolidation

#### B. 歧义检测(Ambiguity Detection)
- Flag vague adjectives (fast, scalable, secure, intuitive, robust) lacking measurable criteria
- Flag unresolved placeholders (TODO, TKTK, ???, `<placeholder>`, etc.)

#### C. 规格不足(Underspecification)
- Requirements with verbs but missing object or measurable outcome
- User stories missing acceptance criteria alignment
- Tasks referencing files or components not defined in spec/plan

#### D. 宪法一致性(Constitution Alignment)
- Any requirement or plan element conflicting with a MUST principle
- Missing mandated sections or quality gates from constitution

#### E. 覆盖缺口(Coverage Gaps)
- Requirements with zero associated tasks
- Tasks with no mapped requirement/story
- Success Criteria requiring buildable work (performance, security, availability) not reflected in tasks

#### F. 不一致(Inconsistency)
- Terminology drift (same concept named differently across files / 术语漂移，同一概念在不同文件中命名不同)
- Data entities referenced in plan but absent in spec (or vice versa)
- Task ordering contradictions (e.g., integration tasks before foundational setup tasks without dependency note)
- Conflicting requirements (e.g., one requires Next.js while other specifies Vue)

### 5. 严重度分级(Severity Assignment)

Use this heuristic to prioritize findings:
- **CRITICAL** (严重): Violates constitution MUST, missing core spec artifact, or requirement with zero coverage that blocks baseline functionality
- **HIGH** (高): Duplicate or conflicting requirement, ambiguous security/performance attribute, untestable acceptance criterion
- **MEDIUM** (中): Terminology drift, missing non-functional task coverage, underspecified edge case
- **LOW** (低): Style/wording improvements, minor redundancy not affecting execution order

### 6. 生成紧凑分析报告(Produce Compact Analysis Report)

Output a Markdown report (no file writes) with the following structure:

## 规格分析报告(Specification Analysis Report)

| ID | Category | Severity | Location(s) | Summary | Recommendation |
|----|----------|----------|-------------|---------|----------------|
| A1 | Duplication | HIGH | spec.md:L120-134 | Two similar requirements ... | Merge phrasing; keep clearer version |

(Add one row per finding; generate stable IDs prefixed by category initial.)

**Coverage Summary Table (覆盖汇总表):**

| Requirement Key | Has Task? | Task IDs | Notes |
|-----------------|-----------|----------|-------|

**Constitution Alignment Issues (宪法一致性问题):** (if any)

**Unmapped Tasks (未映射任务):** (if any)

**Metrics (指标):**
- Total Requirements (需求总数)
- Total Tasks (任务总数)
- Coverage % (requirements with >=1 task / 覆盖率)
- Ambiguity Count (歧义数)
- Duplication Count (重复数)
- Critical Issues Count (严重问题数)

### 7. 提供下一步动作(Provide Next Actions)

At end of report, output a concise Next Actions block:
- If CRITICAL issues exist: Recommend resolving before `/speckit.implement`
- If only LOW/MEDIUM: User may proceed, but provide improvement suggestions
- Provide explicit command suggestions: e.g., "Run /speckit.specify with refinement", "Run /speckit.plan to adjust architecture", "Manually edit tasks.md to add coverage for 'performance-metrics'"

### 8. 提供修复建议(Offer Remediation)

Ask the user: "Would you like me to suggest concrete remediation edits for the top N issues?" (Do NOT apply them automatically.)

### 9. 检查扩展钩子(Check for extension hooks)

After reporting, check if `.specify/extensions.yml` exists in the project root.
- If it exists, read it and look for entries under the `hooks.after_analyze` key
- If the YAML cannot be parsed or is invalid, skip hook checking silently and continue normally
- Filter out hooks where `enabled` is explicitly `false`. Treat hooks without an `enabled` field as enabled by default.
- For each remaining hook, do **not** attempt to interpret or evaluate hook `condition` expressions:
  - If the hook has no `condition` field, or it is null/empty, treat it as executable
  - If the hook defines a non-empty `condition`, skip the hook and leave condition evaluation to the HookExecutor implementation
- For each executable hook, output the following based on its `optional` flag:
  - **Optional hook** (`optional: true` / 可选钩子):
    ```
    ## 扩展钩子(Extension Hooks)

    **Optional Hook**: {extension}
    Command: `/{command}`
    Description: {description}

    Prompt: {prompt}
    To execute: `/{command}`
    ```
  - **Mandatory hook** (`optional: false` / 强制钩子):
    ```
    ## 扩展钩子(Extension Hooks)

    **Automatic Hook**: {extension}
    Executing: `/{command}`
    EXECUTE_COMMAND: {command}
    ```
    After emitting the block above you MUST actually invoke the hook and wait for it to finish before continuing. Run it the same way you would run the command yourself in this agent/session (the invocation may differ from the literal `{command}` id shown above, e.g. a skills-mode agent runs it as `/skill:speckit-...` or `$speckit-...`). Emitting the block alone does not run the hook.
- If no hooks are registered or `.specify/extensions.yml` does not exist, skip silently

## 运行原则(Operating Principles)

### 上下文效率(Context Efficiency)
- **Minimal high-signal tokens** (最小高信号令牌): Focus on actionable findings, not exhaustive documentation
- **Progressive disclosure** (渐进式披露): Load artifacts incrementally; don't dump all content into analysis
- **Token-efficient output** (令牌高效输出): Limit findings table to 50 rows; summarize overflow
- **Deterministic results** (确定性结果): Rerunning without changes should produce consistent IDs and counts

### 分析指南(Analysis Guidelines)
- **NEVER modify files** (绝不修改文件, this is read-only analysis / 只读分析)
- **NEVER hallucinate missing sections** (绝不臆造缺失章节, if absent, report them accurately / 若缺失则准确报告)
- **Prioritize constitution violations** (优先处理宪法违规, these are always CRITICAL / 始终为严重级)
- **Use examples over exhaustive rules** (以示例代替穷举规则, cite specific instances, not generic patterns / 引用具体实例而非泛化模式)
- **Report zero issues gracefully** (无问题时优雅报告, emit success report with coverage statistics / 输出带覆盖统计的成功报告)

## 上下文(Context)

$ARGUMENTS
