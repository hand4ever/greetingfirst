# Makefile Target Contracts

**Feature**: 009-deploy-usr-shortcut
**Date**: 2026-07-16
**Version**: 2.0

> 本文档基于 007-optimize-makefile 的契约文档，记录本次变更涉及的目标契约。仅列出变更的目标；未变更目标（`help`、`rundev`、`fmt`、`lint`、`build`、`build-linux`、`test`、`clean`）的契约参见 [specs/007-optimize-makefile/contracts/makefile-targets.md](../../007-optimize-makefile/contracts/makefile-targets.md)。

---

## Contract: `deploy-qa` (MODIFIED)

**Call**: `make deploy-qa [DEPLOY_HOST=<host>] [DEPLOY_USR=<user>] [DEPLOY_PATH=<path>] [DEPLOY_SUPERVISOR=<name>]`

**Defaults**:
| Variable | Default |
|----------|---------|
| `DEPLOY_HOST` | `111.229.4.203` |
| `DEPLOY_USR` | `ubuntu` |
| `DEPLOY_PATH` | `/opt/project/greeting` |
| `DEPLOY_SUPERVISOR` | `greeting` |

**Preconditions**: ssh/scp access configured for `$(DEPLOY_USR)@$(DEPLOY_HOST)`; all variables have defaults so zero-arg invocation works

**Behavior**:
1. Validate required variables (`DEPLOY_HOST`, `DEPLOY_SUPERVISOR`)
   - Triggered only if user explicitly unsets them; with defaults this should not occur in normal use
2. Invoke `build-linux` (implicit dependency)
3. Print: `Uploading to $(DEPLOY_HOST):$(DEPLOY_PATH)/ ...`
4. Execute: `scp -O $(BIN_DIR)/$(BIN_NAME) $(DEPLOY_USR)@$(DEPLOY_HOST):$(DEPLOY_PATH)/`
5. Print: `Restarting service $(DEPLOY_SUPERVISOR) on $(DEPLOY_HOST)...`
6. Execute: `ssh $(DEPLOY_USR)@$(DEPLOY_HOST) "supervisorctl restart $(DEPLOY_SUPERVISOR)"`
7. Print: `Deploy complete!`

**Changes from 007**:

| Component | 007 (Before) | 009 (After) |
|-----------|-------------|-------------|
| `DEPLOY_HOST` default | _(empty)_ | `111.229.4.203` |
| `DEPLOY_USR` variable | _(not exist)_ | `ubuntu` |
| `DEPLOY_PATH` default | `/opt/src/main` | `/opt/project/greeting` |
| `DEPLOY_SUPERVISOR` default | _(empty)_ | `greeting` |
| scp command | `root@$(DEPLOY_HOST)` | `$(DEPLOY_USR)@$(DEPLOY_HOST)` |
| ssh command | `root@$(DEPLOY_HOST)` | `$(DEPLOY_USR)@$(DEPLOY_HOST)` |
| Usage hint | `DEPLOY_HOST=<host> DEPLOY_SUPERVISOR=<name> [DEPLOY_PATH=<path>]` | `[DEPLOY_HOST=<host>] [DEPLOY_USR=<user>] [DEPLOY_PATH=<path>] [DEPLOY_SUPERVISOR=<name>]` |

**Parallel Execution**: 不提供锁机制（同 007）。

**Exit codes**:

| Code | Meaning |
|------|---------|
| 0 | Deploy successful |
| 1 | Pre-check failed or scp/ssh failed |

---

## Contract: `runqa` (NEW)

**Call**: `make runqa [DEPLOY_HOST=<host>] [DEPLOY_USR=<user>] [DEPLOY_PATH=<path>] [DEPLOY_SUPERVISOR=<name>]`

**Preconditions**: Same as `deploy-qa` (all variables have QA defaults)

**Behavior**:
1. Invoke `deploy-qa` directly (no extra logic)
2. All variable defaults and override behavior inherited from `deploy-qa`

**Equivalent to** (when no overrides):
```bash
make deploy-qa
```

**Override example**:
```bash
# Override host only, keep other QA defaults
make runqa DEPLOY_HOST=other-qa-server

# Full custom deployment
make runqa DEPLOY_HOST=prod.example.com DEPLOY_USR=root DEPLOY_PATH=/app DEPLOY_SUPERVISOR=app
```

**Exit codes**: Same as `deploy-qa` (inherited).

---

## Summary of Changes

| Contract | Status |
|----------|--------|
| `deploy-qa` | MODIFIED — variable defaults updated + 2 command lines (`root@` → `$(DEPLOY_USR)@`) + usage hint updated |
| `runqa` | NEW — alias for `deploy-qa` |
| All others | UNCHANGED — see 007 contracts |
