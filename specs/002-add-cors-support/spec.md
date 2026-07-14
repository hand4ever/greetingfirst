# Feature Specification: 添加 CORS 跨域支持

**Feature Branch**: `002-add-cors-support`

**Created**: 2026-07-14

**Status**: Draft

**Input**: User description: "为接口添加cors跨域"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 前端跨域调用 API (Priority: P1)

前端应用（运行在不同域名或端口）需要通过 XMLHttpRequest 或 Fetch API 调用本服务的 REST API。浏览器会先发送 OPTIONS 预检请求，验证服务器是否允许跨域访问，确认允许后才发送实际请求。

**Why this priority**: 这是 CORS 的核心功能，没有它前端无法从不同源调用任何 API。

**Independent Test**: 启动服务后，使用 curl 模拟浏览器发送 OPTIONS 预检请求和带 Origin 头的 GET 请求，验证响应包含正确的 CORS 头。

**Acceptance Scenarios**:

1. **Given** 服务已启动并启用 CORS，**When** 浏览器发送 OPTIONS 预检请求（含 `Origin`、`Access-Control-Request-Method` 头），**Then** 响应包含 `Access-Control-Allow-Origin`、`Access-Control-Allow-Methods` 头，HTTP 状态码为 200 或 204。
2. **Given** 服务已启动并启用 CORS，**When** 浏览器发送带 `Origin` 头的 GET 请求，**Then** 响应头包含 `Access-Control-Allow-Origin`，允许前端读取响应数据。
3. **Given** 服务已启动并启用 CORS，**When** 浏览器发送带 `Origin` 头的 POST/PUT/DELETE 等非简单请求，**Then** 先通过 OPTIONS 预检，再正常处理实际请求。

---

### User Story 2 - 自定义允许的来源域名 (Priority: P2)

运维或开发人员可以根据部署环境配置允许跨域访问的来源域名列表，而非写死允许所有来源（`*`），以提高安全性。

**Why this priority**: 安全配置是生产环境部署的必要条件，但在开发和演示阶段允许所有来源即可满足需求。

**Independent Test**: 修改 CORS 允许的来源配置，验证只有配置中的来源可以成功跨域访问，未配置的来源返回无 CORS 头或错误。

**Acceptance Scenarios**:

1. **Given** CORS 配置允许来源 `https://example.com`，**When** 浏览器从 `https://example.com` 发起跨域请求，**Then** 响应包含 `Access-Control-Allow-Origin: https://example.com`。
2. **Given** CORS 配置允许来源 `https://example.com`，**When** 浏览器从 `https://evil.com` 发起跨域请求，**Then** 响应不包含 `Access-Control-Allow-Origin` 头，浏览器阻止响应。

---

### Edge Cases

- 不带 `Origin` 头的同源请求如何处理？→ 正常处理，不添加 CORS 响应头
- OPTIONS 请求不包含 `Access-Control-Request-Method` 头时如何处理？→ 正常返回 200，不添加 CORS 响应头
- 请求头 `Access-Control-Request-Headers` 包含非标准自定义头时如何处理？→ 根据配置决定是否允许，默认允许常见头
- 服务未启用 CORS（默认关闭）时跨域请求如何处理？→ 不添加 CORS 响应头，浏览器默认阻止

## Clarifications

### Session 2026-07-14

- Q: CORS 中间件是否支持携带凭证（Cookie/Authorization 头）的跨域请求？ → A: Option B — 默认关闭凭证模式，仅支持无凭证的简单跨域请求

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统 MUST 提供 CORS 中间件，可在服务启动时选择启用或禁用
- **FR-002**: 启用 CORS 后，系统 MUST 对 OPTIONS 预检请求做出正确响应，返回允许的方法和头部信息
- **FR-003**: 启用 CORS 后，系统 MUST 在响应中添加 `Access-Control-Allow-Origin` 头，值可配置（支持指定域名列表或通配符 `*`）
- **FR-004**: 系统 MUST 支持配置允许的 HTTP 方法列表（如 GET、POST、PUT、DELETE、OPTIONS），默认为常用方法
- **FR-005**: 系统 MUST 支持配置允许的请求头列表（如 Content-Type、Authorization），默认为常用头
- **FR-006**: CORS 中间件 MUST 默认不启用凭证模式（`Access-Control-Allow-Credentials` 为 false），仅支持无凭证的简单跨域请求
- **FR-007**: CORS 中间件 MUST 在全局中间件链中按正确顺序注册（应在 Recover 和 RequestID 之后，业务路由之前）
- **FR-008**: CORS 中间件 MUST 对非跨域请求（无 Origin 头）保持透明，不添加多余的 CORS 响应头
- **FR-009**: CORS 配置 MUST 支持通过常量或配置文件方式修改，保持项目作为模板的可复制性

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 启用 CORS 后，来自不同源的浏览器请求可以成功获取 API 响应内容，不再出现 CORS 错误
- **SC-002**: OPTIONS 预检请求响应时间与普通请求一致，无明显延迟增加（在 50ms 以内）
- **SC-003**: CORS 中间件不对非跨域请求产生任何影响（响应内容、响应时间、响应头与未启用 CORS 时一致）
- **SC-004**: CORS 配置可以通过修改一个常量或配置变量来调整允许的来源、方法和头部，无需改动中间件逻辑代码

## Assumptions

- 目标用户是前端开发者，需要从不同域名/端口访问本服务 API
- 默认配置允许所有来源（`*`）和常用 HTTP 方法，适合开发环境
- 使用 Echo v5 框架内置的 CORS 中间件（`github.com/labstack/echo/v5/middleware`），无需引入新的第三方依赖
- CORS 作为全局中间件应用于所有路由，不做按路由粒度的 CORS 配置
- 生产环境部署时由运维人员根据实际需要修改 CORS 配置常量
- 默认不启用凭证模式（`Access-Control-Allow-Credentials: false`），若需要携带 Cookie 或 Authorization 头的跨域请求，需由使用方自行开启
