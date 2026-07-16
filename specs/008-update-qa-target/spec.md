# Feature Specification: 更新 QA 部署目标机器

**Feature Branch**: `008-update-qa-target`

**Created**: 2026-07-16

**Status**: Draft

**Input**: User description: "qa的目标机器更换为ssh ubuntu@111.229.4.203,目录为/opt/project/greeting，supervisor的name为greeting"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 部署到新的 QA 服务器 (Priority: P1)

作为一名运维/后端开发者，我希望通过 `make buildqa` 将编译产物部署到新的 QA 目标服务器 `111.229.4.203`，使用 `ubuntu` 用户和正确的项目目录 `/opt/project/greeting` 以及 supervisor 名称 `greeting`。

**Why this priority**: 这是唯一的用户场景，QA 部署目标是开发者日常部署的核心配置。

**Independent Test**: 配置完成后执行 `make buildqa`，二进制上传至 `ubuntu@111.229.4.203:/opt/project/greeting/` 并正确重启 `greeting` 进程。

**Acceptance Scenarios**:

1. **Given** 项目代码已就绪, **When** 开发者执行 `make buildqa`, **Then** 编译 Linux 二进制，以 `ubuntu` 用户上传至 `111.229.4.203:/opt/project/greeting/`，通过 supervisor 重启 `greeting` 进程
2. **Given** 目标服务器不可达（网络问题或 SSH 未配置）, **When** 开发者执行 `make buildqa`, **Then** 在 scp 步骤失败并显示连接错误，不执行后续重启步骤

---

### Edge Cases

- 目标机器 `/opt/project/greeting` 目录不存在时，scp 是否会报错？是否需要预先创建？
- supervisor 进程名 `greeting` 在目标机器上不存在时，重启命令会失败

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: `buildqa` 目标的 SSH 连接 MUST 使用 `ubuntu@111.229.4.203` 替代原有地址
- **FR-002**: `buildqa` 目标的部署路径 MUST 使用 `/opt/project/greeting` 替代原有路径
- **FR-003**: `buildqa` 目标的 supervisor 进程名 MUST 使用 `greeting` 替代原有占位名称
- **FR-004**: 所有涉及的 SSH/scp 命令 MUST 使用 `ubuntu` 用户（而非原有 `root`）

### Key Entities

不涉及数据实体。仅修改 Makefile 中硬编码的部署参数。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 执行 `make buildqa` 后，二进制成功上传至 `111.229.4.203:/opt/project/greeting/`
- **SC-002**: supervisor 的 `greeting` 进程在部署后成功重启
- **SC-003**: 部署全流程（编译→上传→重启）在 60 秒内完成

## Assumptions

- 目标服务器 `111.229.4.203` 已配置 SSH 免密登录（或开发者持有对应私钥）
- 目标服务器已安装并配置 supervisor，进程名为 `greeting`
- 目标路径 `/opt/project/greeting` 已存在
- 此变更仅涉及 Makefile 中的硬编码参数替换，不涉及架构重构（架构重构属于 `007-optimize-makefile`）
- 保留 `buildqa` 目标名称，不在此次变更中重命名为 `deploy-qa`
