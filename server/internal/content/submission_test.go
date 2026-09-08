package content

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubmitVideoValidation(t *testing.T) {
	handler := NewSubmissionHandler(nil, nil)

	// 1. Missing fields
	body, _ := json.Marshal(map[string]string{
		"title": "",
		"url":   "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
	})
	req := httptest.NewRequest(http.MethodPost, "/content/submit/video", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.HandleSubmitVideo(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code 422 for missing fields, got %d", rr.Code)
	}

	// 2. Invalid Video URL
	body, _ = json.Marshal(map[string]string{
		"title":        "My Submission Video",
		"url":          "https://invalid-video-host.com/video123",
		"categoryId":   "cat-123",
		"districtId":   "dist-123",
		"authorUserId": "user-123",
	})
	req = httptest.NewRequest(http.MethodPost, "/content/submit/video", bytes.NewReader(body))
	rr = httptest.NewRecorder()

	handler.HandleSubmitVideo(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code 422 for unsupported video platform URL, got %d", rr.Code)
	}
}
