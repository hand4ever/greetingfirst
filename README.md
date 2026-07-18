# Greeting

基于 [Echo v5](https://github.com/labstack/echo) 的 Go Web 服务示例项目，采用分层架构搭建，内置统一响应格式、请求耗时统计、请求链路追踪等基础能力。

## 技术栈

| 项目 | 说明 |
|------|------|
| 语言 | Go 1.26.3 |
| Web 框架 | Echo v5.2.1 |
| ORM | GORM v1.31.2 (MySQL) |
| 配置文件 | TOML (`config.toml`) |
| 模块名 | `greeting.first` |
| 监听端口 | `:1323` |

## 目录结构

```
greeting/
├── main.go            # 程序入口，初始化 Echo、注册中间件、启动服务
├── go.mod             # 模块定义与依赖
├── Makefile           # 本地运行 / 测试服部署脚本
├── config.toml        # 全局配置文件（TOML 格式）
├── config/            # 配置加载包（全局 Cfg 实例）
├── router/            # 路由分组与注册
├── handler/           # 请求处理（控制器层）
├── entity/            # 请求参数 / 数据实体定义
├── model/             # 数据库映射模型（GORM），全局 DB 实例
├── response/          # 统一 JSON 响应格式封装
└── middle/            # 自定义中间件
```

## 分层职责

| 目录 | 职责 |
|------|------|
| `router/` | 路由分组与注册，按业务模块划分 |
| `handler/` | 接收请求、参数绑定、调用响应封装 |
| `entity/` | 请求参数结构体（查询参数 / 路径参数等） |
| `model/` | 数据库映射模型（GORM），通过 `model.DB` 访问全局实例 |
| `response/` | 统一的成功 / 错误响应结构 |
| `middle/` | 自定义中间件（耗时统计等） |
| `config/` | 全局配置加载（TOML），通过 `config.Cfg` 访问 |

## 中间件

全局注册的中间件（按执行顺序）：

| 中间件 | 来源 | 作用 |
|--------|------|------|
| `RequestLogger` | Echo 内置 | 请求日志 |
| `Recover` | Echo 内置 | panic 恢复，避免进程崩溃 |
| `CORS` | Echo 内置 | 跨域资源共享，允许前端从不同源访问 API |
| `RequestID` | Echo 内置 | 生成请求追踪 ID（`X-Request-ID`） |
| `CostTime` | 自定义 | 记录请求耗时，并写入响应 `cost` 字段 |

> 注：自定义错误处理器 `CustomHTTPErrorHandler` 暂未启用。

## CORS 配置

CORS 跨域支持通过 Echo v5 内置中间件实现，配置变量定义在 `main.go` 中的 `corsConfig`：

```go
var corsConfig = middleware.CORSConfig{
    AllowOrigins:     []string{"*"},                                                      // 允许的来源域名
    AllowMethods:     []string{http.MethodGet, http.MethodPost, ...},                     // 允许的 HTTP 方法
    AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With", ...}, // 允许的请求头
    AllowCredentials: false,                                                               // 是否允许携带凭证（Cookie/Authorization）
    MaxAge:           86400,                                                               // 预检缓存时间（秒）
}
```

**自定义配置**：

| 场景 | 修改方式 |
|------|----------|
| 限制特定域名 | 将 `AllowOrigins` 改为 `[]string{"https://myapp.com"}` |
| 允许凭证模式 | 将 `AllowCredentials` 改为 `true`（此时 `AllowOrigins` 不能为 `*`） |
| 限制 HTTP 方法 | 修改 `AllowMethods` 列表 |
| 调整预检缓存 | 修改 `MaxAge` 值（单位：秒） |

开发环境默认允许所有来源（`*`），生产环境部署时建议按需修改。

## 配置文件

项目使用 TOML 格式的全局配置文件 `config.toml`，在服务启动时自动加载。若文件缺失或格式错误，服务将打印错误信息并拒绝启动（遵循宪法原则 VI：错误及时抛出）。

```toml
[app]
name = "Greeting"
version = "0.1.0"
build_time = "2026-07-14"

[server]
port = ":1323"

[database.mysql]
dsn = "root:@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
# NOTE: changelog is kept in a separate file `changelog.toml` (not here).
```

**配置项说明**：

| 段 | 字段 | 说明 |
|----|------|------|
| `[app]` | `name` | 应用名称 |
| `[app]` | `version` | 应用版本号 |
| `[app]` | `build_time` | 构建时间 |
| `[server]` | `port` | 服务监听端口 |
| `[database.mysql]` | `dsn` | MySQL 连接串 |
| `changelog.toml` | `date` | 更新日志日期（独立文件，不写入 config.toml） |
| `changelog.toml` | `content` | 更新日志内容 |

修改配置后重启服务即可生效。

## 统一响应格式

所有接口统一返回如下 JSON：

```json
{
  "code": 0,
  "message": "",
  "data": {},
  "trace_id": "请求 ID",
  "cost": "处理耗时，如 1.234ms",
  "extra": "可选扩展字段"
}
```

业务错误码定义（见 `response/message.go`）：

| 常量 | 值 | 说明 |
|------|------|------|
| `ErrCodeOk` | `0` | 成功 |
| `ErrCodeCustom` | `100001` | 通用业务错误 |
| `ErrCodeNetwork` | `100002` | 网络错误 |

## API 列表

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/demo/search?tag=a&tag=b` | 接收多值查询参数 `tag` |
| GET | `/demo/err/debug/:str` | 接收路径参数 `str` 并返回 |
| GET | `/demo/sha256?text=hello` | 计算输入文本的 SHA256 哈希值，返回 `{input, hash}` |
| GET | `/demo/user/phone` | 按手机号查询用户，不存在则创建（测试用） |
| POST | `/demo/usr` | 创建 MySQL 用户（phone, realname, username, age） |
| GET | `/demo/usr/:id` | 查询单个 MySQL 用户 |
| PUT | `/demo/usr/:id` | 更新 MySQL 用户（部分字段） |
| DELETE | `/demo/usr/:id` | 软删除 MySQL 用户 |
| GET | `/demo/usrs` | 查询 MySQL 用户列表 |
| POST | `/sqlite/testuser` | 创建 SQLite 测试用户（name, phone, age） |
| GET | `/sqlite/testuser/:id` | 查询单个 SQLite 测试用户 |
| PUT | `/sqlite/testuser/:id` | 更新 SQLite 测试用户（部分字段） |
| DELETE | `/sqlite/testuser/:id` | 软删除 SQLite 测试用户 |
| GET | `/sqlite/testusers` | 查询 SQLite 测试用户列表 |
| GET | `/common/version` | 返回应用版本信息（版本号、构建时间、Go 版本） |
| GET | `/common/changelog` | 返回应用更新日志列表 |
| GET | `/common/setting` | 返回应用配置信息 |

## 本地开发

```bash
# 本地运行
make rundev
# 等价命令
go run main.go
```

启动后访问 `http://localhost:1323`。

## 测试服部署

```bash
make buildqa
```

该命令会：交叉编译为 Linux amd64 二进制，scp 到测试服务器，并通过 `supervisorctl` 重启服务。
（服务器地址与部署路径见 `Makefile`。）

---

## 更新日志

> 本区块用于持续记录功能迭代，后续新增功能请在此追加。

| 添加时间 | 说明 |
|----------|------|
| 2026-06-30 | 初始化项目骨架：分层架构、统一响应封装、Echo 中间件栈 |
| 2026-07-04 | 实现请求耗时统计（`middle.CostTime` + `response.getCost`），修复 `cost` 字段取值 panic |
| 2026-07-13 | 引入 GORM + SQLite：创建 `model/` 目录，全局 `model.DB` 实例，启动时自动初始化 |
| 2026-07-13 | 新增 `/demo/sha256` 接口：接收 `text` 查询参数，返回 SHA256 哈希值及原始输入 |
| 2026-07-14 | 新增 CORS 跨域支持：使用 Echo v5 内置中间件，默认允许所有来源，可配置域名、方法、请求头等 |
| 2026-07-14 | 新增 `/common/*` 公共路由组：`/common/version`、`/common/changelog`、`/common/setting` |
| 2026-07-14 | 新增全局配置文件 `config.toml`（TOML 格式），version/changelog/setting 改为从配置读取 |
| 2026-07-15 | 新增 MySQL 数据库支持：`/demo/usr` CRUD 接口，独立 SQL 迁移脚本 |
| 2026-07-15 | User 模型审计：Name 拆分为 Realname+Username，新增 PasswordHash，Phone 长度修正为 varchar(20) |
| 2026-07-17 | 新增独立 SQLite 实例（`model.SQLiteDB`）与 `/sqlite/testuser` CRUD 测试接口，用户自管 `test_user` 表（宪法原则 VII：应用零自动建表） |

---

## TODO / 待完善

- [ ] 启用自定义错误处理器 `CustomHTTPErrorHandler`，补充 404/500 等错误页
- [ ] 补充更多业务接口与路由分组
- [ ] 增加单元测试
- [ ] 完善 CI / 部署流程
