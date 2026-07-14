package middle

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// logOK prints test success output. Uses t.Log in verbose mode, fmt.Println otherwise.
func logOK(t *testing.T, format string, args ...any) {
	t.Helper()
	msg := fmt.Sprintf(format, args...)
	if testing.Verbose() {
		t.Log(msg)
	} else {
		fmt.Println(msg)
	}
}

// newCORSApp creates an Echo instance with CORS middleware and a simple test handler.
func newCORSApp(allowOrigins []string) *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch, http.MethodHead},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With", "Accept", "Origin"},
		AllowCredentials: false,
		MaxAge:           86400,
	}))
	e.GET("/test", func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.POST("/test", func(c *echo.Context) error {
		return c.String(http.StatusOK, "created")
	})
	e.OPTIONS("/test", func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	return e
}

// TestCORSPreflight verifies OPTIONS preflight returns correct CORS headers.
// Status 204 is the standard Echo CORS preflight response.
func TestCORSPreflight(t *testing.T) {
	e := newCORSApp([]string{"*"})

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK && rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 200 or 204, got %d", rec.Code)
	}

	allowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "*" {
		t.Errorf("expected Access-Control-Allow-Origin '*', got '%s'", allowOrigin)
	}

	allowMethods := rec.Header().Get("Access-Control-Allow-Methods")
	if allowMethods == "" {
		t.Error("expected non-empty Access-Control-Allow-Methods")
	}

	logOK(t, "TestCORSPreflight PASS: status=%d, Allow-Origin=%s, Allow-Methods=%s",
		rec.Code, allowOrigin, allowMethods)
}

// TestCORSGETWithOrigin verifies GET with Origin header returns Access-Control-Allow-Origin.
func TestCORSGETWithOrigin(t *testing.T) {
	e := newCORSApp([]string{"*"})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	allowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "*" {
		t.Errorf("expected Access-Control-Allow-Origin '*', got '%s'", allowOrigin)
	}

	if rec.Body.String() != "ok" {
		t.Errorf("expected body 'ok', got '%s'", rec.Body.String())
	}

	logOK(t, "TestCORSGETWithOrigin PASS: status=%d, Allow-Origin=%s, body=%s",
		rec.Code, allowOrigin, rec.Body.String())
}

// TestCORSSameOrigin verifies request without Origin header does NOT add CORS headers.
func TestCORSSameOrigin(t *testing.T) {
	e := newCORSApp([]string{"*"})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// No Origin header set – simulating same-origin request
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	allowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "" {
		t.Errorf("expected no Access-Control-Allow-Origin for same-origin request, got '%s'", allowOrigin)
	}

	allowMethods := rec.Header().Get("Access-Control-Allow-Methods")
	if allowMethods != "" {
		t.Errorf("expected no Access-Control-Allow-Methods for same-origin request, got '%s'", allowMethods)
	}

	logOK(t, "TestCORSSameOrigin PASS: status=%d, no CORS headers added", rec.Code)
}

// TestCORSOptionsWithoutRequestMethod verifies OPTIONS without Access-Control-Request-Method
// is still handled gracefully by Echo CORS middleware (returns CORS headers as a safe default).
func TestCORSOptionsWithoutRequestMethod(t *testing.T) {
	e := newCORSApp([]string{"*"})

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	// No Access-Control-Request-Method header
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Echo CORS middleware handles OPTIONS with Origin even without ACRM as a safe default.
	if rec.Code != http.StatusOK && rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 200 or 204, got %d", rec.Code)
	}

	// Echo CORS middleware still adds CORS headers as a safe default behavior.
	allowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin == "" {
		t.Error("expected Access-Control-Allow-Origin to be present")
	}

	logOK(t, "TestCORSOptionsWithoutRequestMethod PASS: status=%d, Allow-Origin=%s",
		rec.Code, allowOrigin)
}

// TestCORSSpecificOriginAllowed verifies that a matched specific origin receives CORS headers.
func TestCORSSpecificOriginAllowed(t *testing.T) {
	e := newCORSApp([]string{"https://myapp.com"})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://myapp.com")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	allowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "https://myapp.com" {
		t.Errorf("expected Access-Control-Allow-Origin 'https://myapp.com', got '%s'", allowOrigin)
	}

	logOK(t, "TestCORSSpecificOriginAllowed PASS: status=%d, Allow-Origin=%s",
		rec.Code, allowOrigin)
}

// TestCORSSpecificOriginDisallowed verifies that an unmatched origin does NOT receive CORS headers.
func TestCORSSpecificOriginDisallowed(t *testing.T) {
	e := newCORSApp([]string{"https://myapp.com"})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://evil.com")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	allowOrigin := rec.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin != "" {
		t.Errorf("expected no Access-Control-Allow-Origin for disallowed origin, got '%s'", allowOrigin)
	}

	logOK(t, "TestCORSSpecificOriginDisallowed PASS: status=%d, no CORS headers (origin not allowed)",
		rec.Code)
}
