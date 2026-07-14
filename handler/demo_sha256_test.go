package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"greeting.first/response"
)

func TestSha256Compute_Basic(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/demo/sha256?text=hello", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := Sha256.Compute(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var resp response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != response.ErrCodeOk {
		t.Fatalf("expected code 0, got %d", resp.Code)
	}

	data, ok := resp.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected data to be map, got %T", resp.Data)
	}

	if data["input"] != "hello" {
		t.Errorf("expected input 'hello', got '%v'", data["input"])
	}

	expectedHash := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if data["hash"] != expectedHash {
		t.Errorf("expected hash '%s', got '%v'", expectedHash, data["hash"])
	}

	logOK(t, "TestSha256Compute_Basic 响应内容: %s", rec.Body.String())
}

func TestSha256Compute_Chinese(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/demo/sha256?text=%E4%BD%A0%E5%A5%BD", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := Sha256.Compute(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var resp response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != response.ErrCodeOk {
		t.Fatalf("expected code 0, got %d", resp.Code)
	}

	data, ok := resp.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected data to be map, got %T", resp.Data)
	}

	if data["input"] != "你好" {
		t.Errorf("expected input '你好', got '%v'", data["input"])
	}

	// Cross-verify: echo -n "你好" | sha256sum
	expectedHash := "670d9743542cae3ea7ebe36af56bd53648b0a1126162e78d81a32934a711302e"
	if data["hash"] != expectedHash {
		t.Errorf("expected hash '%s', got '%v'", expectedHash, data["hash"])
	}

	logOK(t, "TestSha256Compute_Chinese 响应内容: %s", rec.Body.String())
}

func TestSha256Compute_EmptyString(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/demo/sha256?text=", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := Sha256.Compute(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var resp response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != response.ErrCodeOk {
		t.Fatalf("expected code 0, got %d", resp.Code)
	}

	data, ok := resp.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected data to be map, got %T", resp.Data)
	}

	expectedHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	if data["hash"] != expectedHash {
		t.Errorf("expected hash '%s', got '%v'", expectedHash, data["hash"])
	}

	if data["input"] != "" {
		t.Errorf("expected empty input, got '%v'", data["input"])
	}

	logOK(t, "TestSha256Compute_EmptyString 响应内容: %s", rec.Body.String())
}

func TestSha256Compute_MissingParam(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/demo/sha256", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := Sha256.Compute(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var resp response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code == response.ErrCodeOk {
		t.Fatal("expected non-zero error code for missing parameter")
	}

	if resp.Message == "" {
		t.Error("expected non-empty error message")
	}

	logOK(t, "TestSha256Compute_MissingParam 响应内容: %s", rec.Body.String())
}

func TestSha256Compute_SpecialChars(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/demo/sha256?text=%F0%9F%98%80", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := Sha256.Compute(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var resp response.ErrMsg
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Code != response.ErrCodeOk {
		t.Fatalf("expected code 0, got %d", resp.Code)
	}

	data, ok := resp.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected data to be map, got %T", resp.Data)
	}

	// cross-verify: echo -n "😀" | sha256sum
	expectedHash := "f0443a342c5ef54783a111b51ba56c938e474c32324d90c3a60c9c8e3a37e2d9"
	if data["hash"] != expectedHash {
		t.Errorf("expected hash '%s', got '%v'", expectedHash, data["hash"])
	}

	logOK(t, "TestSha256Compute_SpecialChars 响应内容: %s", rec.Body.String())
}
