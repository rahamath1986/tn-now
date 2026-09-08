package social

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSocialHandlerValidation(t *testing.T) {
	handler := NewHandler(nil, nil)

	// 1. Like Missing Fields
	body, _ := json.Marshal(map[string]string{
		"contentId": "",
		"userId":    "user-123",
	})
	req := httptest.NewRequest(http.MethodPost, "/social/like", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.HandleLike(rr, req)
	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code 422 for missing contentId, got %d", rr.Code)
	}

	// 2. Comment Empty Body
	bodyComment, _ := json.Marshal(map[string]string{
		"contentId": "content-123",
		"userId":    "user-123",
		"body":      "   ",
	})
	reqComment := httptest.NewRequest(http.MethodPost, "/social/comment", bytes.NewReader(bodyComment))
	rrComment := httptest.NewRecorder()

	handler.HandleComment(rrComment, reqComment)
	if rrComment.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code 422 for empty comment body, got %d", rrComment.Code)
	}
}
