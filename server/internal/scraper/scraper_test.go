package scraper

import (
	"context"
	"testing"
)

func TestParseHTMLPage(t *testing.T) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Madurai Meenakshi Amman Temple Grand Chithirai Festival Begins</title>
    <meta property="og:title" content="Madurai Meenakshi Amman Temple Grand Chithirai Festival Begins" />
    <meta property="og:description" content="Thousands of devotees gathered in Madurai for the annual celebrations." />
    <meta property="og:image" content="https://example.com/madurai_temple.jpg" />
</head>
<body>
    <iframe src="https://www.youtube.com/embed/dQw4w9WgXcQ"></iframe>
</body>
</html>`

	item := parseHTMLPage("https://example.com/news/madurai", html)
	if item == nil {
		t.Fatal("Expected item to be parsed, got nil")
	}
	if item.District != "Madurai" {
		t.Errorf("Expected district Madurai, got %s", item.District)
	}
	if item.ContentType != "VIDEO_LINK" {
		t.Errorf("Expected VIDEO_LINK content type, got %s", item.ContentType)
	}
	if item.VideoID != "dQw4w9WgXcQ" {
		t.Errorf("Expected video ID dQw4w9WgXcQ, got %s", item.VideoID)
	}
	if len(item.ImageURLs) == 0 || item.ImageURLs[0] != "https://example.com/madurai_temple.jpg" {
		t.Errorf("Expected og:image, got %v", item.ImageURLs)
	}
	if item.Status != "PENDING" {
		t.Errorf("Expected staged status PENDING, got %s", item.Status)
	}
}

func TestScrapeAndStageValidation(t *testing.T) {
	_, err := ScrapeAndStage(context.Background(), nil, "")
	if err == nil {
		t.Error("Expected error for empty URL, got nil")
	}
}
