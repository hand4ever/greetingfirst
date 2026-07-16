<!--
  ============================================================================
  同步影响报告(Sync Impact Report)

  版本变更(Version change): 1.3.2 → 1.3.3
  原因(Reason): (1) 澄清原则 I 中间件规则——原来"不阻断请求链"过于绝对，中间件
  遇致命错误时可通过不调用 next(c) 并直接返回错误响应来阻断请求链；(2) 中文化
  原则 VI 和 VII 正文中的 feature、spec、pause-and-continue 为双语标注格式。
  无新增或修改原则，PATCH 版本升级。

  修改的原则(Modified principles):
    - I. 分层架构 (Layered Architecture): 中间件规则从"不阻断请求链"改为允许
      遇致命错误时阻断。
    - VI. 错误及时抛出 (Fail Fast): feature→功能(feature)、spec→规格说明(spec)、
      pause-and-continue→暂停继续(pause-and-continue)。
    - VII. 数据库表由用户创建 (User-Owned Schema): 同上中文化修正。
  新增章节(Added sections): 无。
  移除章节(Removed sections): 无。
  需更新的模板(Templates requiring updates):
    - .specify/templates/plan-template.md       ✅ 无需更新
    - .specify/templates/spec-template.md        ✅ 无需更新
    - .specify/templates/tasks-template.md       ✅ 无需更新
    - .specify/templates/constitution-template.md ✅ 无需更新
  后续待办(Follow-up TODOs): 无。
  ============================================================================
-->
<!--
  ============================================================================
  同步影响报告(Sync Impact Report)（历史记录）

  版本变更(Version change): 1.3.1 → 1.3.2
  原因(Reason): (1) 在开发流程中新增"文档模板规范"规则；(2) 中文化顶层标题
  （Greeting Constitution → 中英双语、Core Principles → 核心原则、
  Governance → 治理）以对齐模板格式。无新增或修改原则，PATCH 版本升级。

  修改的原则(Modified principles): 无。
  新增章节(Added sections): 在开发流程中新增"文档模板规范"规则。
  移除章节(Removed sections): 无。
  需更新的模板(Templates requiring updates):
    - .specify/templates/plan-template.md       ✅ 无需更新
    - .specify/templates/spec-template.md        ✅ 无需更新
    - .specify/templates/tasks-template.md       ✅ 无需更新
    - .specify/templates/constitution-template.md ✅ 无需更新
  后续待办(Follow-up TODOs): 无。
  ============================================================================
-->
<!--
  ============================================================================
  同步影响报告(Sync Impact Report)（历史记录）

  版本变更(Version change): 1.3.0 → 1.3.1
  原因(Reason): 更新技术栈表格，将 MySQL 标记为当前数据库，同时保留多库架构灵活性。
  从原则 VII 中移除 SQLite 相关描述（测试环境描述）。无新增、移除或重定义原则，
  PATCH 版本升级，仅澄清措辞。

  修改的原则(Modified principles):
    - VII. 数据库表由用户创建 (User-Owned Schema): 移除"（`:memory:`
      SQLite 等）"硬编码引用，替换为通用描述。
  新增章节(Added sections): 无。
  移除章节(Removed sections): 无。
  修改的技术栈(Modified tech stack): 更新数据库行。
  需更新的模板(Templates requiring updates):
    - .specify/templates/plan-template.md       ✅ 无需更新
    - .specify/templates/spec-template.md        ✅ 无需更新
    - .specify/templates/tasks-template.md       ✅ 无需更新
    - .specify/templates/constitution-template.md ✅ 无需更新
  后续待办(Follow-up TODOs): 无。
  ============================================================================
-->
<!--
  ============================================================================
  同步影响报告(Sync Impact Report)（历史记录）

  版本变更(Version change): 1.1.0 → 1.2.0
  原因(Reason): 新增原则 VII（数据库表由用户创建）—— 数据库表必须由用户通过 SQL
  迁移脚本创建；应用禁止在启动时自动创建或迁移表结构。此项实质性细化了原则 III
  （移除 AutoMigrate-from-modules 条款）和原则 V（测试环境现使用 SQL schema 脚本
  替代 AutoMigrate）。MINOR 版本升级：新增原则 + 实质性细化已有原则。
  （此前 1.0.0→1.1.0：新增原则 VI，无重定义。）

  修改的原则(Modified principles):
    - III. 可复制为模板: 移除"（业务表 AutoMigrate 由各模块自行负责）"，
      现 InitDB 禁止执行任何建表或迁移操作。
    - V. 测试覆盖: 将"TestMain MUST call AutoMigrate"替换为
      "TestMain MUST run a SQL schema script to init test tables"。
  新增章节(Added sections):
    - VII. 数据库表由用户创建 (User-Owned Schema)
  移除章节(Removed sections): 无。
  需更新的模板(Templates requiring updates):
    - .specify/templates/plan-template.md       ✅ 无需更新 （通用门禁）
    - .specify/templates/spec-template.md        ✅ 无需更新
    - .specify/templates/tasks-template.md       ✅ 无需更新
    - .specify/templates/constitution-template.md ✅ 无需更新
  后续待办(Follow-up TODOs)（手动，实现阶段）：
    - specs/004-mysql-support/spec.md: FR-008 和 User Story 4 已更新为
      用户自有 schema（本次变更已完成）。
    - model/user_test.go:14 和 handler/demo_test.go:22 中的 TestMain 必须
      从 DB.AutoMigrate 迁移到 SQL schema 脚本（按原则 V）。
    - 应添加 schema SQL 文件（如 migrations/001_user.sql）作为
      规范的、用户管理的表定义。
  ============================================================================
-->

# Greeting 项目宪法(Constitution)

## 核心原则(Core Principles)

### I. 分层架构 (Layered Architecture)

所有代码 MUST 严格遵循以下分层，各层职责明确、不可跨层调用：

```
router/    → 路由分组与注册，按业务模块拆分文件
handler/   → 请求处理层，调用 response 封装返回
entity/    → 请求参数 / 数据实体定义，按模块分子目录
model/     → 数据库映射模型（GORM），通过全局 model.DB 访问
response/  → 统一 JSON 响应格式封装
middle/    → 自定义中间件
```

- Handler 层 MUST 使用包级变量 `var Xxx = &_Xxx{}` 暴露实例
- Entity 结构体 tag MUST 使用 `query`、`param`、`json`
- 中间件默认使用 `next(c)` 将请求传递给后续处理器；遇到致命错误（如鉴权失败、参数校验不通过）时，可通过不调用 `next(c)` 并直接返回错误响应的方式阻断请求链
- 新增模块时 MUST 在 `router/router.go` 中集中注册

**设计理由(Rationale)**: 分层架构确保代码可读性、可测试性和可复用性，也是本项目作为后续项目模板的核心价值所在。

### II. 统一响应格式 (Unified Response Format)

所有 API 接口 MUST 使用 `response` 包统一返回，结构如下：

```go
type ErrMsg struct {
    Code    Code   `json:"code"`
    Message string `json:"message"`
    Data    any    `json:"data"`
    TraceID string `json:"trace_id"`
    Cost    string `json:"cost"`
    Extra   string `json:"extra,omitempty"`
}
```

- 成功：MUST 使用 `response.Ok(c, data)`
- 错误：MUST 使用 `response.NotOk(c, "message")` 或 `response.NotOkWithCode(c, "message", code)`
- 禁止 handler 中直接使用 `c.JSON()` 绕过统一封装
- 时间字段 MUST 使用 `model.LocalTime` 类型，输出格式 `2006-01-02 15:04:05`

**设计理由(Rationale)**: 统一的响应格式让前端/调用方无需适配不同接口的返回结构，降低对接成本。

### III. 可复制为模板 (Copy-Ready Template)

本项目 MUST 保持自包含和最小化，确保可以直接复制到新项目使用：

- 配置与业务逻辑分离，禁止硬编码项目特定值
- `InitDB` MUST 仅连接数据库，不执行任何建表或迁移操作（表结构由用户通过 SQL 脚本手动创建，见原则 VII）
- 数据库连接 MUST 使用环境变量或配置文件注入 DSN，不使用硬编码路径
- 项目结构中的基础能力（中间件、响应封装）MUST 保持通用，不包含业务特定逻辑
- 新增依赖时 MUST 评估是否增加复制负担，优先使用标准库或轻量依赖

**设计理由(Rationale)**: 作为学习型和基础架构项目，可复制性是其最重要的非功能性需求。

### IV. 英文代码产物 (English-Only Code Artifacts)

所有代码相关的文字产物 MUST 使用英文：

- 代码注释 MUST 使用英文
- Git commit message MUST 使用英文
- 导出函数/类型 MUST 有英文注释，格式为 `// Name describes ...`
- Commit message 格式：`type: brief description`（如 `feat:`, `fix:`, `refactor:`, `docs:`, `chore:`）
- 文档、README、接口说明等可使用中文

**设计理由(Rationale)**: 英文代码产物确保国际化和团队协作的通用性，同时中英文文档兼顾本地开发效率。

### V. 测试覆盖 (Test Coverage)

每一个 handler 和 model 方法 MUST 有对应的单元测试：

- 测试文件命名：`xxx_test.go`，与源文件同目录
- Handler 测试：使用 `httptest.NewRequest` + `echo.New().NewContext` 构造请求
- Model 测试：直接调用模型方法，使用内存数据库 `:memory:`
- 每个测试文件 MUST 在 `TestMain` 中通过 SQL 建表脚本（如 `xxx_schema.sql`）初始化测试表，禁止使用 `AutoMigrate` 自动建表（见原则 VII）
- 测试输出 MUST 使用 `logOK` 辅助函数，确保无 `-v` 时也能看到响应内容
- VSCode 测试配置 MUST 包含 `"go.testFlags": ["-v"]`
- 运行命令：`go test -v ./... -count=1`

**设计理由(Rationale)**: 测试是代码质量的最后防线，完整的测试覆盖让项目模板更可靠、更值得信赖。

### VI. 错误及时抛出 (Fail Fast)

所有错误处理 MUST 以显式抛出和告警为默认策略，静默降级为显式设计的例外：

- 配置文件缺失或加载失败时，MUST 打印明确错误信息并以非零退出码终止启动
- 外部依赖（数据库、外部服务等）连接失败时，MUST 显式报错而非使用默认值继续运行
- 表结构缺失（目标数据库中不存在所需表）时，默认 MUST 显式报错退出，禁止自动建表补全（见原则 VII）。例外：某功能(feature)的规格说明(spec)可显式约定「暂停继续(pause-and-continue)」模式——应用打印提醒并暂停轮询，待用户建表后自动继续；该模式 MUST 仍禁止自动建表，且连接失败仍按默认快速失败(fail-fast)退出，不受此例外影响
- 降级措施 ONLY 允许在以下条件下使用：
  - 设计文档（spec.md）中明确约定了降级策略
  - 降级行为有对应的测试用例覆盖
- 未明确定义降级策略的错误场景，DEFAULT 按抛出错误 / 日志告警处理
- 禁止使用 `_` 忽略 error 返回值；禁止无日志的静默吞错

**设计理由(Rationale)**: 错误被静默吞掉会导致线上行为不可预期、运维排障困难。明确报错让问题在最早暴露、最容易定位的阶段被处理，是项目可靠性的基石。降级看似"健壮",实则隐藏了真实故障。

### VII. 数据库表由用户创建 (User-Owned Schema)

数据库表结构 MUST 由用户通过 SQL 脚本（迁移脚本）手动创建，应用 MUST NOT 在启动时自动创建或迁移表结构：

- 禁止在 `InitDB` 或任何初始化路径中调用 `AutoMigrate` / 执行 `CREATE TABLE`
- 数据库 schema 以独立的 `.sql` 迁移文件管理，作为项目可复制资产的一部分
- 应用启动只需连接已存在表结构的数据库；表缺失默认 MUST 按原则 VI 显式报错，而非自动补全。若某功能(feature)的规格说明(spec)显式约定暂停继续(pause-and-continue)模式，则按该模式处理（仍禁止自动补全）
- 测试环境 MUST 在 `TestMain` 中执行与生产环境相同的 schema SQL 脚本，保持 schema 一致（见原则 V）

**设计理由(Rationale)**: 表结构属于数据资产，交由用户显式管理可避免 schema 漂移、隐式变更与环境不一致，也契合 原则 III（可复制为模板）与原则 VI（错误及时抛出）——schema 的缺失应在最早阶段被明确暴露，而非被框架静默抹平。

## 技术栈约束

本项目技术栈 MUST 在以下范围内选择，新增技术需评估必要性：

| 类别 | 选择 | 约束 |
|------|------|------|
| 语言 | Go | 版本 ≥ 1.22 |
| Web 框架 | Echo v5 | 禁止混用其他 Web 框架 |
| ORM | GORM v2 | 通过 `model.DB` 全局实例访问 |
| 数据库 | MySQL / 可替换 | 通过 DSN 注入切换，架构支持多库并存 |
| 模块名 | `greeting.first` | 待稳定后调整 |

- 禁止引入与现有分层架构冲突的框架（如全栈框架替换 Echo）
- 新增依赖前 MUST 评估必要性，优先使用 Go 标准库
- 格式化 MUST 使用 `go fmt` 或 `gofumpt`，提交前执行

## 开发流程

文档模板规范：

- 所有 spec.md / plan.md / tasks.md / research.md / quickstart.md / contracts/ 等规格产物 MUST 严格遵循 `.specify/templates/` 中对应模板的标题格式、层级结构和语言（中文标题 + 双语标注），禁止自行替换为纯英文标题或改变模板定义的章节结构

新增接口 MUST 遵循以下流程：

1. 在 `entity/<模块>/` 中定义请求参数结构体
2. 在 `handler/` 中创建处理器文件
3. 在 `router/` 中注册路由
4. 在 `api.http` 中添加 REST Client 测试用例
5. 在 `README.md` 的「API 列表」和「更新日志」中记录

代码质量要求：

- 单文件 MUST NOT 超过 500 行
- 单函数 MUST NOT 超过 80 行，复杂逻辑拆分子函数
- 控制流缩进 MUST NOT 超过 3 层
- 错误 MUST 始终检查并及早返回，禁止使用 `_` 忽略 error
- Panic MUST NOT 用于常规业务错误，仅用于不可恢复的初始化失败
- 方法接收器 MUST 使用类型首字母小写

## 治理(Governance)

本宪法是 Greeting 项目的最高行为准则，所有代码变更和架构决策 MUST 以此为基准：

- **修订流程**：任何原则变更 MUST 通过 PR 提交宪法修订，经审查后合并；重大变更（MAJOR 版本）需额外说明迁移方案
- **版本策略**：遵循语义化版本（SemVer）—— MAJOR 为不兼容的原则移除/重定义，MINOR 为新增原则或章节，PATCH 为措辞澄清或修正
- **合规审查**：每次 `/speckit.plan` 执行时 MUST 检查 宪法检查(Constitution Check)门禁，违规需在复杂度追踪(Complexity Tracking)中说明理由和替代方案
- **运行时指导**：日常开发细节（命名、错误处理、注释规范等）详见 `.codebuddy/rules/GO_STYLE.mdc`

**版本(Version)**: 1.3.3 | **批准日期(Ratified)**: 2026-07-13 | **最后修订(Last Amended)**: 2026-07-16
