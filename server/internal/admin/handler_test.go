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

func TestAdminSMTPSettingsHandlers(t *testing.T) {
	cfg := &config.Config{
		EmailProvider:   "resend",
		ResendAPIKey:    "re_test_123456",
		ResendFromEmail: "TN24 News <onboarding@resend.dev>",
		SMTPHost:        "smtp.test.example",
		SMTPPort:        "587",
		SMTPUser:        "testuser",
		SMTPPass:        "testpass",
		SMTPFrom:        "noreply@test.example",
	}
	handler := NewHandler(nil, nil, nil, nil, cfg)

	// 1. GET /admin/api/settings/smtp
	reqGet := httptest.NewRequest(http.MethodGet, "/admin/api/settings/smtp", nil)
	rrGet := httptest.NewRecorder()
	handler.HandleSMTPSettings(rrGet, reqGet)
	if rrGet.Code != http.StatusOK {
		t.Errorf("Expected 200 for GET smtp settings, got %d", rrGet.Code)
	}

	var res map[string]interface{}
	if err := json.NewDecoder(rrGet.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to parse json: %v", err)
	}
	if res["success"] != true {
		t.Error("Expected success: true in smtp response")
	}
	data, ok := res["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected data map in response")
	}
	if data["provider"] != "resend" {
		t.Errorf("Expected provider 'resend', got %v", data["provider"])
	}
	if data["hasResendKey"] != true {
		t.Errorf("Expected hasResendKey true, got %v", data["hasResendKey"])
	}

	// 2. POST /admin/api/settings/smtp (updating email settings)
	postBody := `{"provider":"resend","resendApiKey":"re_newkey_789","resendFromEmail":"TN24 <alerts@tn24.in>","host":"smtp.gmail.com","port":"465","user":"user@gmail.com","pass":"apppassword","from":"TN24 <alerts@tn24.in>"}`
	reqPost := httptest.NewRequest(http.MethodPost, "/admin/api/settings/smtp", strings.NewReader(postBody))
	rrPost := httptest.NewRecorder()
	handler.HandleSMTPSettings(rrPost, reqPost)
	if rrPost.Code != http.StatusOK {
		t.Errorf("Expected 200 for POST email settings, got %d", rrPost.Code)
	}
	if cfg.ResendAPIKey != "re_newkey_789" {
		t.Errorf("Expected cfg.ResendAPIKey updated to re_newkey_789, got %s", cfg.ResendAPIKey)
	}

	// 3. Test email validation in HandleTestSMTP
	reqTestInvalid := httptest.NewRequest(http.MethodPost, "/admin/api/settings/smtp/test", strings.NewReader(`{"testEmail":"notanemail"}`))
	rrTestInvalid := httptest.NewRecorder()
	handler.HandleTestSMTP(rrTestInvalid, reqTestInvalid)
	if rrTestInvalid.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid test email, got %d", rrTestInvalid.Code)
	}
}

func TestAdminImageKitSettingsHandlers(t *testing.T) {
	cfg := &config.Config{
		ImageKitURLEndpoint: "https://ik.imagekit.io/tn24test",
		ImageKitPublicKey:   "public_test_key_123",
		ImageKitPrivateKey:  "private_test_key_456",
		ImageKitFolder:      "articles_test",
	}
	handler := NewHandler(nil, nil, nil, nil, cfg)

	// 1. GET /admin/api/settings/imagekit
	reqGet := httptest.NewRequest(http.MethodGet, "/admin/api/settings/imagekit", nil)
	rrGet := httptest.NewRecorder()
	handler.HandleImageKitSettings(rrGet, reqGet)
	if rrGet.Code != http.StatusOK {
		t.Errorf("Expected 200 for GET imagekit settings, got %d", rrGet.Code)
	}

	var res map[string]interface{}
	if err := json.NewDecoder(rrGet.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to parse json: %v", err)
	}
	if res["success"] != true {
		t.Error("Expected success: true in imagekit response")
	}
	data, ok := res["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected data object in imagekit response")
	}
	if data["urlEndpoint"] != "https://ik.imagekit.io/tn24test" {
		t.Errorf("Expected urlEndpoint https://ik.imagekit.io/tn24test, got %v", data["urlEndpoint"])
	}
	if data["publicKey"] != "public_test_key_123" {
		t.Errorf("Expected publicKey public_test_key_123, got %v", data["publicKey"])
	}
	if data["hasPrivateKey"] != true {
		t.Errorf("Expected hasPrivateKey to be true, got %v", data["hasPrivateKey"])
	}

	// 2. Test invalid method for TestImageKit (e.g. GET instead of POST)
	reqTestGet := httptest.NewRequest(http.MethodGet, "/admin/api/settings/imagekit/test", nil)
	rrTestGet := httptest.NewRecorder()
	handler.HandleTestImageKit(rrTestGet, reqTestGet)
	if rrTestGet.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 for GET on test endpoint, got %d", rrTestGet.Code)
	}
}


