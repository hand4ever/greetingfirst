# 快速入门(Quickstart): 部署用户变量化与 QA 快捷目标

**功能(Feature)**: 009-deploy-usr-shortcut
**日期(Date)**: 2026-07-16

> 本文档提供变更后的验证场景。由于 Makefile 变更不涉及 Go 代码，以下场景聚焦于 Makefile 行为验证。详细信息参见 [contracts/makefile-targets.md](./contracts/makefile-targets.md) 和 [data-model.md](./data-model.md)。

---

## 前置条件(Prerequisites)

- GNU Make（`make --version`）
- Go ≥ 1.22（`go version`，用于 `build-linux` 编译）
- （可选）实际 QA 服务器访问权限（用于完整部署验证）

---

## 验证场景(Validation Scenarios)

### 1. `deploy-qa` 零参数使用 QA 默认值（FR-001, SC-001）

```bash
# 测试(Test): deploy-qa 不传任何参数 — 应使用 QA 默认值
make deploy-qa
```

**预期结果(Expected)**:
- `DEPLOY_HOST=111.229.4.203`、`DEPLOY_USR=ubuntu`、`DEPLOY_PATH=/opt/project/greeting`、`DEPLOY_SUPERVISOR=greeting`
- scp 命令显示 `ubuntu@111.229.4.203`
- 编译、上传并重启 greeting 服务

---

### 2. `DEPLOY_USR` 通过命令行覆盖（FR-001, SC-002）

```bash
# 测试(Test): 覆盖 SSH 用户为 root
make deploy-qa DEPLOY_USR=root 2>&1 | head -5
```

**预期结果(Expected)**:
- scp 命令中显示 `root@111.229.4.203`（用户被覆盖，HOST 保持默认）

---

### 3. 任意参数可覆盖（FR-001, FR-006）

```bash
# 测试(Test): 完全自定义部署
make deploy-qa DEPLOY_HOST=prod.example.com DEPLOY_USR=root DEPLOY_PATH=/app DEPLOY_SUPERVISOR=myapp 2>&1 | head -5
```

**预期结果(Expected)**:
- scp 命令使用 `root@prod.example.com:/app/`
- 所有默认值均被覆盖

---

### 4. `check_required` 用法提示已更新（FR-003）

```bash
# 测试(Test): 空变量触发 check_required 错误并显示更新后的提示
make deploy-qa DEPLOY_HOST=
```

**预期结果(Expected)**:
```
Error: DEPLOY_HOST is not set.
  Usage: make deploy-qa [DEPLOY_HOST=<host>] [DEPLOY_USR=<user>] [DEPLOY_PATH=<path>] [DEPLOY_SUPERVISOR=<name>]
```

---

### 5. `runqa` 目标等价于 `deploy-qa`（FR-004, SC-003）

```bash
# 测试(Test): runqa 不传任何参数 — 应完全等价于 deploy-qa
make runqa
```

**预期结果(Expected)**:
- 等价于执行 `make deploy-qa`
- 使用所有 QA 默认值，编译上传并重启 greeting 服务

---

### 6. `runqa` 参数可覆盖（FR-006）

```bash
# 测试(Test): 在 runqa 中覆盖 DEPLOY_HOST
make runqa DEPLOY_HOST=other-server 2>&1 | head -5
```

**预期结果(Expected)**:
- scp 命令使用 `ubuntu@other-server`（仅 HOST 被覆盖，其余保持默认）

---

### 7. `runqa` 已在 `.PHONY` 中声明（FR-005）

```bash
# 测试(Test): 验证 runqa 声明为 PHONY
grep -E '^\.PHONY.*runqa' Makefile
```

**预期结果(Expected)**: 存在包含 `runqa` 的 `.PHONY` 声明行

---

### 8. `make help` 显示 `runqa` 目标（FR-007, SC-004）

```bash
make help
```

**预期结果(Expected)**:
- 在「部署 (Deploy)」分组下显示：
  ```
  make runqa         QA 一键部署 (Quick deploy to QA server)
  ```

---

### 9. QA 服务器不可达时的错误处理（用户故事(User Story) 2, 场景(Scenario) 2）

```bash
# 测试(Test): 模拟 QA 服务器不可达（替换为不存在的主机）
make deploy-qa DEPLOY_HOST=192.0.2.1 2>&1
```

**预期结果(Expected)**:
- scp 步骤失败并显示连接错误
- 不执行后续 supervisorctl restart 步骤
- 退出码非零

---

## 完整验证清单(Complete Verification Checklist)

- [ ] `deploy-qa` 零参数使用 QA 默认值 (`ubuntu@111.229.4.203`)
- [ ] `DEPLOY_USR` 可通过命令行覆盖 (`make deploy-qa DEPLOY_USR=root`)
- [ ] 所有变量可同时被覆盖
- [ ] `check_required` 用法提示展示所有参数
- [ ] `runqa` 等价于 `deploy-qa`
- [ ] `runqa` 参数可通过命令行覆盖
- [ ] `runqa` 已声明 `.PHONY`
- [ ] `make help` 显示 `runqa` 并附有双语描述
- [ ] 服务器不可达时部署优雅失败
