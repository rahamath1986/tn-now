package portal

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"tn-now/server/internal/db"
	"tn-now/server/internal/scraper"
)

type PortalHandler struct {
	q    *db.Queries
	conn *pgxpool.Pool
}

func NewPortalHandler(q *db.Queries, conn *pgxpool.Pool) *PortalHandler {
	return &PortalHandler{q: q, conn: conn}
}

type ArticleItem struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ContentType string    `json:"contentType"`
	SourceURL   string    `json:"sourceUrl"`
	District    string    `json:"district"`
	Category    string    `json:"category"`
	Language    string    `json:"language"`
	Status      string    `json:"status"`
	VideoID     string    `json:"videoId"`
	VideoURL    string    `json:"videoUrl"`
	VideoType   string    `json:"videoType"`
	Thumbnail   string    `json:"thumbnail"`
	CreatedAt   time.Time `json:"createdAt"`
	IsViral     bool      `json:"isViral"`
	Author      string    `json:"author"`
}

type AdBannerSlot struct {
	Type     string       `json:"type"` // "brand" or "article"
	Article  *ArticleItem `json:"article,omitempty"`
	AssetURL string       `json:"assetUrl,omitempty"`
}

type EventItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	District    string `json:"district"`
	Venue       string `json:"venue"`
	Badge       string `json:"badge"`
	BadgeColor  string `json:"badgeColor"`
	DateDisplay string `json:"dateDisplay"`
	Countdown   string `json:"countdown"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type PortalFeedResponse struct {
	Hero           *ArticleItem             `json:"hero"`
	HeroTeasers    []ArticleItem            `json:"heroTeasers"`
	LeftFeed       []ArticleItem            `json:"leftFeed"`
	PressReleases  []ArticleItem            `json:"pressReleases"`
	CenterArticles []ArticleItem            `json:"centerArticles"`
	MostRead       []ArticleItem            `json:"mostRead"`
	Districts      []string                 `json:"districts"`
	Categories     []string                 `json:"categories"`
	TotalCount     int                      `json:"totalCount"`
	ViralCount     int                      `json:"viralCount"`
	Banners        map[string]*AdBannerSlot `json:"banners,omitempty"`
	Events         []EventItem              `json:"events,omitempty"`
}

func (h *PortalHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/portal", h.HandlePortalPage)
	mux.HandleFunc("/api/portal/feed", h.HandlePortalFeed)
	mux.HandleFunc("/api/portal/post", h.HandleGetSinglePost)
	mux.HandleFunc("/portal/assets/brand/", h.HandleBrandAssets)
	mux.HandleFunc("/portal/assets/", h.HandleBrandAssets)
	mux.HandleFunc("/assets/", h.HandleBrandAssets)
	mux.HandleFunc("/admin/assets/", h.HandleBrandAssets)

	// SEO Crawler & Indexing Endpoints
	mux.HandleFunc("/robots.txt", h.HandleRobotsTxt)
	mux.HandleFunc("/sitemap.xml", h.HandleSitemapXML)
	mux.HandleFunc("/sitemap-news.xml", h.HandleNewsSitemapXML)

	// AdSense-required static pages
	mux.HandleFunc("/privacy", h.HandlePrivacyPolicy)
	mux.HandleFunc("/about", h.HandleAboutUs)
	mux.HandleFunc("/contact", h.HandleContactUs)

	// ads.txt — required by Google AdSense to authorize ad sellers
	mux.HandleFunc("/ads.txt", h.HandleAdsTxt)
}

func (h *PortalHandler) HandlePortalPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")

	postID := strings.TrimSpace(r.URL.Query().Get("post"))
	pageHTML := RenderPortalPage()

	if postID != "" && h.conn != nil {
		pageHTML = h.injectPostMetadata(r.Context(), pageHTML, postID, r)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(pageHTML))
}

func (h *PortalHandler) HandleGetSinglePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Method not allowed"})
		return
	}

	postID := strings.TrimSpace(r.URL.Query().Get("id"))
	if postID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Missing article id"})
		return
	}

	if h.conn == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Database not connected"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var it ArticleItem
	var cid pgtype.UUID
	err := h.conn.QueryRow(ctx, `
		SELECT c.id, c.title, COALESCE(c.description, ''), c.content_type, COALESCE(c.source_url, ''),
		       COALESCE(d.name, 'Tamil Nadu'), COALESCE(cat.name, 'News'), c.status,
		       COALESCE(vl.external_video_id, ''), COALESCE(vl.canonical_url, ''), COALESCE(vl.platform, ''),
		       COALESCE(vl.thumbnail_url, s.single_photo_url, (p.photo_urls)[1], ''),
		       COALESCE(c.published_at, c.created_at), COALESCE(c.is_viral, false)
		FROM content c
		LEFT JOIN districts d ON c.district_id = d.id
		LEFT JOIN categories cat ON c.category_id = cat.id
		LEFT JOIN video_links vl ON c.id = vl.content_id
		LEFT JOIN stories s ON c.id = s.content_id
		LEFT JOIN photos p ON c.id = p.content_id
		WHERE c.id::text = $1
		LIMIT 1
	`, postID).Scan(&cid, &it.Title, &it.Description, &it.ContentType, &it.SourceURL,
		&it.District, &it.Category, &it.Status, &it.VideoID, &it.VideoURL, &it.VideoType, &it.Thumbnail, &it.CreatedAt, &it.IsViral)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Article not found"})
		return
	}

	it.ID = fmtUUID(cid)
	it.Title = scraper.CleanHTML(it.Title)
	it.Description = sanitizeEditorialDescription(scraper.CleanHTML(it.Description), it.Title, it.District)
	it.Language = scraper.DetectLanguage(it.Title + " " + it.Description)
	if strings.TrimSpace(it.Thumbnail) == "" {
		it.Thumbnail = scraper.GetFallbackImageWithPerson(it.Title, it.District, it.Category)
	}
	it.Author = "TN24 செய்திக் குழு"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    it,
	})
}

func (h *PortalHandler) injectPostMetadata(ctx context.Context, baseHTML string, postID string, r *http.Request) string {
	var title, desc, thumb, district, category, lang string
	var pubAt, createdAt time.Time
	err := h.conn.QueryRow(ctx, `
		SELECT c.title, COALESCE(c.description, ''),
		       COALESCE(vl.thumbnail_url, s.single_photo_url, (p.photo_urls)[1], ''),
		       COALESCE(d.name, 'Tamil Nadu'), COALESCE(cat.name, 'News'),
		       COALESCE(c.language, 'ta'),
		       COALESCE(c.published_at, c.created_at), c.created_at
		FROM content c
		LEFT JOIN video_links vl ON c.id = vl.content_id
		LEFT JOIN stories s ON c.id = s.content_id
		LEFT JOIN photos p ON c.id = p.content_id
		LEFT JOIN districts d ON c.district_id = d.id
		LEFT JOIN categories cat ON c.category_id = cat.id
		WHERE c.id::text = $1
		LIMIT 1
	`, postID).Scan(&title, &desc, &thumb, &district, &category, &lang, &pubAt, &createdAt)

	if err != nil || strings.TrimSpace(title) == "" {
		return baseHTML
	}

	cleanTitle := scraper.CleanHTML(title)
	cleanTitleEscaped := html.EscapeString(cleanTitle)

	cleanDesc := scraper.CleanHTML(desc)
	cleanDesc = strings.Join(strings.Fields(cleanDesc), " ")
	runes := []rune(cleanDesc)
	if len(runes) > 180 {
		cleanDesc = string(runes[:180]) + "..."
	}
	cleanDescEscaped := html.EscapeString(cleanDesc)

	scheme := "https"
	host := r.Host
	if host == "" {
		host = "www.tn24.in"
	}
	fullPostURL := fmt.Sprintf("%s://%s/portal?post=%s", scheme, host, url.QueryEscape(postID))

	if strings.TrimSpace(thumb) == "" {
		thumb = scraper.GetFallbackImageWithPerson(cleanTitle, district, category)
	}
	if strings.HasPrefix(thumb, "/") {
		thumb = fmt.Sprintf("%s://%s%s", scheme, host, thumb)
	}

	// 1. Replace <title>
	reTitle := regexp.MustCompile(`(?i)<title>.*?</title>`)
	baseHTML = reTitle.ReplaceAllString(baseHTML, fmt.Sprintf("<title>%s &mdash; TN24 | Tamil Nadu News | தமிழ் செய்திகள்</title>", cleanTitleEscaped))

	// 2. Replace meta description
	reDesc := regexp.MustCompile(`(?i)<meta name="description" content=".*?">`)
	baseHTML = reDesc.ReplaceAllString(baseHTML, fmt.Sprintf(`<meta name="description" content="%s — TN24 (tn 24) Tamil Nadu news live 24x7. தமிழ் செய்திகள் உடனுக்குடன்.">`, cleanDescEscaped))

	// 3. Replace og:title
	reOgTitle := regexp.MustCompile(`(?i)<meta property="og:title" content=".*?">`)
	baseHTML = reOgTitle.ReplaceAllString(baseHTML, fmt.Sprintf(`<meta property="og:title" content="%s | TN24 — Tamil Nadu News | தமிழ் செய்திகள்">`, cleanTitleEscaped))

	// 4. Replace og:description
	reOgDesc := regexp.MustCompile(`(?i)<meta property="og:description" content=".*?">`)
	baseHTML = reOgDesc.ReplaceAllString(baseHTML, fmt.Sprintf(`<meta property="og:description" content="%s">`, cleanDescEscaped))

	// 5. Replace og:url
	reOgURL := regexp.MustCompile(`(?i)<meta property="og:url" content=".*?">`)
	baseHTML = reOgURL.ReplaceAllString(baseHTML, fmt.Sprintf(`<meta property="og:url" content="%s">`, fullPostURL))

	// 6. Replace og:image
	reOgImg := regexp.MustCompile(`(?i)<meta property="og:image" content=".*?">`)
	baseHTML = reOgImg.ReplaceAllString(baseHTML, fmt.Sprintf(`<meta property="og:image" content="%s">`, html.EscapeString(thumb)))

	// 7. Replace twitter tags
	reTwTitle := regexp.MustCompile(`(?i)<meta name="twitter:title" content=".*?">`)
	baseHTML = reTwTitle.ReplaceAllString(baseHTML, fmt.Sprintf(`<meta name="twitter:title" content="%s — TN24 Tamil Nadu News">`, cleanTitleEscaped))

	reTwDesc := regexp.MustCompile(`(?i)<meta name="twitter:description" content=".*?">`)
	baseHTML = reTwDesc.ReplaceAllString(baseHTML, fmt.Sprintf(`<meta name="twitter:description" content="%s">`, cleanDescEscaped))

	reTwImg := regexp.MustCompile(`(?i)<meta name="twitter:image" content=".*?">`)
	baseHTML = reTwImg.ReplaceAllString(baseHTML, fmt.Sprintf(`<meta name="twitter:image" content="%s">`, html.EscapeString(thumb)))

	// 8. Inject NewsArticle structured data & window.INITIAL_POST_ID
	articleSchemaJSON, _ := json.Marshal(map[string]interface{}{
		"@context": "https://schema.org",
		"@type":    "NewsArticle",
		"mainEntityOfPage": map[string]string{
			"@type": "WebPage",
			"@id":   fullPostURL,
		},
		"headline":      cleanTitle,
		"description":   cleanDesc,
		"image":         []string{thumb},
		"datePublished": pubAt.Format(time.RFC3339),
		"dateModified":  createdAt.Format(time.RFC3339),
		"inLanguage":    lang,
		"keywords": []string{
			"news tamil nadu", "tamil nadu news", "tn 24", "news tamil 24x7 live",
			"news live tamilnadu", "தமிழ் செய்திகள்", "today news in tamil",
			"news tamil today", "tamil news online", "latest tamil news", "tamil nadu news in tamil",
		},
		"articleSection": category,
		"contentLocation": map[string]string{
			"@type": "AdministrativeArea",
			"name":  district,
		},
		"publisher": map[string]interface{}{
			"@type": "NewsMediaOrganization",
			"name":  "TN24 — Tamil Nadu News",
			"url":   fmt.Sprintf("%s://%s/portal", scheme, host),
			"logo": map[string]string{
				"@type": "ImageObject",
				"url":   fmt.Sprintf("%s://%s/portal/assets/brand/tn24-profile.jpg", scheme, host),
			},
		},
	})

	injectedHead := fmt.Sprintf(`<script type="application/ld+json">%s</script><script>window.INITIAL_POST_ID = %q;</script></head>`, string(articleSchemaJSON), postID)
	baseHTML = strings.Replace(baseHTML, "</head>", injectedHead, 1)

	return baseHTML
}

func (h *PortalHandler) HandlePortalFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"success":false,"message":"Method not allowed"}`))
		return
	}

	if h.conn == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"success":false,"message":"Database not connected"}`))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	districtFilter := strings.TrimSpace(r.URL.Query().Get("district"))
	categoryFilter := strings.TrimSpace(r.URL.Query().Get("category"))
	searchQuery := strings.TrimSpace(r.URL.Query().Get("q"))
	viralFilter := strings.TrimSpace(r.URL.Query().Get("viral"))
	langFilter := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("lang")))

	// Portal displays ONLY APPROVED (PUBLISHED) content within the active 72-hour window
	whereClauses := []string{"c.status = 'PUBLISHED'", "COALESCE(c.published_at, c.created_at) >= NOW() - interval '72 hours'"}
	var args []interface{}
	argIdx := 1

	if viralFilter == "true" {
		whereClauses = append(whereClauses, "c.is_viral = TRUE")
	}

	if langFilter == "ta" {
		whereClauses = append(whereClauses, "(c.title ~ '[\\u0B80-\\u0BFF]' OR COALESCE(c.description, '') ~ '[\\u0B80-\\u0BFF]')")
	} else if langFilter == "en" {
		whereClauses = append(whereClauses, "(c.title !~ '[\\u0B80-\\u0BFF]' AND COALESCE(c.description, '') !~ '[\\u0B80-\\u0BFF]')")
	}

	if districtFilter != "" && districtFilter != "All Districts" && districtFilter != "ALL" && districtFilter != "TN-ALL / REGIONAL" {
		aliases, hasAliases := getDistrictAliases(districtFilter)
		if hasAliases && len(aliases) > 0 {
			var distParts []string
			distParts = append(distParts, fmt.Sprintf("d.name ILIKE $%d", argIdx))
			args = append(args, "%"+districtFilter+"%")
			argIdx++
			for _, alias := range aliases {
				distParts = append(distParts, fmt.Sprintf("(c.title ILIKE $%d OR COALESCE(c.description, '') ILIKE $%d OR d.name ILIKE $%d)", argIdx, argIdx, argIdx))
				args = append(args, "%"+alias+"%")
				argIdx++
			}
			whereClauses = append(whereClauses, "("+strings.Join(distParts, " OR ")+")")
		} else {
			whereClauses = append(whereClauses, fmt.Sprintf("(d.name ILIKE $%d OR c.title ILIKE $%d OR COALESCE(c.description, '') ILIKE $%d)", argIdx, argIdx, argIdx))
			args = append(args, "%"+districtFilter+"%")
			argIdx++
		}
	}

	if categoryFilter != "" && categoryFilter != "All" && categoryFilter != "ALL" {
		catLower := strings.ToLower(categoryFilter)
		switch catLower {
		case "sports":
			whereClauses = append(whereClauses, "(cat.name ILIKE '%Sports%' OR c.title ILIKE '%கிரிக்கெட்%' OR c.title ILIKE '%கால்பந்து%' OR c.title ILIKE '%துலீப்%' OR c.title ILIKE '%ஐபிஎல்%' OR c.title ILIKE '%cricket%' OR c.title ILIKE '%football%' OR c.title ILIKE '%ipl%' OR c.title ILIKE '%ரோகித்%' OR c.title ILIKE '%தோனி%' OR c.title ILIKE '%கோலி%' OR c.title ILIKE '%செஸ்%' OR c.title ILIKE '%கபடி%' OR c.title ILIKE '%தடகள%' OR c.title ILIKE '%விளையாட்டு செய்தி%' OR c.title ILIKE '%விளையாட்டு அரங்கம்%') AND c.title NOT ILIKE '%NOTE THIS POINT%' AND c.title NOT ILIKE '%அரசியல்%' AND c.title NOT ILIKE '%கூட்டாட்சி%'")
		case "crime":
			whereClauses = append(whereClauses, "(cat.name ILIKE '%Crime%' OR c.title ILIKE '%crime%' OR c.title ILIKE '%குற்றம்%' OR c.title ILIKE '%கைது%' OR c.title ILIKE '%போலீஸ்%' OR c.title ILIKE '%கொலை%' OR c.title ILIKE '%மோசடி%' OR c.title ILIKE '%police%' OR c.title ILIKE '%arrest%')")
		case "politics":
			whereClauses = append(whereClauses, "(cat.name ILIKE '%Politics%' OR c.title ILIKE '%அரசியல்%' OR c.title ILIKE '%அரசு%' OR c.title ILIKE '%கட்சி%' OR c.title ILIKE '%தேர்தல்%' OR c.title ILIKE '%முதல்வர்%' OR c.title ILIKE '%அமைச்சர்%' OR c.title ILIKE '%திமுக%' OR c.title ILIKE '%அதிமுக%' OR c.title ILIKE '%பாஜக%' OR c.title ILIKE '%தவெக%' OR c.title ILIKE '%MLA%' OR c.title ILIKE '%MP%' OR c.title ILIKE '%politics%' OR c.title ILIKE '%minister%' OR c.title ILIKE '%government%' OR COALESCE(c.description, '') ILIKE '%அரசியல்%' OR COALESCE(c.description, '') ILIKE '%அரசு%' OR COALESCE(c.description, '') ILIKE '%முதல்வர்%' OR COALESCE(c.description, '') ILIKE '%அமைச்சர்%')")
		case "technical", "tech":
			whereClauses = append(whereClauses, "(cat.name ILIKE '%Tech%' OR c.title ILIKE '%tech%' OR c.title ILIKE '%ai%' OR c.title ILIKE '%crypto%' OR c.title ILIKE '%bitcoin%' OR c.title ILIKE '%software%' OR c.title ILIKE '%இஸ்ரோ%' OR c.title ILIKE '%விண்கலம்%' OR c.title ILIKE '%செயற்கைக்கோள்%' OR c.title ILIKE '%தொழில்நுட்ப%' OR c.title ILIKE '%சாப்ட்வேர்%' OR c.title ILIKE '%digital%' OR c.title ILIKE '%science%' OR COALESCE(c.description, '') ILIKE '%tech%' OR COALESCE(c.description, '') ILIKE '%தொழில்நுட்ப%' OR COALESCE(c.description, '') ILIKE '%இஸ்ரோ%')")
		case "business":
			whereClauses = append(whereClauses, "(cat.name ILIKE '%Business%' OR cat.name ILIKE '%Finance%' OR c.title ILIKE '%business%' OR c.title ILIKE '%market%' OR c.title ILIKE '%finance%' OR c.title ILIKE '%bank%' OR c.title ILIKE '%வங்கி%' OR c.title ILIKE '%ரூபாய்%' OR c.title ILIKE '%பொருளாதார%' OR c.title ILIKE '%பங்கு%' OR c.title ILIKE '%தங்கம்%' OR c.title ILIKE '%விலை%' OR COALESCE(c.description, '') ILIKE '%finance%' OR COALESCE(c.description, '') ILIKE '%வங்கி%' OR COALESCE(c.description, '') ILIKE '%ரூபாய்%' OR COALESCE(c.description, '') ILIKE '%business%')")
		case "entertainment", "cinema":
			whereClauses = append(whereClauses, "(cat.name ILIKE '%Entertainment%' OR cat.name ILIKE '%Cinema%' OR c.title ILIKE '%cinema%' OR c.title ILIKE '%movie%' OR c.title ILIKE '%actor%' OR c.title ILIKE '%திரைப்பட%' OR c.title ILIKE '%நடிக%' OR c.title ILIKE '%Bigg Boss%' OR c.title ILIKE '%படம்%' OR c.title ILIKE '%விஜய்%' OR c.title ILIKE '%அஜித்%' OR c.title ILIKE '%ரஜினி%' OR c.title ILIKE '%கமல்%' OR c.title ILIKE '%பாடல்%' OR COALESCE(c.description, '') ILIKE '%cinema%' OR COALESCE(c.description, '') ILIKE '%திரைப்பட%' OR COALESCE(c.description, '') ILIKE '%நடிக%')")
		case "news", "civic":
			whereClauses = append(whereClauses, "(cat.name ILIKE '%News%' OR cat.name ILIKE '%Civic%' OR cat.name ILIKE '%General%' OR c.title ILIKE '%செய்தி%' OR c.title ILIKE '%மக்கள்%' OR c.title ILIKE '%மாவட்டம்%' OR COALESCE(c.description, '') ILIKE '%செய்தி%')")
		default:
			whereClauses = append(whereClauses, fmt.Sprintf("(cat.name ILIKE $%d OR c.title ILIKE $%d OR COALESCE(c.description, '') ILIKE $%d)", argIdx, argIdx, argIdx))
			args = append(args, "%"+categoryFilter+"%", "%"+categoryFilter+"%", "%"+categoryFilter+"%")
			argIdx++
		}
	}

	if searchQuery != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(c.title ILIKE $%d OR c.description ILIKE $%d OR d.name ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+searchQuery+"%")
		argIdx++
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	var mainBannerID, adHeaderID, adSidebarID, adInfeedID, adSquareID string
	_ = h.conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'portal_main_banner_id'").Scan(&mainBannerID)
	_ = h.conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'portal_ad_banner_header_id'").Scan(&adHeaderID)
	_ = h.conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'portal_ad_banner_sidebar_id'").Scan(&adSidebarID)
	_ = h.conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'portal_ad_banner_infeed_id'").Scan(&adInfeedID)
	_ = h.conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'portal_ad_banner_square_id'").Scan(&adSquareID)

	query := fmt.Sprintf(`
		SELECT c.id, c.title, COALESCE(c.description, ''), c.content_type, COALESCE(c.source_url, ''),
		       COALESCE(d.name, 'Tamil Nadu'), COALESCE(cat.name, 'News'), c.status,
		       COALESCE(vl.external_video_id, ''), COALESCE(vl.canonical_url, ''), COALESCE(vl.platform, ''),
		       COALESCE(vl.thumbnail_url, s.single_photo_url, (p.photo_urls)[1], ''),
		       COALESCE(c.published_at, c.created_at), COALESCE(c.is_viral, false)
		FROM content c
		LEFT JOIN districts d ON c.district_id = d.id
		LEFT JOIN categories cat ON c.category_id = cat.id
		LEFT JOIN video_links vl ON c.id = vl.content_id
		LEFT JOIN stories s ON c.id = s.content_id
		LEFT JOIN photos p ON c.id = p.content_id
		WHERE %s
		ORDER BY COALESCE(c.published_at, c.created_at) DESC, c.id DESC
		LIMIT 80
	`, whereSQL)

	rows, err := h.conn.Query(ctx, query, args...)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	defer rows.Close()

	var allItems []ArticleItem
	viralCount := 0

	authorPool := []string{"Adam Smith", "P. Senthil Nathan", "K. Radhakrishnan", "Meera Ramanathan", "TN24 Editorial", "S. Vijayakumar", "Ananya Sundaram"}

	for rows.Next() {
		var it ArticleItem
		var cid pgtype.UUID
		err := rows.Scan(&cid, &it.Title, &it.Description, &it.ContentType, &it.SourceURL,
			&it.District, &it.Category, &it.Status, &it.VideoID, &it.VideoURL, &it.VideoType, &it.Thumbnail, &it.CreatedAt, &it.IsViral)
		if err == nil {
			it.ID = fmtUUID(cid)
			it.Title = scraper.CleanHTML(it.Title)
			it.Description = sanitizeEditorialDescription(scraper.CleanHTML(it.Description), it.Title, it.District)
			it.Language = scraper.DetectLanguage(it.Title + " " + it.Description)
			if it.IsViral {
				viralCount++
			}
			if strings.TrimSpace(it.Thumbnail) == "" {
				it.Thumbnail = scraper.GetFallbackImageWithPerson(it.Title, it.District, it.Category)
			}
			idx := 0
			if len(it.ID) > 0 {
				idx = int(it.ID[len(it.ID)-1]) % len(authorPool)
			}
			it.Author = authorPool[idx]
			allItems = append(allItems, it)
		}
	}

	// In-memory duplicate removal: retain latest one (allItems is already sorted by latest first)
	var distinctItems []ArticleItem
	seenNormTitles := make(map[string]bool)
	seenCanonURLs := make(map[string]bool)
	seenVideoIDs := make(map[string]bool)

	for _, it := range allItems {
		normT := scraper.NormalizeTitle(it.Title)
		canonU := scraper.CanonicalizeURL(it.SourceURL)

		if len([]rune(normT)) >= 8 {
			if seenNormTitles[normT] {
				continue
			}
			seenNormTitles[normT] = true
		}
		if canonU != "" && canonU != "http://" && canonU != "https://" {
			if seenCanonURLs[canonU] {
				continue
			}
			seenCanonURLs[canonU] = true
		}
		if it.VideoID != "" {
			if seenVideoIDs[it.VideoID] {
				continue
			}
			seenVideoIDs[it.VideoID] = true
		}
		distinctItems = append(distinctItems, it)
	}

	// Strictly ensure allItems is sorted newest first by CreatedAt / PublishedAt
	sort.SliceStable(distinctItems, func(i, j int) bool {
		return distinctItems[i].CreatedAt.After(distinctItems[j].CreatedAt)
	})
	allItems = distinctItems

	resp := PortalFeedResponse{
		HeroTeasers:    make([]ArticleItem, 0),
		LeftFeed:       make([]ArticleItem, 0),
		PressReleases:  make([]ArticleItem, 0),
		CenterArticles: make([]ArticleItem, 0),
		MostRead:       make([]ArticleItem, 0),
		Districts:      []string{"Chennai", "Coimbatore", "Madurai", "Tiruchirappalli", "Salem", "Tirunelveli", "Erode", "Vellore", "Thanjavur", "Kanyakumari", "Ranipet", "Dindigul"},
		Categories:     []string{"News", "Politics", "Civic", "Business", "Technical", "Entertainment", "Crime", "Sports"},
		TotalCount:     len(allItems),
		ViralCount:     viralCount,
		Banners:        make(map[string]*AdBannerSlot),
	}

	// Determine Hero banner
	var heroItem *ArticleItem
	if mainBannerID != "" && districtFilter == "" && categoryFilter == "" && searchQuery == "" {
		// Check if already in allItems
		for i, it := range allItems {
			if it.ID == mainBannerID {
				heroCopy := it
				heroItem = &heroCopy
				allItems = append(allItems[:i], allItems[i+1:]...)
				break
			}
		}
		// If not in top 60, fetch directly
		if heroItem == nil {
			var it ArticleItem
			var cid pgtype.UUID
			singleQ := `
				SELECT c.id, c.title, COALESCE(c.description, ''), c.content_type, COALESCE(c.source_url, ''),
				       COALESCE(d.name, 'Tamil Nadu'), COALESCE(cat.name, 'News'), c.status,
				       COALESCE(vl.external_video_id, ''), COALESCE(vl.canonical_url, ''), COALESCE(vl.platform, ''),
				       COALESCE(vl.thumbnail_url, s.single_photo_url, (p.photo_urls)[1], ''),
				       COALESCE(c.published_at, c.created_at), COALESCE(c.is_viral, false)
				FROM content c
				LEFT JOIN districts d ON c.district_id = d.id
				LEFT JOIN categories cat ON c.category_id = cat.id
				LEFT JOIN video_links vl ON c.id = vl.content_id
				LEFT JOIN stories s ON c.id = s.content_id
				LEFT JOIN photos p ON c.id = p.content_id
				WHERE c.id::text = $1 AND c.status = 'PUBLISHED'
				LIMIT 1
			`
			err := h.conn.QueryRow(ctx, singleQ, mainBannerID).Scan(
				&cid, &it.Title, &it.Description, &it.ContentType, &it.SourceURL,
				&it.District, &it.Category, &it.Status, &it.VideoID, &it.VideoURL, &it.VideoType, &it.Thumbnail, &it.CreatedAt, &it.IsViral,
			)
			if err == nil {
				it.ID = fmtUUID(cid)
				it.Title = scraper.CleanHTML(it.Title)
				it.Description = sanitizeEditorialDescription(scraper.CleanHTML(it.Description), it.Title, it.District)
				if strings.TrimSpace(it.Thumbnail) == "" {
					it.Thumbnail = scraper.GetFallbackImageWithPerson(it.Title, it.District, it.Category)
				}
				it.Author = "TN24 Chief Correspondent"
				heroItem = &it
			}
		}
	}

	if heroItem == nil && len(allItems) > 0 {
		bestIdx := 0
		// Pick from the top 5 newest articles the first one that has a clean thumbnail
		topCheck := 5
		if len(allItems) < topCheck {
			topCheck = len(allItems)
		}
		for i := 0; i < topCheck; i++ {
			it := allItems[i]
			hasRealImage := strings.TrimSpace(it.Thumbnail) != "" &&
				!strings.Contains(it.Thumbnail, "upload.wikimedia.org") &&
				!strings.Contains(it.Thumbnail, "/admin/api/maps/svg") &&
				!strings.Contains(it.Thumbnail, "dummy.svg")
			if hasRealImage {
				bestIdx = i
				break
			}
		}

		chosen := allItems[bestIdx]
		heroItem = &chosen
		allItems = append(allItems[:bestIdx], allItems[bestIdx+1:]...)
	}
	resp.Hero = heroItem

	// Helper to find or load item for Ad slot
	findOrLoadItem := func(slotID string) *ArticleItem {
		if slotID == "" {
			return nil
		}
		for _, it := range allItems {
			if it.ID == slotID {
				cp := it
				return &cp
			}
		}
		var it ArticleItem
		var cid pgtype.UUID
		singleQ := `
			SELECT c.id, c.title, COALESCE(c.description, ''), c.content_type, COALESCE(c.source_url, ''),
			       COALESCE(d.name, 'Tamil Nadu'), COALESCE(cat.name, 'News'), c.status,
			       COALESCE(vl.external_video_id, ''), COALESCE(vl.canonical_url, ''), COALESCE(vl.platform, ''),
			       COALESCE(vl.thumbnail_url, s.single_photo_url, (p.photo_urls)[1], ''),
			       COALESCE(c.published_at, c.created_at), COALESCE(c.is_viral, false)
			FROM content c
			LEFT JOIN districts d ON c.district_id = d.id
			LEFT JOIN categories cat ON c.category_id = cat.id
			LEFT JOIN video_links vl ON c.id = vl.content_id
			LEFT JOIN stories s ON c.id = s.content_id
			LEFT JOIN photos p ON c.id = p.content_id
			WHERE c.id::text = $1 AND c.status = 'PUBLISHED'
			LIMIT 1
		`
		err := h.conn.QueryRow(ctx, singleQ, slotID).Scan(
			&cid, &it.Title, &it.Description, &it.ContentType, &it.SourceURL,
			&it.District, &it.Category, &it.Status, &it.VideoID, &it.VideoURL, &it.VideoType, &it.Thumbnail, &it.CreatedAt, &it.IsViral,
		)
		if err == nil {
			it.ID = fmtUUID(cid)
			it.Title = scraper.CleanHTML(it.Title)
			it.Description = scraper.CleanHTML(it.Description)
			it.Language = scraper.DetectLanguage(it.Title + " " + it.Description)
			if strings.TrimSpace(it.Thumbnail) == "" {
				it.Thumbnail = scraper.GetFallbackImageWithPerson(it.Title, it.District, it.Category)
			}
			it.Author = "TN24 Editorial"
			return &it
		}
		return nil
	}

	assignSlot := func(slotName, slotID string) {
		if slotID != "" {
			if it := findOrLoadItem(slotID); it != nil {
				resp.Banners[slotName] = &AdBannerSlot{
					Type:    "article",
					Article: it,
				}
				return
			}
		}
		resp.Banners[slotName] = &AdBannerSlot{
			Type:     "brand",
			AssetURL: "/portal/assets/brand/tn24-" + slotName + ".svg?v=20260908c",
		}
	}

	assignSlot("header", adHeaderID)
	assignSlot("sidebar", adSidebarID)
	assignSlot("infeed", adInfeedID)
	assignSlot("square", adSquareID)

	isFiltered := (districtFilter != "" || categoryFilter != "" || viralFilter == "true" || searchQuery != "")
	if isFiltered {
		// When a filter is active (e.g. Sports, Viral Radar, or specific District/Search),
		// all matched articles must be presented directly in the Center Editorial feed!
		if heroItem != nil {
			resp.CenterArticles = append(resp.CenterArticles, *heroItem)
		}
		resp.CenterArticles = append(resp.CenterArticles, allItems...)
		for i := 0; i < len(allItems) && i < 3; i++ {
			resp.HeroTeasers = append(resp.HeroTeasers, allItems[i])
		}
		for i := 0; i < len(allItems) && i < 5; i++ {
			resp.LeftFeed = append(resp.LeftFeed, allItems[i])
		}
	} else {
		for i := 0; i < len(allItems) && i < 3; i++ {
			resp.HeroTeasers = append(resp.HeroTeasers, allItems[i])
		}

		for i := 0; i < len(allItems) && i < 6; i++ {
			resp.LeftFeed = append(resp.LeftFeed, allItems[i])
		}

		for i := 0; i < len(allItems) && i < 6; i++ {
			resp.PressReleases = append(resp.PressReleases, allItems[i])
		}

		// Center editorial grid shows all articles in strictly newest-first order
		resp.CenterArticles = append(resp.CenterArticles, allItems...)
	}

	// Guarantee strict newest-first ordering on the Center Articles feed
	sort.SliceStable(resp.CenterArticles, func(i, j int) bool {
		return resp.CenterArticles[i].CreatedAt.After(resp.CenterArticles[j].CreatedAt)
	})

	for i := 0; i < len(allItems) && i < 5; i++ {
		resp.MostRead = append(resp.MostRead, allItems[i])
	}

	resp.Events = GenerateRealtimeTNEvents(time.Now())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    resp,
		"message": "Portal feed fetched successfully",
	})
}

func GenerateRealtimeTNEvents(now time.Time) []EventItem {
	d1Start := now.AddDate(0, 0, 4)
	d1End := now.AddDate(0, 0, 6)
	d2Start := now.AddDate(0, 0, 9)
	d2End := now.AddDate(0, 0, 11)
	d3Start := now.AddDate(0, 0, 16)
	d3End := now.AddDate(0, 0, 19)
	d4Start := now.AddDate(0, 0, 24)
	d4End := now.AddDate(0, 0, 26)
	d5Start := now.AddDate(0, 1, 2)
	d5End := now.AddDate(0, 1, 5)
	d6Start := now.AddDate(0, 1, 10)
	d6End := now.AddDate(0, 1, 12)
	d7Start := now.AddDate(0, 1, 18)
	d7End := now.AddDate(0, 1, 20)
	d8Start := now.AddDate(0, 2, 5)
	d8End := now.AddDate(0, 2, 7)

	fmtRange := func(s, e time.Time) string {
		if s.Month() == e.Month() {
			return fmt.Sprintf("%s %02d - %02d, %d", strings.ToUpper(s.Format("Jan")), s.Day(), e.Day(), s.Year())
		}
		return fmt.Sprintf("%s %02d - %s %02d, %d", strings.ToUpper(s.Format("Jan")), s.Day(), strings.ToUpper(e.Format("Jan")), e.Day(), s.Year())
	}

	return []EventItem{
		{
			ID:          "event-1",
			Title:       "TN Global Investors Meet & Innovation Summit",
			Category:    "Governance & Industry",
			District:    "Chennai",
			Venue:       "Chennai Trade Centre • Nandambakkam",
			Badge:       "TN",
			BadgeColor:  "#0284c7",
			DateDisplay: fmtRange(d1Start, d1End),
			Countdown:   "In 4 Days",
			Status:      "REGISTRATION OPEN",
			Description: "Premier global summit hosting international delegates, tech entrepreneurs, policy makers, and industrial leaders to accelerate Tamil Nadu's $1 Trillion economy roadmap.",
		},
		{
			ID:          "event-2",
			Title:       "Tamil Nadu Agri-Tech & Modern Farmers Summit",
			Category:    "Agriculture & Innovation",
			District:    "Madurai",
			Venue:       "Madurai Convention Centre • Madurai",
			Badge:       "AG",
			BadgeColor:  "#16a34a",
			DateDisplay: fmtRange(d2Start, d2End),
			Countdown:   "In 9 Days",
			Status:      "UPCOMING",
			Description: "State-level agricultural exhibition featuring drone farming, precision irrigation, organic horticulture value chain, and direct farmer-to-market trade pavilions.",
		},
		{
			ID:          "event-3",
			Title:       "Coimbatore International EV & Automation Expo",
			Category:    "Tech & Manufacturing",
			District:    "Coimbatore",
			Venue:       "CODISSIA Trade Fair Complex • Coimbatore",
			Badge:       "EV",
			BadgeColor:  "#8b5cf6",
			DateDisplay: fmtRange(d3Start, d3End),
			Countdown:   "In 2 Weeks",
			Status:      "FEATURED",
			Description: "South India's flagship electric mobility, battery technology, and smart factory robotics showcase with over 350+ global manufacturers and research institutes.",
		},
		{
			ID:          "event-4",
			Title:       "Tamil Nadu Cultural & International Book Fair",
			Category:    "Arts & Literature",
			District:    "Tiruchirappalli",
			Venue:       "YMCA International Grounds • Tiruchirappalli",
			Badge:       "CL",
			BadgeColor:  "#ea580c",
			DateDisplay: fmtRange(d4Start, d4End),
			Countdown:   "In 3 Weeks",
			Status:      "OPEN TO PUBLIC",
			Description: "Celebrating classical and contemporary Tamil literature, digital publishing, multilingual authors panel, and cultural performing arts pavilions.",
		},
		{
			ID:          "event-5",
			Title:       "TN State Youth & CM Trophy Sports Championship",
			Category:    "Sports & Athletics",
			District:    "Chennai",
			Venue:       "Jawaharlal Nehru Stadium • Chennai",
			Badge:       "SP",
			BadgeColor:  "#ec4899",
			DateDisplay: fmtRange(d5Start, d5End),
			Countdown:   "Next Month",
			Status:      "CONFIRMED",
			Description: "State-wide inter-district athletic championships, swimming, football, kabaddi, and track events featuring top emerging athletes from all 38 districts.",
		},
		{
			ID:          "event-6",
			Title:       "Salem Global Textile & Sustainable Handloom Conclave",
			Category:    "Commerce & Textile",
			District:    "Salem",
			Venue:       "Salem Trade Centre • Shevapet",
			Badge:       "TX",
			BadgeColor:  "#0d9488",
			DateDisplay: fmtRange(d6Start, d6End),
			Countdown:   "In 6 Weeks",
			Status:      "BUYER PASSES",
			Description: "International textile buyer-seller meet focusing on technical textiles, sustainable organic dyes, modern spinning innovations, and export corridors.",
		},
		{
			ID:          "event-7",
			Title:       "Hosur Advanced Aerospace & Defense Manufacturing Expo",
			Category:    "Aerospace & Defense",
			District:    "Krishnagiri",
			Venue:       "Hosur Industrial Growth Centre • Hosur",
			Badge:       "AD",
			BadgeColor:  "#6366f1",
			DateDisplay: fmtRange(d7Start, d7End),
			Countdown:   "In 7 Weeks",
			Status:      "INDUSTRY PASS",
			Description: "Showcasing Tamil Nadu Defense Industrial Corridor components, drone propulsion systems, precision machining, and avionics supply chains.",
		},
		{
			ID:          "event-8",
			Title:       "Thanjavur Heritage Tourism & Chola Architecture Conclave",
			Category:    "Tourism & Heritage",
			District:    "Thanjavur",
			Venue:       "Royal Palace Heritage Complex • Thanjavur",
			Badge:       "TR",
			BadgeColor:  "#eab308",
			DateDisplay: fmtRange(d8Start, d8End),
			Countdown:   "Next Quarter",
			Status:      "DELEGATE REGISTRATION",
			Description: "UNESCO world heritage celebration, temple art preservation conference, classical Thanjavur arts festival, and eco-tourism promotion panels.",
		},
	}
}

func (h *PortalHandler) HandleBrandAssets(w http.ResponseWriter, r *http.Request) {
	rawPath := r.URL.Path
	parts := strings.Split(rawPath, "/")
	name := ""
	if len(parts) > 0 {
		name = parts[len(parts)-1]
	}
	name = strings.ToLower(strings.TrimSpace(name))

	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	switch name {
	case "tn24-logo.svg":
		_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 260 52" width="260" height="52">
			<defs>
				<linearGradient id="logoCapsuleBg" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" stop-color="#0b1329" />
					<stop offset="50%" stop-color="#070d1c" />
					<stop offset="100%" stop-color="#030712" />
				</linearGradient>
				<linearGradient id="capsuleBorder" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" stop-color="#38bdf8" />
					<stop offset="50%" stop-color="#6366f1" />
					<stop offset="100%" stop-color="#f43f5e" />
				</linearGradient>
				<linearGradient id="emblemBg" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" stop-color="#1e293b" />
					<stop offset="100%" stop-color="#0f172a" />
				</linearGradient>
				<linearGradient id="badge24" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" stop-color="#ef4444" />
					<stop offset="100%" stop-color="#f97316" />
				</linearGradient>
			</defs>
			<rect x="1" y="1" width="258" height="50" rx="10" fill="url(#logoCapsuleBg)" stroke="url(#capsuleBorder)" stroke-width="1.2" />
			<rect x="7" y="7" width="36" height="36" rx="8" fill="url(#emblemBg)" stroke="#38bdf8" stroke-width="1.2" />
			<path d="M 23 13 L 30 13 L 23 23 L 31 23 L 19 37 L 23 27 L 17 27 Z" fill="#38bdf8" />
			<text x="52" y="32" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="900" font-size="24" fill="#ffffff" letter-spacing="0.5">TN</text>
			<rect x="92" y="12" width="42" height="24" rx="5" fill="url(#badge24)" />
			<text x="113" y="29.5" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="900" font-size="16" fill="#ffffff" text-anchor="middle" letter-spacing="0.5">24</text>
			<rect x="142" y="15" width="38" height="18" rx="4" fill="#0284c7" />
			<circle cx="150" cy="24" r="2.5" fill="#f43f5e" />
			<text x="165" y="27.5" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="8.5" fill="#ffffff" text-anchor="middle" letter-spacing="0.5">LIVE</text>
			<text x="53" y="44" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="700" font-size="7.5" fill="#94a3b8" letter-spacing="1.2">24/7 TAMIL NADU NEWS</text>
		</svg>`))
	case "tn24-profile.jpg", "profile.jpg":
		w.Header().Set("Content-Type", "image/jpeg")
		http.ServeFile(w, r, "/Users/rahamathalikhan/Documents/tn-now/tn24_social_profile_logo.jpg")
		return
	case "tn24-profile.svg", "profile.svg":
		_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" width="512" height="512">
			<defs>
				<radialGradient id="profileBg" cx="50%" cy="50%" r="50%">
					<stop offset="0%" stop-color="#0b1329" />
					<stop offset="80%" stop-color="#030712" />
					<stop offset="100%" stop-color="#010308" />
				</radialGradient>
				<linearGradient id="profileRing" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" stop-color="#38bdf8" />
					<stop offset="50%" stop-color="#6366f1" />
					<stop offset="100%" stop-color="#f43f5e" />
				</linearGradient>
				<linearGradient id="shieldFill" x1="0%" y1="0%" x2="0%" y2="100%">
					<stop offset="0%" stop-color="#0f172a" />
					<stop offset="100%" stop-color="#020617" />
				</linearGradient>
				<linearGradient id="badge24Profile" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" stop-color="#ef4444" />
					<stop offset="100%" stop-color="#f97316" />
				</linearGradient>
				<linearGradient id="livePillBg" x1="0%" y1="0%" x2="100%" y2="0%">
					<stop offset="0%" stop-color="#0369a1" />
					<stop offset="100%" stop-color="#0284c7" />
				</linearGradient>
				<filter id="glowG" x="-20%" y="-20%" width="140%" height="140%">
					<feGaussianBlur stdDeviation="6" result="blur" />
					<feComposite in="SourceGraphic" in2="blur" operator="over" />
				</filter>
			</defs>
			<!-- Circular Profile Canvas -->
			<circle cx="256" cy="256" r="250" fill="url(#profileBg)" stroke="url(#profileRing)" stroke-width="6" />
			
			<!-- Background Glow Emblem -->
			<path d="M 256 95 L 365 145 L 365 270 C 365 345 256 395 256 395 C 256 395 147 345 147 270 L 147 145 Z" fill="url(#shieldFill)" stroke="#38bdf8" stroke-width="3.5" filter="url(#glowG)" opacity="0.9" />
			
			<!-- Inner Shield Border -->
			<path d="M 256 112 L 350 155 L 350 265 C 350 330 256 375 256 375 C 256 375 162 330 162 265 L 162 155 Z" fill="none" stroke="rgba(244,63,94,0.4)" stroke-width="2" />
			
			<!-- Center Brand Typography -->
			<g transform="translate(0, -10)">
				<text x="180" y="272" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="900" font-size="115" fill="#ffffff" text-anchor="middle" letter-spacing="-2">TN</text>
				<rect x="256" y="180" width="135" height="100" rx="22" fill="url(#badge24Profile)" filter="url(#glowG)" />
				<text x="323" y="260" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="900" font-size="78" fill="#ffffff" text-anchor="middle">24</text>
			</g>

			<!-- 24/7 Live Capsule -->
			<g transform="translate(0, 15)">
				<rect x="156" y="320" width="200" height="42" rx="21" fill="url(#livePillBg)" stroke="#38bdf8" stroke-width="2" />
				<circle cx="180" cy="341" r="5.5" fill="#4ade80" />
				<text x="260" y="347" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="900" font-size="16" fill="#ffffff" text-anchor="middle" letter-spacing="2.5">24/7 LIVE</text>
				<circle cx="332" cy="341" r="5.5" fill="#ef4444" />
			</g>

			<!-- Bottom State Tagline -->
			<text x="256" y="442" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="14" fill="#38bdf8" text-anchor="middle" letter-spacing="4">TAMIL NADU NEWS</text>
		</svg>`))
	case "tn24-icon.svg":
		_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64" width="64" height="64">
			<defs>
				<linearGradient id="iconBg" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" stop-color="#0b1329" />
					<stop offset="100%" stop-color="#030712" />
				</linearGradient>
				<linearGradient id="iconBorder" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" stop-color="#38bdf8" />
					<stop offset="50%" stop-color="#6366f1" />
					<stop offset="100%" stop-color="#f43f5e" />
				</linearGradient>
				<linearGradient id="badge24Sm" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" stop-color="#ef4444" />
					<stop offset="100%" stop-color="#f97316" />
				</linearGradient>
			</defs>
			<rect x="2" y="2" width="60" height="60" rx="14" fill="url(#iconBg)" stroke="url(#iconBorder)" stroke-width="2" />
			<text x="32" y="31" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="900" font-size="21" fill="#ffffff" text-anchor="middle">TN</text>
			<rect x="14" y="37" width="36" height="18" rx="4" fill="url(#badge24Sm)" />
			<text x="32" y="51" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="900" font-size="13" fill="#ffffff" text-anchor="middle">24</text>
		</svg>`))
	case "tn24-header.svg", "tn24-leaderboard.svg":
		_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 728 90" width="728" height="90">
			<defs>
				<linearGradient id="bgH" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" stop-color="#091326" />
					<stop offset="100%" stop-color="#0f172a" />
				</linearGradient>
				<pattern id="gridH" width="20" height="20" patternUnits="userSpaceOnUse">
					<line x1="0" y1="0" x2="20" y2="0" stroke="rgba(56,189,248,0.07)" stroke-width="1" />
					<line x1="0" y1="0" x2="0" y2="20" stroke="rgba(56,189,248,0.07)" stroke-width="1" />
				</pattern>
				<linearGradient id="cyanPill" x1="0%" y1="0%" x2="100%" y2="0%">
					<stop offset="0%" stop-color="#0284c7" />
					<stop offset="100%" stop-color="#38bdf8" />
				</linearGradient>
			</defs>
			<rect width="728" height="90" rx="8" fill="url(#bgH)" stroke="#38bdf8" stroke-width="1.5" stroke-dasharray="6,4" />
			<rect width="728" height="90" rx="8" fill="url(#gridH)" />
			<rect x="20" y="18" width="138" height="22" rx="4" fill="url(#cyanPill)"/>
			<text x="89" y="33" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="10" fill="#091326" text-anchor="middle" letter-spacing="0.5">📢 ADVERTISE WITH US</text>
			<text x="172" y="34" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="14" fill="#f8fafc" letter-spacing="0.3">உங்கள் வணிக விளம்பரங்களுக்கு தொடர்பு கொள்ளவும்</text>
			<text x="20" y="58" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="500" font-size="11.5" fill="#94a3b8">Reach 1.85M+ Daily Readers across 38 Districts in Tamil Nadu • Digital Showcase Leaderboard</text>
			<text x="20" y="74" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="700" font-size="11.5" fill="#38bdf8">Email: tn24now@gmail.com</text>
			<rect x="560" y="22" width="148" height="46" rx="6" fill="#1e293b" stroke="#38bdf8" stroke-width="1.5"/>
			<text x="634" y="42" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="11.5" fill="#38bdf8" text-anchor="middle">தொடர்புக்கு ↗</text>
			<text x="634" y="57" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="700" font-size="9" fill="#94a3b8" text-anchor="middle">tn24now@gmail.com</text>
		</svg>`))
	case "tn24-sidebar.svg":
		_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 250 250" width="250" height="250">
			<defs>
				<linearGradient id="bgS" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" stop-color="#091326" />
					<stop offset="100%" stop-color="#0f172a" />
				</linearGradient>
				<pattern id="gridS" width="20" height="20" patternUnits="userSpaceOnUse">
					<line x1="0" y1="0" x2="20" y2="0" stroke="rgba(56,189,248,0.07)" stroke-width="1" />
					<line x1="0" y1="0" x2="0" y2="20" stroke="rgba(56,189,248,0.07)" stroke-width="1" />
				</pattern>
				<linearGradient id="cyanPillS" x1="0%" y1="0%" x2="100%" y2="0%">
					<stop offset="0%" stop-color="#0284c7" />
					<stop offset="100%" stop-color="#38bdf8" />
				</linearGradient>
			</defs>
			<rect width="250" height="250" rx="10" fill="url(#bgS)" stroke="#38bdf8" stroke-width="1.5" stroke-dasharray="6,4"/>
			<rect width="250" height="250" rx="10" fill="url(#gridS)" />
			<rect x="18" y="16" width="138" height="22" rx="4" fill="url(#cyanPillS)"/>
			<text x="87" y="31" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="9.5" fill="#091326" text-anchor="middle" letter-spacing="0.5">📢 ADVERTISE WITH US</text>
			<rect x="18" y="48" width="214" height="130" rx="8" fill="#131c31" stroke="rgba(56,189,248,0.35)" stroke-width="1"/>
			<circle cx="125" cy="80" r="16" fill="rgba(56,189,248,0.12)"/>
			<text x="125" y="86" font-family="sans-serif" font-size="16" fill="#38bdf8" text-anchor="middle">📢</text>
			<text x="125" y="112" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="12" fill="#f8fafc" text-anchor="middle">விளம்பரம் செய்ய தொடர்பு கொள்ளவும்</text>
			<text x="125" y="132" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="700" font-size="11" fill="#38bdf8" text-anchor="middle">tn24now@gmail.com</text>
			<text x="125" y="152" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="500" font-size="10" fill="#94a3b8" text-anchor="middle">தமிழ்நாடு முழுவதும் பிராண்ட் விளம்பரம்</text>
			<text x="125" y="168" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="600" font-size="9" fill="#64748b" text-anchor="middle">250 × 250 SIDEBAR SPOT</text>
			<rect x="18" y="190" width="214" height="38" rx="6" fill="#1e293b" stroke="#38bdf8" stroke-width="1.2"/>
			<text x="125" y="214" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="11" fill="#38bdf8" text-anchor="middle">தொடர்புக்கு: tn24now@gmail.com ↗</text>
		</svg>`))
	case "tn24-square.svg":
		_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 200 200" width="200" height="200">
			<defs>
				<linearGradient id="bgSq" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" stop-color="#091326" />
					<stop offset="100%" stop-color="#0f172a" />
				</linearGradient>
				<pattern id="gridSq" width="16" height="16" patternUnits="userSpaceOnUse">
					<line x1="0" y1="0" x2="16" y2="0" stroke="rgba(56,189,248,0.07)" stroke-width="1" />
					<line x1="0" y1="0" x2="0" y2="16" stroke="rgba(56,189,248,0.07)" stroke-width="1" />
				</pattern>
				<linearGradient id="cyanPillSq" x1="0%" y1="0%" x2="100%" y2="0%">
					<stop offset="0%" stop-color="#0284c7" />
					<stop offset="100%" stop-color="#38bdf8" />
				</linearGradient>
			</defs>
			<rect width="200" height="200" rx="10" fill="url(#bgSq)" stroke="#38bdf8" stroke-width="1.5" stroke-dasharray="5,4"/>
			<rect width="200" height="200" rx="10" fill="url(#gridSq)" />
			<rect x="16" y="14" width="130" height="20" rx="4" fill="url(#cyanPillSq)"/>
			<text x="81" y="28" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="8.5" fill="#091326" text-anchor="middle" letter-spacing="0.5">📢 ADVERTISE WITH US</text>
			<circle cx="100" cy="65" r="15" fill="rgba(56,189,248,0.12)"/>
			<text x="100" y="71" font-family="sans-serif" font-size="14" fill="#38bdf8" text-anchor="middle">📢</text>
			<text x="100" y="96" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="11" fill="#f8fafc" text-anchor="middle">விளம்பரம் செய்ய</text>
			<text x="100" y="112" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="700" font-size="10" fill="#38bdf8" text-anchor="middle">tn24now@gmail.com</text>
			<text x="100" y="128" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="500" font-size="8.5" fill="#94a3b8" text-anchor="middle">தொடர்பு கொள்ளவும்</text>
			<rect x="16" y="148" width="168" height="34" rx="5" fill="#1e293b" stroke="#38bdf8" stroke-width="1"/>
			<text x="100" y="170" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="9" fill="#38bdf8" text-anchor="middle">tn24now@gmail.com ↗</text>
		</svg>`))
	case "tn24-infeed.svg":
		_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 728 90" width="728" height="90">
			<defs>
				<linearGradient id="bgInf" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" stop-color="#091326" />
					<stop offset="100%" stop-color="#0f172a" />
				</linearGradient>
				<pattern id="gridInf" width="20" height="20" patternUnits="userSpaceOnUse">
					<line x1="0" y1="0" x2="20" y2="0" stroke="rgba(168,85,247,0.08)" stroke-width="1" />
					<line x1="0" y1="0" x2="0" y2="20" stroke="rgba(168,85,247,0.08)" stroke-width="1" />
				</pattern>
				<linearGradient id="purplePill" x1="0%" y1="0%" x2="100%" y2="0%">
					<stop offset="0%" stop-color="#7c3aed" />
					<stop offset="100%" stop-color="#a855f7" />
				</linearGradient>
			</defs>
			<rect width="728" height="90" rx="8" fill="url(#bgInf)" stroke="#a855f7" stroke-width="1.5" stroke-dasharray="6,4" />
			<rect width="728" height="90" rx="8" fill="url(#gridInf)" />
			<rect x="20" y="18" width="144" height="22" rx="4" fill="url(#purplePill)"/>
			<text x="92" y="33" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="10" fill="#ffffff" text-anchor="middle" letter-spacing="0.5">📢 SPONSORED ADS</text>
			<text x="176" y="34" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="14" fill="#f8fafc" letter-spacing="0.3">உங்கள் பிராண்ட் விளம்பரங்களுக்கு தொடர்பு கொள்ளவும்</text>
			<text x="20" y="58" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="500" font-size="11.5" fill="#cbd5e1">TN24 In-Feed Native Stream Placement • Highest Engagement Advertising Space</text>
			<text x="20" y="74" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="700" font-size="11.5" fill="#c084fc">Email: tn24now@gmail.com</text>
			<rect x="560" y="22" width="148" height="46" rx="6" fill="#1e1b4b" stroke="#a855f7" stroke-width="1.5"/>
			<text x="634" y="42" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="11.5" fill="#c084fc" text-anchor="middle">தொடர்புக்கு ↗</text>
			<text x="634" y="57" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="700" font-size="9" fill="#cbd5e1" text-anchor="middle">tn24now@gmail.com</text>
		</svg>`))
	default:
		// Universal Fallback Vector Placeholder (handles dummy.svg, missing images, placeholders)
		_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 800 450" width="100%" height="100%">
			<defs>
				<linearGradient id="dummyBg" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" stop-color="#091326" />
					<stop offset="50%" stop-color="#0f172a" />
					<stop offset="100%" stop-color="#1e293b" />
				</linearGradient>
				<pattern id="dummyGrid" width="30" height="30" patternUnits="userSpaceOnUse">
					<path d="M 30 0 L 0 0 0 30" fill="none" stroke="rgba(56,189,248,0.06)" stroke-width="1"/>
				</pattern>
			</defs>
			<rect width="800" height="450" fill="url(#dummyBg)" />
			<rect width="800" height="450" fill="url(#dummyGrid)" />
			<circle cx="400" cy="200" r="44" fill="rgba(56,189,248,0.12)" stroke="#38bdf8" stroke-width="1.5" stroke-dasharray="4,4"/>
			<text x="400" y="212" font-family="system-ui, sans-serif" font-size="34" text-anchor="middle">🏛️</text>
			<text x="400" y="280" font-family="system-ui, sans-serif" font-weight="800" font-size="22" fill="#f8fafc" text-anchor="middle" letter-spacing="1">TN24 MEDIA ASSET</text>
			<text x="400" y="310" font-family="system-ui, sans-serif" font-weight="500" font-size="14" fill="#94a3b8" text-anchor="middle">Official Tamil Nadu Regional &amp; District Coverage</text>
			<rect x="300" y="340" width="200" height="32" rx="6" fill="#0f172a" stroke="#38bdf8" stroke-width="1"/>
			<text x="400" y="361" font-family="monospace" font-weight="700" font-size="12" fill="#38bdf8" text-anchor="middle">tn24.news • 2026</text>
		</svg>`))
	}
}

func fmtUUID(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	src := u.Bytes
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		src[0], src[1], src[2], src[3],
		src[4], src[5],
		src[6], src[7],
		src[8], src[9],
		src[10], src[11], src[12], src[13], src[14], src[15],
	)
}

func getDistrictAliases(district string) ([]string, bool) {
	m := map[string][]string{
		"chennai":         {"சென்னை", "chennai", "தாம்பரம்", "ஆவடி"},
		"coimbatore":      {"கோவை", "கோயம்புத்தூர்", "coimbatore", "பொள்ளாச்சி"},
		"madurai":         {"மதுரை", "madurai", "மேலூர்"},
		"tiruchirappalli": {"திருச்சி", "திருச்சிராப்பள்ளி", "trichy", "tiruchirappalli", "ஸ்ரீரங்கம்"},
		"salem":           {"சேலம்", "salem", "ஆத்தூர்", "மேட்டூர்"},
		"tirunelveli":     {"திருநெல்வேலி", "நெல்லை", "tirunelveli"},
		"erode":           {"ஈரோடு", "erode", "பவானி", "கோபிசெட்டிபாளையம்"},
		"vellore":         {"வேலூர்", "காட்பாடி", "குடியாத்தம்", "vellore"},
		"thanjavur":       {"தஞ்சாவூர்", "தஞ்சை", "thanjavur", "கும்பகோணம்"},
		"kanyakumari":     {"கன்னியாகுமரி", "நாகர்கோவில்", "kanyakumari"},
		"ranipet":         {"ராணிப்பேட்டை", "ஆற்காடு", "ranipet"},
		"dindigul":        {"திண்டுக்கல்", "பழனி", "கொடைக்கானல்", "dindigul"},
		"tiruppur":        {"திருப்பூர்", "tiruppur", "அவிநாசி"},
		"kanchipuram":     {"காஞ்சிபுரம்", "காஞ்சி", "kanchipuram"},
		"chengalpattu":    {"செங்கல்பட்டு", "மகாபலிபுரம்", "chengalpattu"},
		"tiruvallur":      {"திருவள்ளூர்", "tiruvallur"},
		"cuddalore":       {"கடலூர்", "சிதம்பரம்", "neyveli", "cuddalore"},
		"dharmapuri":      {"தருமபுரி", "dharmapuri"},
		"krishnagiri":     {"கிருஷ்ணகிரி", "ஓசூர்", "hosur", "krishnagiri"},
		"theni":           {"தேனி", "theni"},
		"thoothukudi":     {"தூத்துக்குடி", "thoothukudi", "tuticorin", "கோவில்பட்டி"},
		"sivaganga":       {"சிவகங்கை", "காரைக்குடி", "sivaganga"},
		"ramanathapuram":  {"ராமநாதபுரம்", "ராமேஸ்வரம்", "ramanathapuram"},
		"virudhunagar":    {"விருதுநகர்", "சிவகாசி", "virudhunagar"},
		"nilgiris":        {"நீலகிரி", "ஊட்டி", "ooty", "nilgiris"},
		"namakkal":        {"நாமக்கல்", "namakkal"},
		"tiruvannamalai":  {"திருவண்ணாமலை", "tiruvannamalai"},
		"tiruvarur":       {"திருவாரூர்", "tiruvarur"},
		"nagapattinam":    {"நாகப்பட்டினம்", "நாகை", "nagapattinam"},
		"pudukkottai":     {"புதுக்கோட்டை", "pudukkottai"},
		"karur":           {"கரூர்", "karur"},
		"ariyalur":        {"அரியலூர்", "ariyalur"},
		"perambalur":      {"பெரம்பலூர்", "perambalur"},
		"viluppuram":      {"விழுப்புரம்", "viluppuram"},
		"kallakurichi":    {"கள்ளக்குறிச்சி", "kallakurichi"},
		"tenkasi":         {"தென்காசி", "குற்றாலம்", "tenkasi"},
		"tirupattur":      {"திருப்பத்தூர்", "ஆம்பூர்", "வாணியம்பாடி", "tirupattur"},
		"mayiladuthurai":  {"மயிலாடுதுறை", "mayiladuthurai"},
	}
	aliases, ok := m[strings.ToLower(strings.TrimSpace(district))]
	return aliases, ok
}

func sanitizeEditorialDescription(desc, title, district string) string {
	clean := strings.TrimSpace(desc)
	if clean == "" || strings.HasPrefix(clean, "Report from ") || strings.Contains(clean, "Full coverage: http") {
		if strings.ContainsAny(title, "அஆஇஈஉஊஎஏஐஒஓஔகஙசஞடணதநபமயரலவழளறன") {
			dist := district
			if dist == "" {
				dist = "தமிழ்நாடு"
			}
			return fmt.Sprintf("%s - %s மற்றும் தமிழ்நாடு வட்டார முக்கிய நிகழ்வுகள் குறித்த விரிவான கள நிலவரம் மற்றும் நேரடி செய்தி தொகுப்பு.", title, dist)
		}
		dist := district
		if dist == "" {
			dist = "Tamil Nadu"
		}
		return fmt.Sprintf("%s. Comprehensive on-ground news coverage and latest regional updates from %s.", title, dist)
	}
	return clean
}

// HandleAdsTxt serves /ads.txt — required by Google AdSense to authorize
// Google as a DIRECT seller of ads on this domain.
// Without this file AdSense dashboard shows "Ads.txt: Not found" and
// refuses to serve ads.
func (h *PortalHandler) HandleAdsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	// Format: <ad network domain>, <publisher ID>, <relationship>, <certification authority ID>
	_, _ = w.Write([]byte("google.com, pub-1894301748406603, DIRECT, f08c47fec0942fa0\n"))
}

func (h *PortalHandler) HandleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400")

	scheme := "https"
	host := r.Host
	if host == "" {
		host = "www.tn24.in"
	}

	robotsContent := fmt.Sprintf(`User-agent: *
Allow: /
Allow: /portal
Allow: /portal/*
Allow: /api/portal/feed
Allow: /privacy
Allow: /about
Allow: /contact
Disallow: /admin
Disallow: /admin/*
Disallow: /api/scraper/*
Disallow: /api/auth/*

User-agent: Googlebot-News
Allow: /
Allow: /portal
Allow: /portal/*
Allow: /sitemap-news.xml

Sitemap: %s://%s/sitemap.xml
Sitemap: %s://%s/sitemap-news.xml
`, scheme, host, scheme, host)
	_, _ = w.Write([]byte(robotsContent))
}

func (h *PortalHandler) HandleSitemapXML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	scheme := "https"
	host := r.Host
	if host == "" {
		host = "www.tn24.in"
	}
	base := fmt.Sprintf("%s://%s", scheme, host)
	now := time.Now().Format("2006-01-02")

	// Target SEO keyword search paths for crawler discovery
	keywordUrls := []string{
		"tn+24",
		"news+tamil+24x7+live",
		"tamil+nadu+news",
		"news+live+tamilnadu",
		"today+news+in+tamil",
		"news+tamil+today",
		"tamil+news+online",
		"latest+tamil+news",
		"tamil+nadu+news+in+tamil",
		"news+tamil+nadu",
	}

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	// Core Static Pages
	sb.WriteString(fmt.Sprintf(`    <url>
        <loc>%s/portal</loc>
        <lastmod>%s</lastmod>
        <changefreq>always</changefreq>
        <priority>1.0</priority>
    </url>
    <url>
        <loc>%s/</loc>
        <lastmod>%s</lastmod>
        <changefreq>always</changefreq>
        <priority>0.95</priority>
    </url>
    <url>
        <loc>%s/about</loc>
        <lastmod>%s</lastmod>
        <changefreq>monthly</changefreq>
        <priority>0.7</priority>
    </url>
    <url>
        <loc>%s/privacy</loc>
        <lastmod>%s</lastmod>
        <changefreq>monthly</changefreq>
        <priority>0.6</priority>
    </url>
    <url>
        <loc>%s/contact</loc>
        <lastmod>%s</lastmod>
        <changefreq>monthly</changefreq>
        <priority>0.6</priority>
    </url>
`, base, now, base, now, base, now, base, now, base, now))

	// Target Keyword Query URLs
	for _, kw := range keywordUrls {
		sb.WriteString(fmt.Sprintf(`    <url>
        <loc>%s/portal?q=%s</loc>
        <lastmod>%s</lastmod>
        <changefreq>hourly</changefreq>
        <priority>0.85</priority>
    </url>
`, base, kw, now))
	}

	// Categories & Key Districts
	cats := []string{"News", "Politics", "Sports", "Business", "Entertainment", "Technical"}
	for _, c := range cats {
		sb.WriteString(fmt.Sprintf(`    <url>
        <loc>%s/portal?cat=%s</loc>
        <lastmod>%s</lastmod>
        <changefreq>hourly</changefreq>
        <priority>0.8</priority>
    </url>
`, base, c, now))
	}

	districts := []string{"Chennai", "Madurai", "Coimbatore", "Salem", "Tiruchirappalli", "Tirunelveli"}
	for _, d := range districts {
		sb.WriteString(fmt.Sprintf(`    <url>
        <loc>%s/portal?district=%s</loc>
        <lastmod>%s</lastmod>
        <changefreq>hourly</changefreq>
        <priority>0.8</priority>
    </url>
`, base, d, now))
	}

	// Dynamic Published Articles (Top 100 recent articles)
	if h.conn != nil {
		rows, err := h.conn.Query(r.Context(), `
			SELECT c.id::text, COALESCE(c.published_at, c.created_at)
			FROM content c
			WHERE c.status = 'PUBLISHED'
			ORDER BY COALESCE(c.published_at, c.created_at) DESC
			LIMIT 100
		`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var aID string
				var aMod time.Time
				if err := rows.Scan(&aID, &aMod); err == nil {
					sb.WriteString(fmt.Sprintf(`    <url>
        <loc>%s/portal?post=%s</loc>
        <lastmod>%s</lastmod>
        <changefreq>daily</changefreq>
        <priority>0.75</priority>
    </url>
`, base, url.QueryEscape(aID), aMod.Format("2006-01-02")))
				}
			}
		}
	}

	sb.WriteString("</urlset>\n")
	_, _ = w.Write([]byte(sb.String()))
}

// HandleNewsSitemapXML serves Google News specific sitemap (/sitemap-news.xml)
// with articles from the last 48 hours per Google News sitemap specifications.
func (h *PortalHandler) HandleNewsSitemapXML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=1800")

	scheme := "https"
	host := r.Host
	if host == "" {
		host = "www.tn24.in"
	}
	base := fmt.Sprintf("%s://%s", scheme, host)

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
        xmlns:news="http://www.google.com/schemas/sitemap-news/0.9">` + "\n")

	if h.conn != nil {
		rows, err := h.conn.Query(r.Context(), `
			SELECT c.id::text, c.title, COALESCE(c.language, 'ta'), COALESCE(c.published_at, c.created_at)
			FROM content c
			WHERE c.status = 'PUBLISHED'
			  AND COALESCE(c.published_at, c.created_at) >= NOW() - INTERVAL '48 hours'
			ORDER BY COALESCE(c.published_at, c.created_at) DESC
			LIMIT 250
		`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var aID, aTitle, aLang string
				var aPub time.Time
				if err := rows.Scan(&aID, &aTitle, &aLang, &aPub); err == nil {
					cleanTitle := scraper.CleanHTML(aTitle)
					cleanTitle = html.EscapeString(cleanTitle)
					if aLang == "" {
						aLang = "ta"
					}
					sb.WriteString(fmt.Sprintf(`    <url>
        <loc>%s/portal?post=%s</loc>
        <news:news>
            <news:publication>
                <news:name>TN24 — Tamil Nadu News</news:name>
                <news:language>%s</news:language>
            </news:publication>
            <news:publication_date>%s</news:publication_date>
            <news:title>%s</news:title>
        </news:news>
    </url>
`, base, url.QueryEscape(aID), aLang, aPub.Format(time.RFC3339), cleanTitle))
				}
			}
		}
	}

	sb.WriteString("</urlset>\n")
	_, _ = w.Write([]byte(sb.String()))
}

// HandlePrivacyPolicy serves the Privacy Policy page — required for Google AdSense approval.
func (h *PortalHandler) HandlePrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(renderStaticPage("Privacy Policy | TN24 — Tamil Nadu News", "தனியுரிமைக் கொள்கை — Privacy Policy", `
<h2>Privacy Policy</h2>
<p><strong>Effective Date:</strong> September 2026</p>
<p>TN24 (<strong>www.tn24.in</strong>) is committed to protecting your privacy. This Privacy Policy explains how we collect, use, and safeguard your information when you visit our website.</p>

<h3>1. Information We Collect</h3>
<ul>
  <li><strong>Log Data:</strong> Browser type, IP address, operating system, and pages visited.</li>
  <li><strong>Cookies:</strong> Essential session cookies and analytics cookies to enhance user experience.</li>
  <li><strong>Advertising Cookies:</strong> Google AdSense uses cookies to serve personalized or non-personalized ads based on prior visits. Users may opt out of personalized advertising by visiting <a href="https://www.google.com/settings/ads" target="_blank" rel="noopener">Google Ads Settings</a>.</li>
</ul>

<h3>2. How We Use Information</h3>
<ul>
  <li>To provide and maintain the TN24 Tamil news portal.</li>
  <li>To analyze site traffic, popular content, and improve reader experience.</li>
  <li>To display advertisements that support our free news publishing operations.</li>
</ul>

<h3>3. Third-Party Services</h3>
<p>We work with trusted third-party providers including:</p>
<ul>
  <li><strong>Google AdSense:</strong> Advertising network (<a href="https://policies.google.com/privacy" target="_blank" rel="noopener">Google Privacy Policy</a>).</li>
  <li><strong>Google Analytics:</strong> Audience measurement and aggregated reporting.</li>
</ul>

<h3>4. Data Retention</h3>
<p>News articles published on TN24 are automatically retained for 24 hours before being removed or archived per our editorial retention policy.</p>

<h3>5. Your Rights</h3>
<p>You have the right to access, correct, or request deletion of your personal data. Contact us at <a href="mailto:admin@tn24.in">admin@tn24.in</a>.</p>

<h3>6. Grievance Officer</h3>
<p>In accordance with the Information Technology (Intermediary Guidelines and Digital Media Ethics Code) Rules, 2021:</p>
<p><strong>Grievance Officer:</strong> Editorial Desk, TN24 Media<br>
Email: <a href="mailto:admin@tn24.in">admin@tn24.in</a><br>
Tamil Nadu, India</p>

<h3>7. Contact Us</h3>
<p>If you have questions about this Privacy Policy, please contact us at <a href="mailto:admin@tn24.in">admin@tn24.in</a>.</p>
`)))
}

// HandleAboutUs serves the About Us page — required for Google AdSense approval.
func (h *PortalHandler) HandleAboutUs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(renderStaticPage("About Us | TN24 — Tamil Nadu News | தமிழ் செய்திகள்", "எங்களைப் பற்றி — About Us", `
<h2>About TN24</h2>
<p><strong>TN24</strong> (www.tn24.in) is Tamil Nadu's leading 24/7 digital news platform, dedicated to delivering accurate, fast, and comprehensive news coverage across all 38 districts of Tamil Nadu.</p>

<h3>Our Mission</h3>
<p>To provide unbiased, community-verified news in Tamil and English, empowering citizens with timely information about governance, civic issues, sports, entertainment, and cultural events.</p>

<h3>Key Features</h3>
<ul>
  <li><strong>38 District Coverage:</strong> Dedicated hyper-local feeds for Chennai, Coimbatore, Madurai, Salem, Trichy, Tirunelveli, and every district of Tamil Nadu.</li>
  <li><strong>24/7 Live Updates:</strong> Breaking news alerts, live tickers, and real-time updates as events unfold.</li>
  <li><strong>Editorial Integrity:</strong> Independent reporting with source attribution and factual verification.</li>
  <li><strong>Civic Focus:</strong> Highlighting citizen reports, public announcements, and administrative updates.</li>
</ul>

<h3>Editorial Team</h3>
<p>TN24 follows strict editorial guidelines. All news content is verified before publication. We are committed to accuracy, fairness, and transparency in all reporting.</p>

<h3>Contact &amp; Corporate Information</h3>
<p>Email: <a href="mailto:admin@tn24.in">admin@tn24.in</a></p>
<p>Website: <a href="https://www.tn24.in">www.tn24.in</a></p>
<p>Operating Region: Tamil Nadu, India</p>
`)))
}

// HandleContactUs serves the Contact page — required for Google AdSense approval.
func (h *PortalHandler) HandleContactUs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(renderStaticPage("Contact Us | TN24 — Tamil Nadu News", "தொடர்பு கொள்ளுங்கள் — Contact Us", `
<h2>Contact TN24</h2>
<p>In accordance with <strong>IT Rules 2021</strong>, you may submit a content grievance via our <a href="/portal">portal grievance form</a> or by writing to:</p>
<p><a href="mailto:admin@tn24.in">admin@tn24.in</a></p>
<p>We acknowledge grievances within <strong>24 hours</strong> and resolve them within <strong>15 days</strong> as required by law.</p>

<h3>🌐 Website</h3>
<p><a href="https://www.tn24.in">www.tn24.in</a></p>
`)))
}

// renderStaticPage renders a simple, SEO-friendly HTML page for Privacy Policy, About Us, Contact.
func renderStaticPage(title, heading, bodyHTML string) string {
	return `<!DOCTYPE html>
<html lang="ta">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>` + title + `</title>
<meta name="description" content="TN24 (tn 24) — Tamil Nadu news, latest tamil news &amp; தமிழ் செய்திகள் live 24x7. Breaking updates across 38 districts.">
<meta name="keywords" content="tn 24, news tamil 24x7 live, tamil nadu news, news live tamilnadu, தமிழ் செய்திகள், today news in tamil, news tamil today, tamil news online, latest tamil news, tamil nadu news in tamil, news tamil nadu">
<meta name="robots" content="index, follow">
<link rel="canonical" href="https://www.tn24.in">
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body { font-family: 'Segoe UI', Arial, sans-serif; background: #0f0f0f; color: #e0e0e0; line-height: 1.8; }
  header { background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%); padding: 16px 24px; display: flex; align-items: center; gap: 16px; border-bottom: 2px solid #e53e3e; }
  header a { text-decoration: none; color: #fff; font-size: 1.6rem; font-weight: 800; letter-spacing: -0.5px; }
  header span { color: #e53e3e; }
  nav { background: #161616; padding: 10px 24px; display: flex; gap: 20px; flex-wrap: wrap; }
  nav a { color: #aaa; text-decoration: none; font-size: 0.9rem; } nav a:hover { color: #fff; }
  .container { max-width: 860px; margin: 40px auto; padding: 0 24px 60px; }
  h1 { font-size: 1.5rem; color: #e53e3e; margin-bottom: 24px; border-bottom: 1px solid #333; padding-bottom: 12px; }
  h2 { font-size: 1.3rem; color: #fff; margin: 28px 0 12px; }
  h3 { font-size: 1.05rem; color: #ccc; margin: 20px 0 8px; }
  p { color: #bbb; margin-bottom: 12px; }
  ul { margin: 8px 0 12px 20px; color: #bbb; }
  li { margin-bottom: 6px; }
  a { color: #e53e3e; }
  footer { text-align: center; padding: 24px; background: #111; color: #555; font-size: 0.82rem; border-top: 1px solid #222; margin-top: 40px; }
</style>
<!-- Google AdSense -->
<script async src="https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=ca-pub-1894301748406603" crossorigin="anonymous"></script>
</head>
<body>
<header>
  <a href="/portal">TN<span>24</span></a>
  <span style="color:#aaa;font-size:0.85rem;">TN24 &mdash; Tamil Nadu News | தமிழ் செய்திகள் | 24x7 Live</span>
</header>
<nav>
  <a href="/portal">🏠 Home</a>
  <a href="/about">About</a>
  <a href="/privacy">Privacy</a>
  <a href="/contact">Contact</a>
</nav>
<div class="container">
  <h1>` + heading + `</h1>
  ` + bodyHTML + `
</div>
<footer>
  &copy; 2026 TN24 &mdash; www.tn24.in &nbsp;|&nbsp; <a href="/privacy">Privacy Policy</a> &nbsp;|&nbsp; <a href="/about">About Us</a> &nbsp;|&nbsp; <a href="/contact">Contact</a>
</footer>
</body>
</html>`
}

