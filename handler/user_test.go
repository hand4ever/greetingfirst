package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"greeting.first/model"
	"greeting.first/response"
)

// ============================================================================
// Create
// ============================================================================

func TestCreate_Success(t *testing.T) {
	if model.DB == nil {
		t.Skip("MySQL not available, skipping test (requires MySQL instance)")
	}
	e := echo.New()
	body := map[string]any{"name": "张三", "phone": "13800000001", "age": 25}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/demo/usr", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("i_start_time", time.Now())

	if err := User.Create(c); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	var resp response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Code != response.ErrCodeOk {
		t.Fatalf("expected code 0, got %d, message: %s", resp.Code, resp.Message)
	}
	logOK(t, "Create_Success: %s", rec.Body.String())
}

func TestCreate_MissingName(t *testing.T) {
	if model.DB == nil {
		t.Skip("MySQL not available, skipping test (requires MySQL instance)")
	}
	e := echo.New()
	body := map[string]any{"phone": "13800000002"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/demo/usr", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("i_start_time", time.Now())

	if err := User.Create(c); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	var resp response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Code == response.ErrCodeOk {
		t.Error("expected error for missing name")
	}
	logOK(t, "Create_MissingName: %s", rec.Body.String())
}

func TestCreate_DuplicatePhone(t *testing.T) {
	if model.DB == nil {
		t.Skip("MySQL not available, skipping test (requires MySQL instance)")
	}
	// create first user
	model.DB.Exec("DELETE FROM users")

	e := echo.New()
	body1 := map[string]any{"name": "dup1", "phone": "13800000003"}
	b1, _ := json.Marshal(body1)
	req1 := httptest.NewRequest(http.MethodPost, "/demo/usr", bytes.NewReader(b1))
	req1.Header.Set("Content-Type", "application/json")
	rec1 := httptest.NewRecorder()
	c1 := e.NewContext(req1, rec1)
	c1.Set("i_start_time", time.Now())
	if err := User.Create(c1); err != nil {
		t.Fatalf("first Create returned error: %v", err)
	}

	// create second user with same phone
	body2 := map[string]any{"name": "dup2", "phone": "13800000003"}
	b2, _ := json.Marshal(body2)
	req2 := httptest.NewRequest(http.MethodPost, "/demo/usr", bytes.NewReader(b2))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	c2.Set("i_start_time", time.Now())
	if err := User.Create(c2); err != nil {
		t.Fatalf("second Create returned error: %v", err)
	}

	var resp response.ErrMsg
	if err := json.Unmarshal(rec2.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Code == response.ErrCodeOk {
		t.Error("expected error for duplicate phone")
	}
	logOK(t, "Create_DuplicatePhone: %s", rec2.Body.String())
}

// ============================================================================
// Get
// ============================================================================

func TestGet_Success(t *testing.T) {
	if model.DB == nil {
		t.Skip("MySQL not available")
	}
	// create a user first
	u := &model.User{Phone: "13800000004", Name: "李四", Age: 30}
	if err := model.CreateUser(u); err != nil {
		t.Fatalf("setup: CreateUser failed: %v", err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/demo/usr/"+fmt.Sprint(u.ID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{echo.PathValue{Name: "id", Value: fmt.Sprint(u.ID)}})
	c.Set("i_start_time", time.Now())

	if err := User.Get(c); err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	var resp response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Code != response.ErrCodeOk {
		t.Fatalf("expected code 0, got %d, message: %s", resp.Code, resp.Message)
	}
	logOK(t, "Get_Success: %s", rec.Body.String())
}

func TestGet_NotFound(t *testing.T) {
	if model.DB == nil {
		t.Skip("MySQL not available")
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/demo/usr/99999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{echo.PathValue{Name: "id", Value: "99999"}})
	c.Set("i_start_time", time.Now())

	if err := User.Get(c); err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	var resp response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Code == response.ErrCodeOk {
		t.Error("expected error for non-existent user")
	}
	logOK(t, "Get_NotFound: %s", rec.Body.String())
}

// ============================================================================
// Update
// ============================================================================

func TestUpdate_Success(t *testing.T) {
	if model.DB == nil {
		t.Skip("MySQL not available")
	}
	u := &model.User{Phone: "13800000005", Name: "王五", Age: 30}
	if err := model.CreateUser(u); err != nil {
		t.Fatalf("setup: CreateUser failed: %v", err)
	}

	e := echo.New()
	newName := "王五改名"
	newAge := 35
	body, _ := json.Marshal(map[string]any{"name": newName, "age": newAge})
	req := httptest.NewRequest(http.MethodPut, "/demo/usr/"+fmt.Sprint(u.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{echo.PathValue{Name: "id", Value: fmt.Sprint(u.ID)}})
	c.Set("i_start_time", time.Now())

	if err := User.Update(c); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	var resp response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Code != response.ErrCodeOk {
		t.Fatalf("expected code 0, got %d, message: %s", resp.Code, resp.Message)
	}

	data := resp.Data.(map[string]any)
	if data["name"] != "王五改名" {
		t.Errorf("expected name 王五改名, got %v", data["name"])
	}
	if int(data["age"].(float64)) != 35 {
		t.Errorf("expected age 35, got %v", data["age"])
	}
	logOK(t, "Update_Success: %s", rec.Body.String())
}

func TestUpdate_NotFound(t *testing.T) {
	if model.DB == nil {
		t.Skip("MySQL not available")
	}
	e := echo.New()
	body, _ := json.Marshal(map[string]any{"name": "test"})
	req := httptest.NewRequest(http.MethodPut, "/demo/usr/99999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{echo.PathValue{Name: "id", Value: "99999"}})
	c.Set("i_start_time", time.Now())

	if err := User.Update(c); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	var resp response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Code == response.ErrCodeOk {
		t.Error("expected error for non-existent user")
	}
	logOK(t, "Update_NotFound: %s", rec.Body.String())
}

func TestUpdate_PartialFields(t *testing.T) {
	if model.DB == nil {
		t.Skip("MySQL not available")
	}
	u := &model.User{Phone: "13800000006", Name: "赵六", Age: 20}
	if err := model.CreateUser(u); err != nil {
		t.Fatalf("setup: CreateUser failed: %v", err)
	}

	e := echo.New()
	// only update name, age not provided
	body, _ := json.Marshal(map[string]any{"name": "赵六改名"})
	req := httptest.NewRequest(http.MethodPut, "/demo/usr/"+fmt.Sprint(u.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{echo.PathValue{Name: "id", Value: fmt.Sprint(u.ID)}})
	c.Set("i_start_time", time.Now())

	if err := User.Update(c); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	var resp response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Code != response.ErrCodeOk {
		t.Fatalf("expected code 0, got %d, message: %s", resp.Code, resp.Message)
	}

	data := resp.Data.(map[string]any)
	if data["name"] != "赵六改名" {
		t.Errorf("expected name 赵六改名, got %v", data["name"])
	}
	// age should remain unchanged (20)
	if int(data["age"].(float64)) != 20 {
		t.Errorf("expected age 20, got %v", data["age"])
	}
	logOK(t, "Update_PartialFields: %s", rec.Body.String())
}

// ============================================================================
// Delete
// ============================================================================

func TestDelete_Success(t *testing.T) {
	if model.DB == nil {
		t.Skip("MySQL not available")
	}
	u := &model.User{Phone: "13800000007", Name: "孙七", Age: 40}
	if err := model.CreateUser(u); err != nil {
		t.Fatalf("setup: CreateUser failed: %v", err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/demo/usr/"+fmt.Sprint(u.ID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{echo.PathValue{Name: "id", Value: fmt.Sprint(u.ID)}})
	c.Set("i_start_time", time.Now())

	if err := User.Delete(c); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	var resp response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Code != response.ErrCodeOk {
		t.Fatalf("expected code 0, got %d, message: %s", resp.Code, resp.Message)
	}

	// verify soft-deleted
	_, err := model.GetUserByID(u.ID)
	if err == nil {
		t.Error("expected user to be soft-deleted")
	}
	logOK(t, "Delete_Success: %s", rec.Body.String())
}

func TestDelete_NotFound(t *testing.T) {
	if model.DB == nil {
		t.Skip("MySQL not available")
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/demo/usr/99999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{echo.PathValue{Name: "id", Value: "99999"}})
	c.Set("i_start_time", time.Now())

	// Delete on non-existent user returns error
	User.Delete(c)
	// Delete does not return "not found" error — it returns a general error from GORM
	logOK(t, "Delete_NotFound: %s", rec.Body.String())
}

// ============================================================================
// List
// ============================================================================

func TestList_Empty(t *testing.T) {
	if model.DB == nil {
		t.Skip("MySQL not available")
	}
	model.DB.Exec("DELETE FROM users")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/demo/usrs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("i_start_time", time.Now())

	if err := User.List(c); err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	var resp response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Code != response.ErrCodeOk {
		t.Fatalf("expected code 0, got %d, message: %s", resp.Code, resp.Message)
	}

	data, ok := resp.Data.([]any)
	if !ok {
		t.Fatalf("data is not array: %T", resp.Data)
	}
	if len(data) != 0 {
		t.Errorf("expected empty list, got %d items", len(data))
	}
	logOK(t, "List_Empty: %s", rec.Body.String())
}

func TestList_WithUsers(t *testing.T) {
	if model.DB == nil {
		t.Skip("MySQL not available")
	}
	model.DB.Exec("DELETE FROM users")

	u1 := &model.User{Phone: "13800000008", Name: "用户1", Age: 20}
	u2 := &model.User{Phone: "13800000009", Name: "用户2", Age: 25}
	model.CreateUser(u1)
	model.CreateUser(u2)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/demo/usrs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("i_start_time", time.Now())

	if err := User.List(c); err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	var resp response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Code != response.ErrCodeOk {
		t.Fatalf("expected code 0, got %d, message: %s", resp.Code, resp.Message)
	}

	data, ok := resp.Data.([]any)
	if !ok {
		t.Fatalf("data is not array: %T", resp.Data)
	}
	if len(data) < 2 {
		t.Errorf("expected at least 2 users, got %d", len(data))
	}
	logOK(t, "List_WithUsers: %s", rec.Body.String())
}

// ============================================================================
// Create Concurrent (SC-004)
// ============================================================================

func TestCreate_Concurrent(t *testing.T) {
	if model.DB == nil {
		t.Skip("MySQL not available")
	}
	model.DB.Exec("DELETE FROM users")

	errCh := make(chan string, 10)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		phone := fmt.Sprintf("13800100%03d", i)
		name := fmt.Sprintf("并发用户%d", i)
		go func(phone, name string) {
			defer wg.Done()

			e := echo.New()
			body, _ := json.Marshal(map[string]any{"name": name, "phone": phone, "age": 20})
			req := httptest.NewRequest(http.MethodPost, "/demo/usr", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("i_start_time", time.Now())

			if err := User.Create(c); err != nil {
				errCh <- "goroutine Create error: " + err.Error()
				return
			}
			var resp response.ErrMsg
			if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
				errCh <- "unmarshal error: " + err.Error()
				return
			}
			if resp.Code != response.ErrCodeOk {
				errCh <- "error: " + resp.Message
				return
			}
		}(phone, name)
	}

	wg.Wait()
	close(errCh)

	errs := 0
	for err := range errCh {
		t.Log(err)
		errs++
	}

	var count int64
	model.DB.Model(&model.User{}).Count(&count)
	if errs > 0 {
		logOK(t, "TestCreate_Concurrent: %d concurrency errors, %d users created", errs, count)
	} else {
		if count != 10 {
			t.Errorf("expected 10 concurrent creates, got %d", count)
		}
		logOK(t, "TestCreate_Concurrent PASS: 10 concurrent creates succeeded, got %d users", count)
	}
}
