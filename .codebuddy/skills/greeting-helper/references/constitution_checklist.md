# 宪法合规核对清单(Constitution Compliance Checklist)

**用途**：配合 `greeting-helper` 技能，对 Greeting 项目改动做合规自检。
**来源**：`.specify/memory/constitution.md`（当前版本见文件末尾 Version 字段）。
**用法**：仅勾选本次改动受影响的原则条目；每条给出 通过 / 不通过，不通过需附 文件:行号 + 整改建议。

---

## I. 分层架构(Layered Architecture)

- [ ] 代码仅落在 `router/ handler/ entity/ model/ response/ middle/` 分层，无跨层越权调用
- [ ] handler 使用包级变量 `var Xxx = &_Xxx{}` 暴露实例
- [ ] entity 结构体 tag 使用 `query` / `param` / `json`
- [ ] 中间件默认 `next(c)` 链式传递；遇致命错误可不调用 `next(c)` 直接返回错误响应
- [ ] 新增模块在 `router/router.go` 中集中注册

## II. 统一响应格式(Unified Response Format)

- [ ] 成功返回使用 `response.Ok(c, data)`
- [ ] 错误返回使用 `response.NotOk(c, "message")` 或 `response.NotOkWithCode(c, "message", code)`
- [ ] handler 中禁止直接使用 `c.JSON()` 绕过统一封装
- [ ] 时间字段使用 `model.LocalTime`，输出格式 `2006-01-02 15:04:05`

## III. 可复制为模板(Copy-Ready Template)

- [ ] 配置与业务逻辑分离，无硬编码项目特定值
- [ ] `InitDB` 仅连接数据库，不执行任何建表或迁移
- [ ] 数据库连接使用环境变量或配置文件注入 DSN，无硬编码路径
- [ ] 基础能力（中间件、响应封装）保持通用，不含业务特定逻辑

## IV. 英文代码产物(English-Only Code Artifacts)

- [ ] 代码注释使用英文
- [ ] git commit message 使用英文，格式 `type: brief description`（feat/fix/refactor/docs/chore）
- [ ] 导出函数 / 类型有英文注释，格式 `// Name describes ...`
- [ ] 文档、README、接口说明可使用中文

## V. 测试覆盖(Test Coverage)

- [ ] 每个 handler 与 model 方法都有对应单元测试（`xxx_test.go` 同目录）
- [ ] handler 测试使用 `httptest.NewRequest` + `echo.New().NewContext`
- [ ] model 测试直接调用模型方法，使用 `:memory:` 内存库
- [ ] `TestMain` 通过 SQL 建表脚本（如 `xxx_schema.sql`）初始化，**禁止 `AutoMigrate` 自动建表**
- [ ] 测试输出使用 `logOK` 辅助函数
- [ ] 提交前运行 `go test -v ./... -count=1` 全绿

## VI. 错误及时抛出(Fail Fast)

- [ ] 配置文件缺失 / 加载失败 → 打印明确错误并以非零退出码终止启动
- [ ] 外部依赖（数据库、外部服务）连接失败 → 显式报错，不用默认值继续
- [ ] 表结构缺失 → 默认显式报错退出，禁止自动建表补全
- [ ] 禁止用 `_` 忽略 error；禁止无日志静默吞错
- [ ] 降级 ONLY 允许在 spec 明确约定策略且有测试覆盖时；否则默认抛出

## VII. 数据库表由用户创建(User-Owned Schema)

- [ ] `InitDB` 及任何初始化路径禁止 `AutoMigrate` / `CREATE TABLE`
- [ ] 数据库 schema 由独立 `.sql` 迁移文件管理（项目可复制资产）
- [ ] 测试环境 `TestMain` 执行与生产相同的 schema SQL，保持 schema 一致

## VIII. 中文交互(Chinese-First Interaction)

- [ ] 面向用户的思考、分析、回答使用简体中文
- [ ] 错误提示、解释说明、状态报告、验证结论使用中文
- [ ] 规格文档（spec/plan/tasks 等）叙述与标题使用中文（遵循 `.specify/templates` 双语格式）
- [ ] 不覆盖原则 IV：代码注释、commit message、导出符号注释仍用英文

---

## 开发流程(Development Workflow)

### 新增接口流程

- [ ] 在 `entity/<模块>/` 定义请求参数结构体
- [ ] 在 `handler/` 创建处理器文件
- [ ] 在 `router/` 注册路由
- [ ] 在 `api.http` 添加 REST Client 测试用例
- [ ] 在 `README.md` 的「API 列表」与「更新日志」中记录

### 代码质量要求

- [ ] 单文件不超过 500 行
- [ ] 单函数不超过 80 行，复杂逻辑拆分子函数
- [ ] 控制流缩进不超过 3 层
- [ ] 错误始终检查并及早返回，禁止 `_` 忽略 error
- [ ] 不使用 panic 处理常规业务错误（仅用于不可恢复的初始化失败）
- [ ] 方法接收器使用类型首字母小写

### Changelog 强制登记

- [ ] 所有更改（新增功能、缺陷修复、配置调整、重构、文档/spec 变更等）与新增任务，MUST 在 `changelog.toml` 的 `[[changelog]]` 追加一条（`date=YYYY-MM-DD` 当天，`content` 简述）
- [ ] 若 `config.go` 的 `defaultConfig()` 默认 changelog 也需同步，避免配置文件缺失时接口内容不一致
- [ ] 未登记 changelog 的改动视为未完成

---

## 备注(Notes)

- 违反原则时，在复杂度追踪(Complexity Tracking)中说明理由与替代方案
- 清单条目为"受影响才核对"，避免对非相关改动做无意义扫描
