package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tn-now/server/internal/config"
)

func TestGoogleOAuthInitiation(t *testing.T) {
	cfg := &config.Config{
		GoogleClientID:    "test-google-client-id.apps.googleusercontent.com",
		GoogleRedirectURL: "http://localhost:8080/admin/auth/google/callback",
	}
	copilot := NewCopilotEngine(nil, nil, nil, nil, cfg)

	req := httptest.NewRequest(http.MethodGet, "/admin/auth/google", nil)
	rr := httptest.NewRecorder()

	copilot.HandleInitiateGoogleOAuth(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Fatalf("Expected 307 Temporary Redirect to Google, got %d", rr.Code)
	}

	location := rr.Header().Get("Location")
	if !strings.Contains(location, "accounts.google.com/o/oauth2/v2/auth") {
		t.Errorf("Expected redirect to accounts.google.com, got %s", location)
	}
	if !strings.Contains(location, "client_id=test-google-client-id") {
		t.Errorf("Expected client_id in redirect location, got %s", location)
	}
}

func TestGoogleOAuthInitiation_Unconfigured(t *testing.T) {
	cfg := &config.Config{}
	copilot := NewCopilotEngine(nil, nil, nil, nil, cfg)

	req := httptest.NewRequest(http.MethodGet, "/admin/auth/google", nil)
	rr := httptest.NewRecorder()

	copilot.HandleInitiateGoogleOAuth(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 setup prompt, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Google OAuth 2.0 Credentials Required") {
		t.Error("Expected setup instruction page when Google Client ID is unconfigured")
	}
}

func TestGoogleIDTokenHandler_Validation(t *testing.T) {
	copilot := NewCopilotEngine(nil, nil, nil, nil, nil)

	// Missing credential
	req := httptest.NewRequest(http.MethodPost, "/admin/api/auth/google-id-token", strings.NewReader(`{}`))
	rr := httptest.NewRecorder()
	copilot.HandleVerifyGoogleIDToken(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request for missing credential, got %d", rr.Code)
	}
}

func TestAgentPromptHandler(t *testing.T) {
	copilot := NewCopilotEngine(nil, nil, nil, nil, nil)

	promptPayload := strings.NewReader(`{"prompt":"Run full system telemetry and diagnostics"}`)
	req := httptest.NewRequest(http.MethodPost, "/admin/api/agent/prompt", promptPayload)
	rr := httptest.NewRecorder()

	copilot.HandleAgentPrompt(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", rr.Code)
	}

	var res map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to parse agent response: %v", err)
	}

	if res["success"] != true {
		t.Error("Expected success: true in agent response")
	}

	data := res["data"].(map[string]interface{})
	reply, ok := data["reply"].(string)
	if !ok || !strings.Contains(strings.ToLower(reply), "antigravity") {
		t.Errorf("Expected Antigravity engine reference in reply, got %s", reply)
	}

	thoughtTrace, ok := data["thoughtTrace"].([]interface{})
	if !ok || len(thoughtTrace) == 0 {
		t.Error("Expected thought trace steps from Antigravity engine")
	}
}

func TestAgentExecuteHandler(t *testing.T) {
	copilot := NewCopilotEngine(nil, nil, nil, nil, nil)

	payload := strings.NewReader(`{"actionType":"GENERIC_EVALUATE","payload":"Test submission title"}`)
	req := httptest.NewRequest(http.MethodPost, "/admin/api/agent/execute", payload)
	rr := httptest.NewRecorder()

	copilot.HandleAgentExecute(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", rr.Code)
	}

	var res map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatalf("Failed to parse execute response: %v", err)
	}

	if res["success"] != true {
		t.Error("Expected success: true in execute response")
	}
}
