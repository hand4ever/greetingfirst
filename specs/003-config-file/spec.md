# Feature Specification: 全局配置文件

**Feature Branch**: `003-config-file`

**Created**: 2026-07-14

**Status**: Draft

**Input**: User description: "将version和changelog里的数据改为从配置文件里读取，故需要新增全局配置文件，在项目启动的时候就读取，配置文件的格式为toml或yaml，给出两者的优劣，取其一。"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 通过配置文件管理版本和更新日志 (Priority: P1)

运维或开发人员通过编辑项目根目录的 `config.toml` 文件来修改应用版本号、更新日志等内容，无需重新编译代码即可让 `/common/version`、`/common/changelog`、`/common/setting` 三个接口返回最新信息。

**Why this priority**: 这是配置文件的核心价值——将运行时数据与代码分离，实现不重新编译即可变更。

**Independent Test**: 修改 `config.toml` 中的版本号后重启服务，访问 `/common/version` 验证返回新版本号。

**Acceptance Scenarios**:

1. **Given** 服务已启动且 `config.toml` 中 `version = "1.0.0"`，**When** 调用 `GET /common/version`，**Then** 返回 `{"version": "1.0.0", ...}`。
2. **Given** 编辑 `config.toml` 将 `version` 改为 `"2.0.0"` 并重启服务，**When** 再次调用 `GET /common/version`，**Then** 返回 `{"version": "2.0.0", ...}`。
3. **Given** `config.toml` 中定义了 changelog 条目，**When** 调用 `GET /common/changelog`，**Then** 返回配置文件中定义的所有条目。

---

### User Story 2 - 配置文件异常时阻止启动 (Priority: P2)

当配置文件缺失或 TOML 格式错误时，系统 MUST 打印明确错误信息并拒绝启动，不允许静默降级或使用默认值继续运行。

**Why this priority**: 运维安全——配置文件异常意味着部署意图未正确表达，静默降级会掩盖配置错误，导致线上行为不符合预期。

**Independent Test**: 删除 `config.toml` 后启动服务，验证进程非零退出码并 stderr 包含明确错误信息。

**Acceptance Scenarios**:

1. **Given** 项目根目录不存在 `config.toml`，**When** 启动服务，**Then** 进程退出码非零，stderr 打印 "config.toml not found" 错误信息。
2. **Given** `config.toml` 存在但格式错误（如缺少引号），**When** 启动服务，**Then** 进程退出码非零，stderr 打印 TOML 解析错误信息。

---

### Edge Cases

- 配置文件存在但为空文件时如何处理？→ 空文件 TOML 解析失败，阻止启动并报错
- changelog 数组为空时 `/common/changelog` 返回什么？→ 返回空数组 `[]`
- 配置项缺失部分字段时如何处理？→ 配置项缺失字段使用 Go 零值（文件合法且已成功加载，仅部分字段未填写）
- 配置文件编码非 UTF-8 时如何处理？→ TOML 标准要求 UTF-8，非 UTF-8 解析失败按格式错误处理，阻止启动

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统 MUST 在启动时自动加载项目根目录的 `config.toml` 配置文件
- **FR-002**: 配置文件 MUST 使用 TOML 格式（经对比，TOML 更适合 Go 生态且不易出错）
- **FR-003**: 配置文件缺失或 TOML 解析失败时，系统 MUST 打印错误信息并以非零退出码终止启动，不允许静默降级
- **FR-004**: `GET /common/version` MUST 从配置文件的 `[app]` 段读取 `version` 和 `build_time`
- **FR-005**: `GET /common/changelog` MUST 从配置文件的 `[[changelog]]` 数组读取条目
- **FR-006**: `GET /common/setting` MUST 从配置文件的 `[app]`、`[server]`、`[database]` 段读取设置项
- **FR-007**: 配置结构体 MUST 定义为包级导出类型，方便其他模块读取
- **FR-008**: 配置加载逻辑 MUST 提供 `InitConfig(configPath string)` 函数，支持通过参数指定路径（便于测试）
- **FR-009**: (Deferred) 运行时通过命令行 `-config` flag 或环境变量 `CONFIG_PATH` 覆盖默认路径 — 当前硬编码 `"config.toml"`，后续迭代实现
- **FR-010**: 服务启动监听端口 MUST 从配置文件 `[server] port` 读取并用于 `e.Start`，不得硬编码端口号

### Key Entities

- **Config**: 顶层配置结构体，包含 App、Server、Database、Changelog 四个子段
- **AppConfig**: 应用元信息（名称、版本、构建时间）
- **ServerConfig**: 服务端配置（监听端口）
- **DatabaseConfig**: 数据库配置（类型、DSN）
- **ChangelogEntry**: 更新日志条目（日期、内容）——复用 `entity/common.ChangelogEntry`

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 修改 `config.toml` 后重启服务，3 个 common 接口均返回配置中的新值
- **SC-002**: 缺失 `config.toml` 或格式错误时，进程以非零退出码终止，stderr 包含明确错误信息
- **SC-003**: 全量单元测试通过，新增配置相关的测试用例
- **SC-004**: 项目保持零外部依赖增长（TOML 库为标准 Go 模块生态，增量可控）
- **SC-005**: 修改 `config.toml` 的 `[server] port` 后重启服务，服务监听新端口（验证 FR-010）

## Clarifications

### Session 2026-07-14

- Q: 配置文件路径应如何灵活指定？ → A: 暂不处理，当前硬编码 `"config.toml"`；命令行参数和环境变量覆盖留待后续迭代。
- Q: 配置文件不存在或格式错误时应如何处理？ → A: 统一抛错阻止启动（文件不存在 OR 格式错误均视为致命错误），不允许静默降级。

## Assumptions

- 配置文件默认路径为项目根目录的 `config.toml`，可通过 `InitConfig` 参数自定义；运行时路径覆盖（命令行/环境变量）延后处理
- 配置在启动时一次性加载到内存，运行期间不热更新（后续可扩展）
- TOML 序列化库使用 `github.com/BurntSushi/toml`（Go 生态最成熟的 TOML 库）
- 配置文件不包含敏感信息（数据库 DSN 仅示例用途）
- 配置文件异常（缺失/格式错误）一律视为致命错误，拒绝启动——不提供默认值降级路径
