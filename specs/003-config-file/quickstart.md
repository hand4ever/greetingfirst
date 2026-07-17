# Quickstart: 全局配置文件

## Prerequisites

- Go 1.22+
- Project dependencies: `go mod tidy`
- `config.toml` in project root (shipped with defaults)

## Validation Scenarios

### Scenario 1: Normal startup with config file

**Setup**: Ensure `config.toml` exists with valid TOML content.

**Run**:
```bash
go run main.go
```

**Expected**:
- Log: `[config] loaded config from config.toml`
- Service starts on the port from `config.toml` `[server] port` (default `:1323`)
- `GET /common/version` returns `version` from config
- `GET /common/changelog` returns `[[changelog]]` entries from config
- `GET /common/setting` returns `[app]`, `[server]`, `[database]` values from config

**Verify**:
```bash
# Version matches config.toml [app] section
curl -s http://localhost:1323/common/version | jq '.data.version'
# → "0.1.0"

# Changelog count matches config.toml [[changelog]] entries
curl -s http://localhost:1323/common/changelog | jq '.data | length'
# → 5
```

### Scenario 2: Missing config file (degraded mode)

**Setup**: Rename or delete `config.toml`.

**Run**:
```bash
mv config.toml config.toml.bak
go run main.go
```

**Expected**:
- Log: `[config] config file not found at config.toml, using defaults`
- Service starts normally (no panic)
- All three common endpoints return built-in defaults (same values as shipped `config.toml`)

**Cleanup**:
```bash
mv config.toml.bak config.toml
```

### Scenario 3: Invalid config file (degraded mode)

**Setup**: Corrupt `config.toml` (e.g., missing closing quote).

**Run**:
```bash
echo 'version = "bad' > config.toml.bak
mv config.toml config.toml.good
mv config.toml.bak config.toml
go run main.go
```

**Expected**:
- Log: `[config] failed to parse config.toml: ... using defaults`
- Service starts normally (no panic)
- All endpoints return defaults

**Cleanup**:
```bash
mv config.toml.good config.toml
```

### Scenario 4: Config changes take effect after restart

**Setup**: Edit `config.toml` to change version.

**Run**:
```bash
# Edit version in config.toml to "2.0.0"
# Restart service
go run main.go &
sleep 1
curl -s http://localhost:1323/common/version | jq '.data.version'
# → "2.0.0"
```

**Expected**: Version reflects the edited value after restart.

### Scenario 5: Run all tests

```bash
go test -v ./... -count=1
```

**Expected**: All tests pass, including any config-related tests in `config/`.

## API Contracts

See `specs/003-config-file/data-model.md` for entity definitions.

| Endpoint | Method | Response Data |
|----------|--------|---------------|
| `/common/version` | GET | `VersionResponse` { version, build_time, go_version } |
| `/common/changelog` | GET | `[]ChangelogEntry` [{ date, content }] |
| `/common/setting` | GET | `[]SettingItem` [{ key, value, description }] |
