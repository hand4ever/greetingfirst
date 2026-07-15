# 006: Audit User Model Against Actual Database Schema

## Overview

将 `model.User` 结构体及关联代码与数据库实际建表语句对齐。

### 实际数据库建表语句

```sql
create table users
(
    id            int auto_increment primary key,
    phone         varchar(20)  not null comment 'cellphone number',
    realname      varchar(100) null,
    username      varchar(20)  null,
    age           int          null,
    password_hash varchar(200) null,
    created_at    datetime     null,
    updated_at    datetime     null,
    deleted_at    datetime     null
);
```

## 差异分析

| 字段 | 当前 model.User | 实际建表 | 差异 |
|------|----------------|----------|------|
| `Name` | `gorm:"type:varchar(64);not null" json:"name"` | 不存在 | 应拆分为 `realname` + `username` |
| `realname` | 缺失 | `varchar(100) null` | 新增 |
| `username` | 缺失 | `varchar(20) null` | 新增 |
| `password_hash` | 缺失 | `varchar(200) null` | 新增 |
| `Phone` | `type:varchar(32)` | `varchar(20)` | 长度不一致 |
| `Age` | `default:0` | `null` | 默认值和可空性不一致 |
| `DeletedAt` | `gorm.DeletedAt` | `datetime` | 改为 `*time.Time`，移除 GORM 软删除 |
| `ID` | `uint` | `int` | 改为 `int`，与 DB 类型完全一致 |

## 需要变更的文件

### 1. `model/user.go`
- User 结构体字段重定义：`Name` → `Realname` + `Username`，新增 `PasswordHash`
- `Phone` 长度从 32 → 20
- `ID` 从 `uint` → `int`
- `DeletedAt` 从 `gorm.DeletedAt` → `*time.Time`
- 删除相关 CRUD 函数需改为手动软删除：`DeleteUser` 更新 `deleted_at`，`GetAllUsers`/`GetUserByID`/`GetUserByPhone` 添加 `WHERE deleted_at IS NULL` 过滤

### 2. `entity/user/user.go`
- `UserCreateReq.Name` → `Realname` + `Username`
- `UserUpdateReq.Name` → `Realname` + `Username`

### 3. `handler/user.go`
- 所有引用 `.Name` 的地方改为 `.Realname`（或按语义拆分）

### 4. `handler/demo.go`
- `GetUserByPhoneTest` 中的 `Name: "test_user"` → 对应字段调整

### 5. `migrations/001_user.mysql.sql`
- 重建为与实际建表一致的 DDL

### 6. 测试文件
- `model/user_test.go`、`handler/user_test.go` 中所有 `.Name` 引用改为 `.Realname`

### 7. `README.md` / `api.http`
- API 文档更新字段名

## Clarifications

### Session 2026-07-15

- Q: 数据库迁移策略：原 Name 字段改为 Realname + Username，旧数据如何处理？→ A: 无需迁移脚本，数据库表已与实际建表语句一致，表中不存在 name 列（DBA 已提前改好表结构为 realname + username）。
- Q: `DeletedAt` 类型：DB 为 `datetime`，模型用 `gorm.DeletedAt`，是否对齐？→ A: 改为 `*time.Time`，不使用 GORM 软删除封装。
- Q: `password_hash` 字段使用范围：是否需要 CRUD 接口支持？→ A: 仅表结构对齐，CRUD 接口暂不支持读写 password_hash，后续按需扩展。
- Q: `ID` 类型：DB 用 `int`，模型用 `uint`，是否改为 `int`？→ A: 改为 `int`，与 DB 类型完全一致。
- Q: 删除行为：`DeletedAt` 改为 `*time.Time` 后，删除策略？→ A: 手动软删除，更新 `deleted_at = NOW()`，查询加 `WHERE deleted_at IS NULL`。
- Q: `Name` 字段如何拆分？原业务中 `Name` 对应的是真实姓名还是用户名？→ A: 删除 Name，重新添加 Realname 和 Username 两个独立字段。
- Q: `Age` 是否应改为 `*int`（可空指针）以匹配表定义？→ A: 保持 int，不改为指针。
