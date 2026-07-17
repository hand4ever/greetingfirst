package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"greeting.first/model"
	"greeting.first/response"
)

// TestMain initializes an in-memory SQLite instance with the user-owned schema
// script so the SQLite handler tests run without an external database.
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
	sqlDB.SetMaxOpenConns(1)

	schema, err := os.ReadFile("../migrations/002_test_user.sql")
	if err != nil {
		panic("failed to read migration sql: " + err.Error())
	}
	if err := db.Exec(string(schema)).Error; err != nil {
		panic("failed to exec migration sql: " + err.Error())
	}
	model.SQLiteDB = db

	code := m.Run()
	os.Exit(code)
}

// newSqliteCtx builds an echo context for a direct (non-routed) handler call.
// Suitable for endpoints without path parameters (e.g. Create).
func newSqliteCtx(method, path string, body any) (*echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	var reader *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reader = bytes.NewReader(b)
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("i_start_time", time.Now())
	return c, rec
}

// serveSqlite registers the SQLite routes and dispatches the request through
// Echo so path parameters (e.g. :id) are parsed for the handler under test.
func serveSqlite(method, path string, body any) *httptest.ResponseRecorder {
	e := echo.New()
	g := e.Group("/sqlite/testuser")
	g.POST("", SqliteUser.Create)
	g.GET("/:id", SqliteUser.Get)
	g.PUT("/:id", SqliteUser.Update)
	g.DELETE("/:id", SqliteUser.Delete)
	e.GET("/sqlite/testusers", SqliteUser.List)

	var reader *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reader = bytes.NewReader(b)
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

// sqliteRespCode parses the unified response from a recorder.
func sqliteRespCode(t *testing.T, rec *httptest.ResponseRecorder) response.ErrMsg {
	t.Helper()
	var resp response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	return resp
}

func itoa(n int) string {
	return strconv.Itoa(n)
}

func TestSqliteUser_Create_Success(t *testing.T) {
	c, rec := newSqliteCtx(http.MethodPost, "/sqlite/testuser", map[string]any{
		"name": "张三", "phone": "13900138000", "age": 28,
	})
	if err := SqliteUser.Create(c); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	resp := sqliteRespCode(t, rec)
	if resp.Code != response.ErrCodeOk {
		t.Fatalf("expected code 0, got %d, message: %s", resp.Code, resp.Message)
	}
	logOK(t, "Create_Success: %s", rec.Body.String())
}

func TestSqliteUser_Create_MissingName(t *testing.T) {
	c, rec := newSqliteCtx(http.MethodPost, "/sqlite/testuser", map[string]any{
		"phone": "13900138001",
	})
	if err := SqliteUser.Create(c); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	resp := sqliteRespCode(t, rec)
	if resp.Code == response.ErrCodeOk {
		t.Fatalf("expected error for missing name, got code 0")
	}
	logOK(t, "Create_MissingName: %s", rec.Body.String())
}

func TestSqliteUser_Create_MissingPhone(t *testing.T) {
	c, rec := newSqliteCtx(http.MethodPost, "/sqlite/testuser", map[string]any{
		"name": "李四",
	})
	if err := SqliteUser.Create(c); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	resp := sqliteRespCode(t, rec)
	if resp.Code == response.ErrCodeOk {
		t.Fatalf("expected error for missing phone, got code 0")
	}
	logOK(t, "Create_MissingPhone: %s", rec.Body.String())
}

func TestSqliteUser_Create_DuplicatePhone(t *testing.T) {
	c1, _ := newSqliteCtx(http.MethodPost, "/sqlite/testuser", map[string]any{
		"name": "A", "phone": "13900138002",
	})
	if err := SqliteUser.Create(c1); err != nil {
		t.Fatalf("first Create returned error: %v", err)
	}
	rec2 := serveSqlite(http.MethodPost, "/sqlite/testuser", map[string]any{
		"name": "B", "phone": "13900138002",
	})
	resp := sqliteRespCode(t, rec2)
	if resp.Code == response.ErrCodeOk {
		t.Fatalf("expected duplicate-phone error, got code 0")
	}
	logOK(t, "Create_DuplicatePhone: %s", rec2.Body.String())
}

func TestSqliteUser_Get_Success(t *testing.T) {
	c, rec := newSqliteCtx(http.MethodPost, "/sqlite/testuser", map[string]any{
		"name": "王五", "phone": "13900138003", "age": 30,
	})
	if err := SqliteUser.Create(c); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	var created response.ErrMsg
	_ = json.Unmarshal(rec.Body.Bytes(), &created)
	id := int(created.Data.(map[string]any)["id"].(float64))

	rec2 := serveSqlite(http.MethodGet, "/sqlite/testuser/"+itoa(id), nil)
	resp := sqliteRespCode(t, rec2)
	if resp.Code != response.ErrCodeOk {
		t.Fatalf("expected code 0, got %d, message: %s", resp.Code, resp.Message)
	}
	logOK(t, "Get_Success: %s", rec2.Body.String())
}

func TestSqliteUser_Get_NotFound(t *testing.T) {
	rec := serveSqlite(http.MethodGet, "/sqlite/testuser/999999", nil)
	resp := sqliteRespCode(t, rec)
	if resp.Code == response.ErrCodeOk {
		t.Fatalf("expected not-found error, got code 0")
	}
	logOK(t, "Get_NotFound: %s", rec.Body.String())
}

func TestSqliteUser_Get_InvalidID(t *testing.T) {
	rec := serveSqlite(http.MethodGet, "/sqlite/testuser/abc", nil)
	resp := sqliteRespCode(t, rec)
	if resp.Code == response.ErrCodeOk {
		t.Fatalf("expected invalid-id error, got code 0")
	}
	logOK(t, "Get_InvalidID: %s", rec.Body.String())
}

func TestSqliteUser_Update_Success(t *testing.T) {
	c, rec := newSqliteCtx(http.MethodPost, "/sqlite/testuser", map[string]any{
		"name": "赵六", "phone": "13900138004", "age": 20,
	})
	if err := SqliteUser.Create(c); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	var created response.ErrMsg
	_ = json.Unmarshal(rec.Body.Bytes(), &created)
	id := int(created.Data.(map[string]any)["id"].(float64))

	rec2 := serveSqlite(http.MethodPut, "/sqlite/testuser/"+itoa(id), map[string]any{
		"name": "赵六丰", "age": 25,
	})
	resp := sqliteRespCode(t, rec2)
	if resp.Code != response.ErrCodeOk {
		t.Fatalf("expected code 0, got %d, message: %s", resp.Code, resp.Message)
	}
	data := resp.Data.(map[string]any)
	if data["name"] != "赵六丰" || int(data["age"].(float64)) != 25 {
		t.Errorf("expected 赵六丰/25, got %v/%v", data["name"], data["age"])
	}
	logOK(t, "Update_Success: %s", rec2.Body.String())
}

func TestSqliteUser_Delete_Success(t *testing.T) {
	c, rec := newSqliteCtx(http.MethodPost, "/sqlite/testuser", map[string]any{
		"name": "钱七", "phone": "13900138005",
	})
	if err := SqliteUser.Create(c); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	var created response.ErrMsg
	_ = json.Unmarshal(rec.Body.Bytes(), &created)
	id := int(created.Data.(map[string]any)["id"].(float64))

	recDel := serveSqlite(http.MethodDelete, "/sqlite/testuser/"+itoa(id), nil)
	respDel := sqliteRespCode(t, recDel)
	if respDel.Code != response.ErrCodeOk {
		t.Fatalf("expected code 0, got %d, message: %s", respDel.Code, respDel.Message)
	}

	// subsequent get must return not found (soft delete)
	recGet := serveSqlite(http.MethodGet, "/sqlite/testuser/"+itoa(id), nil)
	respGet := sqliteRespCode(t, recGet)
	if respGet.Code == response.ErrCodeOk {
		t.Fatalf("expected not-found after delete, got code 0")
	}
	logOK(t, "Delete_Success: delete ok, subsequent get -> %s", recGet.Body.String())
}

func TestSqliteUser_List_AfterCreate(t *testing.T) {
	c, _ := newSqliteCtx(http.MethodPost, "/sqlite/testuser", map[string]any{
		"name": "孙八", "phone": "13900138006",
	})
	if err := SqliteUser.Create(c); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	rec := serveSqlite(http.MethodGet, "/sqlite/testusers", nil)
	resp := sqliteRespCode(t, rec)
	if resp.Code != response.ErrCodeOk {
		t.Fatalf("expected code 0, got %d, message: %s", resp.Code, resp.Message)
	}
	logOK(t, "List_AfterCreate: %s", rec.Body.String())
}
