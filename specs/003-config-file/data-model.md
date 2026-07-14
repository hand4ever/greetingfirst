# Data Model: 全局配置文件

## Entities

### Config (top-level)

| Field | Type | TOML Key | Required | Default | Description |
|-------|------|----------|----------|---------|-------------|
| App | `AppConfig` | `[app]` | No | See below | Application metadata |
| Server | `ServerConfig` | `[server]` | No | See below | Server configuration |
| Database | `DatabaseConfig` | `[database]` | No | See below | Database configuration |
| Changelog | `[]ChangelogConfig` | `[[changelog]]` | No | See below | Changelog entries (array of tables) |

**Identity**: Singleton — one instance per process, stored in `config.Cfg`.

**Lifecycle**:
```
package init  →  Cfg = defaultConfig()     (safe defaults)
main()        →  config.InitConfig(path)   (overlay from TOML file or keep defaults)
runtime       →  read-only access via config.Cfg
```

### AppConfig

| Field | Type | TOML Key | Required | Default |
|-------|------|----------|----------|---------|
| Name | `string` | `name` | No | `"Greeting"` |
| Version | `string` | `version` | No | `"0.1.0"` |
| BuildTime | `string` | `build_time` | No | `"unknown"` |

### ServerConfig

| Field | Type | TOML Key | Required | Default |
|-------|------|----------|----------|---------|
| Port | `string` | `port` | No | `":1323"` |

### DatabaseConfig

| Field | Type | TOML Key | Required | Default |
|-------|------|----------|----------|---------|
| Type | `string` | `type` | No | `"sqlite"` |
| DSN | `string` | `dsn` | No | `"greeting.db"` |

### ChangelogConfig

| Field | Type | TOML Key | Required | Default |
|-------|------|----------|----------|---------|
| Date | `string` | `date` | No | `""` |
| Content | `string` | `content` | No | `""` |

**Note**: The `ChangelogConfig` type in `config/` is separate from `entity/common.ChangelogEntry`. The handler converts between them when building the API response, isolating the config representation from the API contract.

## Relationships

```
Config (1) ─── (1) AppConfig
Config (1) ─── (1) ServerConfig
Config (1) ─── (1) DatabaseConfig
Config (1) ─── (N) ChangelogConfig
```

No database persistence — all data is in-memory, sourced from the TOML file at startup.

## Validation Rules

- All fields have defaults; no field is strictly required
- TOML parse errors are caught by `toml.Unmarshal` → triggers fallback to defaults
- Empty config file → `toml.Unmarshal` succeeds with zero values → defaults apply
- Missing sections (e.g., no `[server]` block) → Go zero values for that sub-struct
- Changelog array may be empty (`[[changelog]]` omitted) → API returns `[]`
- Non-UTF-8 encoding → `os.ReadFile` + `toml.Unmarshal` will fail → fallback to defaults
