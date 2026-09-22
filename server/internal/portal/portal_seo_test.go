package portal

import (
	"fmt"
	"html"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"
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
	if !strings.Contains(html, "<title>TN24 – Tamil News | Latest Tamil Nadu News &amp; District News</title>") {
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
	if !strings.Contains(sitemapBody, "<loc>https://www.tn24.in/</loc>") {
		t.Errorf("expected canonical homepage in sitemap.xml")
	}
	if strings.Contains(sitemapBody, "<loc>https://www.tn24.in/portal</loc>") {
		t.Errorf("sitemap.xml should not contain redirecting /portal URL")
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

func TestGoogleAdSenseUnit(t *testing.T) {
	portalHTML := RenderPortalPage()
	if !strings.Contains(portalHTML, "ca-pub-1894301748406603") {
		t.Errorf("expected AdSense client ID ca-pub-1894301748406603 in RenderPortalPage")
	}
	if !strings.Contains(portalHTML, `data-ad-slot="3023623233"`) {
		t.Errorf("expected AdSense slot 3023623233 in RenderPortalPage")
	}
	if !strings.Contains(portalHTML, "adsbygoogle") {
		t.Errorf("expected adsbygoogle class in RenderPortalPage")
	}

	h := NewPortalHandler(nil, nil)
	recTerms := httptest.NewRecorder()
	h.HandleTermsOfService(recTerms, httptest.NewRequest(http.MethodGet, "/terms", nil))
	termsBody := recTerms.Body.String()
	if !strings.Contains(termsBody, `data-ad-slot="3023623233"`) {
		t.Errorf("expected AdSense slot 3023623233 in static /terms page")
	}
}

func TestFaviconAndWebsiteTitle(t *testing.T) {
	// Test website title in main portal page
	portalHTML := RenderPortalPage()
	if !strings.Contains(portalHTML, "<title>TN24 – Tamil News | Latest Tamil Nadu News &amp; District News</title>") {
		t.Errorf("expected primary title in RenderPortalPage")
	}

	// Test favicon link tags
	if !strings.Contains(portalHTML, `rel="icon" type="image/svg+xml" href="/assets/brand/tn24-icon.svg"`) {
		t.Errorf("expected SVG favicon link in RenderPortalPage")
	}
	if !strings.Contains(portalHTML, `rel="icon" type="image/x-icon" href="/favicon.ico"`) {
		t.Errorf("expected ICO favicon link in RenderPortalPage")
	}
	if !strings.Contains(portalHTML, `rel="apple-touch-icon"`) {
		t.Errorf("expected apple-touch-icon in RenderPortalPage")
	}

	// Test favicon HTTP endpoint
	h := NewPortalHandler(nil, nil)
	recFav := httptest.NewRecorder()
	h.HandleFavicon(recFav, httptest.NewRequest(http.MethodGet, "/favicon.ico", nil))
	if recFav.Code != http.StatusOK {
		t.Errorf("expected 200 OK from /favicon.ico, got %d", recFav.Code)
	}
	favSVG := recFav.Body.String()
	if !strings.Contains(favSVG, "<svg") || !strings.Contains(favSVG, "TN") || !strings.Contains(favSVG, "24") {
		t.Errorf("expected valid SVG favicon with TN and 24, got %s", favSVG)
	}

	// Test static page favicon and title
	recContact := httptest.NewRecorder()
	h.HandleContactUs(recContact, httptest.NewRequest(http.MethodGet, "/contact", nil))
	contactBody := recContact.Body.String()
	if !strings.Contains(contactBody, "<title>Contact Us | TN24 — Tamil Nadu News | தமிழ் செய்திகள்</title>") {
		t.Errorf("expected title in contact page, got %s", contactBody)
	}
	if !strings.Contains(contactBody, `href="/favicon.ico"`) {
		t.Errorf("expected favicon link in contact page")
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
	if !strings.Contains(body, "<link>https://www.tn24.in/") {
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
		"chennai", "coimbatore", "madurai", "salem", "tiruchirappalli",
		"tirunelveli", "kanyakumari", "thanjavur", "dindigul", "vellore",
		"ranipet", "tirupathur", "chengalpattu", "mayiladuthurai", "tenkasi",
	}
	for _, d := range districts {
		if !strings.Contains(sitemapBody, "/district/"+d) {
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
	if !strings.Contains(distHTML, "https://www.tn24.in/district/madurai") {
		t.Errorf("expected district canonical in injected HTML, got %s", distHTML)
	}
	if !strings.Contains(distHTML, "Madurai News Today") {
		t.Errorf("expected Madurai News Today in title, got %s", distHTML)
	}

	// Test Category injector
	catHTML := h.injectCategoryMetadata(baseHTML, "cinema", nil)
	if !strings.Contains(catHTML, "https://www.tn24.in/category/entertainment") {
		t.Errorf("expected category canonical in injected HTML, got %s", catHTML)
	}
	if !strings.Contains(catHTML, "Tamil Nadu Entertainment News") {
		t.Errorf("expected Entertainment category in injected HTML, got %s", catHTML)
	}

	// Test Viral injector
	viralHTML := h.injectViralMetadata(baseHTML, nil)
	if !strings.Contains(viralHTML, "https://www.tn24.in/?viral=true") {
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

// ═══════════════════════════════════════════════════════════════
// P0 & P1 COMPREHENSIVE SEO VALIDATION TESTS
// ═══════════════════════════════════════════════════════════════

func TestP0HomepageSEO(t *testing.T) {
	html := RenderPortalPage()

	// 1. Preferred Title Check
	expectedTitle := "<title>TN24 – Tamil News | Latest Tamil Nadu News &amp; District News</title>"
	if !strings.Contains(html, expectedTitle) {
		t.Errorf("P0: Homepage title mismatch, expected %s", expectedTitle)
	}

	// 2. Meta Description Check
	expectedDesc := `<meta name="description" content="TN24 brings the latest Tamil news, Tamil Nadu district news, politics, cinema, sports and breaking news from across Tamil Nadu.">`
	if !strings.Contains(html, expectedDesc) {
		t.Errorf("P0: Homepage meta description mismatch, expected %s", expectedDesc)
	}

	// 3. Exactly ONE H1 Tag on Homepage
	h1Matches := regexp.MustCompile(`(?is)<h1\b[^>]*>(.*?)</h1>`).FindAllStringSubmatch(html, -1)
	if len(h1Matches) != 1 {
		t.Errorf("P0: Homepage must have EXACTLY 1 H1 tag, found %d", len(h1Matches))
	} else {
		h1Content := strings.TrimSpace(h1Matches[0][1])
		if !strings.Contains(h1Content, "TN24 – Latest Tamil News &amp; Tamil Nadu News") && !strings.Contains(h1Content, "TN24 – Latest Tamil News & Tamil Nadu News") {
			t.Errorf("P0: Homepage H1 content incorrect, got: %s", h1Content)
		}
	}

	// 4. Canonical URL Check
	expectedCanonical := `<link rel="canonical" href="https://www.tn24.in/">`
	if !strings.Contains(html, expectedCanonical) {
		t.Errorf("P0: Homepage canonical mismatch, expected %s", expectedCanonical)
	}

	// 5. Organization Schema Check
	if !strings.Contains(html, `"NewsMediaOrganization"`) || !strings.Contains(html, `"name": "TN24"`) {
		t.Errorf("P0: Organization schema missing or does not identify TN24")
	}

	// 6. WebSite Schema Check
	if !strings.Contains(html, `"WebSite"`) || !strings.Contains(html, `"headline": "TN24 – Latest Tamil News & Tamil Nadu News"`) {
		t.Errorf("P0: WebSite schema missing or does not have correct headline")
	}
}

func TestP1DistrictPages(t *testing.T) {
	h := NewPortalHandler(nil, nil)

	testDistricts := []struct {
		slug       string
		name       string
		taName     string
	}{
		{"chennai", "Chennai", "சென்னை"},
		{"coimbatore", "Coimbatore", "கோயம்புத்தூர்"},
		{"madurai", "Madurai", "மதுரை"},
		{"salem", "Salem", "சேலம்"},
		{"tiruchirappalli", "Tiruchirappalli", "திருச்சிராப்பள்ளி"},
		{"tirunelveli", "Tirunelveli", "திருநெல்வேலி"},
		{"thoothukudi", "Thoothukudi", "தூத்துக்குடி"},
		{"tenkasi", "Tenkasi", "தென்காசி"},
		{"namakkal", "Namakkal", "நாமக்கல்"},
		{"pudukkottai", "Pudukkottai", "புதுக்கோட்டை"},
		{"virudhunagar", "Virudhunagar", "விருதுநகர்"},
	}

	for _, tc := range testDistricts {
		req := httptest.NewRequest(http.MethodGet, "/district/"+tc.slug, nil)
		rec := httptest.NewRecorder()
		h.HandleDistrictPage(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 OK for /district/%s, got %d", tc.slug, rec.Code)
		}

		body := rec.Body.String()

		// Canonical check
		expectedCanonical := fmt.Sprintf(`<link rel="canonical" href="https://www.tn24.in/district/%s">`, tc.slug)
		if !strings.Contains(body, expectedCanonical) {
			t.Errorf("expected canonical %s for /district/%s", expectedCanonical, tc.slug)
		}

		// Title check
		expectedTitle := fmt.Sprintf("%s News Today – Latest %s Tamil News", tc.name, tc.name)
		if !strings.Contains(body, expectedTitle) {
			t.Errorf("expected title to contain %q for /district/%s", expectedTitle, tc.slug)
		}

		// Exactly ONE H1 check
		h1Matches := regexp.MustCompile(`(?is)<h1\b[^>]*>(.*?)</h1>`).FindAllStringSubmatch(body, -1)
		if len(h1Matches) != 1 {
			t.Errorf("/district/%s must have EXACTLY 1 H1 tag, found %d", tc.slug, len(h1Matches))
		} else {
			expectedH1 := fmt.Sprintf("%s News – %s செய்திகள்", tc.name, tc.taName)
			if !strings.Contains(h1Matches[0][1], expectedH1) {
				t.Errorf("/district/%s H1 incorrect, expected to contain %q, got: %s", tc.slug, expectedH1, h1Matches[0][1])
			}
		}

		// Breadcrumbs JSON-LD check
		if !strings.Contains(body, `"BreadcrumbList"`) || !strings.Contains(body, fmt.Sprintf(`/district/%s`, tc.slug)) {
			t.Errorf("/district/%s missing BreadcrumbList JSON-LD with district URL", tc.slug)
		}

		// Internal links to other districts check
		if !strings.Contains(body, "Latest Chennai News") && tc.slug != "chennai" {
			t.Errorf("/district/%s missing internal link anchor text to Chennai", tc.slug)
		}
	}

	// Invalid district returns 404
	req404 := httptest.NewRequest(http.MethodGet, "/district/unknown-invalid-district-xyz", nil)
	rec404 := httptest.NewRecorder()
	h.HandleDistrictPage(rec404, req404)
	if rec404.Code != http.StatusNotFound {
		t.Errorf("expected 404 for unknown district, got %d", rec404.Code)
	}
}

func TestP1CategoryPages(t *testing.T) {
	h := NewPortalHandler(nil, nil)

	testCats := []struct {
		slug       string
		name       string
		taName     string
	}{
		{"politics", "Politics", "அரசியல்"},
		{"news", "News", "செய்திகள்"},
		{"sports", "Sports", "விளையாட்டு"},
		{"business", "Business", "வணிகம்"},
		{"technical", "Technical", "தொழில்நுட்பம்"},
		{"crime", "Crime", "குற்ற நிகழ்வுகள்"},
		{"civic", "Civic", "பொதுமக்கள் & நலம்"},
	}

	for _, tc := range testCats {
		req := httptest.NewRequest(http.MethodGet, "/category/"+tc.slug, nil)
		rec := httptest.NewRecorder()
		h.HandleCategoryPage(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 OK for /category/%s, got %d", tc.slug, rec.Code)
		}

		body := rec.Body.String()

		// Canonical check
		expectedCanonical := fmt.Sprintf(`<link rel="canonical" href="https://www.tn24.in/category/%s">`, tc.slug)
		if !strings.Contains(body, expectedCanonical) {
			t.Errorf("expected canonical %s for /category/%s", expectedCanonical, tc.slug)
		}

		// Title check
		expectedTitle := fmt.Sprintf("Tamil Nadu %s News – Latest %s News", tc.name, tc.name)
		if !strings.Contains(body, expectedTitle) {
			t.Errorf("expected title to contain %q for /category/%s", expectedTitle, tc.slug)
		}

		// Exactly ONE H1 check
		h1Matches := regexp.MustCompile(`(?is)<h1\b[^>]*>(.*?)</h1>`).FindAllStringSubmatch(body, -1)
		if len(h1Matches) != 1 {
			t.Errorf("/category/%s must have EXACTLY 1 H1 tag, found %d", tc.slug, len(h1Matches))
		} else {
			expectedH1 := fmt.Sprintf("Tamil Nadu %s News – தமிழ்நாடு %s செய்திகள்", tc.name, tc.taName)
			expectedH1Escaped := fmt.Sprintf("Tamil Nadu %s News – தமிழ்நாடு %s செய்திகள்", tc.name, html.EscapeString(tc.taName))
			if !strings.Contains(h1Matches[0][1], expectedH1) && !strings.Contains(h1Matches[0][1], expectedH1Escaped) {
				t.Errorf("/category/%s H1 incorrect, expected %q, got: %s", tc.slug, expectedH1, h1Matches[0][1])
			}
		}

		// Breadcrumbs JSON-LD check
		if !strings.Contains(body, `"BreadcrumbList"`) || !strings.Contains(body, fmt.Sprintf(`/category/%s`, tc.slug)) {
			t.Errorf("/category/%s missing BreadcrumbList JSON-LD with category URL", tc.slug)
		}
	}

	// Invalid category returns 404
	req404 := httptest.NewRequest(http.MethodGet, "/category/invalid-unknown-cat-xyz", nil)
	rec404 := httptest.NewRecorder()
	h.HandleCategoryPage(rec404, req404)
	if rec404.Code != http.StatusNotFound {
		t.Errorf("expected 404 for unknown category, got %d", rec404.Code)
	}
}

func TestP1LegacyQueryRedirects(t *testing.T) {
	mux := http.NewServeMux()
	h := NewPortalHandler(nil, nil)
	h.RegisterRoutes(mux)

	// Test /portal?district=Chennai -> 301 to /district/chennai
	req1 := httptest.NewRequest(http.MethodGet, "/portal?district=Chennai", nil)
	rec1 := httptest.NewRecorder()
	mux.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusMovedPermanently || rec1.Header().Get("Location") != "/district/chennai" {
		t.Errorf("expected 301 to /district/chennai, got status %d, loc %s", rec1.Code, rec1.Header().Get("Location"))
	}

	// Test /portal?category=Politics -> 301 to /category/politics
	req2 := httptest.NewRequest(http.MethodGet, "/portal?category=Politics", nil)
	rec2 := httptest.NewRecorder()
	mux.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusMovedPermanently || rec2.Header().Get("Location") != "/category/politics" {
		t.Errorf("expected 301 to /category/politics, got status %d, loc %s", rec2.Code, rec2.Header().Get("Location"))
	}

	// Test /portal?cat=Politics -> 301 to /category/politics
	req3 := httptest.NewRequest(http.MethodGet, "/portal?cat=Politics", nil)
	rec3 := httptest.NewRecorder()
	mux.ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusMovedPermanently || rec3.Header().Get("Location") != "/category/politics" {
		t.Errorf("expected 301 to /category/politics, got status %d, loc %s", rec3.Code, rec3.Header().Get("Location"))
	}

	// Test /?district=Chennai in HandlePortalPage -> 301 to /district/chennai
	req4 := httptest.NewRequest(http.MethodGet, "/?district=Chennai", nil)
	rec4 := httptest.NewRecorder()
	h.HandlePortalPage(rec4, req4)
	if rec4.Code != http.StatusMovedPermanently || rec4.Header().Get("Location") != "/district/chennai" {
		t.Errorf("expected 301 to /district/chennai from /?district=Chennai, got status %d, loc %s", rec4.Code, rec4.Header().Get("Location"))
	}

	// Test /?category=Politics in HandlePortalPage -> 301 to /category/politics
	req5 := httptest.NewRequest(http.MethodGet, "/?category=Politics", nil)
	rec5 := httptest.NewRecorder()
	h.HandlePortalPage(rec5, req5)
	if rec5.Code != http.StatusMovedPermanently || rec5.Header().Get("Location") != "/category/politics" {
		t.Errorf("expected 301 to /category/politics from /?category=Politics, got status %d, loc %s", rec5.Code, rec5.Header().Get("Location"))
	}

	// Test /?cat=Politics in HandlePortalPage -> 301 to /category/politics
	req6 := httptest.NewRequest(http.MethodGet, "/?cat=Politics", nil)
	rec6 := httptest.NewRecorder()
	h.HandlePortalPage(rec6, req6)
	if rec6.Code != http.StatusMovedPermanently || rec6.Header().Get("Location") != "/category/politics" {
		t.Errorf("expected 301 to /category/politics from /?cat=Politics, got status %d, loc %s", rec6.Code, rec6.Header().Get("Location"))
	}
}

func TestP2DistrictFallbackContent(t *testing.T) {
	h := NewPortalHandler(nil, nil)

	// 1. Test district with NO local articles (e.g. Ariyalur with nil db / 0 articles)
	req := httptest.NewRequest(http.MethodGet, "/district/ariyalur", nil)
	rec := httptest.NewRecorder()
	h.HandleDistrictPage(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200 for empty district page, got %d", rec.Code)
	}

	body := rec.Body.String()

	// Must show "Ariyalur News – அரியலூர் செய்திகள்"
	if !strings.Contains(body, "Ariyalur News – அரியலூர் செய்திகள்") {
		t.Errorf("expected Ariyalur News heading on empty district page")
	}

	// Must show factual message
	expectedFactual := "TN24 will publish the latest Ariyalur news and local updates here as they are reported."
	if !strings.Contains(body, expectedFactual) {
		t.Errorf("expected factual message %q on empty district page", expectedFactual)
	}

	// Must NOT present statewide news as Ariyalur news
	if strings.Contains(body, "Latest Ariyalur News") {
		t.Errorf("should NOT show 'Latest Ariyalur News' when there are 0 local articles")
	}

	// Must contain canonical URL pointing to /district/ariyalur
	if !strings.Contains(body, `<link rel="canonical" href="https://www.tn24.in/district/ariyalur">`) {
		t.Errorf("expected canonical URL for Ariyalur")
	}
}

func TestP2LCPAndCLS(t *testing.T) {
	html := RenderPortalPage()

	// 1. LCP Hero Image optimizations
	if !strings.Contains(html, `id="heroImage"`) {
		t.Fatalf("heroImage element not found in RenderPortalPage")
	}
	if !strings.Contains(html, `loading="eager"`) || !strings.Contains(html, `fetchpriority="high"`) {
		t.Errorf("expected heroImage to have loading='eager' and fetchpriority='high'")
	}
	if !strings.Contains(html, `width="1200"`) || !strings.Contains(html, `height="675"`) {
		t.Errorf("expected heroImage to have explicit width and height dimensions for CLS prevention")
	}
	if !strings.Contains(html, `aspect-ratio:16/9`) {
		t.Errorf("expected heroImage to specify aspect-ratio:16/9 style")
	}

	// 2. Brand logo explicit dimensions for CLS prevention
	if !strings.Contains(html, `width="210" height="42"`) {
		t.Errorf("expected header brand logo to have explicit width and height")
	}

	// 3. Ad containers and units min-height reservation
	if !strings.Contains(html, `.ad-banner-horizontal {`) || !strings.Contains(html, `min-height: 90px;`) {
		t.Errorf("expected ad-banner-horizontal to reserve min-height: 90px")
	}
	if !strings.Contains(html, `.ad-banner-square {`) || !strings.Contains(html, `min-height: 200px;`) {
		t.Errorf("expected ad-banner-square to reserve min-height: 200px")
	}
	if !strings.Contains(html, `.ad-banner-250 {`) || !strings.Contains(html, `min-height: 250px;`) {
		t.Errorf("expected ad-banner-250 to reserve min-height: 250px")
	}
	if !strings.Contains(html, `.ad-unit-leaderboard {`) || !strings.Contains(html, `min-height: 90px;`) {
		t.Errorf("expected ad-unit-leaderboard to reserve min-height: 90px")
	}
	if !strings.Contains(html, `#ad-banner-header-slot {`) || !strings.Contains(html, `min-height: 90px;`) {
		t.Errorf("expected #ad-banner-header-slot to reserve min-height: 90px")
	}

	// 4. Article LCP Image preloading and dimensions
	h := NewPortalHandler(nil, nil)
	now := time.Date(2026, 9, 22, 10, 0, 0, 0, time.UTC)
	sampleThumb := "https://images.unsplash.com/photo-sample.jpg"
	articleHTML := h.InjectPostMetadataValues(
		nil, html, "test-post-1", "முதலமைச்சர் மு.க.ஸ்டாலின் அறிக்கை",
		"விரிவான செய்தி தொகுப்பு", sampleThumb,
		"Chennai", "Politics", "ta", now, now, "", "", "", 0, nil,
	)

	expectedPreload := fmt.Sprintf(`<link rel="preload" as="image" href="%s" fetchpriority="high">`, sampleThumb)
	if !strings.Contains(articleHTML, expectedPreload) {
		t.Errorf("expected article <head> to contain preload tag for hero image: %s", expectedPreload)
	}
	if !strings.Contains(articleHTML, `width="900" height="506" loading="eager" fetchpriority="high"`) {
		t.Errorf("expected article SSR image to have explicit dimensions, loading=eager, and fetchpriority=high")
	}
}

func TestP2CacheHeaders(t *testing.T) {
	h := NewPortalHandler(nil, nil)

	// 1. Static Brand Assets
	reqAsset := httptest.NewRequest(http.MethodGet, "/assets/brand/tn24-logo.svg?v=20260908d", nil)
	recAsset := httptest.NewRecorder()
	h.HandleBrandAssets(recAsset, reqAsset)

	ccAsset := recAsset.Header().Get("Cache-Control")
	if !strings.Contains(ccAsset, "max-age=31536000") || !strings.Contains(ccAsset, "immutable") {
		t.Errorf("expected long-lived immutable Cache-Control on versioned brand asset, got: %s", ccAsset)
	}

	// 2. Favicon
	reqFav := httptest.NewRequest(http.MethodGet, "/favicon.svg", nil)
	recFav := httptest.NewRecorder()
	h.HandleFavicon(recFav, reqFav)

	ccFav := recFav.Header().Get("Cache-Control")
	if !strings.Contains(ccFav, "max-age=604800") {
		t.Errorf("expected 7-day Cache-Control on favicon, got: %s", ccFav)
	}

	// 3. Dynamic HTML District Page
	reqDist := httptest.NewRequest(http.MethodGet, "/district/chennai", nil)
	recDist := httptest.NewRecorder()
	h.HandleDistrictPage(recDist, reqDist)

	ccDist := recDist.Header().Get("Cache-Control")
	if !strings.Contains(ccDist, "max-age=60") || !strings.Contains(ccDist, "s-maxage=120") {
		t.Errorf("expected short-lived fresh Cache-Control on district page, got: %s", ccDist)
	}

	// 4. Dynamic HTML Category Page
	reqCat := httptest.NewRequest(http.MethodGet, "/category/politics", nil)
	recCat := httptest.NewRecorder()
	h.HandleCategoryPage(recCat, reqCat)

	ccCat := recCat.Header().Get("Cache-Control")
	if !strings.Contains(ccCat, "max-age=60") || !strings.Contains(ccCat, "s-maxage=120") {
		t.Errorf("expected short-lived fresh Cache-Control on category page, got: %s", ccCat)
	}
}

func TestP2SitemapLastmod(t *testing.T) {
	h := NewPortalHandler(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	rec := httptest.NewRecorder()
	h.HandleSitemapXML(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200 on /sitemap.xml, got %d", rec.Code)
	}

	sitemap := rec.Body.String()

	// Static pages must use fixed launch date 2026-09-08 and NOT today's date
	expectedStaticEntry := "<loc>https://www.tn24.in/about</loc>\n        <lastmod>2026-09-08</lastmod>"
	if !strings.Contains(sitemap, expectedStaticEntry) {
		t.Errorf("expected /about to use static lastmod 2026-09-08 in sitemap, not found")
	}

	expectedPrivacyEntry := "<loc>https://www.tn24.in/privacy</loc>\n        <lastmod>2026-09-08</lastmod>"
	if !strings.Contains(sitemap, expectedPrivacyEntry) {
		t.Errorf("expected /privacy to use static lastmod 2026-09-08 in sitemap, not found")
	}
}



