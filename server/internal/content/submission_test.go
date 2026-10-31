package content

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestSubmitStoryValidation(t *testing.T) {
	handler := NewSubmissionHandler(nil, nil)

	// Missing title and body
	body, _ := json.Marshal(map[string]string{
		"title": "",
		"body":  "",
	})
	req := httptest.NewRequest(http.MethodPost, "/content/submit/story", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.HandleSubmitStory(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code 422 for empty story title and body, got %d", rr.Code)
	}
}

func TestHandleUploadImageValidation(t *testing.T) {
	handler := NewSubmissionHandler(nil, nil)

	// 1. Not a multipart request
	req := httptest.NewRequest(http.MethodPost, "/content/upload/image", strings.NewReader("bad payload"))
	rr := httptest.NewRecorder()
	handler.HandleUploadImage(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for non-multipart request, got %d", rr.Code)
	}

	// 2. Multipart request without "image" field
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("some_other_field", "value")
	writer.Close()

	req2 := httptest.NewRequest(http.MethodPost, "/content/upload/image", body)
	req2.Header.Set("Content-Type", writer.FormDataContentType())
	rr2 := httptest.NewRecorder()
	handler.HandleUploadImage(rr2, req2)
	if rr2.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for missing image field, got %d", rr2.Code)
	}

	// 3. Multipart with non-image file (e.g., text/plain)
	body3 := &bytes.Buffer{}
	writer3 := multipart.NewWriter(body3)
	part, _ := writer3.CreateFormFile("image", "test.txt")
	part.Write([]byte("this is plain text"))
	writer3.Close()

	req3 := httptest.NewRequest(http.MethodPost, "/content/upload/image", body3)
	req3.Header.Set("Content-Type", writer3.FormDataContentType())
	rr3 := httptest.NewRecorder()
	handler.HandleUploadImage(rr3, req3)
	if rr3.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for text/plain file, got %d", rr3.Code)
	}
}


