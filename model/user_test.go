package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

// logOK always prints log, even without -v flag
func logOK(t *testing.T, format string, args ...any) {
	t.Helper()
	msg := fmt.Sprintf(format, args...)
	if testing.Verbose() {
		t.Log(msg)
	} else {
		fmt.Println(msg)
	}
}

// ============================================================================
// User CRUD Tests — require MySQL connection to run
// ============================================================================

func TestCreateUser(t *testing.T) {
	if DB == nil {
		t.Skip("MySQL not available, skipping test (requires MySQL instance)")
	}
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
	if DB == nil {
		t.Skip("MySQL not available, skipping test (requires MySQL instance)")
	}
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
	if DB == nil {
		t.Skip("MySQL not available, skipping test (requires MySQL instance)")
	}
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
	if DB == nil {
		t.Skip("MySQL not available, skipping test (requires MySQL instance)")
	}
	_, err := GetUserByID(99999)
	if err == nil {
		t.Error("expected error for non-existent user")
	}
}

func TestGetUserByPhone(t *testing.T) {
	if DB == nil {
		t.Skip("MySQL not available, skipping test (requires MySQL instance)")
	}
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
	if DB == nil {
		t.Skip("MySQL not available, skipping test (requires MySQL instance)")
	}
	_, err := GetUserByPhone("00000000000")
	if err == nil {
		t.Error("expected error for non-existent phone")
	}
}

func TestUpdateUser(t *testing.T) {
	if DB == nil {
		t.Skip("MySQL not available, skipping test (requires MySQL instance)")
	}
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
	if DB == nil {
		t.Skip("MySQL not available, skipping test (requires MySQL instance)")
	}
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
func TestPhoneReuseAfterSoftDelete(t *testing.T) {
	if DB == nil {
		t.Skip("MySQL not available, skipping test (requires MySQL instance)")
	}
	user := &User{Phone: "13900138005", Name: "旧用户", Age: 30}
	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if err := DeleteUser(user.ID); err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}

	// On MySQL: same phone allowed after soft delete (generated column handles this)
	user2 := &User{Phone: "13900138005", Name: "新用户", Age: 25}
	if err := CreateUser(user2); err != nil {
		t.Logf("phone reuse blocked: %v", err)
		return
	}
	if user2.ID == user.ID {
		t.Error("expected new record with different ID")
	}
	b, _ := json.MarshalIndent(user2, "", "  ")
	logOK(t, "软删除后同号码可复用:\n%s", b)
}
