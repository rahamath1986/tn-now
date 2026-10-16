package admin

import (
	"context"
	"strings"
	"testing"
)

func TestAdminRewriterIntegration(t *testing.T) {
	rewriter := NewArticleRewriter("")

	req := ArticleRewriteRequest{
		ID:          "cont_test_456",
		Title:       "சென்னை மெட்ரோ ரயில் விரிவாக்கம் - The Hindu",
		Description: "சென்னை மெட்ரோ ரயில் 2-ம் கட்ட விரிவாக்கப் பணிகள் தீவிரமாக நடைபெற்று வருகின்றன.\n\nFollow us on Twitter for live updates.\n\nஅடுத்த ஆண்டு இறுதிக்குள் பணிகள் நிறைவடையும் என அதிகாரிகள் தெரிவித்துள்ளனர்.",
		District:    "Chennai",
		Category:    "Civic",
		Language:    "ta",
	}

	res, err := rewriter.RewriteArticle(context.Background(), req)
	if err != nil {
		t.Fatalf("RewriteArticle failed: %v", err)
	}

	if res.Title != "சென்னை மெட்ரோ ரயில் விரிவாக்கம்" {
		t.Errorf("Expected branding stripped title, got: %s", res.Title)
	}
	if strings.Contains(res.Description, "Twitter") {
		t.Errorf("Expected Twitter CTA to be removed from description, got: %s", res.Description)
	}
	if res.Slug == "" {
		t.Errorf("Expected generated slug to be non-empty")
	}
	if res.Language != "ta" {
		t.Errorf("Expected language 'ta', got: %s", res.Language)
	}
}
