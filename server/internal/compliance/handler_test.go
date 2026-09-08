package compliance

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestComplianceHandlerValidation(t *testing.T) {
	handler := NewHandler(nil, nil)

	// 1. Missing Grievance Email
	body, _ := json.Marshal(map[string]string{
		"complainantEmail": "",
		"category":         "DEFAMATION",
		"description":      "Objectionable content report",
	})
	req := httptest.NewRequest(http.MethodPost, "/compliance/grievance", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.HandleSubmitGrievance(rr, req)
	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code 422 for missing complainantEmail, got %d", rr.Code)
	}

	// 2. Missing Legal Takedown Order Reference
	bodyTakedown, _ := json.Marshal(map[string]string{
		"contentId":   "content-123",
		"orderRef":    "",
		"adminUserId": "admin-123",
	})
	reqTakedown := httptest.NewRequest(http.MethodPost, "/compliance/takedown", bytes.NewReader(bodyTakedown))
	rrTakedown := httptest.NewRecorder()

	handler.HandleLegalTakedown(rrTakedown, reqTakedown)
	if rrTakedown.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code 422 for missing orderRef, got %d", rrTakedown.Code)
	}
}
