package scraper

import (
	"strings"
	"testing"
)

func TestGenerateArticleSlug(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		id       string
		mustHave []string
	}{
		{
			name:     "Pure Tamil title with keywords",
			title:    "சென்னையில் கனமழை: பள்ளிகளுக்கு நாளை விடுமுறை அறிவிப்பு",
			id:       "3f82a1bc-79f2-4e89-8d19-411a7a28e8d2",
			mustHave: []string{"chennai", "kanamazhai", "schools", "holiday", "3f82a1bc"},
		},
		{
			name:     "Pure English title",
			title:    "Tamil Nadu Government Announces Free Laptop Scheme For College Students",
			id:       "c8b18a12-45e2-4567-89ab-cdef01234567",
			mustHave: []string{"tamil", "nadu", "government", "announces", "c8b18a12"},
		},
		{
			name:     "Mixed Tamil and English",
			title:    "IPL 2026: CSK அணியின் புதிய கேப்டன் யார்?",
			id:       "a1b2c3d4-0000-0000-0000-000000000000",
			mustHave: []string{"ipl", "2026", "csk", "a1b2c3d4"},
		},
		{
			name:     "Special characters, emojis and quotes",
			title:    "🔥 'பட்ஜெட் 2026' - முதலமைச்சர் மு.க.ஸ்டாலின் முக்கிய உரை! [முழு விவரம்]",
			id:       "e5f6a7b8-1111-2222-3333-444455556666",
			mustHave: []string{"budget", "2026", "cm", "stalin", "e5f6a7b8"},
		},
		{
			name:     "Empty title fallback",
			title:    "",
			id:       "99887766-aaaa-bbbb-cccc-ddddeeeeffff",
			mustHave: []string{"news-99887766"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			slug := GenerateArticleSlug(tc.title, tc.id)
			if slug == "" {
				t.Fatalf("expected non-empty slug, got empty")
			}
			// Slug must be valid URL path (only lowercase alphanumeric and hyphens)
			for _, r := range slug {
				if !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9') && r != '-' {
					t.Errorf("slug contains invalid URL character %q: %s", r, slug)
				}
			}
			// Check for expected keywords
			for _, expected := range tc.mustHave {
				if !strings.Contains(slug, expected) {
					t.Errorf("slug %q does not contain expected substring %q", slug, expected)
				}
			}
			// No consecutive hyphens
			if strings.Contains(slug, "--") {
				t.Errorf("slug contains consecutive hyphens: %s", slug)
			}
			// No leading or trailing hyphen
			if strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
				t.Errorf("slug has leading/trailing hyphen: %s", slug)
			}
		})
	}
}
