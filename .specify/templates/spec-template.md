# 功能规格说明(Feature Specification): [功能名称]

**功能分支(Feature Branch)**: `[###-feature-name]`

**创建日期(Created)**: [DATE]

**状态(Status)**: 草稿(Draft)

**输入(Input)**: 用户描述(User description): "$ARGUMENTS"

## 用户场景与测试(User Scenarios & Testing) *(必填)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.

  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### 用户故事(User Story) 1 - [简要标题] (优先级(Priority): P1)

[用平实的语言描述这个用户旅程]

**优先级理由(Why this priority)**: [解释这个优先级的价值和理由]

**独立测试(Independent Test)**: [描述如何独立测试 - 例如："可以通过[具体操作]完整测试，并交付[具体价值]"]

**验收场景(Acceptance Scenarios)**:

1. **假设(Given)** [初始状态], **当(When)** [操作], **那么(Then)** [预期结果]
2. **假设(Given)** [初始状态], **当(When)** [操作], **那么(Then)** [预期结果]

---

### 用户故事(User Story) 2 - [简要标题] (优先级(Priority): P2)

[用平实的语言描述这个用户旅程]

**优先级理由(Why this priority)**: [解释这个优先级的价值和理由]

**独立测试(Independent Test)**: [描述如何独立测试]

**验收场景(Acceptance Scenarios)**:

1. **假设(Given)** [初始状态], **当(When)** [操作], **那么(Then)** [预期结果]

---

### 用户故事(User Story) 3 - [简要标题] (优先级(Priority): P3)

[用平实的语言描述这个用户旅程]

**优先级理由(Why this priority)**: [解释这个优先级的价值和理由]

**独立测试(Independent Test)**: [描述如何独立测试]

**验收场景(Acceptance Scenarios)**:

1. **假设(Given)** [初始状态], **当(When)** [操作], **那么(Then)** [预期结果]

---

[根据需要添加更多用户故事，每个都分配优先级]

### 边界情况(Edge Cases)

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right edge cases.
-->

- 当[边界条件]发生时，会发生什么？
- 系统如何处理[错误场景]？

## 需求(Requirements) *(必填)*

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### 功能需求(Functional Requirements)

- **FR-001**: 系统必须(System MUST) [具体能力，例如："允许用户创建账号"]
- **FR-002**: 系统必须(System MUST) [具体能力，例如："验证邮箱地址"]
- **FR-003**: 用户必须能够(Users MUST be able to) [关键交互，例如："重置密码"]
- **FR-004**: 系统必须(System MUST) [数据需求，例如："持久化用户偏好"]
- **FR-005**: 系统必须(System MUST) [行为，例如："记录所有安全事件"]

*标注不明确需求的示例:*

- **FR-006**: 系统必须通过[需要澄清(NEEDS CLARIFICATION): 未指定认证方式 - 邮箱/密码、SSO、OAuth？]对用户进行认证
- **FR-007**: 系统必须保留用户数据[需要澄清(NEEDS CLARIFICATION): 未指定保留期限]

### 关键实体(Key Entities) *(涉及数据时填写)*

- **[实体(Entity) 1]**: [代表什么，关键属性（不含实现细节）]
- **[实体(Entity) 2]**: [代表什么，与其他实体的关系]

## 成功标准(Success Criteria) *(必填)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### 可衡量成果(Measurable Outcomes)

- **SC-001**: [可衡量指标，例如："用户能在 2 分钟内完成账号创建"]
- **SC-002**: [可衡量指标，例如："系统支持 1000 并发用户无性能下降"]
- **SC-003**: [用户满意度指标，例如："90% 用户首次尝试即成功完成主要任务"]
- **SC-004**: [业务指标，例如："减少与[X]相关的支持工单 50%"]

## 假设(Assumptions)

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right assumptions based on reasonable defaults
  chosen when the feature description did not specify certain details.
-->

- [关于目标用户的假设，例如："用户拥有稳定的互联网连接"]
- [关于范围边界的假设，例如："v1 不支持移动端"]
- [关于数据/环境的假设，例如："将复用现有认证系统"]
- [对现有系统/服务的依赖，例如："需要访问现有的用户资料 API"]
