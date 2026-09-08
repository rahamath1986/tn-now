package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestIDMiddleware(t *testing.T) {
	// Simple handler that reads request ID from context
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := GetRequestID(r.Context())
		if reqID == "" {
			t.Error("RequestID was not injected in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	respHeader := rr.Header().Get("X-Request-ID")
	if respHeader == "" {
		t.Error("X-Request-ID header was not written to response")
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	// Handler that explicitly triggers a panic
	handler := Recovery(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("intentional crash test")
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 500, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Internal Server Error") {
		t.Errorf("Expected response body to be formatted JSON error, got %s", body)
	}
}

func TestCORSMiddleware(t *testing.T) {
	allowedOrigins := []string{"http://example.com"}
	handler := CORS(allowedOrigins)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Match origin
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Header().Get("Access-Control-Allow-Origin") != "http://example.com" {
		t.Errorf("Expected Access-Control-Allow-Origin to match, got %s", rr.Header().Get("Access-Control-Allow-Origin"))
	}

	// Mismatched origin
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://malicious.com")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("Expected Access-Control-Allow-Origin to be empty, got %s", rr.Header().Get("Access-Control-Allow-Origin"))
	}
}
