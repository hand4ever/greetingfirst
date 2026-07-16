# 实施计划(Implementation Plan): [功能名称]

**分支(Branch)**: `[###-feature-name]` | **日期(Date)**: [DATE] | **规格(Spec)**: [link]

**输入(Input)**: Feature specification from `/specs/[###-feature-name]/spec.md`

**说明(Note)**: This template is filled in by the `/speckit.plan` command; its definition describes the execution workflow.

## 概述(Summary)

[从功能规格中提取：主要需求 + 研究中的技术方案]

## 技术上下文(Technical Context)

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**语言/版本(Language/Version)**: [例如：Python 3.11、Swift 5.9、Rust 1.75 或 需要澄清(NEEDS CLARIFICATION)]

**主要依赖(Primary Dependencies)**: [例如：FastAPI、UIKit、LLVM 或 需要澄清(NEEDS CLARIFICATION)]

**存储(Storage)**: [如适用，例如：PostgreSQL、CoreData、文件 或 N/A]

**测试框架(Testing)**: [例如：pytest、XCTest、cargo test 或 需要澄清(NEEDS CLARIFICATION)]

**目标平台(Target Platform)**: [例如：Linux 服务器、iOS 15+、WASM 或 需要澄清(NEEDS CLARIFICATION)]

**项目类型(Project Type)**: [例如：library/cli/web-service/mobile-app/compiler/desktop-app 或 需要澄清(NEEDS CLARIFICATION)]

**性能目标(Performance Goals)**: [领域相关，例如：1000 req/s、10k lines/sec、60 fps 或 需要澄清(NEEDS CLARIFICATION)]

**约束(Constraints)**: [领域相关，例如：<200ms p95、<100MB 内存、离线可用 或 需要澄清(NEEDS CLARIFICATION)]

**规模/范围(Scale/Scope)**: [领域相关，例如：10k 用户、1M 行代码、50 个页面 或 需要澄清(NEEDS CLARIFICATION)]

## 宪法检查(Constitution Check)

*门禁(GATE): Must pass before Phase 0 research. Re-check after Phase 1 design.*

[Gates determined based on constitution file]

## 项目结构(Project Structure)

### 文档(Documentation)（本特性）

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### 源代码(Source Code)（仓库根目录）
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
# [如未使用请删除] Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# [如未使用请删除] Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# [如未使用请删除] Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure: feature modules, UI flows, platform tests]
```

**结构决策(Structure Decision)**: [记录选定的结构并引用上面捕获的真实目录]

## 复杂度追踪(Complexity Tracking)

> **仅当宪法检查有违规且需要说明理由时填写**

| 违规(Violation) | 必要性(Why Needed) | 拒绝更简单方案的原因(Simpler Alternative Rejected Because) |
|-----------|------------|-------------------------------------|
| [例如：第4个项目] | [当前需求] | [为什么3个项目不够] |
| [例如：Repository 模式] | [具体问题] | [为什么直接访问数据库不够] |
