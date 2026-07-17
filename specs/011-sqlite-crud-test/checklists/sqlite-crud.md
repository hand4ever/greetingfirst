# 需求质量检查清单(Requirements Quality Checklist): 新增 SQLite 实例与 CRUD 测试接口

**目的(Purpose)**: 以「英文单元测试」视角校验本特性需求文档（spec/plan/data-model/contracts/quickstart/tasks）的完整性、清晰性、一致性、可测性与覆盖度；并将项目宪法（原则 V/VII/VI/II/I/III）合规与设计重点「双库数据隔离 SC-006」作为显式门禁项。

**创建日期(Created)**: 2026-07-17
**特性(Feature)**: [spec.md](../spec.md) | [plan.md](../plan.md) | [data-model.md](../data-model.md) | [contracts/api.md](../contracts/api.md)

**说明(Note)**: 此清单由 `/speckit.checklist` 生成，仅检验「需求是否被清晰、完整、一致、可测地描述」，不检验实现行为。范围=完整含宪法合规；深度=标准；强制门禁重点=双库数据隔离 SC-006。

---

## 需求完整性(Requirement Completeness)

- [ ] CHK001 是否将 5 个 CRUD 操作（创建/按 ID 查询/更新/删除/列表）均作为独立需求显式定义并各自可追踪到 FR？[Completeness, Spec §FR-005]
- [ ] CHK002 `age` 缺省值（未传时=0）是否在需求层（而非仅在 data-model）被显式规定？[Completeness, Spec §FR-006]
- [ ] CHK003 请求字段容量约束（name≤100、phone≤20）是否作为 API 层校验需求存在，还是仅停留在 DB 列类型定义？[Gap, data-model §1]
- [ ] CHK004 是否显式规定了 phone 的格式/长度校验（如纯数字、位数），或明确将其排除在范围外？[Gap, Spec §FR-006]
- [ ] CHK005 列表接口的分页（pagination）是否明确声明为「本期不做（exclusion）」，还是需求遗漏？[Gap, contracts §5]

## 需求清晰性(Requirement Clarity)

- [ ] CHK006 "test_user 表由用户预先创建并存在" 是否清晰说明「由谁」通过「哪个产物」（migrations/002_test_user.sql）创建？[Clarity, Spec §FR-010]
- [ ] CHK007 默认 DSN 位置 "项目同级目录" 是否为精确、无歧义的文件系统路径描述？[Clarity, Spec Assumptions]
- [ ] CHK008 FR-009 "可按创建时间倒序" 是否明确该倒序是「强制」还是「可选」？[Ambiguity, Spec §FR-009]
- [ ] CHK009 "活跃记录内唯一（phone）" 是否明确定义为「排除已软删除记录」？[Clarity, Spec Edge Cases]
- [ ] CHK010 各接口的错误消息（如 "phone already exists"）是否被规定为「必须满足的响应契约」，而非示例？[Clarity, contracts §1]

## 需求一致性(Requirement Consistency)

- [x] CHK011 列表端点路由在 plan/tasks 与 contracts 之间现已一致（统一为 `/sqlite/testuser` + `/sqlite/testusers`），无冲突。[Conflict, resolved via /speckit.clarify 2026-07-17]
- [ ] CHK012 SC-004（并发创建）在 spec「成功标准」与 quickstart「验证检查清单」之间是否被一致对待（必需 vs 可选）？[Consistency, Spec §SC-004 vs quickstart]
- [ ] CHK013 FR-009（列表"可按…倒序"）与 contracts §5（"默认按 created_at DESC"）在排序行为上是否对齐？[Consistency, Spec §FR-009 vs contracts §5]

## 验收标准可测性(Acceptance Criteria Quality)

- [ ] CHK014 SC-006（双库隔离）是否以「双向可验证」准则表述：SQLite→MySQL 不可见 且 MySQL→SQLite 不可见？[Measurability, Spec §SC-006]
- [ ] CHK015 是否存在可度量的验收准则来证明应用「绝不自动建表/迁移 test_user」（FR-010 / 原则 VII）？[Measurability, Spec §FR-010]
- [ ] CHK016 每个 FR 是否均可映射到具体 SC/US 验收场景，或存在缺少验收的 FR？[Traceability, Spec §Requirements]
- [ ] CHK017 SC-002（"仅需配置 [database.sqlite] 即可启用"）是否可被客观验证？[Measurability, Spec §SC-002]

## 场景覆盖(Scenario Coverage)

- [ ] CHK018 备选流程（部分更新、未传字段保持不变）是否由显式需求 + 验收场景共同覆盖？[Coverage, Spec §FR-007]
- [ ] CHK019 全部 5 个端点的异常流（id 非法、not found、phone 重复、DB 失败）是否在 contracts 中一致定义？[Coverage, contracts §1–§5]
- [ ] CHK020 当 MySQL 成功而 SQLite 初始化失败（或反之）时，是否定义了恢复/中止需求（而非静默单库运行）？[Coverage, Gap, Spec §FR-003]
- [ ] CHK021 非功能并发需求（SC-004）是否定义了明确的并发量级与成功判定谓词？[Coverage, Spec §SC-004]

## 边界与异常覆盖(Edge Case Coverage)

- [ ] CHK022 是否规定 `age` 取负值/越界（如 -1）时的处理（拒绝或接受）？[Gap, Edge Case]
- [ ] CHK023 创建/更新时 `name`/`phone` 超出列容量（DB 报错路径）的行为是否被定义？[Gap, data-model §1]
- [ ] CHK024 是否规定 phone 含前后空格/格式化的归一化或拒绝策略？[Gap, Spec §FR-006]
- [ ] CHK025 `:memory:` 测试模式与生产文件型 DB 之间「互不污染」是否作为显式需求？[Coverage, Spec Assumptions]

## 非功能需求(Non-Functional Requirements)

- [ ] CHK026 性能目标（SC-001 连接 3 秒、SC-004 并发）是否覆盖所有关键路径，还是仅启动阶段？[Completeness, Spec §SC-001/SC-004]
- [ ] CHK027 这些测试接口是否需要访问控制/鉴权需求，或明确声明「开放测试接口、无需鉴权」？[Gap, Spec §FR-005]
- [ ] CHK028 fail-fast 时效（SC-005 "5 秒内退出"）是否以可度量的触发条件表述？[Measurability, Spec §SC-005]

## 依赖与假设(Dependencies & Assumptions)

- [ ] CHK029 假设 "glebarez/sqlite 纯 Go、无 CGO" 是否作为需求被记录/校验（契合原则 III 可复制）？[Assumption, Spec Assumptions / Principle III]
- [ ] CHK030 对既有 `config`/`model`/`response`/`router` 模块的依赖是否显式列为带版本的先决条件？[Dependency, Spec §Dependencies]
- [ ] CHK031 假设 "greeting.db 中 test_user 已存在" 是否被校验，并定义「若不存在」的失败模式（FR-010 fail-fast）？[Assumption, Spec §FR-010]

## 宪法合规门禁(Constitution Compliance Gate)

- [ ] CHK032 需求是否显式禁止在任何初始化路径对 test_user 调用 AutoMigrate / CREATE TABLE（原则 VII）？[Constitution, Principle VII / FR-010]
- [ ] CHK033 测试需求是否强制 TestMain 执行 migrations/002_test_user.sql（非 AutoMigrate）来初始化 :memory: 表（原则 V）？[Constitution, Principle V / tasks T010]
- [ ] CHK034 SQLite 连接/Ping 失败与表缺失的 fail-fast 需求是否与原则 VI（明确报错 + 非零退出）一致？[Constitution, Principle VI / FR-003]
- [ ] CHK035 是否强制全部 5 个端点使用 response.Ok/NotOk/NotOkWithCode（原则 II）？[Constitution, Principle II / FR-011]
- [ ] CHK036 新模块是否落在正确分层（entity/sqliteusr、handler、router、model）（原则 I）？[Constitution, Principle I]

## 数据隔离(强制门禁重点 / Mandatory Gate: SC-006)

- [ ] CHK037 SC-006 是否要求「经 /sqlite/testuser 写入的数据在 MySQL `users` 中不可查询」（方向 A）？[Mandatory Gate, Spec §SC-006]
- [ ] CHK038 SC-006 是否要求「经 MySQL /demo/usr 写入的数据在 SQLite test_user 中不可查询」（方向 B）？[Mandatory Gate, Spec §SC-006]
- [ ] CHK039 是否存在验证「双向隔离」的验收场景，而非仅验证单一方向？[Mandatory Gate, Coverage, quickstart §Scenario D]
- [ ] CHK040 隔离的实现机制（独立 DB 实例 + 独立表）是否被显式陈述，使隔离成为「设计保证」而非偶然？[Mandatory Gate, Clarity, plan §Technical Context]

## 备注(Notes)

- 完成后标记：`[x]`
- 重点门禁：CHK037–CHK040（双库数据隔离 SC-006）为合并前必须全部通过的强制项
- 宪法门禁：CHK032–CHK036 须与 plan.md「宪法检查」结论一致；若任一不通过需记入复杂度追踪
- 项目按顺序编号，便于引用
