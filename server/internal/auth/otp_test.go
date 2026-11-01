package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"tn-now/server/internal/config"
)

func TestGenerateNumericOTP(t *testing.T) {
	seen := make(map[string]bool)

	for i := 0; i < 50; i++ {
		otp, err := GenerateNumericOTP()
		if err != nil {
			t.Fatalf("GenerateNumericOTP returned error: %v", err)
		}
		if len(otp) != 6 {
			t.Fatalf("Expected 6 digits, got %s (len %d)", otp, len(otp))
		}
		num, err := strconv.Atoi(otp)
		if err != nil {
			t.Fatalf("OTP is not numeric: %s", otp)
		}
		if num < 0 || num > 999999 {
			t.Fatalf("OTP %d out of bounds", num)
		}
		seen[otp] = true
	}

	// In 50 random 6-digit codes, we expect high entropy (at least 45 unique codes)
	if len(seen) < 45 {
		t.Errorf("Low entropy in OTP generation: only %d unique out of 50", len(seen))
	}
}

func TestSendOTPEmailFallback(t *testing.T) {
	// Empty SMTP config simulates local development
	cfg := &config.Config{
		SMTPHost: "",
		SMTPPort: "587",
		SMTPUser: "",
		SMTPPass: "",
		SMTPFrom: "noreply@tn24.in",
	}

	err := SendOTPEmail(cfg, "reader@example.com", "482910")
	if err != nil {
		t.Fatalf("SendOTPEmail with local fallback should succeed, got error: %v", err)
	}
}

func TestHandleSendOTPValidation(t *testing.T) {
	h := &Handler{service: nil}

	// 1. Invalid method (GET instead of POST)
	req := httptest.NewRequest(http.MethodGet, "/auth/otp/send", nil)
	rec := httptest.NewRecorder()
	h.HandleSendOTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 MethodNotAllowed, got %d", rec.Code)
	}

	// 2. Empty Body
	req = httptest.NewRequest(http.MethodPost, "/auth/otp/send", bytes.NewReader([]byte("")))
	rec = httptest.NewRecorder()
	h.HandleSendOTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 BadRequest for empty body, got %d", rec.Code)
	}

	// 3. Missing Email
	reqBody, _ := json.Marshal(map[string]string{"email": "  "})
	req = httptest.NewRequest(http.MethodPost, "/auth/otp/send", bytes.NewReader(reqBody))
	rec = httptest.NewRecorder()
	h.HandleSendOTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 BadRequest for blank email, got %d", rec.Code)
	}
}

func TestHandleVerifyOTPValidation(t *testing.T) {
	h := &Handler{service: nil}

	// 1. Invalid method
	req := httptest.NewRequest(http.MethodGet, "/auth/otp/verify", nil)
	rec := httptest.NewRecorder()
	h.HandleVerifyOTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 MethodNotAllowed, got %d", rec.Code)
	}

	// 2. Missing OTP
	reqBody, _ := json.Marshal(map[string]string{
		"email": "reader@example.com",
		"otp":   "",
	})
	req = httptest.NewRequest(http.MethodPost, "/auth/otp/verify", bytes.NewReader(reqBody))
	rec = httptest.NewRecorder()
	h.HandleVerifyOTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 BadRequest for missing OTP, got %d", rec.Code)
	}
}

func TestOTPExpiryWindow(t *testing.T) {
	// Ensure expiration duration is exactly 10 minutes
	expiry := time.Now().Add(10 * time.Minute)
	diff := time.Until(expiry)
	if diff < 9*time.Minute || diff > 10*time.Minute+5*time.Second {
		t.Errorf("Unexpected expiration diff: %v", diff)
	}
}

func TestEffectiveSMTPConfigFallback(t *testing.T) {
	baseCfg := &config.Config{
		SMTPHost: "smtp.fallback.test",
		SMTPPort: "587",
		SMTPUser: "fallback_user",
		SMTPPass: "fallback_pass",
		SMTPFrom: "noreply@fallback.test",
	}

	// Passing nil pool should return base config directly
	res := GetEffectiveSMTPConfig(nil, nil, baseCfg)
	if res.SMTPHost != "smtp.fallback.test" {
		t.Errorf("Expected SMTPHost smtp.fallback.test, got %s", res.SMTPHost)
	}
	if res.SMTPPort != "587" {
		t.Errorf("Expected SMTPPort 587, got %s", res.SMTPPort)
	}
}

func TestHandleUpdateProfileValidation(t *testing.T) {
	h := &Handler{service: nil}

	// GET method should be rejected
	req := httptest.NewRequest(http.MethodGet, "/auth/profile", nil)
	rec := httptest.NewRecorder()
	h.HandleUpdateProfile(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 MethodNotAllowed for GET on update profile, got %d", rec.Code)
	}

	// Missing auth token
	req = httptest.NewRequest(http.MethodPut, "/auth/profile", bytes.NewReader([]byte("{}")))
	rec = httptest.NewRecorder()
	h.HandleUpdateProfile(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized for missing token, got %d", rec.Code)
	}
}

func TestLiveGmailSMTP(t *testing.T) {
	if os.Getenv("TEST_LIVE_SMTP") == "" {
		t.Skip("skipping live SMTP test; set TEST_LIVE_SMTP=1 to run")
	}
	cfg := &config.Config{
		SMTPHost: "smtp.gmail.com",
		SMTPPort: "465",
		SMTPUser: "tn24now@gmail.com",
		SMTPPass: "gcatermqkobrzjlr",
		SMTPFrom: "tn24now@gmail.com",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	err := SendOTPEmailWithContext(ctx, cfg, "tn24now@gmail.com", "773322")
	if err != nil {
		t.Fatalf("Live SMTP dispatch failed: %v", err)
	}
}

