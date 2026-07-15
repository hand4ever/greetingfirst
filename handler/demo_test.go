package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"greeting.first/response"
)

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

func TestSearch_MultipleTags(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/demo/search?tag=go&tag=web&tag=api", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("i_start_time", time.Now())

	if err := Demo.Search(c); err != nil {
		t.Fatalf("Search returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var body response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if body.Code != 0 {
		t.Errorf("expected code 0, got %d", body.Code)
	}

	data, ok := body.Data.(map[string]any)
	if !ok {
		t.Fatalf("data is not map: %T", body.Data)
	}
	tags, ok := data["tag"].([]any)
	if !ok {
		t.Fatalf("tag is not array: %T", data["tag"])
	}
	if len(tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(tags))
	}
	logOK(t, "响应内容: %s", rec.Body.String())
}

func TestSearch_NoTags(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/demo/search", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("i_start_time", time.Now())

	if err := Demo.Search(c); err != nil {
		t.Fatalf("Search returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	logOK(t, "响应内容: %s", rec.Body.String())
}

func TestErrDebug(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/demo/err/debug/hello", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{echo.PathValue{Name: "str", Value: "hello"}})
	c.Set("i_start_time", time.Now())

	if err := Demo.ErrDebug(c); err != nil {
		t.Fatalf("ErrDebug returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var body response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if body.Code != 0 {
		t.Errorf("expected code 0, got %d", body.Code)
	}

	data, ok := body.Data.(map[string]any)
	if !ok {
		t.Fatalf("data is not map: %T", body.Data)
	}
	if str, ok := data["str"].(string); !ok || str != "hello" {
		t.Errorf("expected str=hello, got %v", data["str"])
	}
	logOK(t, "响应内容: %s", rec.Body.String())
}

func TestErrDebug_ChineseStr(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/demo/err/debug/你好世界", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{echo.PathValue{Name: "str", Value: "你好世界"}})
	c.Set("i_start_time", time.Now())

	if err := Demo.ErrDebug(c); err != nil {
		t.Fatalf("ErrDebug with Chinese returned error: %v", err)
	}

	var body response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	data, ok := body.Data.(map[string]any)
	if !ok {
		t.Fatalf("data is not map: %T", body.Data)
	}
	if str, ok := data["str"].(string); !ok || str != "你好世界" {
		t.Errorf("expected str=你好世界, got %v", data["str"])
	}
	logOK(t, "响应内容: %s", rec.Body.String())
}
