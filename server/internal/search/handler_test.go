package search

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchValidation(t *testing.T) {
	handler := NewHandler(nil)

	// Missing q parameter
	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	rr := httptest.NewRecorder()

	handler.HandleSearch(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code 400 for missing q query parameter, got %d", rr.Code)
	}
}
