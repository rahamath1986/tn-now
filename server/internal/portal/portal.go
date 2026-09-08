package portal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	mux.HandleFunc("/portal/assets/brand/", h.HandleBrandAssets)
	mux.HandleFunc("/portal/assets/", h.HandleBrandAssets)
	mux.HandleFunc("/assets/", h.HandleBrandAssets)
	mux.HandleFunc("/admin/assets/", h.HandleBrandAssets)

	// SEO Crawler & Indexing Endpoints
	mux.HandleFunc("/robots.txt", h.HandleRobotsTxt)
	mux.HandleFunc("/sitemap.xml", h.HandleSitemapXML)
}

func (h *PortalHandler) HandlePortalPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(RenderPortalPage()))
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

	// Portal displays ONLY APPROVED (PUBLISHED) content within the active 48-hour window
	whereClauses := []string{"c.status = 'PUBLISHED'", "c.created_at >= NOW() - interval '48 hours'"}
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
		       c.created_at, COALESCE(c.is_viral, false)
		FROM content c
		LEFT JOIN districts d ON c.district_id = d.id
		LEFT JOIN categories cat ON c.category_id = cat.id
		LEFT JOIN video_links vl ON c.id = vl.content_id
		LEFT JOIN stories s ON c.id = s.content_id
		LEFT JOIN photos p ON c.id = p.content_id
		WHERE %s
		ORDER BY c.is_viral DESC, COALESCE(c.updated_at, c.created_at) DESC
		LIMIT 60
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
				       c.created_at, COALESCE(c.is_viral, false)
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
		bestIdx := -1
		bestScore := -1

		for i, it := range allItems {
			score := 0
			hasRealImage := strings.TrimSpace(it.Thumbnail) != "" &&
				!strings.Contains(it.Thumbnail, "upload.wikimedia.org") &&
				!strings.Contains(it.Thumbnail, "/admin/api/maps/svg") &&
				!strings.Contains(it.Thumbnail, "dummy.svg")

			descLen := len(strings.TrimSpace(it.Description))

			if it.IsViral {
				score += 1000
			}
			if hasRealImage {
				score += 500
			}
			if descLen > 1000 {
				score += 400
			} else if descLen > 500 {
				score += 250
			} else if descLen > 200 {
				score += 100
			}

			if score > bestScore {
				bestScore = score
				bestIdx = i
			}
		}

		if bestIdx >= 0 {
			chosen := allItems[bestIdx]
			heroItem = &chosen
			allItems = append(allItems[:bestIdx], allItems[bestIdx+1:]...)
		} else {
			heroItem = &allItems[0]
			allItems = allItems[1:]
		}
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
			       c.created_at, COALESCE(c.is_viral, false)
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

		for i := 3; i < len(allItems) && i < 7; i++ {
			resp.LeftFeed = append(resp.LeftFeed, allItems[i])
		}

		for i := 7; i < len(allItems) && i < 12; i++ {
			resp.PressReleases = append(resp.PressReleases, allItems[i])
		}

		for i := 12; i < len(allItems); i++ {
			resp.CenterArticles = append(resp.CenterArticles, allItems[i])
		}
	}

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
			<text x="20" y="74" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="700" font-size="11.5" fill="#38bdf8">Email: tn24now@gmail.com  |  Call / WhatsApp: +91 81243 95082</text>
			<rect x="560" y="22" width="148" height="46" rx="6" fill="#1e293b" stroke="#38bdf8" stroke-width="1.5"/>
			<text x="634" y="42" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="11.5" fill="#38bdf8" text-anchor="middle">தொடர்புக்கு ↗</text>
			<text x="634" y="57" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="700" font-size="9.5" fill="#94a3b8" text-anchor="middle">8124395082</text>
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
			<text x="125" y="214" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="11" fill="#38bdf8" text-anchor="middle">தொடர்புக்கு: +91 81243 95082 ↗</text>
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
			<text x="100" y="170" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="9.5" fill="#38bdf8" text-anchor="middle">+91 81243 95082 ↗</text>
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
			<text x="20" y="74" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="700" font-size="11.5" fill="#c084fc">Email: tn24now@gmail.com  |  Call: +91 81243 95082</text>
			<rect x="560" y="22" width="148" height="46" rx="6" fill="#1e1b4b" stroke="#a855f7" stroke-width="1.5"/>
			<text x="634" y="42" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="800" font-size="11.5" fill="#c084fc" text-anchor="middle">தொடர்புக்கு ↗</text>
			<text x="634" y="57" font-family="-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif" font-weight="700" font-size="9.5" fill="#cbd5e1" text-anchor="middle">8124395082</text>
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

func (h *PortalHandler) HandleRobotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	robotsContent := `User-agent: *
Allow: /
Allow: /portal
Allow: /portal/*
Allow: /api/portal/feed
Disallow: /admin
Disallow: /admin/*
Disallow: /api/scraper/*
Disallow: /api/auth/*

Sitemap: http://localhost:8080/sitemap.xml
`
	_, _ = w.Write([]byte(robotsContent))
}

func (h *PortalHandler) HandleSitemapXML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	now := time.Now().Format("2006-01-02")
	sitemap := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <url>
        <loc>http://localhost:8080/portal</loc>
        <lastmod>%s</lastmod>
        <changefreq>always</changefreq>
        <priority>1.0</priority>
    </url>
    <url>
        <loc>http://localhost:8080/</loc>
        <lastmod>%s</lastmod>
        <changefreq>always</changefreq>
        <priority>0.9</priority>
    </url>
    <url>
        <loc>http://localhost:8080/portal?cat=News</loc>
        <lastmod>%s</lastmod>
        <changefreq>hourly</changefreq>
        <priority>0.8</priority>
    </url>
    <url>
        <loc>http://localhost:8080/portal?cat=Politics</loc>
        <lastmod>%s</lastmod>
        <changefreq>hourly</changefreq>
        <priority>0.8</priority>
    </url>
    <url>
        <loc>http://localhost:8080/portal?cat=Sports</loc>
        <lastmod>%s</lastmod>
        <changefreq>hourly</changefreq>
        <priority>0.8</priority>
    </url>
    <url>
        <loc>http://localhost:8080/portal?cat=Technical</loc>
        <lastmod>%s</lastmod>
        <changefreq>hourly</changefreq>
        <priority>0.8</priority>
    </url>
    <url>
        <loc>http://localhost:8080/portal?cat=Business</loc>
        <lastmod>%s</lastmod>
        <changefreq>hourly</changefreq>
        <priority>0.8</priority>
    </url>
    <url>
        <loc>http://localhost:8080/portal?cat=Entertainment</loc>
        <lastmod>%s</lastmod>
        <changefreq>hourly</changefreq>
        <priority>0.8</priority>
    </url>
    <url>
        <loc>http://localhost:8080/portal?district=Madurai</loc>
        <lastmod>%s</lastmod>
        <changefreq>hourly</changefreq>
        <priority>0.8</priority>
    </url>
    <url>
        <loc>http://localhost:8080/portal?district=Chennai</loc>
        <lastmod>%s</lastmod>
        <changefreq>hourly</changefreq>
        <priority>0.8</priority>
    </url>
    <url>
        <loc>http://localhost:8080/portal?district=Coimbatore</loc>
        <lastmod>%s</lastmod>
        <changefreq>hourly</changefreq>
        <priority>0.8</priority>
    </url>
</urlset>`, now, now, now, now, now, now, now, now, now, now, now)

	_, _ = w.Write([]byte(sitemap))
}

