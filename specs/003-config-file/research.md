# Research: 全局配置文件

## Decision 1: TOML Configuration Format

**Decision**: Use TOML as the configuration file format.

**Rationale**:
- TOML is explicitly designed for configuration files (unlike YAML which is a general-purpose serialization format)
- No indentation sensitivity — eliminates the most common source of YAML parse errors
- `github.com/BurntSushi/toml` is the de facto standard Go TOML library, well-maintained and widely used
- The project's config structure is flat-to-moderately nested (app/server/database/changelog), which TOML handles naturally
- Learning curve is near zero — `[section]` and `key = "value"` are self-explanatory

**Alternatives considered**:

| Alternative | Pros | Cons | Verdict |
|-------------|------|------|---------|
| YAML (`go-yaml/yaml.v3`) | K8s ecosystem familiarity, supports complex nesting | Indentation-sensitive, anchors/aliases add cognitive overhead, heavier library | Rejected — overkill for flat config |
| JSON | No external dependency, every language supports it | No comments, verbose syntax, not human-friendly for editing | Rejected — poor DX for config editing |
| Environment variables | No file needed, 12-factor app compliant | No structure, hard to manage many settings, no arrays | Rejected — insufficient for changelog arrays |
| Go flags (`flag` package) | Standard library, no dependency | Only scalar values, verbose for many settings, not file-based | Rejected — complementary, not alternative |

## Decision 2: Degraded Mode on Config Failure

**Decision**: Return `nil` (no error) when config file is missing or TOML parse fails; print a warning to stdout and continue with built-in defaults.

**Rationale**:
- Aligns with "Copy-Ready Template" principle — new project clones should start without a config file
- Prevents startup panic from a missing/malformed config (only truly unrecoverable errors should panic per constitution)
- Built-in `defaultConfig()` provides the same values as the shipped `config.toml`, ensuring consistency
- The `InitConfig()` signature `error` return allows callers (like tests) to handle errors differently if needed

**Alternatives considered**:
- Panic on config failure → rejected: violates constitution ("Panic MUST NOT 用于常规业务错误") and breaks copy-ready guarantee
- Return error and let `main.go` decide → partially adopted: `main.go` panics only on non-parse read errors (e.g., permission denied), while missing/invalid files are silent

## Decision 3: Singleton Global Config

**Decision**: Expose `var Cfg *Config` as a package-level variable, initialized with defaults before `InitConfig()` runs.

**Rationale**:
- Consistent with existing patterns in the project (`model.DB` is a similar global singleton)
- Zero-copy access: handlers dereference `config.Cfg` directly without parameter passing
- Thread-safe by design: config is set once at startup before the HTTP server starts, never mutated afterward
- `defaultConfig()` initializes `Cfg` at package init time, so it's never nil

**Alternatives considered**:
- Dependency injection → rejected: adds complexity without benefit for a process-lifetime singleton
- `sync.Once` lazy init → rejected: unnecessary — config is always loaded synchronously at startup

## Decision 4: Config Path Resolution (Deferred)

**Decision**: Hardcode `"config.toml"` in `main.go`. Command-line flag (`-config`) and environment variable (`CONFIG_PATH`) are deferred to a future iteration per specification clarification.

**Rationale**: Simplicity first; `InitConfig(configPath string)` already accepts a custom path for testing flexibility. Runtime path override will be added when deployment scenarios demand it.

## Dependency Impact

Single new dependency added:
- `github.com/BurntSushi/toml` v1.6.0 — pure Go, no transitive dependencies beyond what's already in the module graph
