package model

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// jsonMarshal formats v for human-readable test output.
func jsonMarshal(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// TestMain initializes an in-memory SQLite instance with the user-owned schema
// script (principle V: TestMain runs SQL, never AutoMigrate) for model tests.
func TestMain(m *testing.M) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to open :memory: sqlite: " + err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get sql.DB: " + err.Error())
	}
	// Single connection keeps the in-memory database alive across queries.
	sqlDB.SetMaxOpenConns(1)

	schema, err := os.ReadFile("../migrations/002_test_user.sql")
	if err != nil {
		panic("failed to read migration sql: " + err.Error())
	}
	if err := db.Exec(string(schema)).Error; err != nil {
		panic("failed to exec migration sql: " + err.Error())
	}
	SQLiteDB = db

	code := m.Run()
	os.Exit(code)
}

func TestCreateTestUser(t *testing.T) {
	tu := &TestUser{Name: "张三", Phone: "13900138000", Age: 28}
	if err := CreateTestUser(tu); err != nil {
		t.Fatalf("CreateTestUser failed: %v", err)
	}
	if tu.ID == 0 {
		t.Error("expected ID to be set after create")
	}
	b, _ := jsonMarshal(tu)
	logOK(t, "created:\n%s", b)
}

func TestGetTestUserByID(t *testing.T) {
	tu := &TestUser{Name: "李四", Phone: "13900138001"}
	if err := CreateTestUser(tu); err != nil {
		t.Fatalf("CreateTestUser failed: %v", err)
	}
	got, err := GetTestUserByID(tu.ID)
	if err != nil {
		t.Fatalf("GetTestUserByID failed: %v", err)
	}
	if got.Name != "李四" || got.Phone != "13900138001" {
		t.Errorf("expected 李四/13900138001, got %s/%s", got.Name, got.Phone)
	}
	logOK(t, "get by id %d ok", tu.ID)
}

func TestGetTestUserByID_NotFound(t *testing.T) {
	_, err := GetTestUserByID(999999)
	if err == nil {
		t.Fatal("expected error for missing id, got nil")
	}
	logOK(t, "not found returns error as expected: %v", err)
}

func TestUpdateTestUser(t *testing.T) {
	tu := &TestUser{Name: "王五", Phone: "13900138002", Age: 20}
	if err := CreateTestUser(tu); err != nil {
		t.Fatalf("CreateTestUser failed: %v", err)
	}
	tu.Name = "王五丰"
	tu.Age = 22
	if err := UpdateTestUser(tu); err != nil {
		t.Fatalf("UpdateTestUser failed: %v", err)
	}
	got, err := GetTestUserByID(tu.ID)
	if err != nil {
		t.Fatalf("GetTestUserByID failed: %v", err)
	}
	if got.Name != "王五丰" || got.Age != 22 {
		t.Errorf("expected 王五丰/22, got %s/%d", got.Name, got.Age)
	}
	logOK(t, "updated id %d ok", tu.ID)
}

func TestDeleteTestUser_SoftDelete(t *testing.T) {
	tu := &TestUser{Name: "赵六", Phone: "13900138003"}
	if err := CreateTestUser(tu); err != nil {
		t.Fatalf("CreateTestUser failed: %v", err)
	}
	if err := DeleteTestUser(tu.ID); err != nil {
		t.Fatalf("DeleteTestUser failed: %v", err)
	}
	_, err := GetTestUserByID(tu.ID)
	if err == nil {
		t.Error("expected not found after soft delete")
	}
	logOK(t, "soft delete id %d ok", tu.ID)
}

func TestListTestUsers(t *testing.T) {
	// clean slate
	if err := SQLiteDB.Exec("DELETE FROM test_user").Error; err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
	var ids []int
	for i, phone := range []string{"13900138010", "13900138011", "13900138012"} {
		tu := &TestUser{Name: "U", Phone: phone, Age: i}
		if err := CreateTestUser(tu); err != nil {
			t.Fatalf("CreateTestUser failed: %v", err)
		}
		ids = append(ids, tu.ID)
	}
	// soft-delete one, it must not appear in list
	if err := DeleteTestUser(ids[0]); err != nil {
		t.Fatalf("DeleteTestUser failed: %v", err)
	}
	ts, err := ListTestUsers()
	if err != nil {
		t.Fatalf("ListTestUsers failed: %v", err)
	}
	if len(ts) != 2 {
		t.Errorf("expected 2 active users, got %d", len(ts))
	}
	b, _ := jsonMarshal(ts)
	logOK(t, "list:\n%s", b)
}

func TestPhoneActiveExists(t *testing.T) {
	if err := SQLiteDB.Exec("DELETE FROM test_user").Error; err != nil {
		t.Fatalf("cleanup failed: %v", err)
	}
	exists, err := PhoneActiveExists("13900138020")
	if err != nil {
		t.Fatalf("PhoneActiveExists failed: %v", err)
	}
	if exists {
		t.Error("expected false before creation")
	}
	tu := &TestUser{Name: "A", Phone: "13900138020"}
	if err := CreateTestUser(tu); err != nil {
		t.Fatalf("CreateTestUser failed: %v", err)
	}
	exists, err = PhoneActiveExists("13900138020")
	if err != nil {
		t.Fatalf("PhoneActiveExists failed: %v", err)
	}
	if !exists {
		t.Error("expected true after creation")
	}
	// soft-delete then the active check must allow reuse
	if err := DeleteTestUser(tu.ID); err != nil {
		t.Fatalf("DeleteTestUser failed: %v", err)
	}
	exists, err = PhoneActiveExists("13900138020")
	if err != nil {
		t.Fatalf("PhoneActiveExists failed: %v", err)
	}
	if exists {
		t.Error("expected false after soft delete (reuse allowed)")
	}
	logOK(t, "phone active uniqueness with soft-delete reuse ok")
}
