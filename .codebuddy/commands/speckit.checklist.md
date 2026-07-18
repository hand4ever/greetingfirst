---
description: Generate a custom checklist for the current feature based on user requirements.
---

## 检查清单目的(Checklist Purpose)：「需求编写的单元测试」("Unit Tests for English")

**CRITICAL CONCEPT** (核心概念): Checklists are **UNIT TESTS FOR REQUIREMENTS WRITING** (需求编写的单元测试) - they validate the quality, clarity, and completeness of requirements in a given domain.

**NOT for verification/testing** (不用于验证/测试):
- ❌ NOT "Verify the button clicks correctly"
- ❌ NOT "Test error handling works"
- ❌ NOT "Confirm the API returns 200"
- ❌ NOT checking if code/implementation matches the spec

**FOR requirements quality validation** (用于需求质量校验):
- ✅ "Are visual hierarchy requirements defined for all card types?" (completeness / 完整性)
- ✅ "Is 'prominent display' quantified with specific sizing/positioning?" (clarity / 清晰度)
- ✅ "Are hover state requirements consistent across all interactive elements?" (consistency / 一致性)
- ✅ "Are accessibility requirements defined for keyboard navigation?" (coverage / 覆盖度)
- ✅ "Does the spec define what happens when logo image fails to load?" (edge cases / 边界情况)

**Metaphor** (隐喻): If your spec is code written in English, the checklist is its unit test suite. You're testing whether the requirements are well-written, complete, unambiguous, and ready for implementation - NOT whether the implementation works.

## 用户输入(User Input)

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## 执行前检查(Pre-Execution Checks)

**Check for extension hooks (before checklist generation / 检查清单生成前的扩展钩子检查)**:
- Check if `.specify/extensions.yml` exists in the project root.
- If it exists, read it and look for entries under the `hooks.before_checklist` key
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

    Wait for the result of the hook command before proceeding to the Execution Steps.
    ```
    After emitting the block above you MUST actually invoke the hook and wait for it to finish before continuing. Run it the same way you would run the command yourself in this agent/session (the invocation may differ from the literal `{command}` id shown above, e.g. a skills-mode agent runs it as `/skill:speckit-...` or `$speckit-...`). Emitting the block alone does not run the hook.
- If no hooks are registered or `.specify/extensions.yml` does not exist, skip silently

## 执行步骤(Execution Steps)

1. **Setup** (初始化): Run `.specify/scripts/bash/check-prerequisites.sh --json` from repo root and parse JSON for FEATURE_DIR and AVAILABLE_DOCS list.
   - All file paths must be absolute.
   - For single quotes in args like "I'm Groot", use escape syntax: e.g 'I'\''m Groot' (or double-quote if possible: "I'm Groot").

2. **IF EXISTS**: Load `.specify/memory/constitution.md` for project principles and governance constraints.

3. **Clarify intent (dynamic)** (澄清意图，动态): Derive up to THREE initial contextual clarifying questions (no pre-baked catalog). They MUST:
   - Be generated from the user's phrasing + extracted signals from spec/plan/tasks
   - Only ask about information that materially changes checklist content
   - Be skipped individually if already unambiguous in `$ARGUMENTS`
   - Prefer precision over breadth

   **Generation algorithm** (生成算法):
   1. Extract signals: feature domain keywords (e.g., auth, latency, UX, API), risk indicators ("critical", "must", "compliance"), stakeholder hints ("QA", "review", "security team"), and explicit deliverables ("a11y", "rollback", "contracts").
   2. Cluster signals into candidate focus areas (max 4) ranked by relevance.
   3. Identify probable audience & timing (author, reviewer, QA, release) if not explicit.
   4. Detect missing dimensions: scope breadth, depth/rigor, risk emphasis, exclusion boundaries, measurable acceptance criteria.
   5. Formulate questions chosen from these archetypes:
      - Scope refinement (范围细化, e.g., "Should this include integration touchpoints with X and Y or stay limited to local module correctness?")
      - Risk prioritization (风险优先级, e.g., "Which of these potential risk areas should receive mandatory gating checks?")
      - Depth calibration (深度校准, e.g., "Is this a lightweight pre-commit sanity list or a formal release gate?")
      - Audience framing (受众定位, e.g., "Will this be used by the author only or peers during PR review?")
      - Boundary exclusion (边界排除, e.g., "Should we explicitly exclude performance tuning items this round?")
      - Scenario class gap (场景类别缺口, e.g., "No recovery flows detected—are rollback / partial failure paths in scope?")

   **Question formatting rules** (问题格式规则):
   - If presenting options, generate a compact table with columns: Option | Candidate | Why It Matters
   - Limit to A–E options maximum; omit table if a free-form answer is clearer
   - Never ask the user to restate what they already said
   - Avoid speculative categories (no hallucination). If uncertain, ask explicitly: "Confirm whether X belongs in scope."

   **Defaults when interaction impossible** (无法交互时的默认值):
   - Depth: Standard (标准)
   - Audience: Reviewer (PR) if code-related; Author otherwise
   - Focus: Top 2 relevance clusters

   Output the questions (label Q1/Q2/Q3). After answers: if ≥2 scenario classes (Alternate / Exception / Recovery / Non-Functional domain) remain unclear, you MAY ask up to TWO more targeted follow‑ups (Q4/Q5) with a one-line justification each (e.g., "Unresolved recovery path risk"). Do not exceed five total questions. Skip escalation if user explicitly declines more.

4. **Understand user request** (理解用户请求): Combine `$ARGUMENTS` + clarifying answers:
   - Derive checklist theme (e.g., security, review, deploy, ux / 如安全、评审、部署、用户体验)
   - Consolidate explicit must-have items mentioned by user
   - Map focus selections to category scaffolding
   - Infer any missing context from spec/plan/tasks (do NOT hallucinate)

5. **Load feature context** (加载功能上下文): Read from FEATURE_DIR:
   - spec.md: Feature requirements and scope (功能需求与范围)
   - plan.md (if exists): Technical details, dependencies (技术细节、依赖)
   - tasks.md (if exists): Implementation tasks (实现任务)

   **Context Loading Strategy** (上下文加载策略):
   - Load only necessary portions relevant to active focus areas (avoid full-file dumping)
   - Prefer summarizing long sections into concise scenario/requirement bullets
   - Use progressive disclosure: add follow-on retrieval only if gaps detected
   - If source docs are large, generate interim summary items instead of embedding raw text

6. **Generate checklist** (生成检查清单) - Create "Unit Tests for Requirements" (需求单元测试):
   - Create `FEATURE_DIR/checklists/` directory if it doesn't exist
   - Generate unique checklist filename:
     - Use short, descriptive name based on domain (e.g., `ux.md`, `api.md`, `security.md`)
     - Format: `[domain].md`
   - File handling behavior:
     - If file does NOT exist: Create new file and number items starting from CHK001
     - If file exists: Append new items to existing file, continuing from the last CHK ID (e.g., if last item is CHK015, start new items at CHK016)
   - Never delete or replace existing checklist content - always preserve and append

   **CORE PRINCIPLE** (核心原则) - Test the Requirements, Not the Implementation (测试需求，而非实现):
   Every checklist item MUST evaluate the REQUIREMENTS THEMSELVES for:
   - **Completeness** (完整性): Are all necessary requirements present?
   - **Clarity** (清晰度): Are requirements unambiguous and specific?
   - **Consistency** (一致性): Do requirements align with each other?
   - **Measurability** (可衡量性): Can requirements be objectively verified?
   - **Coverage** (覆盖度): Are all scenarios/edge cases addressed?

   **Category Structure** (类别结构) - Group items by requirement quality dimensions:
   - **Requirement Completeness** (需求完整性, Are all necessary requirements documented?)
   - **Requirement Clarity** (需求清晰度, Are requirements specific and unambiguous?)
   - **Requirement Consistency** (需求一致性, Do requirements align without conflicts?)
   - **Acceptance Criteria Quality** (验收标准质量, Are success criteria measurable?)
   - **Scenario Coverage** (场景覆盖, Are all flows/cases addressed?)
   - **Edge Case Coverage** (边界情况覆盖, Are boundary conditions defined?)
   - **Non-Functional Requirements** (非功能需求, Performance, Security, Accessibility, etc. - are they specified?)
   - **Dependencies & Assumptions** (依赖与假设, Are they documented and validated?)
   - **Ambiguities & Conflicts** (歧义与冲突, What needs clarification?)

   **HOW TO WRITE CHECKLIST ITEMS** (如何编写检查项) - "Unit Tests for English" (需求编写的单元测试):

   ❌ **WRONG** (错误, Testing implementation / 测试实现):
   - "Verify landing page displays 3 episode cards"
   - "Test hover states work on desktop"
   - "Confirm logo click navigates home"

   ✅ **CORRECT** (正确, Testing requirements quality / 测试需求质量):
   - "Are the exact number and layout of featured episodes specified?" [Completeness]
   - "Is 'prominent display' quantified with specific sizing/positioning?" [Clarity]
   - "Are hover state requirements consistent across all interactive elements?" [Consistency]
   - "Are keyboard navigation requirements defined for all interactive UI?" [Coverage]
   - "Is the fallback behavior specified when logo image fails to load?" [Edge Cases]
   - "Are loading states defined for asynchronous episode data?" [Completeness]
   - "Does the spec define visual hierarchy for competing UI elements?" [Clarity]

   **ITEM STRUCTURE** (条目结构):
   Each item should follow this pattern:
   - Question format asking about requirement quality (询问需求质量的问题格式)
   - Focus on what's WRITTEN (or not written) in the spec/plan
   - Include quality dimension in brackets [Completeness/Clarity/Consistency/etc.]
   - Reference spec section `[Spec §X.Y]` when checking existing requirements
   - Use `[Gap]` marker when checking for missing requirements

   **EXAMPLES BY QUALITY DIMENSION** (按质量维度的示例):

   Completeness (完整性):
   - "Are error handling requirements defined for all API failure modes? [Gap]"
   - "Are accessibility requirements specified for all interactive elements? [Completeness]"
   - "Are mobile breakpoint requirements defined for responsive layouts? [Gap]"

   Clarity (清晰度):
   - "Is 'fast loading' quantified with specific timing thresholds? [Clarity, Spec §NFR-2]"
   - "Are 'related episodes' selection criteria explicitly defined? [Clarity, Spec §FR-5]"
   - "Is 'prominent' defined with measurable visual properties? [Ambiguity, Spec §FR-4]"

   Consistency (一致性):
   - "Do navigation requirements align across all pages? [Consistency, Spec §FR-10]"
   - "Are card component requirements consistent between landing and detail pages? [Consistency]"

   Coverage (覆盖度):
   - "Are requirements defined for zero-state scenarios (no episodes)? [Coverage, Edge Case]"
   - "Are concurrent user interaction scenarios addressed? [Coverage, Gap]"
   - "Are requirements specified for partial data loading failures? [Coverage, Exception Flow]"

   Measurability (可衡量性):
   - "Are visual hierarchy requirements measurable/testable? [Acceptance Criteria, Spec §FR-1]"
   - "Can 'balanced visual weight' be objectively verified? [Measurability, Spec §FR-2]"

   **Scenario Classification & Coverage** (场景分类与覆盖, Requirements Quality Focus):
   - Check if requirements exist for: Primary, Alternate, Exception/Error, Recovery, Non-Functional scenarios
   - For each scenario class, ask: "Are [scenario type] requirements complete, clear, and consistent?"
   - If scenario class missing: "Are [scenario type] requirements intentionally excluded or missing? [Gap]"
   - Include resilience/rollback when state mutation occurs: "Are rollback requirements defined for migration failures? [Gap]"

   **Traceability Requirements** (可追溯性要求):
   - MINIMUM: ≥80% of items MUST include at least one traceability reference
   - Each item should reference: spec section `[Spec §X.Y]`, or use markers: `[Gap]`, `[Ambiguity]`, `[Conflict]`, `[Assumption]`
   - If no ID system exists: "Is a requirement & acceptance criteria ID scheme established? [Traceability]"

   **Surface & Resolve Issues** (暴露并解决问题, Requirements Quality Problems):
   Ask questions about the requirements themselves:
   - Ambiguities: "Is the term 'fast' quantified with specific metrics? [Ambiguity, Spec §NFR-1]"
   - Conflicts: "Do navigation requirements conflict between §FR-10 and §FR-10a? [Conflict]"
   - Assumptions: "Is the assumption of 'always available podcast API' validated? [Assumption]"
   - Dependencies: "Are external podcast API requirements documented? [Dependency, Gap]"
   - Missing definitions: "Is 'visual hierarchy' defined with measurable criteria? [Gap]"

   **Content Consolidation** (内容整合):
   - Soft cap: If raw candidate items > 40, prioritize by risk/impact
   - Merge near-duplicates checking the same requirement aspect
   - If >5 low-impact edge cases, create one item: "Are edge cases X, Y, Z addressed in requirements? [Coverage]"

   **🚫 ABSOLUTELY PROHIBITED** (绝对禁止) - These make it an implementation test, not a requirements test (这些会使其成为实现测试而非需求测试):
   - ❌ Any item starting with "Verify", "Test", "Confirm", "Check" + implementation behavior
   - ❌ References to code execution, user actions, system behavior
   - ❌ "Displays correctly", "works properly", "functions as expected"
   - ❌ "Click", "navigate", "render", "load", "execute"
   - ❌ Test cases, test plans, QA procedures
   - ❌ Implementation details (frameworks, APIs, algorithms)

   **✅ REQUIRED PATTERNS** (必用模式) - These test requirements quality (这些测试需求质量):
   - ✅ "Are [requirement type] defined/specified/documented for [scenario]?"
   - ✅ "Is [vague term] quantified/clarified with specific criteria?"
   - ✅ "Are requirements consistent between [section A] and [section B]?"
   - ✅ "Can [requirement] be objectively measured/verified?"
   - ✅ "Are [edge cases/scenarios] addressed in requirements?"
   - ✅ "Does the spec define [missing aspect]?"

7. **Structure Reference** (结构参考): Generate the checklist following the canonical template in `.specify/templates/checklist-template.md` for title, meta section, category headings, and ID formatting. If template is unavailable, use: H1 title, purpose/created meta lines, `##` category sections containing `- [ ] CHK### <requirement item>` lines with globally incrementing IDs starting at CHK001.

8. **Report** (报告): Output full path to checklist file, item count, and summarize whether the run created a new file or appended to an existing one. Summarize:
   - Focus areas selected (所选关注领域)
   - Depth level (深度级别)
   - Actor/timing (执行者/时机)
   - Any explicit user-specified must-have items incorporated (纳入的用户明确必含项)

**Important** (重要): Each `/speckit.checklist` command invocation uses a short, descriptive checklist filename and either creates a new file or appends to an existing one. This allows:
- Multiple checklists of different types (e.g., `ux.md`, `test.md`, `security.md`)
- Simple, memorable filenames that indicate checklist purpose
- Easy identification and navigation in the `checklists/` folder

To avoid clutter, use descriptive types and clean up obsolete checklists when done.

## 检查清单类型与示例(Example Checklist Types & Sample Items)

**UX Requirements Quality** (UX 需求质量): `ux.md`

Sample items (testing the requirements, NOT the implementation / 测试需求而非实现):
- "Are visual hierarchy requirements defined with measurable criteria? [Clarity, Spec §FR-1]"
- "Is the number and positioning of UI elements explicitly specified? [Completeness, Spec §FR-1]"
- "Are interaction state requirements (hover, focus, active) consistently defined? [Consistency]"
- "Are accessibility requirements specified for all interactive elements? [Coverage, Gap]"
- "Is fallback behavior defined when images fail to load? [Edge Case, Gap]"
- "Can 'prominent display' be objectively measured? [Measurability, Spec §FR-4]"

**API Requirements Quality** (API 需求质量): `api.md`

Sample items:
- "Are error response formats specified for all failure scenarios? [Completeness]"
- "Are rate limiting requirements quantified with specific thresholds? [Clarity]"
- "Are authentication requirements consistent across all endpoints? [Consistency]"
- "Are retry/timeout requirements defined for external dependencies? [Coverage, Gap]"
- "Is versioning strategy documented in requirements? [Gap]"

**Performance Requirements Quality** (性能需求质量): `performance.md`

Sample items:
- "Are performance requirements quantified with specific metrics? [Clarity]"
- "Are performance targets defined for all critical user journeys? [Coverage]"
- "Are performance requirements under different load conditions specified? [Completeness]"
- "Can performance requirements be objectively measured? [Measurability]"
- "Are degradation requirements defined for high-load scenarios? [Edge Case, Gap]"

**Security Requirements Quality** (安全需求质量): `security.md`

Sample items:
- "Are authentication requirements specified for all protected resources? [Coverage]"
- "Are data protection requirements defined for sensitive information? [Completeness]"
- "Is the threat model documented and requirements aligned to it? [Traceability]"
- "Are security requirements consistent with compliance obligations? [Consistency]"
- "Are security failure/breach response requirements defined? [Gap, Exception Flow]"

## 反例：不该做什么(Anti-Examples: What NOT To Do)

**❌ WRONG** (错误) - These test implementation, not requirements (这些测试实现而非需求):

```markdown
- [ ] CHK001 - Verify landing page displays 3 episode cards [Spec §FR-001]
- [ ] CHK002 - Test hover states work correctly on desktop [Spec §FR-003]
- [ ] CHK003 - Confirm logo click navigates to home page [Spec §FR-010]
- [ ] CHK004 - Check that related episodes section shows 3-5 items [Spec §FR-005]
```

**✅ CORRECT** (正确) - These test requirements quality (这些测试需求质量):

```markdown
- [ ] CHK001 - Are the number and layout of featured episodes explicitly specified? [Completeness, Spec §FR-001]
- [ ] CHK002 - Are hover state requirements consistently defined for all interactive elements? [Consistency, Spec §FR-003]
- [ ] CHK003 - Are navigation requirements clear for all clickable brand elements? [Clarity, Spec §FR-010]
- [ ] CHK004 - Is the selection criteria for related episodes documented? [Gap, Spec §FR-005]
- [ ] CHK005 - Are loading state requirements defined for asynchronous episode data? [Gap]
- [ ] CHK006 - Can "visual hierarchy" requirements be objectively measured? [Measurability, Spec §FR-001]
```

**Key Differences** (关键区别):
- Wrong: Tests if the system works correctly (系统是否正确工作)
- Correct: Tests if the requirements are written correctly (需求是否编写正确)
- Wrong: Verification of behavior (行为验证)
- Correct: Validation of requirement quality (需求质量校验)
- Wrong: "Does it do X?" (它是否做了 X？)
- Correct: "Is X clearly specified?" (X 是否被清晰指定？)

## 后置执行检查(Post-Execution Checks)

**Check for extension hooks (after checklist generation / 检查清单生成后的扩展钩子检查)**:
Check if `.specify/extensions.yml` exists in the project root.
- If it exists, read it and look for entries under the `hooks.after_checklist` key
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
