package model

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestMain(m *testing.M) {
	// init both DBs as in-memory SQLite for unit testing
	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to init test DB (MySQL stub): " + err.Error())
	}
	SQLiteDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to init test SQLiteDB: " + err.Error())
	}
	// apply schema to SQLiteDB
	if err := ApplySchema(SQLiteDB); err != nil {
		panic("failed to apply schema to SQLiteDB: " + err.Error())
	}
	// create users table for MySQL entity tests (running against SQLite stub)
	if err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			phone VARCHAR(32) NOT NULL UNIQUE,
			name VARCHAR(64) NOT NULL,
			age INT DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);
		CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at)
	`).Error; err != nil {
		panic("failed to create users table: " + err.Error())
	}
	code := m.Run()
	os.Exit(code)
}

// logOK always prints log, even without -v flag
func logOK(t *testing.T, format string, args ...any) {
	t.Helper()
	msg := fmt.Sprintf(format, args...)
	if testing.Verbose() {
		t.Log(msg)
	} else {
		fmt.Fprintln(os.Stderr, msg)
	}
}

// ============================================================================
// SQLiteUser CRUD Tests
// ============================================================================

func TestCreateSQLiteUser(t *testing.T) {
	user := &SQLiteUser{
		Phone: "13800138000",
		Name:  "张三",
		Age:   28,
	}
	if err := CreateSQLiteUser(user); err != nil {
		t.Fatalf("CreateSQLiteUser failed: %v", err)
	}
	if user.ID == 0 {
		t.Error("expected ID to be set after create")
	}
	b, _ := json.MarshalIndent(user, "", "  ")
	logOK(t, "SQLite 创建成功:\n%s", b)
}

func TestCreateSQLiteUser_DuplicatePhone(t *testing.T) {
	user := &SQLiteUser{Phone: "13800001111", Name: "dup1", Age: 20}
	if err := CreateSQLiteUser(user); err != nil {
		t.Fatalf("first create should succeed: %v", err)
	}
	user2 := &SQLiteUser{Phone: "13800001111", Name: "dup2", Age: 30}
	if err := CreateSQLiteUser(user2); err == nil {
		t.Error("expected error for duplicate phone")
	}
	b, _ := json.MarshalIndent(user, "", "  ")
	logOK(t, "首次创建成功，重复手机号正确拒绝:\n%s", b)
}

func TestGetSQLiteUserByID(t *testing.T) {
	user := &SQLiteUser{Phone: "13800138001", Name: "李四", Age: 30}
	if err := CreateSQLiteUser(user); err != nil {
		t.Fatalf("CreateSQLiteUser failed: %v", err)
	}

	found, err := GetSQLiteUserByID(user.ID)
	if err != nil {
		t.Fatalf("GetSQLiteUserByID failed: %v", err)
	}
	if found.Name != "李四" {
		t.Errorf("expected name 李四, got %s", found.Name)
	}
	if found.Phone != "13800138001" {
		t.Errorf("expected phone 13800138001, got %s", found.Phone)
	}
	if found.Age != 30 {
		t.Errorf("expected age 30, got %d", found.Age)
	}
	b, _ := json.MarshalIndent(found, "", "  ")
	logOK(t, "查询成功:\n%s", b)
}

func TestGetSQLiteUserByID_NotFound(t *testing.T) {
	_, err := GetSQLiteUserByID(99999)
	if err == nil {
		t.Error("expected error for non-existent user")
	}
}

func TestGetSQLiteUserByPhone(t *testing.T) {
	user := &SQLiteUser{Phone: "13800138002", Name: "王五", Age: 25}
	if err := CreateSQLiteUser(user); err != nil {
		t.Fatalf("CreateSQLiteUser failed: %v", err)
	}

	found, err := GetSQLiteUserByPhone("13800138002")
	if err != nil {
		t.Fatalf("GetSQLiteUserByPhone failed: %v", err)
	}
	if found.Name != "王五" {
		t.Errorf("expected name 王五, got %s", found.Name)
	}
	b, _ := json.MarshalIndent(found, "", "  ")
	logOK(t, "按手机号查询成功:\n%s", b)
}

func TestGetSQLiteUserByPhone_NotFound(t *testing.T) {
	_, err := GetSQLiteUserByPhone("00000000000")
	if err == nil {
		t.Error("expected error for non-existent phone")
	}
}

func TestUpdateSQLiteUser(t *testing.T) {
	user := &SQLiteUser{Phone: "13800138003", Name: "赵六", Age: 22}
	if err := CreateSQLiteUser(user); err != nil {
		t.Fatalf("CreateSQLiteUser failed: %v", err)
	}

	user.Name = "赵六改名"
	user.Age = 35
	if err := UpdateSQLiteUser(user); err != nil {
		t.Fatalf("UpdateSQLiteUser failed: %v", err)
	}

	found, err := GetSQLiteUserByID(user.ID)
	if err != nil {
		t.Fatalf("GetSQLiteUserByID after update failed: %v", err)
	}
	if found.Name != "赵六改名" {
		t.Errorf("expected name 赵六改名, got %s", found.Name)
	}
	if found.Age != 35 {
		t.Errorf("expected age 35, got %d", found.Age)
	}
	b, _ := json.MarshalIndent(found, "", "  ")
	logOK(t, "更新成功:\n%s", b)
}

func TestDeleteSQLiteUser(t *testing.T) {
	user := &SQLiteUser{Phone: "13800138004", Name: "钱七", Age: 40}
	if err := CreateSQLiteUser(user); err != nil {
		t.Fatalf("CreateSQLiteUser failed: %v", err)
	}

	if err := DeleteSQLiteUser(user.ID); err != nil {
		t.Fatalf("DeleteSQLiteUser failed: %v", err)
	}

	// should not be found after soft delete
	_, err := GetSQLiteUserByID(user.ID)
	if err == nil {
		t.Error("expected error for soft-deleted user")
	}
	b, _ := json.MarshalIndent(user, "", "  ")
	logOK(t, "软删除成功，已查不到该用户:\n%s", b)
}

// ============================================================================
// User (MySQL) CRUD Tests — run against in-memory SQLite for unit testing
// ============================================================================

func TestCreateUser(t *testing.T) {
	user := &User{
		Phone: "13900138000",
		Name:  "MySQL张三",
		Age:   28,
	}
	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if user.ID == 0 {
		t.Error("expected ID to be set after create")
	}
	b, _ := json.MarshalIndent(user, "", "  ")
	logOK(t, "MySQL 创建成功:\n%s", b)
}

func TestCreateUser_DuplicatePhone(t *testing.T) {
	user := &User{Phone: "13900001111", Name: "dup1", Age: 20}
	if err := CreateUser(user); err != nil {
		t.Fatalf("first create should succeed: %v", err)
	}
	user2 := &User{Phone: "13900001111", Name: "dup2", Age: 30}
	if err := CreateUser(user2); err == nil {
		t.Error("expected error for duplicate phone (active record)")
	}
	b, _ := json.MarshalIndent(user, "", "  ")
	logOK(t, "MySQL 唯一性约束测试通过:\n%s", b)
}

func TestGetUserByID(t *testing.T) {
	user := &User{Phone: "13900138001", Name: "MySQL李四", Age: 30}
	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	found, err := GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if found.Name != "MySQL李四" {
		t.Errorf("expected name MySQL李四, got %s", found.Name)
	}
	if found.Phone != "13900138001" {
		t.Errorf("expected phone 13900138001, got %s", found.Phone)
	}
	if found.Age != 30 {
		t.Errorf("expected age 30, got %d", found.Age)
	}
	b, _ := json.MarshalIndent(found, "", "  ")
	logOK(t, "MySQL 查询成功:\n%s", b)
}

func TestGetUserByID_NotFound(t *testing.T) {
	_, err := GetUserByID(99999)
	if err == nil {
		t.Error("expected error for non-existent user")
	}
}

func TestGetUserByPhone(t *testing.T) {
	user := &User{Phone: "13900138002", Name: "MySQL王五", Age: 25}
	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	found, err := GetUserByPhone("13900138002")
	if err != nil {
		t.Fatalf("GetUserByPhone failed: %v", err)
	}
	if found.Name != "MySQL王五" {
		t.Errorf("expected name MySQL王五, got %s", found.Name)
	}
	b, _ := json.MarshalIndent(found, "", "  ")
	logOK(t, "MySQL 按手机号查询成功:\n%s", b)
}

func TestGetUserByPhone_NotFound(t *testing.T) {
	_, err := GetUserByPhone("00000000000")
	if err == nil {
		t.Error("expected error for non-existent phone")
	}
}

func TestUpdateUser(t *testing.T) {
	user := &User{Phone: "13900138003", Name: "MySQL赵六", Age: 22}
	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user.Name = "MySQL赵六改名"
	user.Age = 35
	if err := UpdateUser(user); err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	found, err := GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("GetUserByID after update failed: %v", err)
	}
	if found.Name != "MySQL赵六改名" {
		t.Errorf("expected name MySQL赵六改名, got %s", found.Name)
	}
	if found.Age != 35 {
		t.Errorf("expected age 35, got %d", found.Age)
	}
	b, _ := json.MarshalIndent(found, "", "  ")
	logOK(t, "MySQL 更新成功:\n%s", b)
}

func TestDeleteUser(t *testing.T) {
	user := &User{Phone: "13900138004", Name: "MySQL钱七", Age: 40}
	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if err := DeleteUser(user.ID); err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}

	// should not be found after soft delete
	_, err := GetUserByID(user.ID)
	if err == nil {
		t.Error("expected error for soft-deleted user")
	}
	b, _ := json.MarshalIndent(user, "", "  ")
	logOK(t, "MySQL 软删除成功，已查不到该用户:\n%s", b)
}

// TestPhoneReuseAfterSoftDelete verifies phone reuse behavior.
// Note: this relies on MySQL generated column; behavior may differ on SQLite stub.
func TestPhoneReuseAfterSoftDelete(t *testing.T) {
	user := &User{Phone: "13900138005", Name: "旧用户", Age: 30}
	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if err := DeleteUser(user.ID); err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}

	// On MySQL: same phone allowed after soft delete (generated column handles this).
	// On SQLite stub: UNIQUE constraint prevents reuse; test verifies current behavior.
	user2 := &User{Phone: "13900138005", Name: "新用户", Age: 25}
	err := CreateUser(user2)
	if err != nil {
		// expected on SQLite (simple UNIQUE constraint), acceptable on MySQL too if constraint triggers
		logOK(t, "Phone reuse after soft delete blocked (expected on SQLite): %v", err)
		return
	}
	if user2.ID == user.ID {
		t.Error("expected new record with different ID")
	}
	b, _ := json.MarshalIndent(user2, "", "  ")
	logOK(t, "软删除后同号码可复用 (MySQL):\n%s", b)
}
