package portal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserCornerRoutesRegistration(t *testing.T) {
	mux := http.NewServeMux()
	handler := NewPortalHandler(nil, nil)
	handler.RegisterUserCornerRoutes(mux)

	tests := []struct {
		name       string
		path       string
		method     string
		wantStatus int
		contains   []string
	}{
		{
			name:       "Submit Page loads with Quill WYSIWYG editor",
			path:       "/submit",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			contains: []string{
				"வாசகர் அரங்கம்",
				"quill.snow.css",
				"quill.min.js",
				"districtSelect",
				"categorySelect",
				"titleInput",
				"editor-container",
				"Under Review",
			},
		},
		{
			name:       "Front Login Page renders with front theme",
			path:       "/login",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			contains: []string{
				"உள்நுழைக",
				"loginForm",
				"registerForm",
				"tn24-logo.svg",
			},
		},
		{
			name:       "Front Register Page renders with front theme",
			path:       "/register",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			contains: []string{
				"பதிவு செய்க",
				"regFullName",
				"regUsername",
				"regEmail",
				"regPassword",
				"otpVerifyForm",
			},
		},
		{
			name:       "User Corner Hub Page loads with SEO and citizen journalism branding",
			path:       "/user-corner",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
			contains: []string{
				"வாசகர் அரங்கம்",
				"Community Voices",
				"userCornerFeed",
			},
		},
		{
			name:       "My Posts dashboard redirects unauthenticated visitors to login",
			path:       "/user/my-posts",
			method:     http.MethodGet,
			wantStatus: http.StatusSeeOther,
			contains:   []string{},
		},
		{
			name:       "User Profile dashboard redirects unauthenticated visitors to login",
			path:       "/user/profile",
			method:     http.MethodGet,
			wantStatus: http.StatusSeeOther,
			contains:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("%s: got status %d, want %d", tt.name, w.Code, tt.wantStatus)
			}

			body := w.Body.String()
			for _, needle := range tt.contains {
				if !strings.Contains(body, needle) {
					t.Errorf("%s: body missing expected string %q", tt.name, needle)
				}
			}
		})
	}
}

func TestAuthorEEATStructuredData(t *testing.T) {
	html := renderAuthorProfilePageHTML(
		"karthik_mds",
		"Karthikeyan Sundaram",
		"Citizen journalist covering rural education in Madurai",
		"",
		5,
		50,
		95,
		[]ArticleItem{
			{
				ID:       "art-12345",
				Slug:     "madurai-school-renovation-12345",
				Title:    "மதுரை அரசு பள்ளி புதுப்பிப்பு பணி",
				District: "Madurai",
			},
		},
	)

	// Verify Google E-E-A-T Person schema
	expectedSchemaTokens := []string{
		`"@type": "Person"`,
		`"name": "Karthikeyan Sundaram"`,
		`"url": "https://www.tn24.in/author/karthik_mds"`,
		`"jobTitle": "Community Contributor"`,
		`"name": "TN24 — Tamil Nadu News"`,
		`madurai-school-renovation-12345`,
	}

	for _, token := range expectedSchemaTokens {
		if !strings.Contains(html, token) {
			t.Errorf("author profile page missing token %q", token)
		}
	}
}

