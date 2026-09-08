package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tn-now/server/internal/config"
	"tn-now/server/internal/cron"
)

func TestAdminDashboardHandler(t *testing.T) {
	scheduler := cron.NewScheduler(nil)
	handler := NewHandler(nil, nil, nil, scheduler, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	rr := httptest.NewRecorder()

	handler.HandleDashboard(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("Expected text/html content type, got %s", contentType)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "TN24") || !strings.Contains(body, "Control Hub") {
		t.Error("Dashboard HTML missing TN24 branding title")
	}

	// Verify accessibility landmarks
	if !strings.Contains(body, "skip-link") || !strings.Contains(body, "aria-label") {
		t.Error("Dashboard HTML missing accessibility attributes")
	}
}

func TestAdminStatsHandler(t *testing.T) {
	scheduler := cron.NewScheduler(nil)
	handler := NewHandler(nil, nil, nil, scheduler, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/api/stats", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetStats(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	var res map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to parse stats JSON: %v", err)
	}

	if res["success"] != true {
		t.Error("Expected success: true in stats response")
	}
}

func TestAdminCronHandlers(t *testing.T) {
	scheduler := cron.NewScheduler(nil)
	handler := NewHandler(nil, nil, nil, scheduler, nil)

	// Test GET /admin/api/cron
	req := httptest.NewRequest(http.MethodGet, "/admin/api/cron", nil)
	rr := httptest.NewRecorder()
	handler.HandleGetCronJobs(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	// Test POST /admin/api/cron/trigger
	triggerPayload := strings.NewReader(`{"jobId":"dead_link_checker"}`)
	reqTrigger := httptest.NewRequest(http.MethodPost, "/admin/api/cron/trigger", triggerPayload)
	rrTrigger := httptest.NewRecorder()
	handler.HandleTriggerCronJob(rrTrigger, reqTrigger)

	if rrTrigger.Code != http.StatusOK {
		t.Errorf("Expected status code 200 from trigger, got %d", rrTrigger.Code)
	}
}

func TestHandleCreateContentManualValidation(t *testing.T) {
	handler := NewHandler(nil, nil, nil, nil, nil)

	// Test GET method not allowed
	reqGet := httptest.NewRequest(http.MethodGet, "/admin/api/content/create", nil)
	rrGet := httptest.NewRecorder()
	handler.HandleCreateContentManual(rrGet, reqGet)
	if rrGet.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 for GET, got %d", rrGet.Code)
	}

	// Test missing title
	reqEmpty := httptest.NewRequest(http.MethodPost, "/admin/api/content/create", strings.NewReader(`{"title":""}`))
	rrEmpty := httptest.NewRecorder()
	handler.HandleCreateContentManual(rrEmpty, reqEmpty)
	if rrEmpty.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for empty title, got %d", rrEmpty.Code)
	}

	// Test nil DB
	reqValid := httptest.NewRequest(http.MethodPost, "/admin/api/content/create", strings.NewReader(`{"title":"Sample News","description":"Sample Description"}`))
	rrValid := httptest.NewRecorder()
	handler.HandleCreateContentManual(rrValid, reqValid)
	if rrValid.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected 503 for nil DB, got %d", rrValid.Code)
	}
}

func TestAdminAuthFlow(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:     "test-secret-12345678901234567890",
		AdminUsername: "admin",
		AdminPassword: "mypassword123",
	}
	handler := NewHandler(nil, nil, nil, nil, cfg)

	// 1. GET /admin/login renders login page
	reqLoginGet := httptest.NewRequest(http.MethodGet, "/admin/login", nil)
	rrLoginGet := httptest.NewRecorder()
	handler.HandleLogin(rrLoginGet, reqLoginGet)
	if rrLoginGet.Code != http.StatusOK {
		t.Errorf("Expected 200 for GET /admin/login, got %d", rrLoginGet.Code)
	}

	// 2. POST /admin/login with invalid password fails
	reqInvalid := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(`{"username":"admin","password":"wrongpassword"}`))
	reqInvalid.Header.Set("Content-Type", "application/json")
	rrInvalid := httptest.NewRecorder()
	handler.HandleLogin(rrInvalid, reqInvalid)
	if rrInvalid.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for wrong credentials, got %d", rrInvalid.Code)
	}

	// 3. POST /admin/login with correct credentials succeeds and sets cookie
	reqValid := httptest.NewRequest(http.MethodPost, "/admin/login", strings.NewReader(`{"username":"admin","password":"mypassword123"}`))
	reqValid.Header.Set("Content-Type", "application/json")
	rrValid := httptest.NewRecorder()
	handler.HandleLogin(rrValid, reqValid)
	if rrValid.Code != http.StatusOK {
		t.Errorf("Expected 200 for valid login, got %d", rrValid.Code)
	}
	cookies := rrValid.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == AdminSessionCookieName {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil || sessionCookie.Value == "" {
		t.Fatal("Expected admin session cookie to be set")
	}

	// 4. Validate token
	username, ok := ValidateAdminSession(sessionCookie.Value, cfg.JWTSecret)
	if !ok || username != "admin" {
		t.Errorf("Expected valid token for 'admin', got user='%s', ok=%v", username, ok)
	}

	// 5. RequireAuth middleware allows access with valid cookie
	protectedCalled := false
	protectedHandler := handler.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		protectedCalled = true
		w.WriteHeader(http.StatusOK)
	})

	reqProtected := httptest.NewRequest(http.MethodGet, "/admin/api/stats", nil)
	reqProtected.AddCookie(sessionCookie)
	rrProtected := httptest.NewRecorder()
	protectedHandler(rrProtected, reqProtected)
	if !protectedCalled || rrProtected.Code != http.StatusOK {
		t.Errorf("Expected protected handler to be called with 200, got code %d", rrProtected.Code)
	}

	// 6. RequireAuth rejects request without cookie
	reqUnauth := httptest.NewRequest(http.MethodGet, "/admin/api/stats", nil)
	rrUnauth := httptest.NewRecorder()
	protectedHandler(rrUnauth, reqUnauth)
	if rrUnauth.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for unauthenticated API request, got %d", rrUnauth.Code)
	}
}

