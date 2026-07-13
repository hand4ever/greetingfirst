# Feature Specification: SHA256 Demo API

**Feature Branch**: `001-sha256-demo-api`

**Created**: 2026-07-13

**Status**: Draft

**Input**: User description: "实现一个demo新测试接口，用于返回sha256的返回值，输入值在querystring里"

## Clarifications

### Session 2026-07-13

- Q: 响应中是否需要回显原始输入？ → A: 响应中包含 `input`（原始文本）和 `hash`（SHA256 值）两个字段

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Compute SHA256 Hash (Priority: P1)

调用方通过 GET 请求传入一个文本参数，接口返回该文本的 SHA256 哈希值。用于快速验证接口可用性和哈希计算正确性。

**Why this priority**: 这是该接口的唯一核心功能，没有其他子功能。

**Independent Test**: 发送 `GET /demo/sha256?text=hello`，验证返回的 SHA256 值与预期一致（在线工具或命令行 `echo -n "hello" | sha256sum` 可交叉验证）。

**Acceptance Scenarios**:

1. **Given** 接口正常运行，**When** 调用方发送 `GET /demo/sha256?text=hello`，**Then** 返回 HTTP 200，响应体中 `data` 包含 `input: "hello"` 和 `hash: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"`，以及 `code: 0`
2. **Given** 接口正常运行，**When** 调用方发送 `GET /demo/sha256?text=你好`，**Then** 返回 HTTP 200，响应体中 `data` 包含 `input: "你好"` 和正确的 `hash` 值
3. **Given** 接口正常运行，**When** 调用方发送 `GET /demo/sha256`（不传 text 参数），**Then** 返回错误响应，`message` 包含有意义的错误提示

---

### Edge Cases

- 当 `text` 参数为空字符串（`text=`）时，返回空字符串的 SHA256 值（`e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`）
- 当 `text` 参数包含特殊字符（URL 编码字符、emoji 等）时，正确计算其 SHA256
- 当 `text` 参数为超长字符串（如 1MB 文本）时，接口正常完成计算且不超时
- 当请求方法不是 GET 时，返回适当错误

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 接口 MUST 支持 GET 请求，从 query string 中读取名为 `text` 的参数
- **FR-002**: 接口 MUST 对 `text` 参数值计算 SHA256 哈希，并以小写十六进制字符串格式返回；同时 MUST 在响应 `data` 中回显原始输入 `input` 和计算结果 `hash`
- **FR-003**: 接口 MUST 在 `text` 参数缺失时返回错误响应（HTTP 400），附带可读的错误信息
- **FR-004**: 接口 MUST 使用项目统一响应格式（`response.Ok` / `response.NotOk`）返回结果
- **FR-005**: 接口 MUST 路由注册在 `/demo/sha256` 路径下
- **FR-006**: 接口 MUST 处理任意 UTF-8 编码的输入文本

### Key Entities

- **SHA256 Request**: 包含 query 参数 `text`（string），为需要计算哈希的输入文本
- **SHA256 Response Data**: 包含 `input`（string，原始输入文本）和 `hash`（string，小写十六进制 SHA256 哈希值）

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 调用方发送任意有效请求后，在 1 秒内获得正确响应
- **SC-002**: 返回的 SHA256 值与标准 SHA256 算法计算结果 100% 一致（交叉验证通过）
- **SC-003**: 缺少必填参数时，接口返回 HTTP 400 并包含明确错误信息，而非 5xx 或空响应
- **SC-004**: 接口能处理至少 1MB 大小的输入文本而不崩溃

## Assumptions

- 调用方无需认证即可访问该 Demo 接口（纯展示用途）
- 接口仅支持 GET 方法，不支持 POST
- `text` 参数名固定，不区分大小写处理按 Echo 框架默认行为
- 不需要对输入文本做长度硬限制（由 Go 运行时和内存自然限制）
- 不需要持久化请求日志或哈希结果
- 接口使用标准库 `crypto/sha256` 实现，无需额外依赖
