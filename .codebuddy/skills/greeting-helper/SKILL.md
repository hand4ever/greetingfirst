---
name: greeting-helper
description: 本技能是 Greeting 项目的规范校验助手。当用户在 greeting 项目中完成改动、准备提交或执行 /speckit.plan 之前，需要按项目宪法（八项核心原则 + 开发流程）做合规自检时使用；亦可用于代码评审、PR 前核对分层架构、统一响应、changelog 强制登记、英文注释、fail-fast 等约束。
---

# Greeting Helper（问候项目 · 规范校验助手）

## 概述(Overview)

本技能用于在 Greeting 项目发生改动后，依据项目宪法（`.specify/memory/constitution.md`，八项核心原则）与开发流程规范，对改动做结构化合规自检，逐项核对并输出"通过 / 不通过 + 整改建议"的报告。

## 适用场景(When to Use)

- 完成功能开发、准备提交（commit）前
- 执行 `/speckit.plan` 的宪法检查(Constitution Check)门禁前
- 代码评审 / 开 PR 前的自检
- 用户明确要求"按宪法核对改动""检查是否合规"

## 校验工作流(Workflow)

1. **加载宪法**：读取 `.specify/memory/constitution.md`，确认八项原则与开发流程的当前要求（注意版本号与同步影响报告）。
2. **界定范围**：确认本次改动涉及的文件清单（git diff / 用户说明），只核对受影响的原则，避免无谓全量扫描。
3. **逐项核对**：依据 `references/constitution_checklist.md` 的清单，对受影响文件逐条检查；每条给出 通过 / 不通过。
4. **定位问题**：对不合规项，给出 文件:行号 + 违反的原则 + 具体整改建议（引用宪法原文更优）。
5. **输出报告**：汇总为合规报告，包含：通过项计数、不通过项清单（含位置与建议）、整体结论（Ready / Needs Fix）、复杂度追踪中需说明的例外理由（如违反原则时的替代方案）。

## 资源(Resources)

### references/constitution_checklist.md

八项原则 + 开发流程的**可勾选核对清单**，含每条的具体检查点与判定标准。执行第 3 步时按需载入。

---

本技能为"参考/规范"型，主要依赖 SKILL.md 工作流与 references 清单，无需 scripts/ 与 assets/。如已无占位示例，可删除对应空目录。
