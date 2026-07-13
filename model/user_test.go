package model

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	if err := InitDB(":memory:"); err != nil {
		panic("failed to init test db: " + err.Error())
	}
	if err := DB.AutoMigrate(&User{}); err != nil {
		panic("failed to migrate test db: " + err.Error())
	}
	code := m.Run()
	os.Exit(code)
}

// logOK always prints log, even without -v flag
func logOK(t *testing.T, format string, args ...interface{}) {
	t.Helper()
	msg := fmt.Sprintf(format, args...)
	if testing.Verbose() {
		t.Log(msg)
	} else {
		fmt.Fprintln(os.Stderr, msg)
	}
}

func TestCreateUser(t *testing.T) {
	user := &User{
		Phone: "13800138000",
		Name:  "张三",
		Age:   28,
	}
	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if user.ID == 0 {
		t.Error("expected ID to be set after create")
	}
	b, _ := json.MarshalIndent(user, "", "  ")
	logOK(t, "创建成功:\n%s", b)
}

func TestCreateUser_DuplicatePhone(t *testing.T) {
	user := &User{Phone: "13800001111", Name: "dup1", Age: 20}
	if err := CreateUser(user); err != nil {
		t.Fatalf("first create should succeed: %v", err)
	}
	user2 := &User{Phone: "13800001111", Name: "dup2", Age: 30}
	if err := CreateUser(user2); err == nil {
		t.Error("expected error for duplicate phone")
	}
	b, _ := json.MarshalIndent(user, "", "  ")
	logOK(t, "首次创建成功，重复手机号正确拒绝:\n%s", b)
}

func TestGetUserByID(t *testing.T) {
	user := &User{Phone: "13800138001", Name: "李四", Age: 30}
	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	found, err := GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
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

func TestGetUserByID_NotFound(t *testing.T) {
	_, err := GetUserByID(99999)
	if err == nil {
		t.Error("expected error for non-existent user")
	}
}

func TestGetUserByPhone(t *testing.T) {
	user := &User{Phone: "13800138002", Name: "王五", Age: 25}
	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	found, err := GetUserByPhone("13800138002")
	if err != nil {
		t.Fatalf("GetUserByPhone failed: %v", err)
	}
	if found.Name != "王五" {
		t.Errorf("expected name 王五, got %s", found.Name)
	}
	b, _ := json.MarshalIndent(found, "", "  ")
	logOK(t, "按手机号查询成功:\n%s", b)
}

func TestGetUserByPhone_NotFound(t *testing.T) {
	_, err := GetUserByPhone("00000000000")
	if err == nil {
		t.Error("expected error for non-existent phone")
	}
}

func TestUpdateUser(t *testing.T) {
	user := &User{Phone: "13800138003", Name: "赵六", Age: 22}
	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user.Name = "赵六改名"
	user.Age = 35
	if err := UpdateUser(user); err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	found, err := GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("GetUserByID after update failed: %v", err)
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

func TestDeleteUser(t *testing.T) {
	user := &User{Phone: "13800138004", Name: "钱七", Age: 40}
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
	logOK(t, "软删除成功，已查不到该用户:\n%s", b)
}

func TestRestoreUserByPhone(t *testing.T) {
	user := &User{Phone: "13800138005", Name: "孙八", Age: 33}
	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	if err := DeleteUser(user.ID); err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}

	restored, err := RestoreUserByPhone("13800138005")
	if err != nil {
		t.Fatalf("RestoreUserByPhone failed: %v", err)
	}
	if restored.Name != "孙八" {
		t.Errorf("expected name 孙八, got %s", restored.Name)
	}

	// verify it can be found normally after restore
	found, err := GetUserByID(restored.ID)
	if err != nil {
		t.Errorf("GetUserByID after restore should succeed: %v", err)
	}
	if found == nil {
		t.Error("expected to find user after restore")
	}
	b, _ := json.MarshalIndent(found, "", "  ")
	logOK(t, "恢复成功:\n%s", b)
}

func TestRestoreUserByPhone_NotFound(t *testing.T) {
	_, err := RestoreUserByPhone("99999999999")
	if err == nil {
		t.Error("expected error for non-existent phone")
	}
}
