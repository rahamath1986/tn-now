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
	if !strings.Contains(body, "swg-basic.js") || !strings.Contains(body, "CAowvuTHDA:openaccess") {
		t.Errorf("expected Google Reader Revenue Manager script in static pages")
	}
}

func TestReaderRevenueManagerSWG(t *testing.T) {
	portalHTML := RenderPortalPage()
	if !strings.Contains(portalHTML, "https://news.google.com/swg/js/v1/swg-basic.js") {
		t.Errorf("expected SWG basic js script tag in RenderPortalPage")
	}
	if !strings.Contains(portalHTML, "CAowvuTHDA:openaccess") {
		t.Errorf("expected product ID CAowvuTHDA:openaccess in RenderPortalPage")
	}
}

func TestRSSFeed(t *testing.T) {
	h := NewPortalHandler(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "https://www.tn24.in/rss.xml", nil)
	rec := httptest.NewRecorder()
	h.HandleRSSFeed(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK from /rss.xml, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/rss+xml") {
		t.Errorf("expected application/rss+xml content type, got %s", contentType)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `<rss version="2.0"`) {
		t.Errorf("expected RSS 2.0 version attribute")
	}
	if !strings.Contains(body, `xmlns:media="http://search.yahoo.com/mrss/"`) {
		t.Errorf("expected Media RSS namespace")
	}
	if !strings.Contains(body, "<title>TN24") {
		t.Errorf("expected RSS channel title")
	}
	if !strings.Contains(body, "<link>https://www.tn24.in/portal</link>") {
		t.Errorf("expected canonical portal link in RSS channel")
	}
}

func TestSitemap38DistrictsAndRobots(t *testing.T) {
	h := NewPortalHandler(nil, nil)

	// Verify all 38 districts in sitemap.xml
	reqSitemap := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	recSitemap := httptest.NewRecorder()
	h.HandleSitemapXML(recSitemap, reqSitemap)
	sitemapBody := recSitemap.Body.String()

	districts := []string{
		"Chennai", "Coimbatore", "Madurai", "Salem", "Tiruchirappalli",
		"Tirunelveli", "Kanyakumari", "Thanjavur", "Dindigul", "Vellore",
		"Ranipet", "Tirupathur", "Chengalpattu", "Mayiladuthurai", "Tenkasi",
	}
	for _, d := range districts {
		if !strings.Contains(sitemapBody, "portal?district="+d) {
			t.Errorf("expected district %s in sitemap.xml", d)
		}
	}

	// Verify robots.txt allows RSS
	reqRobots := httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	recRobots := httptest.NewRecorder()
	h.HandleRobotsTxt(recRobots, reqRobots)
	robotsBody := recRobots.Body.String()

	if !strings.Contains(robotsBody, "Allow: /rss.xml") {
		t.Errorf("expected robots.txt to allow /rss.xml")
	}
	if !strings.Contains(robotsBody, "Allow: /feed.xml") {
		t.Errorf("expected robots.txt to allow /feed.xml")
	}
}

func TestSEOInjectors(t *testing.T) {
	h := NewPortalHandler(nil, nil)
	baseHTML := "<html><head><title>Original</title><link rel=\"canonical\" href=\"https://www.tn24.in/portal\"></head><body><h1>Heading</h1></body></html>"

	// Test District injector
	distHTML := h.injectDistrictMetadata(baseHTML, "Madurai", nil)
	if !strings.Contains(distHTML, "https://www.tn24.in/portal?district=Madurai") {
		t.Errorf("expected district canonical in injected HTML, got %s", distHTML)
	}
	if !strings.Contains(distHTML, "மதுரை செய்திகள்") {
		t.Errorf("expected Tamil district name in injected HTML")
	}

	// Test Category injector
	catHTML := h.injectCategoryMetadata(baseHTML, "cinema", nil)
	if !strings.Contains(catHTML, "https://www.tn24.in/portal?category=cinema") {
		t.Errorf("expected category canonical in injected HTML, got %s", catHTML)
	}
	if !strings.Contains(catHTML, "சினிமா செய்திகள்") {
		t.Errorf("expected Tamil category name in injected HTML")
	}

	// Test Viral injector
	viralHTML := h.injectViralMetadata(baseHTML, nil)
	if !strings.Contains(viralHTML, "https://www.tn24.in/portal?viral=true") {
		t.Errorf("expected viral canonical in injected HTML, got %s", viralHTML)
	}

	// Test Search injector
	searchHTML := h.injectSearchMetadata(baseHTML, "breaking news", nil)
	if !strings.Contains(searchHTML, "breaking news") {
		t.Errorf("expected search query in title")
	}
	if !strings.Contains(searchHTML, `noindex, follow`) {
		t.Errorf("expected noindex, follow on search page")
	}
}

func TestVideoSitemapXML(t *testing.T) {
	h := NewPortalHandler(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "https://www.tn24.in/sitemap-video.xml", nil)
	rec := httptest.NewRecorder()
	h.HandleVideoSitemapXML(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200 OK from /sitemap-video.xml, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/xml") {
		t.Errorf("expected application/xml content type, got %s", contentType)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `xmlns:video="http://www.google.com/schemas/sitemap-video/1.1"`) {
		t.Errorf("expected Google Video sitemap namespace in /sitemap-video.xml")
	}
	if !strings.Contains(body, `<urlset`) {
		t.Errorf("expected <urlset> in /sitemap-video.xml")
	}
	if !strings.Contains(body, `<url>`) {
		t.Errorf("expected at least one <url> in /sitemap-video.xml to satisfy Google Search Console")
	}
	if !strings.Contains(body, `<video:video>`) {
		t.Errorf("expected <video:video> in /sitemap-video.xml")
	}
}

func TestRobotsTxtGooglebotVideo(t *testing.T) {
	h := NewPortalHandler(nil, nil)

	reqRobots := httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	recRobots := httptest.NewRecorder()
	h.HandleRobotsTxt(recRobots, reqRobots)
	robotsBody := recRobots.Body.String()

	if !strings.Contains(robotsBody, "sitemap-video.xml") {
		t.Errorf("expected robots.txt to mention sitemap-video.xml")
	}
	if !strings.Contains(robotsBody, "User-agent: Googlebot-Video") {
		t.Errorf("expected robots.txt to contain User-agent: Googlebot-Video")
	}
	if !strings.Contains(robotsBody, "Allow: /sitemap-video.xml") {
		t.Errorf("expected robots.txt to allow /sitemap-video.xml")
	}
}



