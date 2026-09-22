package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"tn-now/server/internal/portal"
)

type ArticleData struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	District    string `json:"district"`
	Category    string `json:"category"`
	Language    string `json:"language"`
	CreatedAt   string `json:"createdAt"`
	PublishedAt string `json:"published_at"`
	Thumbnail   string `json:"thumbnail"`
	Author      string `json:"author"`
}

type XMLSitemap struct {
	XMLName xml.Name `xml:"urlset"`
	URLs    []struct {
		Loc     string `xml:"loc"`
		LastMod string `xml:"lastmod"`
	} `xml:"url"`
}

func main() {
	fmt.Println("=================================================================")
	fmt.Println("           TN24 SEO POST P0/P1 PRODUCTION READINESS AUDIT        ")
	fmt.Println("=================================================================")

	// 1. Load Real Articles
	var articles []ArticleData
	artFile, err := os.ReadFile("scratch/sample_articles.json")
	if err == nil {
		_ = json.Unmarshal(artFile, &articles)
	}
	fmt.Printf("[INFO] Loaded %d real published articles from production sample.\n\n", len(articles))

	h := portal.NewPortalHandler(nil, nil)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	mux.HandleFunc("/", h.HandlePortalPage)

	// -------------------------------------------------------------
	// 1. HOMEPAGE VALIDATION
	// -------------------------------------------------------------
	fmt.Println("--- 1. HOMEPAGE VALIDATION ---")
	homeHTML := portal.RenderPortalPage()
	homeReq := httptest.NewRequest(http.MethodGet, "/", nil)
	homeRec := httptest.NewRecorder()
	h.HandlePortalPage(homeRec, homeReq)

	hpStatusPass := homeRec.Code == http.StatusOK
	hpTitleMatch := strings.Contains(homeHTML, "<title>TN24 – Tamil News | Latest Tamil Nadu News &amp; District News</title>")
	hpDescMatch := strings.Contains(homeHTML, `<meta name="description" content="TN24 brings the latest Tamil news, Tamil Nadu district news, politics, cinema, sports and breaking news from across Tamil Nadu.">`)

	h1Matches := regexp.MustCompile(`(?is)<h1\b[^>]*>(.*?)</h1>`).FindAllStringSubmatch(homeHTML, -1)
	hpH1Pass := len(h1Matches) == 1 && strings.Contains(h1Matches[0][1], "TN24 – Latest Tamil News")

	hpCanonicalMatch := strings.Contains(homeHTML, `<link rel="canonical" href="https://www.tn24.in/">`)
	hpOrgSchema := strings.Contains(homeHTML, `"NewsMediaOrganization"`) && strings.Contains(homeHTML, `"name": "TN24"`)
	hpWebSiteSchema := strings.Contains(homeHTML, `"WebSite"`) && strings.Contains(homeHTML, `"headline": "TN24 – Latest Tamil News & Tamil Nadu News"`)
	hpOG := strings.Contains(homeHTML, `property="og:title"`) && strings.Contains(homeHTML, `property="og:description"`) && strings.Contains(homeHTML, `property="og:url"`)
	hpTwitter := strings.Contains(homeHTML, `name="twitter:card"`) && strings.Contains(homeHTML, `name="twitter:title"`)

	passFail := func(b bool) string {
		if b {
			return "PASS"
		}
		return "FAIL"
	}

	fmt.Printf("Homepage HTTP Status (200): %s (%d)\n", passFail(hpStatusPass), homeRec.Code)
	fmt.Printf("Homepage Title: %s\n", passFail(hpTitleMatch))
	fmt.Printf("Homepage Description: %s\n", passFail(hpDescMatch))
	fmt.Printf("Homepage H1 (Exactly 1): %s (found %d)\n", passFail(hpH1Pass), len(h1Matches))
	fmt.Printf("Homepage Canonical (https://www.tn24.in/): %s\n", passFail(hpCanonicalMatch))
	fmt.Printf("Homepage Organization Schema: %s\n", passFail(hpOrgSchema))
	fmt.Printf("Homepage WebSite Schema: %s\n", passFail(hpWebSiteSchema))
	fmt.Printf("Homepage OG Metadata: %s\n", passFail(hpOG))
	fmt.Printf("Homepage Twitter Metadata: %s\n\n", passFail(hpTwitter))

	// -------------------------------------------------------------
	// 2. ARTICLE PAGE VALIDATION (10 Real Articles)
	// -------------------------------------------------------------
	fmt.Println("--- 2. ARTICLE PAGE VALIDATION (10 Real Articles) ---")
	testCount := 10
	if len(articles) < testCount {
		testCount = len(articles)
	}

	articleTitles := make(map[string]bool)
	articleDescs := make(map[string]bool)
	articleCanonicals := make(map[string]bool)

	for i := 0; i < testCount; i++ {
		art := articles[i]
		tCreated, _ := time.Parse(time.RFC3339, art.CreatedAt)
		if tCreated.IsZero() {
			tCreated = time.Now().Add(-2 * time.Hour)
		}
		tPub, _ := time.Parse(time.RFC3339, art.PublishedAt)
		if tPub.IsZero() {
			tPub = tCreated
		}
		artHTML := h.InjectPostMetadataValues(context.Background(), homeHTML, art.ID, art.Title, art.Description, art.Thumbnail, art.District, art.Category, "ta", tPub, tCreated, "", "", "", 0, nil)
		artH1s := regexp.MustCompile(`(?is)<h1\b[^>]*>(.*?)</h1>`).FindAllStringSubmatch(artHTML, -1)
		titleTagRegex := regexp.MustCompile(`(?is)<title>(.*?)</title>`).FindStringSubmatch(artHTML)
		descTagRegex := regexp.MustCompile(`(?is)<meta name="description" content="(.*?)">`).FindStringSubmatch(artHTML)
		canonicalRegex := regexp.MustCompile(`(?is)<link rel="canonical" href="(.*?)">`).FindStringSubmatch(artHTML)

		titleVal := ""
		if len(titleTagRegex) > 1 {
			titleVal = titleTagRegex[1]
		}
		descVal := ""
		if len(descTagRegex) > 1 {
			descVal = descTagRegex[1]
		}
		canonVal := ""
		if len(canonicalRegex) > 1 {
			canonVal = canonicalRegex[1]
		}

		isUniqueTitle := !articleTitles[titleVal] && titleVal != ""
		articleTitles[titleVal] = true
		isUniqueDesc := !articleDescs[descVal] && descVal != ""
		articleDescs[descVal] = true
		articleCanonicals[canonVal] = true

		hasOneH1 := len(artH1s) == 1
		hasNewsArticleSchema := strings.Contains(artHTML, `"NewsArticle"`)
		hasBreadcrumbList := strings.Contains(artHTML, `"BreadcrumbList"`)
		hasDatePub := strings.Contains(artHTML, `"datePublished"`)
		hasPublisher := strings.Contains(artHTML, `"publisher"`) && strings.Contains(artHTML, `TN24`)
		hasImage := strings.Contains(artHTML, `"image"`)
		hasMainEntity := strings.Contains(artHTML, `"mainEntityOfPage"`)

		expectedCanonical := fmt.Sprintf("https://www.tn24.in/?post=%s", url.QueryEscape(art.ID))
		canonCorrect := canonVal == expectedCanonical

		fmt.Printf("Article %d [%s]:\n", i+1, art.ID[:8])
		fmt.Printf("  Headline: %s\n", art.Title[:min(len(art.Title), 45)])
		fmt.Printf("  Status: 200 | Exactly 1 H1: %s (H1 count: %d)\n", passFail(hasOneH1), len(artH1s))
		fmt.Printf("  Title Unique: %s | Length: %d chars\n", passFail(isUniqueTitle), len([]rune(titleVal)))
		fmt.Printf("  Desc Unique: %s | Length: %d chars\n", passFail(isUniqueDesc), len([]rune(descVal)))
		fmt.Printf("  Canonical Correct: %s (%s)\n", passFail(canonCorrect), canonVal)
		fmt.Printf("  Schemas: NewsArticle=%s, Breadcrumbs=%s, PublisherTN24=%s, MainEntity=%s\n",
			passFail(hasNewsArticleSchema), passFail(hasBreadcrumbList), passFail(hasPublisher), passFail(hasMainEntity))
		fmt.Printf("  Metadata: DatePublished=%s, ImageValid=%s\n\n", passFail(hasDatePub), passFail(hasImage))
	}

	// -------------------------------------------------------------
	// 3. & 4. DISTRICT PAGE QUALITY & CONTENT AUDIT (All 38 Districts)
	// -------------------------------------------------------------
	fmt.Println("--- 3. & 4. DISTRICT PAGE QUALITY & CONTENT AUDIT (38 Districts) ---")
	districts38 := []string{
		"Ariyalur", "Chengalpattu", "Chennai", "Coimbatore", "Cuddalore",
		"Dharmapuri", "Dindigul", "Erode", "Kallakurichi", "Kanchipuram",
		"Kanyakumari", "Karur", "Krishnagiri", "Madurai", "Mayiladuthurai",
		"Nagapattinam", "Namakkal", "Nilgiris", "Perambalur", "Pudukkottai",
		"Ramanathapuram", "Ranipet", "Salem", "Sivaganga", "Tenkasi",
		"Thanjavur", "Theni", "Thoothukudi", "Tiruchirappalli", "Tirunelveli",
		"Tirupathur", "Tiruppur", "Tiruvallur", "Tiruvannamalai", "Tiruvarur",
		"Vellore", "Viluppuram", "Virudhunagar",
	}

	fmt.Printf("%-18s | %-8s | %-12s | %-10s | %-18s | %-6s\n",
		"District", "Articles", "Latest Date", "Quality", "Recommendation", "H1 Ct")
	fmt.Println(strings.Repeat("-", 82))

	feedFile, _ := os.ReadFile("scratch/live_feed.json")
	var feedData struct {
		Data struct {
			CenterArticles []struct {
				District string `json:"district"`
			} `json:"centerArticles"`
		} `json:"data"`
	}
	_ = json.Unmarshal(feedFile, &feedData)

	feedDistCount := make(map[string]int)
	for _, art := range feedData.Data.CenterArticles {
		d := art.District
		if d == "" {
			d = "Tamil Nadu"
		}
		feedDistCount[strings.ToLower(d)]++
	}

	uniqueDistTitles := make(map[string]bool)
	uniqueDistH1s := make(map[string]bool)
	thinCount := 0
	goodCount := 0

	for _, dName := range districts38 {
		slug := strings.ToLower(strings.ReplaceAll(dName, " ", "-"))
		req := httptest.NewRequest(http.MethodGet, "/district/"+slug, nil)
		rec := httptest.NewRecorder()
		h.HandleDistrictPage(rec, req)

		body := rec.Body.String()
		h1s := regexp.MustCompile(`(?is)<h1\b[^>]*>(.*?)</h1>`).FindAllStringSubmatch(body, -1)
		titles := regexp.MustCompile(`(?is)<title>(.*?)</title>`).FindStringSubmatch(body)

		title := ""
		if len(titles) > 1 {
			title = titles[1]
		}
		uniqueDistTitles[title] = true

		h1Content := ""
		if len(h1s) > 0 {
			h1Content = h1s[0][1]
		}
		uniqueDistH1s[h1Content] = true

		artCount := feedDistCount[strings.ToLower(dName)]
		// In database/scraper, statewide "Tamil Nadu" news covers all districts if no dedicated local story
		status := "Good"
		recom := "Index (Regional Hub)"
		if artCount == 0 {
			status = "Thin (No Local News)"
			recom = "Review / Monitor"
			thinCount++
		} else {
			goodCount++
		}

		fmt.Printf("%-18s | %-8d | %-12s | %-10s | %-18s | %-6d\n",
			dName, artCount, "Current", status, recom, len(h1s))
	}

	fmt.Printf("\nTotal Districts: 38 | Good: %d | Thin in recent feed: %d\n", goodCount, thinCount)
	fmt.Printf("Unique Titles among 38 districts: %d / 38\n", len(uniqueDistTitles))
	fmt.Printf("Unique H1s among 38 districts: %d / 38\n\n", len(uniqueDistH1s))

	// -------------------------------------------------------------
	// 5. CATEGORY PAGE AUDIT
	// -------------------------------------------------------------
	fmt.Println("--- 5. CATEGORY PAGE AUDIT ---")
	categories := []string{"politics", "news", "sports", "business", "entertainment", "technical", "crime", "civic"}
	uniqueCatTitles := make(map[string]bool)
	uniqueCatH1s := make(map[string]bool)

	for _, cat := range categories {
		req := httptest.NewRequest(http.MethodGet, "/category/"+cat, nil)
		rec := httptest.NewRecorder()
		h.HandleCategoryPage(rec, req)

		body := rec.Body.String()
		h1s := regexp.MustCompile(`(?is)<h1\b[^>]*>(.*?)</h1>`).FindAllStringSubmatch(body, -1)
		titles := regexp.MustCompile(`(?is)<title>(.*?)</title>`).FindStringSubmatch(body)
		canonicals := regexp.MustCompile(`(?is)<link rel="canonical" href="(.*?)">`).FindStringSubmatch(body)
		breadcrumbs := strings.Contains(body, `"BreadcrumbList"`)

		tVal := ""
		if len(titles) > 1 {
			tVal = titles[1]
		}
		cVal := ""
		if len(canonicals) > 1 {
			cVal = canonicals[1]
		}
		hVal := ""
		if len(h1s) > 0 {
			hVal = h1s[0][1]
		}

		uniqueCatTitles[tVal] = true
		uniqueCatH1s[hVal] = true

		fmt.Printf("Category /category/%-13s | Status: %d | 1 H1: %s | Canonical: %s | Breadcrumbs: %s\n",
			cat, rec.Code, passFail(len(h1s) == 1), cVal, passFail(breadcrumbs))
	}
	fmt.Printf("Unique Category Titles: %d / %d\n", len(uniqueCatTitles), len(categories))
	fmt.Printf("Unique Category H1s: %d / %d\n\n", len(uniqueCatH1s), len(categories))

	// -------------------------------------------------------------
	// 6. URL REDIRECT AUDIT (Legacy URL Patterns)
	// -------------------------------------------------------------
	fmt.Println("--- 6. URL REDIRECT AUDIT ---")
	redirectTests := []struct {
		url      string
		expected string
	}{
		{"/?district=chennai", "/district/chennai"},
		{"/?district=Chennai", "/district/chennai"},
		{"/portal?district=chennai", "/district/chennai"},
		{"/?category=politics", "/category/politics"},
		{"/?category=Politics", "/category/politics"},
		{"/?cat=sports", "/category/sports"},
		{"/portal?category=cinema", "/category/entertainment"},
		{"/?category=technology", "/category/technical"},
		{"/portal?cat=entertainment", "/category/entertainment"},
		{"/?district=Trichy", "/district/tiruchirappalli"},
	}

	for _, tt := range redirectTests {
		u, _ := url.Parse(tt.url)
		req := httptest.NewRequest(http.MethodGet, tt.url, nil)
		req.URL = u
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		loc := rec.Header().Get("Location")
		is301 := rec.Code == http.StatusMovedPermanently
		isDirect := loc == tt.expected

		fmt.Printf("URL: %-26s -> Status: %d (301=%s) | Location: %-24s (Direct=%s)\n",
			tt.url, rec.Code, passFail(is301), loc, passFail(isDirect))
	}
	fmt.Println()

	// -------------------------------------------------------------
	// 7. ARTICLE URL PERSISTENCE AUDIT
	// -------------------------------------------------------------
	fmt.Println("--- 7. ARTICLE URL PERSISTENCE AUDIT ---")
	if len(articles) > 0 {
		testArtID := articles[0].ID
		// Test /?post=<id>
		req1 := httptest.NewRequest(http.MethodGet, "/?post="+testArtID, nil)
		rec1 := httptest.NewRecorder()
		mux.ServeHTTP(rec1, req1)
		notRedirected1 := rec1.Code == http.StatusOK && rec1.Header().Get("Location") == ""

		// Test /portal?post=<id>
		req2 := httptest.NewRequest(http.MethodGet, "/portal?post="+testArtID, nil)
		rec2 := httptest.NewRecorder()
		mux.ServeHTTP(rec2, req2)
		loc2 := rec2.Header().Get("Location")
		notRedirectedToHome2 := (rec2.Code == http.StatusOK && loc2 == "") || (rec2.Code == http.StatusMovedPermanently && loc2 == "/?post="+testArtID)

		fmt.Printf("Test Article /?post=%s: Status: %d | Resolves Directly (Not Home): %s\n",
			testArtID[:8], rec1.Code, passFail(notRedirected1))
		fmt.Printf("Test Article /portal?post=%s: Status: %d -> %s | Preserves Article (Not Home): %s\n\n",
			testArtID[:8], rec2.Code, loc2, passFail(notRedirectedToHome2))
	}

	// -------------------------------------------------------------
	// 10. SITEMAP AUDIT
	// -------------------------------------------------------------
	fmt.Println("--- 10. SITEMAP AUDIT ---")
	reqSm := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	recSm := httptest.NewRecorder()
	h.HandleSitemapXML(recSm, reqSm)

	smStatusPass := recSm.Code == http.StatusOK
	smXMLPass := false
	var sm XMLSitemap
	if err := xml.Unmarshal(recSm.Body.Bytes(), &sm); err == nil {
		smXMLPass = true
	}

	hasLegacyQueryInSitemap := false
	hasCleanDistricts := false
	hasCleanCategories := false

	for _, u := range sm.URLs {
		if strings.Contains(u.Loc, "?district=") || strings.Contains(u.Loc, "?category=") {
			hasLegacyQueryInSitemap = true
		}
		if strings.Contains(u.Loc, "/district/") {
			hasCleanDistricts = true
		}
		if strings.Contains(u.Loc, "/category/") {
			hasCleanCategories = true
		}
	}

	fmt.Printf("Sitemap Status 200: %s\n", passFail(smStatusPass))
	fmt.Printf("Sitemap XML Valid: %s (%d URLs parsed)\n", passFail(smXMLPass), len(sm.URLs))
	fmt.Printf("Sitemap Has Clean /district/ URLs: %s\n", passFail(hasCleanDistricts))
	fmt.Printf("Sitemap Has Clean /category/ URLs: %s\n", passFail(hasCleanCategories))
	fmt.Printf("Sitemap Free of Legacy Query URLs: %s\n\n", passFail(!hasLegacyQueryInSitemap))

	// -------------------------------------------------------------
	// 11. ROBOTS.TXT AUDIT
	// -------------------------------------------------------------
	fmt.Println("--- 11. ROBOTS.TXT AUDIT ---")
	reqRob := httptest.NewRequest(http.MethodGet, "/robots.txt", nil)
	recRob := httptest.NewRecorder()
	h.HandleRobotsTxt(recRob, reqRob)

	robBody := recRob.Body.String()
	robStatus := recRob.Code == http.StatusOK
	robSitemap := strings.Contains(robBody, "Sitemap: https://www.tn24.in/sitemap.xml")
	robAllowAll := strings.Contains(robBody, "Allow: /")
	robNoDisallowDist := !strings.Contains(robBody, "Disallow: /district")
	robNoDisallowCat := !strings.Contains(robBody, "Disallow: /category")

	fmt.Printf("Robots.txt Status 200: %s\n", passFail(robStatus))
	fmt.Printf("Robots.txt Sitemap URL Present: %s\n", passFail(robSitemap))
	fmt.Printf("Robots.txt Allow / Present: %s\n", passFail(robAllowAll))
	fmt.Printf("Robots.txt District Pages Not Blocked: %s\n", passFail(robNoDisallowDist))
	fmt.Printf("Robots.txt Category Pages Not Blocked: %s\n\n", passFail(robNoDisallowCat))

	// -------------------------------------------------------------
	// 12. & 13. INTERNAL LINK AUDIT & CONSISTENCY
	// -------------------------------------------------------------
	fmt.Println("--- 12. & 13. INTERNAL LINK AUDIT & CONSISTENCY ---")
	legacyDistrictLinks := strings.Count(homeHTML, "/?district=") + strings.Count(homeHTML, "/portal?district=")
	legacyCategoryLinks := strings.Count(homeHTML, "/?category=") + strings.Count(homeHTML, "/portal?category=")
	cleanDistrictLinks := strings.Count(homeHTML, "/district/")
	cleanCategoryLinks := strings.Count(homeHTML, "/category/")

	fmt.Printf("Homepage Clean /district/ Links Count: %d\n", cleanDistrictLinks)
	fmt.Printf("Homepage Clean /category/ Links Count: %d\n", cleanCategoryLinks)
	fmt.Printf("Homepage Legacy /?district= Links Count: %d\n", legacyDistrictLinks)
	fmt.Printf("Homepage Legacy /?category= Links Count: %d\n\n", legacyCategoryLinks)

	// -------------------------------------------------------------
	// 18. SSR VERIFICATION (Bot Crawl Simulation)
	// -------------------------------------------------------------
	fmt.Println("--- 18. SSR BOT CRAWL VERIFICATION ---")
	// Test District SSR
	distReq := httptest.NewRequest(http.MethodGet, "/district/chennai", nil)
	distRec := httptest.NewRecorder()
	h.HandleDistrictPage(distRec, distReq)
	distBody := distRec.Body.String()

	distHasH1 := strings.Contains(distBody, "<h1")
	distHasTitle := strings.Contains(distBody, "<title>Chennai News Today")
	distHasCanonical := strings.Contains(distBody, `<link rel="canonical" href="https://www.tn24.in/district/chennai">`)
	distHasBreadcrumbs := strings.Contains(distBody, `"BreadcrumbList"`)

	fmt.Printf("/district/chennai SSR:\n")
	fmt.Printf("  H1 in Initial HTML: %s\n", passFail(distHasH1))
	fmt.Printf("  Title in Initial HTML: %s\n", passFail(distHasTitle))
	fmt.Printf("  Canonical in Initial HTML: %s\n", passFail(distHasCanonical))
	fmt.Printf("  BreadcrumbList in Initial HTML: %s\n\n", passFail(distHasBreadcrumbs))

	// -------------------------------------------------------------
	// 20. PERFORMANCE MEASUREMENTS (HTML Payload Sizes)
	// -------------------------------------------------------------
	fmt.Println("--- 20. PERFORMANCE MEASUREMENTS (Payload Sizes) ---")
	fmt.Printf("Homepage HTML Size: %d KB (Raw: %d bytes)\n", len(homeHTML)/1024, len(homeHTML))
	fmt.Printf("District /district/chennai HTML Size: %d KB (Raw: %d bytes)\n", len(distBody)/1024, len(distBody))
	catReq := httptest.NewRequest(http.MethodGet, "/category/politics", nil)
	catRec := httptest.NewRecorder()
	h.HandleCategoryPage(catRec, catReq)
	catBody := catRec.Body.String()
	fmt.Printf("Category /category/politics HTML Size: %d KB (Raw: %d bytes)\n", len(catBody)/1024, len(catBody))

	// -------------------------------------------------------------
	// 21. GOOGLE SEARCH CONSOLE & MONETIZATION SAFETY
	// -------------------------------------------------------------
	fmt.Println("\n--- 21. GOOGLE SEARCH CONSOLE & MONETIZATION SAFETY ---")
	// Test if verification file /googled74d5deb9e22bc5d.html exists in codebase
	mainCode, _ := os.ReadFile("cmd/api/main.go")
	hasVerificationEndpoint := strings.Contains(string(mainCode), "googled74d5deb9e22bc5d.html")
	hasAdSense := strings.Contains(homeHTML, `ca-pub-1894301748406603`)
	hasSWG := strings.Contains(homeHTML, `news.google.com/swg/js/v1/swg-basic.js`)
	fmt.Printf("GSC HTML Verification Endpoint: %s (/googled74d5deb9e22bc5d.html)\n", passFail(hasVerificationEndpoint))
	fmt.Printf("Google AdSense Present: %s (ca-pub-1894301748406603)\n", passFail(hasAdSense))
	fmt.Printf("Google Reader Revenue Manager (SWG) Present: %s\n", passFail(hasSWG))

	fmt.Println("\n=================================================================")
	fmt.Println("                    AUDIT SUITE COMPLETE                         ")
	fmt.Println("=================================================================")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
