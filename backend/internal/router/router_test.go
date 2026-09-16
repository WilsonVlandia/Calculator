package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"calculator/internal/calculator"
	"calculator/internal/handler"
)

func newTestRouter() http.Handler {
	h := handler.New(calculator.New(2))
	return New(h, []string{"http://allowed.example"})
}

type errorResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Code string `json:"code"`
	} `json:"error"`
}

func TestRouterDispatchesKnownRoutes(t *testing.T) {
	mux := newTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/add", bytes.NewBufferString(`{"a":2,"b":3}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body: %s)", rec.Code, rec.Body.String())
	}
}

func TestRouterRejectsWrongMethod(t *testing.T) {
	mux := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/add", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}

	var body errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body.Error.Code != "METHOD_NOT_ALLOWED" {
		t.Errorf("unexpected error code: %s", body.Error.Code)
	}
}

func TestRouterReturnsNotFoundForUnknownRoute(t *testing.T) {
	mux := newTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/unknown", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}

	var body errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body.Error.Code != "NOT_FOUND" {
		t.Errorf("unexpected error code: %s", body.Error.Code)
	}
}

func TestCORSAllowsConfiguredOrigin(t *testing.T) {
	mux := newTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/add", bytes.NewBufferString(`{"a":2,"b":3}`))
	req.Header.Set("Origin", "http://allowed.example")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://allowed.example" {
		t.Errorf("expected allowed origin to be echoed back, got %q", got)
	}
}

func TestCORSRejectsUnknownOrigin(t *testing.T) {
	mux := newTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/add", bytes.NewBufferString(`{"a":2,"b":3}`))
	req.Header.Set("Origin", "http://evil.example")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("expected no CORS header for a disallowed origin, got %q", got)
	}
}

func TestCORSPreflightIsAnsweredDirectly(t *testing.T) {
	mux := newTestRouter()

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/add", nil)
	req.Header.Set("Origin", "http://allowed.example")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Errorf("expected Access-Control-Allow-Methods header to be set")
	}
}

func TestCORSAllowAllWildcard(t *testing.T) {
	h := handler.New(calculator.New(2))
	mux := New(h, []string{"*"})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/add", bytes.NewBufferString(`{"a":2,"b":3}`))
	req.Header.Set("Origin", "http://anything.example")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("expected wildcard origin, got %q", got)
	}
}
