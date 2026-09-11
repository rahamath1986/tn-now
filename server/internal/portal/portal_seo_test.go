package portal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPortalSEOKeywords(t *testing.T) {
	html := RenderPortalPage()

	targetKeywords := []string{
		"tn 24",
		"news tamil 24x7 live",
		"tamil nadu news",
		"news live tamilnadu",
		"தமிழ் செய்திகள்",
		"today news in tamil",
		"news tamil today",
		"tamil news online",
		"latest tamil news",
		"tamil nadu news in tamil",
		"news tamil nadu",
	}

	htmlLower := strings.ToLower(html)

	for _, kw := range targetKeywords {
		if !strings.Contains(htmlLower, strings.ToLower(kw)) {
			t.Errorf("expected RenderPortalPage() to contain keyword %q, but was not found", kw)
		}
	}

	// Verify SEO Title
	if !strings.Contains(html, "<title>TN24 &mdash; Tamil Nadu News | தமிழ் செய்திகள் | Latest Tamil News Today Live 24x7 | TN 24</title>") {
		t.Errorf("expected target title tag not found")
	}

	// Verify Schema.org JSON-LD exists and has target keywords
	if !strings.Contains(html, `application/ld+json`) {
		t.Errorf("expected application/ld+json schema not found")
	}
	if !strings.Contains(html, `"NewsMediaOrganization"`) {
		t.Errorf("expected NewsMediaOrganization schema not found")
	}
	if !strings.Contains(html, `"WebSite"`) {
		t.Errorf("expected WebSite schema not found")
	}

	// Verify OpenGraph and Twitter tags
	if !strings.Contains(html, `property="og:title"`) || !strings.Contains(html, `name="twitter:title"`) {
		t.Errorf("missing OpenGraph or Twitter title tags")
	}

	// Verify Trending Topics in footer
	if !strings.Contains(html, "Trending Tamil News Searches") {
		t.Errorf("expected trending keywords footer section not found")
	}
}

func TestSitemapAndRobots(t *testing.T) {
	h := NewPortalHandler(nil, nil)

	// Test Robots.txt
	reqRobots := httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	recRobots := httptest.NewRecorder()
	h.HandleRobotsTxt(recRobots, reqRobots)
	robotsBody := recRobots.Body.String()

	if !strings.Contains(robotsBody, "sitemap.xml") || !strings.Contains(robotsBody, "sitemap-news.xml") {
		t.Errorf("expected robots.txt to mention sitemap.xml and sitemap-news.xml, got %s", robotsBody)
	}

	// Test Sitemap.xml
	reqSitemap := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	recSitemap := httptest.NewRecorder()
	h.HandleSitemapXML(recSitemap, reqSitemap)
	sitemapBody := recSitemap.Body.String()

	if !strings.Contains(sitemapBody, "<urlset") {
		t.Errorf("expected <urlset> in sitemap, got %s", sitemapBody)
	}
	if !strings.Contains(sitemapBody, "portal?q=tamil+nadu+news") {
		t.Errorf("expected target keyword URL in sitemap.xml")
	}

	// Test Google News Sitemap
	reqNews := httptest.NewRequest(http.MethodGet, "/sitemap-news.xml", nil)
	recNews := httptest.NewRecorder()
	h.HandleNewsSitemapXML(recNews, reqNews)
	newsBody := recNews.Body.String()

	if !strings.Contains(newsBody, "http://www.google.com/schemas/sitemap-news/0.9") {
		t.Errorf("expected Google News sitemap namespace in /sitemap-news.xml")
	}
}

func TestTermsOfService(t *testing.T) {
	h := NewPortalHandler(nil, nil)

	reqTerms := httptest.NewRequest(http.MethodGet, "/terms", nil)
	recTerms := httptest.NewRecorder()
	h.HandleTermsOfService(recTerms, reqTerms)

	if recTerms.Code != http.StatusOK {
		t.Errorf("expected 200 OK from /terms, got %d", recTerms.Code)
	}

	body := recTerms.Body.String()
	if !strings.Contains(body, "Terms of Service") {
		t.Errorf("expected Terms of Service in body")
	}
	if !strings.Contains(body, "Privacy Policy") {
		t.Errorf("expected link to Privacy Policy in Terms of Service page")
	}
}

