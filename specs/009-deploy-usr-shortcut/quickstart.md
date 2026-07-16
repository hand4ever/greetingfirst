# Quickstart: 部署用户变量化与 QA 快捷目标

**Feature**: 009-deploy-usr-shortcut
**Date**: 2026-07-16

> 本文档提供变更后的验证场景。由于 Makefile 变更不涉及 Go 代码，以下场景聚焦于 Makefile 行为验证。详细信息参见 [contracts/makefile-targets.md](./contracts/makefile-targets.md) 和 [data-model.md](./data-model.md)。

---

## Prerequisites

- GNU Make（`make --version`）
- Go ≥ 1.22（`go version`，用于 `build-linux` 编译）
- （可选）实际 QA 服务器访问权限（用于完整部署验证）

---

## Validation Scenarios

### 1. `deploy-qa` 零参数使用 QA 默认值（FR-001, SC-001）

```bash
# Test: deploy-qa without any args — should use QA defaults
make deploy-qa
```

**Expected**:
- `DEPLOY_HOST=111.229.4.203`、`DEPLOY_USR=ubuntu`、`DEPLOY_PATH=/opt/project/greeting`、`DEPLOY_SUPERVISOR=greeting`
- scp 命令显示 `ubuntu@111.229.4.203`
- 编译、上传并重启 greeting 服务

---

### 2. `DEPLOY_USR` 通过命令行覆盖（FR-001, SC-002）

```bash
# Test: override SSH user to root
make deploy-qa DEPLOY_USR=root 2>&1 | head -5
```

**Expected**:
- scp 命令中显示 `root@111.229.4.203`（用户被覆盖，HOST 保持默认）

---

### 3. 任意参数可覆盖（FR-001, FR-006）

```bash
# Test: full custom deployment
make deploy-qa DEPLOY_HOST=prod.example.com DEPLOY_USR=root DEPLOY_PATH=/app DEPLOY_SUPERVISOR=myapp 2>&1 | head -5
```

**Expected**:
- scp 命令使用 `root@prod.example.com:/app/`
- 所有默认值均被覆盖

---

### 4. `check_required` 用法提示已更新（FR-003）

```bash
# Test: empty variable triggers check_required error with updated hint
make deploy-qa DEPLOY_HOST=
```

**Expected**:
```
Error: DEPLOY_HOST is not set.
  Usage: make deploy-qa [DEPLOY_HOST=<host>] [DEPLOY_USR=<user>] [DEPLOY_PATH=<path>] [DEPLOY_SUPERVISOR=<name>]
```

---

### 5. `runqa` 目标等价于 `deploy-qa`（FR-004, SC-003）

```bash
# Test: runqa without any parameters — should behave exactly like deploy-qa
make runqa
```

**Expected**:
- 等价于执行 `make deploy-qa`
- 使用所有 QA 默认值，编译上传并重启 greeting 服务

---

### 6. `runqa` 参数可覆盖（FR-006）

```bash
# Test: override DEPLOY_HOST in runqa
make runqa DEPLOY_HOST=other-server 2>&1 | head -5
```

**Expected**:
- scp 命令使用 `ubuntu@other-server`（仅 HOST 被覆盖，其余保持默认）

---

### 7. `runqa` 已在 `.PHONY` 中声明（FR-005）

```bash
# Test: verify runqa is PHONY
grep -E '^\.PHONY.*runqa' Makefile
```

**Expected**: 存在包含 `runqa` 的 `.PHONY` 声明行

---

### 8. `make help` 显示 `runqa` 目标（FR-007, SC-004）

```bash
make help
```

**Expected**:
- 在「部署 (Deploy)」分组下显示：
  ```
  make runqa         QA 一键部署 (Quick deploy to QA server)
  ```

---

### 9. QA 服务器不可达时的错误处理（User Story 2, Scenario 2）

```bash
# Test: simulate unreachable QA server (replace with non-existent host)
make deploy-qa DEPLOY_HOST=192.0.2.1 2>&1
```

**Expected**:
- scp 步骤失败并显示连接错误
- 不执行后续 supervisorctl restart 步骤
- 退出码非零

---

## Complete Verification Checklist

- [ ] `deploy-qa` with no args uses QA defaults (`ubuntu@111.229.4.203`)
- [ ] `DEPLOY_USR` can be overridden via CLI (`make deploy-qa DEPLOY_USR=root`)
- [ ] All vars can be overridden simultaneously
- [ ] `check_required` usage hint shows all params
- [ ] `runqa` is equivalent to `deploy-qa`
- [ ] `runqa` parameters can be overridden from CLI
- [ ] `runqa` is declared `.PHONY`
- [ ] `make help` displays `runqa` with bilingual description
- [ ] Deploy fails gracefully when server is unreachable
