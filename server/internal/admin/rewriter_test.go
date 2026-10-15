package admin

import (
	"context"
	"strings"
	"testing"
)

func TestValidateLanguageIntegrity(t *testing.T) {
	// Pure Tamil headline and text
	tamilTitle := "சென்னையில் பலத்த மழை: பள்ளிகளுக்கு விடுமுறை அறிவிப்பு"
	tamilBody := "சென்னையில் தொடர்ந்து பெய்து வரும் கனமழை காரணமாக முன்னெச்சரிக்கை நடவடிக்கையாக பள்ளிகளுக்கு விடுமுறை அளிக்கப்பட்டுள்ளது."
	if !validateLanguageIntegrity("ta", tamilTitle, tamilBody) {
		t.Errorf("Expected pure Tamil text to pass language integrity validation")
	}

	// Tanglish text (Tamil words in Latin script) - MUST BE REJECTED FOR 'ta'
	tanglishTitle := "Chennai-la heavy rain: Schools-ku leave announce pannitanga"
	tanglishBody := "Chennai-la romba mazhai peiyuthu. Athunaala students-ku leave."
	if validateLanguageIntegrity("ta", tanglishTitle, tanglishBody) {
		t.Errorf("Expected Tanglish text to FAIL language integrity validation for Tamil 'ta'")
	}

	// Pure English text
	engTitle := "Chennai Rain Alert: Schools Closed Across Coastal Districts"
	engBody := "Heavy rainfall lashed coastal Tamil Nadu on Thursday, prompting district collectors to declare holidays."
	if !validateLanguageIntegrity("en", engTitle, engBody) {
		t.Errorf("Expected pure English text to pass language integrity validation for 'en'")
	}
}

func TestStripBrandingSuffixes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "விவசாயிகளுக்கு புதிய திட்டம் அறிவிப்பு - BBC News தமிழ்",
			expected: "விவசாயிகளுக்கு புதிய திட்டம் அறிவிப்பு",
		},
		{
			input:    "TN Assembly Session Begins Today | Dinamalar",
			expected: "TN Assembly Session Begins Today",
		},
		{
			input:    "போக்குவரத்து நெரிசல் குறைப்பு நடவடிக்கை – மாலை மலர்",
			expected: "போக்குவரத்து நெரிசல் குறைப்பு நடவடிக்கை",
		},
	}

	for _, tc := range tests {
		got := stripBrandingSuffixes(tc.input)
		if got != tc.expected {
			t.Errorf("stripBrandingSuffixes(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestStripBoilerplate(t *testing.T) {
	input := `அமைச்சரவைக் கூட்டம் இன்று தலைமைச் செயலகத்தில் நடைபெற்றது.

Subscribe to our channel for more updates.

முக்கிய முடிவுகள் பல எட்டப்பட்டுள்ளன.

செய்திகளை உடனுக்குடன் பெற வாட்ஸ்அப் சேனலில் இணையவும்.`

	cleaned := stripBoilerplate(input)
	if strings.Contains(cleaned, "Subscribe") {
		t.Errorf("Expected 'Subscribe' junk to be removed, got: %s", cleaned)
	}
	if strings.Contains(cleaned, "வாட்ஸ்அப் சேனலில்") {
		t.Errorf("Expected WhatsApp CTA junk to be removed, got: %s", cleaned)
	}
	if !strings.Contains(cleaned, "அமைச்சரவைக் கூட்டம்") {
		t.Errorf("Expected legitimate content to be preserved, got: %s", cleaned)
	}
}

func TestHeuristicRewrite(t *testing.T) {
	rewriter := NewArticleRewriter("")

	req := ArticleRewriteRequest{
		ID:          "cont_test_123",
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
