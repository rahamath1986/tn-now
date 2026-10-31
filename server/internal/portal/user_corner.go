package portal

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"tn-now/server/internal/scraper"
)

// UserPostItem represents a user's submitted article in their dashboard
type UserPostItem struct {
	ID               string    `json:"id"`
	Slug             string    `json:"slug"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Status           string    `json:"status"`
	ModerationStatus string    `json:"moderationStatus"`
	District         string    `json:"district"`
	Category         string    `json:"category"`
	Thumbnail        string    `json:"thumbnail"`
	CreatedAt        time.Time `json:"createdAt"`
	PublishedAt      time.Time `json:"publishedAt"`
}

// RegisterUserCornerRoutes registers user corner and citizen journalism routes
func (h *PortalHandler) RegisterUserCornerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/submit", h.HandleSubmitPage)
	mux.HandleFunc("/login", h.HandleLoginPage)
	mux.HandleFunc("/register", h.HandleRegisterPage)
	mux.HandleFunc("/user/my-posts", h.HandleMyPostsPage)
	mux.HandleFunc("/user/profile", h.HandleProfilePage)
	mux.HandleFunc("/api/user/my-posts", h.HandleGetMyPostsAPI)
	mux.HandleFunc("/user-corner", h.HandleUserCornerPage)
	mux.HandleFunc("/author/", h.HandleAuthorPage)
}

// getUserFromRequest extracts authenticated user id and username from cookie or authorization header
func (h *PortalHandler) getUserFromRequest(r *http.Request) (string, string, string) {
	if h.conn == nil {
		return "", "", ""
	}

	token := ""
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	} else if cookie, err := r.Cookie("tn24_user_token"); err == nil && cookie.Value != "" {
		token = cookie.Value
	}

	if token == "" {
		return "", "", ""
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var uid pgtype.UUID
	var username, fullName string
	// Find user associated with an active non-revoked session
	err := h.conn.QueryRow(ctx, `
		SELECT u.id, u.username, COALESCE(up.full_name, u.username)
		FROM sessions s
		JOIN users u ON s.user_id = u.id
		LEFT JOIN user_profiles up ON u.id = up.user_id
		WHERE s.is_revoked = false AND (s.expires_at IS NULL OR s.expires_at > NOW())
		ORDER BY s.created_at DESC
		LIMIT 1
	`).Scan(&uid, &username, &fullName)

	if err != nil || !uid.Valid {
		return "", "", ""
	}

	return fmtUUID(uid), username, fullName
}

// HandleSubmitPage renders the front-facing article submission form with Quill WYSIWYG editor
func (h *PortalHandler) HandleSubmitPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.WriteHeader(http.StatusOK)

	userID, username, fullName := h.getUserFromRequest(r)

	pageHTML := renderFrontSubmitPageHTML(userID, username, fullName)
	_, _ = w.Write([]byte(pageHTML))
}

// HandleLoginPage renders the front-themed login page
func (h *PortalHandler) HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(renderFrontAuthPageHTML(true)))
}

// HandleRegisterPage renders the front-themed register page
func (h *PortalHandler) HandleRegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(renderFrontAuthPageHTML(false)))
}

// HandleMyPostsPage renders the author's submission tracking dashboard
func (h *PortalHandler) HandleMyPostsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, username, fullName := h.getUserFromRequest(r)
	if userID == "" {
		http.Redirect(w, r, "/login?next=/user/my-posts", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	pageHTML := renderFrontMyPostsPageHTML(userID, username, fullName)
	_, _ = w.Write([]byte(pageHTML))
}

// HandleProfilePage renders the author's profile editor
func (h *PortalHandler) HandleProfilePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, username, fullName := h.getUserFromRequest(r)
	if userID == "" {
		http.Redirect(w, r, "/login?next=/user/profile", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	pageHTML := renderFrontProfilePageHTML(userID, username, fullName)
	_, _ = w.Write([]byte(pageHTML))
}

// HandleGetMyPostsAPI returns the logged-in user's submissions as JSON
func (h *PortalHandler) HandleGetMyPostsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Method not allowed"})
		return
	}

	userID, _, _ := h.getUserFromRequest(r)
	if userID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Authentication required"})
		return
	}

	if h.conn == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "Database unavailable"})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := h.conn.Query(ctx, `
		SELECT c.id, c.title, COALESCE(c.description, ''), c.status, COALESCE(c.moderation_status, 'UNMODERATED'),
		       COALESCE(d.name, 'Tamil Nadu'), COALESCE(cat.name, 'User Corner'),
		       COALESCE(vl.thumbnail_url, s.single_photo_url, (p.photo_urls)[1], ''),
		       c.created_at, COALESCE(c.published_at, c.created_at), COALESCE(c.slug, '')
		FROM content c
		LEFT JOIN districts d ON c.district_id = d.id
		LEFT JOIN categories cat ON c.category_id = cat.id
		LEFT JOIN stories s ON c.id = s.content_id
		LEFT JOIN video_links vl ON c.id = vl.content_id
		LEFT JOIN photos p ON c.id = p.content_id
		WHERE c.author_user_id::text = $1
		ORDER BY c.created_at DESC
		LIMIT 100
	`, userID)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": err.Error()})
		return
	}
	defer rows.Close()

	var posts []UserPostItem
	for rows.Next() {
		var it UserPostItem
		var cid pgtype.UUID
		err := rows.Scan(&cid, &it.Title, &it.Description, &it.Status, &it.ModerationStatus,
			&it.District, &it.Category, &it.Thumbnail, &it.CreatedAt, &it.PublishedAt, &it.Slug)
		if err == nil {
			it.ID = fmtUUID(cid)
			it.Title = scraper.CleanHTML(it.Title)
			if it.Slug == "" {
				it.Slug = scraper.GenerateArticleSlug(it.Title, it.ID)
			}
			posts = append(posts, it)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    posts,
		"count":   len(posts),
	})
}

// HandleUserCornerPage renders the dedicated User Corner showcase page
func (h *PortalHandler) HandleUserCornerPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=60, s-maxage=120")
	w.WriteHeader(http.StatusOK)

	pageHTML := renderFrontUserCornerHubHTML()
	_, _ = w.Write([]byte(pageHTML))
}

// HandleAuthorPage renders a contributor's public E-E-A-T profile page
func (h *PortalHandler) HandleAuthorPage(w http.ResponseWriter, r *http.Request) {
	username := strings.TrimPrefix(r.URL.Path, "/author/")
	username = strings.Trim(username, "/")
	if username == "" {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	if h.conn == nil {
		http.NotFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var uid pgtype.UUID
	var dbUsername, fullName, bio, avatarUrl string
	var approvedCount, points, trustScore int
	err := h.conn.QueryRow(ctx, `
		SELECT u.id, u.username, COALESCE(up.full_name, u.username), COALESCE(up.bio, ''), COALESCE(up.avatar_url, ''),
		       COALESCE(cp.approved_count, 0), COALESCE(cp.points, 0), COALESCE(cp.trust_score, 10)
		FROM users u
		LEFT JOIN user_profiles up ON u.id = up.user_id
		LEFT JOIN contributor_profiles cp ON u.id = cp.user_id
		WHERE LOWER(u.username) = LOWER($1)
		LIMIT 1
	`, username).Scan(&uid, &dbUsername, &fullName, &bio, &avatarUrl, &approvedCount, &points, &trustScore)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Fetch author's published articles
	rows, err := h.conn.Query(ctx, `
		SELECT c.id, c.title, COALESCE(c.description, ''),
		       COALESCE(d.name, 'Tamil Nadu'), COALESCE(cat.name, 'News'),
		       COALESCE(vl.thumbnail_url, s.single_photo_url, (p.photo_urls)[1], ''),
		       COALESCE(c.published_at, c.created_at), COALESCE(c.slug, '')
		FROM content c
		LEFT JOIN districts d ON c.district_id = d.id
		LEFT JOIN categories cat ON c.category_id = cat.id
		LEFT JOIN stories s ON c.id = s.content_id
		LEFT JOIN video_links vl ON c.id = vl.content_id
		LEFT JOIN photos p ON c.id = p.content_id
		WHERE c.author_user_id = $1 AND c.status = 'PUBLISHED'
		ORDER BY COALESCE(c.published_at, c.created_at) DESC
		LIMIT 50
	`, uid)

	var articles []ArticleItem
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var it ArticleItem
			var cid pgtype.UUID
			if err := rows.Scan(&cid, &it.Title, &it.Description, &it.District, &it.Category, &it.Thumbnail, &it.CreatedAt, &it.Slug); err == nil {
				it.ID = fmtUUID(cid)
				it.Title = scraper.CleanHTML(it.Title)
				if it.Slug == "" {
					it.Slug = scraper.GenerateArticleSlug(it.Title, it.ID)
				}
				it.Author = fullName
				articles = append(articles, it)
			}
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=120, s-maxage=300")
	w.WriteHeader(http.StatusOK)

	pageHTML := renderAuthorProfilePageHTML(dbUsername, fullName, bio, avatarUrl, approvedCount, points, trustScore, articles)
	_, _ = w.Write([]byte(pageHTML))
}

// -----------------------------------------------------------------------------
// FRONT-THEMED HTML BUILDERS (Dark Navy #0d1421, matching Portal UI System)
// -----------------------------------------------------------------------------

func getCommonFrontHead(title, metaDesc, canonicalURL string) string {
	return `<!DOCTYPE html>
<html lang="ta">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + html.EscapeString(title) + `</title>
    <meta name="description" content="` + html.EscapeString(metaDesc) + `">
    <meta name="robots" content="index, follow">
    <link rel="canonical" href="` + html.EscapeString(canonicalURL) + `">
    <link rel="icon" type="image/svg+xml" href="/assets/brand/tn24-icon.svg">
    <link rel="shortcut icon" href="/favicon.ico">
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Mukta+Malar:wght@400;500;600;700;800&family=Inter:wght@400;500;600;700;800&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-base: #080f1c;
            --bg-surface: #0d1421;
            --bg-card: #111827;
            --border-color: #1e293b;
            --border-highlight: #334155;
            --primary-accent: #2563eb;
            --primary-hover: #1d4ed8;
            --brand-red: #ef4444;
            --brand-gold: #f59e0b;
            --text-primary: #f8fafc;
            --text-muted: #94a3b8;
            --font-family: 'Mukta Malar', 'Inter', system-ui, -apple-system, sans-serif;
        }
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: var(--font-family);
            background: var(--bg-base);
            color: var(--text-primary);
            line-height: 1.6;
            min-height: 100vh;
            display: flex;
            flex-direction: column;
        }
        a { color: inherit; text-decoration: none; }
        .front-header {
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 14px 24px;
            background: var(--bg-surface);
            border-bottom: 1px solid var(--border-color);
            position: sticky;
            top: 0;
            z-index: 50;
        }
        .brand-logo-img { height: 38px; width: auto; display: block; }
        .front-nav-pills { display: flex; align-items: center; gap: 12px; }
        .nav-btn {
            padding: 8px 16px;
            border-radius: 8px;
            font-size: 13px;
            font-weight: 700;
            display: inline-flex;
            align-items: center;
            gap: 6px;
            transition: all 0.2s ease;
            cursor: pointer;
            border: 1px solid transparent;
        }
        .nav-btn-outline {
            background: rgba(255,255,255,0.05);
            color: var(--text-primary);
            border-color: var(--border-color);
        }
        .nav-btn-outline:hover { background: rgba(255,255,255,0.1); border-color: var(--text-muted); }
        .nav-btn-primary {
            background: linear-gradient(135deg, #ef4444, #dc2626);
            color: #ffffff;
            box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
        }
        .nav-btn-primary:hover { transform: translateY(-1px); box-shadow: 0 6px 16px rgba(239, 68, 68, 0.4); }
        .front-container { max-width: 1000px; margin: 32px auto; padding: 0 20px; width: 100%; flex: 1; }
        .front-card {
            background: var(--bg-surface);
            border: 1px solid var(--border-color);
            border-radius: 14px;
            padding: 28px;
            box-shadow: 0 8px 30px rgba(0,0,0,0.3);
        }
        .form-label {
            display: block;
            font-size: 14px;
            font-weight: 700;
            margin-bottom: 8px;
            color: #e2e8f0;
        }
        .form-label span.req { color: var(--brand-red); margin-left: 4px; }
        .form-input, .form-select, .form-textarea {
            width: 100%;
            background: var(--bg-card);
            border: 1px solid var(--border-color);
            border-radius: 8px;
            padding: 12px 14px;
            font-size: 14px;
            color: var(--text-primary);
            font-family: inherit;
            outline: none;
            transition: border-color 0.2s;
        }
        .form-input:focus, .form-select:focus, .form-textarea:focus {
            border-color: var(--primary-accent);
            box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.2);
        }
        .form-row { margin-bottom: 20px; }
        .form-grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
        @media (max-width: 640px) { .form-grid-2 { grid-template-columns: 1fr; } }
        .status-badge {
            display: inline-flex;
            align-items: center;
            gap: 6px;
            padding: 4px 10px;
            border-radius: 6px;
            font-size: 12px;
            font-weight: 700;
        }
        .badge-pending { background: rgba(245, 158, 11, 0.15); color: #f59e0b; border: 1px solid rgba(245, 158, 11, 0.3); }
        .badge-published { background: rgba(34, 197, 94, 0.15); color: #22c55e; border: 1px solid rgba(34, 197, 94, 0.3); }
        .badge-rejected { background: rgba(239, 68, 68, 0.15); color: #ef4444; border: 1px solid rgba(239, 68, 68, 0.3); }
        footer {
            background: var(--bg-surface);
            border-top: 1px solid var(--border-color);
            padding: 24px;
            text-align: center;
            color: var(--text-muted);
            font-size: 13px;
            margin-top: 48px;
        }
    </style>`
}

func getCommonFrontHeader(activePage, userID, fullName string) string {
	userBtn := `<a href="/login" class="nav-btn nav-btn-outline">உள்நுழைக / பதிவு</a>`
	if userID != "" {
		displayName := fullName
		if displayName == "" {
			displayName = "எனது பக்கம்"
		}
		userBtn = fmt.Sprintf(`
			<a href="/user/profile" class="nav-btn nav-btn-outline" title="சுயவிவரத்தை மாற்றுக">👤 %s</a>
			<a href="/user/my-posts" class="nav-btn nav-btn-outline" title="எனது கட்டுரைகள்">📋 பதிவுகள்</a>
			<button onclick="handleLogout()" class="nav-btn nav-btn-outline" style="color:#ef4444;" title="வெளியேறு">🚪</button>
		`, html.EscapeString(displayName))
	}

	submitBtn := `<a href="/submit" class="nav-btn nav-btn-primary">✍️ செய்தி எழுதுக</a>`
	if activePage == "submit" {
		submitBtn = `<a href="/user-corner" class="nav-btn nav-btn-outline">📰 வாசகர் அரங்கம்</a>`
	}

	return `
    <header class="front-header">
        <div style="display:flex;align-items:center;gap:16px;">
            <a href="/" title="TN24 முகப்பு">
                <img src="/assets/brand/tn24-logo.svg" alt="TN24 Logo" class="brand-logo-img">
            </a>
            <span style="font-size:12px;color:#94a3b8;border-left:1px solid #1e293b;padding-left:12px;" class="hide-mobile">
                தமிழ்நாடு நம்பகமான செய்திகள்
            </span>
        </div>
        <div class="front-nav-pills">
            <a href="/user-corner" class="nav-btn nav-btn-outline" style="color:#f59e0b;">✍️ வாசகர் அரங்கம்</a>
            ` + submitBtn + `
            ` + userBtn + `
        </div>
    </header>`
}

// renderFrontSubmitPageHTML produces the interactive Quill WYSIWYG article submission form
func renderFrontSubmitPageHTML(userID, username, fullName string) string {
	head := getCommonFrontHead(
		"செய்தி சமர்ப்பிக்கவும் | வாசகர் அரங்கம் — TN24",
		"உங்கள் பகுதி செய்திகள், கட்டுரைகள் மற்றும் பொதுமக்களின் பிரச்சனைகளை வாசகர் அரங்கில் பகிருங்கள்.",
		"https://www.tn24.in/submit",
	)

	return head + `
    <!-- Quill WYSIWYG CSS -->
    <link href="https://cdn.quilljs.com/1.3.6/quill.snow.css" rel="stylesheet">
    <style>
        .ql-toolbar.ql-snow {
            background: #111827;
            border-color: #1e293b !important;
            border-top-left-radius: 8px;
            border-top-right-radius: 8px;
        }
        .ql-container.ql-snow {
            background: #0b1120;
            border-color: #1e293b !important;
            border-bottom-left-radius: 8px;
            border-bottom-right-radius: 8px;
            color: #f8fafc;
            font-family: var(--font-family);
            font-size: 15px;
            min-height: 320px;
        }
        .ql-snow .ql-stroke { stroke: #94a3b8; }
        .ql-snow .ql-fill { fill: #94a3b8; }
        .ql-snow .ql-picker { color: #94a3b8; }
        .ql-editor.ql-blank::before { color: #64748b; font-style: normal; }
        .image-preview-box {
            margin-top: 10px;
            border-radius: 8px;
            overflow: hidden;
            display: none;
            max-height: 280px;
            border: 1px solid #1e293b;
            background: #000;
        }
        .image-preview-box img { width: 100%; height: 280px; object-fit: cover; display: block; }
        .word-meter {
            display: flex;
            align-items: center;
            justify-content: space-between;
            font-size: 12px;
            color: #94a3b8;
            margin-top: 8px;
            padding: 4px 8px;
            background: #111827;
            border-radius: 6px;
        }
        .modal-overlay {
            position: fixed;
            inset: 0;
            background: rgba(0,0,0,0.85);
            backdrop-filter: blur(4px);
            display: none;
            align-items: center;
            justify-content: center;
            z-index: 100;
            padding: 20px;
        }
        .modal-card {
            background: #0d1421;
            border: 1px solid #1e293b;
            border-radius: 16px;
            padding: 32px;
            max-width: 520px;
            width: 100%;
            text-align: center;
            box-shadow: 0 20px 40px rgba(0,0,0,0.6);
        }
    </style>
</head>
<body>
    ` + getCommonFrontHeader("submit", userID, fullName) + `

    <main class="front-container">
        <!-- Breadcrumb & Header -->
        <div style="margin-bottom: 24px;">
            <div style="display:flex;align-items:center;gap:8px;font-size:13px;color:#94a3b8;margin-bottom:8px;">
                <a href="/">முகப்பு</a> <span>&rsaquo;</span> <a href="/user-corner">வாசகர் அரங்கம்</a> <span>&rsaquo;</span> <span style="color:#f8fafc;">செய்தி சமர்ப்பிக்க</span>
            </div>
            <h1 style="font-size: 26px; font-weight: 800; color: #f8fafc; letter-spacing: -0.5px;">
                ✍️ வாசகர் அரங்கம் — உங்கள் செய்தியை சமர்ப்பிக்கவும்
            </h1>
            <p style="color: #94a3b8; font-size: 14px; margin-top: 4px;">
                உங்கள் பகுதியில் நிகழ்ந்த உண்மைச் செய்திகள், மக்கள் பிரச்சனைகள் மற்றும் சமூகக் கட்டுரைகளை TN24 வாசகர் அரங்கில் பகிருங்கள். ஆசிரியர் சரிபார்ப்புக்குப் பின் உடனடியாக வெளியாகும்.
            </p>
        </div>

        <div class="front-card">
            <form id="storySubmitForm" onsubmit="handleStorySubmit(event)">
                <div class="form-grid-2 form-row">
                    <!-- District / Region Selection -->
                    <div>
                        <label class="form-label" for="districtSelect">மாவட்டம் / பகுதி <span class="req">*</span></label>
                        <select id="districtSelect" class="form-select" required>
                            <option value="Tamil Nadu">அனைத்து தமிழ்நாடு (Statewide / Regional)</option>
                            <option value="Chennai">சென்னை (Chennai)</option>
                            <option value="Coimbatore">கோயம்புத்தூர் (Coimbatore)</option>
                            <option value="Madurai">மதுரை (Madurai)</option>
                            <option value="Tiruchirappalli">திருச்சிராப்பள்ளி (Trichy)</option>
                            <option value="Salem">சேலம் (Salem)</option>
                            <option value="Tirunelveli">திருநெல்வேலி (Tirunelveli)</option>
                            <option value="Tiruppur">திருப்பூர் (Tiruppur)</option>
                            <option value="Erode">ஈரோடு (Erode)</option>
                            <option value="Vellore">வேலூர் (Vellore)</option>
                            <option value="Thanjavur">தஞ்சாவூர் (Thanjavur)</option>
                            <option value="Dindigul">திண்டுக்கல் (Dindigul)</option>
                            <option value="Kanyakumari">கன்னியாகுமரி (Kanyakumari)</option>
                            <option value="Thoothukudi">தூத்துக்குடி (Thoothukudi)</option>
                            <option value="Ranipet">ராணிப்பேட்டை (Ranipet)</option>
                            <option value="Virudhunagar">விருதுநகர் (Virudhunagar)</option>
                            <option value="Sivagangai">சிவகங்கை (Sivagangai)</option>
                            <option value="Ramanathapuram">ராமநாதபுரம் (Ramanathapuram)</option>
                            <option value="Theni">தேனி (Theni)</option>
                            <option value="Karur">கரூர் (Karur)</option>
                            <option value="Namakkal">நாமக்கல் (Namakkal)</option>
                            <option value="Nilgiris">நீலகிரி (Nilgiris)</option>
                            <option value="Cuddalore">கடலூர் (Cuddalore)</option>
                            <option value="Villupuram">விழுப்புரம் (Villupuram)</option>
                            <option value="Kallakurichi">கள்ளக்குறிச்சி (Kallakurichi)</option>
                            <option value="Chengalpattu">செங்கல்பட்டு (Chengalpattu)</option>
                            <option value="Kanchipuram">காஞ்சிபுரம் (Kanchipuram)</option>
                            <option value="Tiruvallur">திருவள்ளூர் (Tiruvallur)</option>
                            <option value="Tiruvannamalai">திருவண்ணாமலை (Tiruvannamalai)</option>
                            <option value="Dharmapuri">தருமபுரி (Dharmapuri)</option>
                            <option value="Krishnagiri">கிருஷ்ணகிரி (Krishnagiri)</option>
                            <option value="Pudukkottai">புதுக்கோட்டை (Pudukkottai)</option>
                            <option value="Ariyalur">அரியலூர் (Ariyalur)</option>
                            <option value="Perambalur">பெரம்பலூர் (Perambalur)</option>
                            <option value="Nagapattinam">நாகப்பட்டினம் (Nagapattinam)</option>
                            <option value="Tiruvarur">திருவாரூர் (Tiruvarur)</option>
                            <option value="Mayiladuthurai">மயிலாடுதுறை (Mayiladuthurai)</option>
                            <option value="Tenkasi">தென்காசி (Tenkasi)</option>
                            <option value="Tirupathur">திருப்பத்தூர் (Tirupathur)</option>
                        </select>
                    </div>

                    <!-- Category Selection -->
                    <div>
                        <label class="form-label" for="categorySelect">பிரிவு (Category) <span class="req">*</span></label>
                        <select id="categorySelect" class="form-select" required>
                            <option value="User Corner" selected>✍️ வாசகர் அரங்கம் (Citizen Journalism / Opinion)</option>
                            <option value="News">செய்திகள் & பொது நிகழ்வுகள் (Local News)</option>
                            <option value="Public Updates">பொதுப் பிரச்சனைகள் & கோரிக்கைகள் (Public Grievances)</option>
                            <option value="Politics">அரசியல் & கொள்கை (Politics)</option>
                            <option value="Sports">விளையாட்டு (Sports)</option>
                            <option value="Technical">தொழில்நுட்பம் (Tech & Innovation)</option>
                            <option value="Business">வணிகம் & விவசாயம் (Business & Agri)</option>
                            <option value="Entertainment">கலை & பண்பாடு (Culture & Arts)</option>
                        </select>
                    </div>
                </div>

                <!-- Headline / Title -->
                <div class="form-row">
                    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:8px;">
                        <label class="form-label" for="titleInput" style="margin-bottom:0;">செய்தியின் தலைப்பு (Headline) <span class="req">*</span></label>
                        <span id="titleCounter" style="font-size:12px;color:#94a3b8;">0 / 120 எழுத்துக்கள்</span>
                    </div>
                    <input type="text" id="titleInput" class="form-input" placeholder="கவனத்தை ஈர்க்கும் துல்லியமான தலைப்பை உள்ளிடவும்..." maxlength="150" required oninput="updateTitleCount()">
                </div>

                <!-- Main Cover Image URL -->
                <div class="form-row">
                    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:8px;">
                        <label class="form-label" for="imageInput" style="margin-bottom:0;">முதன்மை புகைப்பட இணைப்பு (Main Image)</label>
                        <span style="font-size:12px;color:#94a3b8;">ImageKit CDN மூலம் உகந்ததாக்கப்படும்</span>
                    </div>
                    <div style="display:flex;gap:10px;align-items:center;flex-wrap:wrap;margin-bottom:10px;">
                        <input type="url" id="imageInput" class="form-input" style="flex:1;min-width:260px;" placeholder="https://example.com/photo.jpg அல்லது சாதனத்திலிருந்து பதிவேற்றுக..." oninput="handleImagePreview()">
                        <input type="file" id="imageFileInput" accept="image/jpeg,image/png,image/webp,image/gif" style="display:none;" onchange="handleImageFileUpload(event)">
                        <button type="button" id="uploadFileBtn" onclick="document.getElementById('imageFileInput').click()" class="nav-btn nav-btn-outline" style="white-space:nowrap;display:inline-flex;align-items:center;gap:6px;border-color:#38bdf8;color:#38bdf8;padding:8px 14px;">
                            <span id="uploadBtnIcon">📷</span> <span id="uploadBtnText">சாதனத்திலிருந்து பதிவேற்று</span>
                        </button>
                    </div>
                    <div id="uploadStatusMsg" style="display:none;font-size:12px;padding:8px 12px;border-radius:6px;margin-bottom:10px;"></div>
                    <div id="imagePreviewBox" class="image-preview-box">
                        <img id="imagePreview" src="" alt="Main Article Image Preview">
                    </div>
                </div>

                <!-- WYSIWYG Editor Body -->
                <div class="form-row">
                    <label class="form-label">முழு செய்தி கட்டுரை (Article Content - WYSIWYG Editor) <span class="req">*</span></label>
                    <div id="editor-container"></div>
                    <div class="word-meter">
                        <span id="wordCountLabel">சொற்கள் எண்ணிக்கை: 0 சொற்கள் (குறைந்தது 50 சொற்கள் பரிந்துரைக்கப்படுகிறது)</span>
                        <span id="readingTimeLabel">⏱️ வாசிக்கும் நேரம்: ~1 நிமிடம்</span>
                    </div>
                </div>

                <!-- Excerpt / Short Summary -->
                <div class="form-row">
                    <label class="form-label" for="descInput">குறுஞ்செய்தி சுருக்கம் (Brief Summary / Excerpt)</label>
                    <textarea id="descInput" class="form-textarea" rows="2" placeholder="கட்டுரையின் முக்கிய சிறப்பம்சங்களை 2 வாக்கியங்களில் சுருக்கமாகத் தருக (கூகுள் தேடல் மற்றும் பகிர்வுக்கு உதவும்)..."></textarea>
                </div>

                <!-- Author & Language Row -->
                <div class="form-grid-2 form-row">
                    <div>
                        <label class="form-label" for="authorNameInput">எழுதியவர் பெயர் (Author Byline)</label>
                        <input type="text" id="authorNameInput" class="form-input" value="` + html.EscapeString(fullName) + `" placeholder="உங்கள் பெயர்">
                    </div>
                    <div>
                        <label class="form-label" for="langSelect">மொழி (Language)</label>
                        <select id="langSelect" class="form-select">
                            <option value="ta" selected>தமிழ் (Tamil)</option>
                            <option value="en">English</option>
                        </select>
                    </div>
                </div>

                <!-- Guideline Agreement -->
                <div class="form-row" style="background:#111827;border:1px solid #1e293b;border-radius:8px;padding:12px 14px;">
                    <label style="display:flex;align-items:flex-start;gap:10px;cursor:pointer;font-size:13px;color:#cbd5e1;">
                        <input type="checkbox" id="termsAgree" required style="margin-top:3px;">
                        <span>
                            இச்செய்தி உண்மையானது என்பதையும், வன்முறை, போலி தகவல் மற்றும் காப்புரிமை மீறல் இல்லாதது என்பதையும் உறுதிப்படுத்துகிறேன். TN24 ஆசிரியர் குழுவின் நெறிமுறைகளுக்கு ஒப்புக்கொள்கிறேன்.
                        </span>
                    </label>
                </div>

                <!-- Submit Button -->
                <div style="display:flex;gap:12px;align-items:center;margin-top:24px;">
                    <button type="submit" id="submitBtn" class="nav-btn nav-btn-primary" style="padding:12px 28px;font-size:15px;border-radius:10px;">
                        🚀 சமர்ப்பிக்கவும் (Submit for Review)
                    </button>
                    <span id="submitSpinner" style="display:none;color:#94a3b8;font-size:13px;">🔄 செய்தி சமர்ப்பிக்கப்படுகிறது...</span>
                </div>
            </form>
        </div>
    </main>

    <!-- Success Modal ("Under Review" Status Notice) -->
    <div id="successModal" class="modal-overlay">
        <div class="modal-card">
            <div style="font-size: 48px; margin-bottom: 12px;">🟡</div>
            <h2 style="font-size: 22px; font-weight: 800; color: #f59e0b; margin-bottom: 8px;">
                பதிவு பெறப்பட்டது — பரிசீலனையில் உள்ளது!
            </h2>
            <p style="color: #94a3b8; font-size: 14px; line-height: 1.6; margin-bottom: 24px;">
                நன்றி! உங்கள் கட்டுரை TN24 ஆசிரியர் குழுவின் நேரடி பரிசீலனைக்கு (<strong>Under Review</strong>) அனுப்பி வைக்கப்பட்டுள்ளது. எடிட்டர்கள் சரிபார்த்து ஒப்புதல் அளித்தவுடன் <strong>வாசகர் அரங்கம்</strong> பகுதியில் நேரடியாக வெளியாகும்.
            </p>
            <div style="display:flex;gap:12px;justify-content:center;">
                <a href="/user/my-posts" class="nav-btn nav-btn-primary" style="padding:10px 20px;">
                    📋 எனது பதிவுகளைக் காண்க
                </a>
                <a href="/" class="nav-btn nav-btn-outline" style="padding:10px 20px;">
                    முகப்பு பக்கம்
                </a>
            </div>
        </div>
    </div>

    <footer>
        &copy; 2026 TN24 — Tamil Nadu's Independent Digital News Platform. வாசகர் அரங்கம்.
    </footer>

    <!-- Quill WYSIWYG JS -->
    <script src="https://cdn.quilljs.com/1.3.6/quill.min.js"></script>
    <script>
        const quill = new Quill('#editor-container', {
            theme: 'snow',
            placeholder: 'உங்கள் செய்தியின் விரிவான உரை, கள நிலவரம், மேற்கோள்கள் மற்றும் தகவல்களை இங்கே தட்டச்சு செய்யவும்...',
            modules: {
                toolbar: [
                    [{ 'header': [2, 3, false] }],
                    ['bold', 'italic', 'underline', 'strike'],
                    ['blockquote'],
                    [{ 'list': 'ordered'}, { 'list': 'bullet' }],
                    ['link', 'clean']
                ]
            }
        });

        // Live Word Counter
        quill.on('text-change', function() {
            const text = quill.getText().trim();
            const words = text ? text.split(/\s+/).length : 0;
            const wLabel = document.getElementById('wordCountLabel');
            if (wLabel) {
                wLabel.textContent = 'சொற்கள் எண்ணிக்கை: ' + words + ' சொற்கள் ' + (words >= 50 ? '✅' : '(குறைந்தது 50 சொற்கள் பரிந்துரைக்கப்படுகிறது)');
                wLabel.style.color = words >= 50 ? '#22c55e' : '#94a3b8';
            }
            const rLabel = document.getElementById('readingTimeLabel');
            if (rLabel) {
                const mins = Math.max(1, Math.ceil(words / 180));
                rLabel.textContent = '⏱️ வாசிக்கும் நேரம்: ~' + mins + ' நிமிடம்';
            }
        });

        function updateTitleCount() {
            const input = document.getElementById('titleInput');
            const counter = document.getElementById('titleCounter');
            counter.textContent = input.value.length + ' / 120 எழுத்துக்கள்';
        }

        function handleImagePreview() {
            const url = document.getElementById('imageInput').value.trim();
            const box = document.getElementById('imagePreviewBox');
            const img = document.getElementById('imagePreview');
            if (url && (url.startsWith('http://') || url.startsWith('https://'))) {
                img.src = url;
                img.onload = () => { box.style.display = 'block'; };
                img.onerror = () => { box.style.display = 'none'; };
            } else {
                box.style.display = 'none';
            }
        }

        async function handleImageFileUpload(e) {
            const file = e.target.files && e.target.files[0];
            if (!file) return;

            if (!file.type.startsWith('image/')) {
                alert('தயவுசெய்து சரியான படக் கோப்பைத் தேர்ந்தெடுக்கவும் (JPG, PNG, WebP, GIF).');
                return;
            }

            if (file.size > 10 * 1024 * 1024) {
                alert('படத்தின் அளவு 10MB-க்கு கீழ் இருக்க வேண்டும்.');
                return;
            }

            const statusMsg = document.getElementById('uploadStatusMsg');
            const uploadBtn = document.getElementById('uploadFileBtn');
            const uploadBtnText = document.getElementById('uploadBtnText');
            const uploadBtnIcon = document.getElementById('uploadBtnIcon');

            statusMsg.style.display = 'block';
            statusMsg.style.background = 'rgba(56, 189, 248, 0.1)';
            statusMsg.style.color = '#38bdf8';
            statusMsg.style.border = '1px solid rgba(56, 189, 248, 0.2)';
            statusMsg.innerHTML = '⏳ படம் ImageKit CDN-ல் பதிவேற்றப்படுகிறது, காத்திருக்கவும்...';

            uploadBtn.disabled = true;
            uploadBtnText.textContent = 'பதிவேற்றுகிறது...';
            uploadBtnIcon.textContent = '⏳';

            const formData = new FormData();
            formData.append('image', file);

            try {
                const res = await fetch('/content/upload/image', {
                    method: 'POST',
                    body: formData
                });
                const data = await res.json().catch(() => ({}));
                if (!res.ok || !data.success) {
                    throw new Error(data.error || 'பதிவேற்றம் தோல்வியடைந்தது. ImageKit அமைப்புகளைச் சரிபார்க்கவும்.');
                }

                document.getElementById('imageInput').value = data.url;
                handleImagePreview();

                statusMsg.style.background = 'rgba(34, 197, 94, 0.1)';
                statusMsg.style.color = '#22c55e';
                statusMsg.style.border = '1px solid rgba(34, 197, 94, 0.2)';
                statusMsg.innerHTML = '✓ படம் வெற்றிகரமாக ImageKit CDN-ல் பதிவேற்றப்பட்டது!';
            } catch (err) {
                statusMsg.style.background = 'rgba(239, 68, 68, 0.1)';
                statusMsg.style.color = '#ef4444';
                statusMsg.style.border = '1px solid rgba(239, 68, 68, 0.2)';
                statusMsg.innerHTML = '⚠️ ' + (err.message || 'பதிவேற்றம் தோல்வியடைந்தது');
            } finally {
                uploadBtn.disabled = false;
                uploadBtnText.textContent = 'சாதனத்திலிருந்து பதிவேற்று';
                uploadBtnIcon.textContent = '📷';
                e.target.value = '';
            }
        }

        async function handleStorySubmit(e) {
            e.preventDefault();
            const title = document.getElementById('titleInput').value.trim();
            const bodyHtml = quill.root.innerHTML;
            const plainText = quill.getText().trim();

            if (!title) {
                alert('செய்தியின் தலைப்பை உள்ளிடவும்.');
                return;
            }
            if (!plainText || plainText.length < 30) {
                alert('கட்டுரையின் உரை மிகவும் சிறியதாக உள்ளது. விரிவாக எழுதவும்.');
                return;
            }

            const district = document.getElementById('districtSelect').value;
            const category = document.getElementById('categorySelect').value;
            const imageUrl = document.getElementById('imageInput').value.trim();
            const desc = document.getElementById('descInput').value.trim();
            const author = document.getElementById('authorNameInput').value.trim();
            const lang = document.getElementById('langSelect').value;

            const submitBtn = document.getElementById('submitBtn');
            const spinner = document.getElementById('submitSpinner');
            submitBtn.disabled = true;
            spinner.style.display = 'inline-block';

            try {
                const res = await fetch('/content/submit/story', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        title: title,
                        headline: title,
                        body: bodyHtml,
                        description: desc,
                        districtId: district,
                        categoryId: category,
                        singlePhotoUrl: imageUrl,
                        authorName: author,
                        language: lang
                    })
                });
                const data = await res.json();
                if (data.success) {
                    document.getElementById('successModal').style.display = 'flex';
                } else {
                    alert('பிழை: ' + (data.message || 'சமர்ப்பிக்க முடியவில்லை'));
                }
            } catch (err) {
                alert('தொடர்பு பிழை: ' + err.message);
            } finally {
                submitBtn.disabled = false;
                spinner.style.display = 'none';
            }
        }

        function handleLogout() {
            fetch('/auth/logout', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: '{}' })
                .then(() => { window.location.reload(); });
        }
    </script>
</body>
</html>`
}

// renderFrontAuthPageHTML produces clean front-themed login/register form supporting traditional registration, Email OTP verification, and password login
func renderFrontAuthPageHTML(isLogin bool) string {
	title := "உள்நுழைக / பதிவு செய்க | TN24 வாசகர் அரங்கம்"
	if !isLogin {
		title = "பதிவு செய்க | TN24 வாசகர் அரங்கம்"
	}

	head := getCommonFrontHead(title, "TN24 செய்திக் களத்தில் உங்கள் கணக்கை நிர்வகிக்கவும் மற்றும் கட்டுரைகளை சமர்ப்பிக்கவும்.", "https://www.tn24.in/login")

	return head + `
    <style>
        .auth-tab-btn {
            flex: 1;
            padding: 12px;
            background: none;
            border: none;
            font-size: 14px;
            font-weight: 700;
            cursor: pointer;
            color: #94a3b8;
            border-bottom: 2px solid transparent;
            font-family: inherit;
            transition: all 0.2s;
        }
        .auth-tab-btn.active {
            color: #ef4444;
            border-bottom-color: #ef4444;
        }
        .otp-box {
            font-size: 26px;
            letter-spacing: 10px;
            text-align: center;
            font-family: monospace;
            font-weight: 800;
            color: #38bdf8;
            background: #0b1120;
        }
    </style>
</head>
<body>
    <header class="front-header">
        <a href="/" title="TN24 முகப்பு">
            <img src="/assets/brand/tn24-logo.svg" alt="TN24 Logo" class="brand-logo-img">
        </a>
        <a href="/" class="nav-btn nav-btn-outline">← முகப்பு திரும்புக</a>
    </header>

    <main class="front-container" style="max-width: 460px; margin-top: 40px; margin-bottom: 60px;">
        <div class="front-card">
            <div style="text-align:center;margin-bottom:20px;">
                <div style="font-size:32px;margin-bottom:6px;">✍️</div>
                <h1 style="font-size:22px;font-weight:800;color:#f8fafc;letter-spacing:-0.5px;">TN24 வாசகர் அரங்கம்</h1>
                <p style="color:#94a3b8;font-size:13px;margin-top:4px;">செய்திகள் எழுத மற்றும் உங்கள் பங்களிப்புகளை நிர்வகிக்க உள்நுழையவும்</p>
            </div>

            <!-- Primary Tabs: Login vs Register -->
            <div id="authTabsContainer" style="display:flex;border-bottom:1px solid #1e293b;margin-bottom:20px;">
                <button type="button" id="tabLogin" onclick="switchAuthTab(true)" class="auth-tab-btn ` + (func() string {
		if isLogin {
			return "active"
		}
		return ""
	}()) + `">உள்நுழைக (Sign In)</button>
                <button type="button" id="tabRegister" onclick="switchAuthTab(false)" class="auth-tab-btn ` + (func() string {
		if !isLogin {
			return "active"
		}
		return ""
	}()) + `">பதிவு செய்க (Sign Up)</button>
            </div>

            <!-- 1. Traditional Password Login Form -->
            <form id="loginForm" style="display:` + (func() string {
		if isLogin {
			return "block"
		}
		return "none"
	}()) + `;" onsubmit="handleLoginSubmit(event)">
                <div class="form-row">
                    <label class="form-label" for="loginEmail">மின்னஞ்சல் (Email) <span class="req">*</span></label>
                    <input type="email" id="loginEmail" class="form-input" placeholder="name@example.com" required autocomplete="email">
                </div>
                <div class="form-row">
                    <label class="form-label" for="loginPassword">கடவுச்சொல் (Password) <span class="req">*</span></label>
                    <input type="password" id="loginPassword" class="form-input" placeholder="••••••••" required autocomplete="current-password">
                </div>
                <button type="submit" id="loginBtn" class="nav-btn nav-btn-primary" style="width:100%;justify-content:center;padding:12px;font-size:15px;margin-top:8px;">
                    உள்நுழைக
                </button>
                <div style="margin-top:16px;text-align:center;">
                    <button type="button" onclick="showQuickOtpLogin()" style="background:none;border:none;color:#38bdf8;font-size:12px;cursor:pointer;text-decoration:underline;">
                        ⚡ கடவுச்சொல் மறந்துவிட்டதா? மின்னஞ்சல் OTP மூலம் உள்நுழைக
                    </button>
                </div>
            </form>

            <!-- 2. Traditional Registration Form -->
            <form id="registerForm" style="display:` + (func() string {
		if !isLogin {
			return "block"
		}
		return "none"
	}()) + `;" onsubmit="handleRegisterSubmit(event)">
                <div class="form-row">
                    <label class="form-label" for="regFullName">முழுப் பெயர் / Display Name <span class="req">*</span></label>
                    <input type="text" id="regFullName" class="form-input" placeholder="எ.கா: கார்த்திகேயன்" required>
                </div>
                <div class="form-row">
                    <label class="form-label" for="regUsername">பயனர் பெயர் (Username) <span class="req">*</span></label>
                    <input type="text" id="regUsername" class="form-input" placeholder="karthik_tn" pattern="[a-zA-Z0-9_]{3,30}" title="எழுத்துக்கள், எண்கள் மற்றும் underscore மட்டுமே" required autocomplete="username">
                </div>
                <div class="form-row">
                    <label class="form-label" for="regEmail">மின்னஞ்சல் (Email) <span class="req">*</span></label>
                    <input type="email" id="regEmail" class="form-input" placeholder="karthik@example.com" required autocomplete="email">
                    <p style="color:#64748b;font-size:12px;margin-top:4px;">💡 பதிவிற்குப் பின் உங்கள் மின்னஞ்சலுக்கு ஒருமுறை சரிபார்ப்புக் குறியீடு (OTP) அனுப்பப்படும்.</p>
                </div>
                <div class="form-row">
                    <label class="form-label" for="regPassword">கடவுச்சொல் (Password) <span class="req">*</span></label>
                    <input type="password" id="regPassword" class="form-input" placeholder="குறைந்தது 6 எழுத்துக்கள்" minlength="6" required autocomplete="new-password">
                </div>
                <button type="submit" id="regBtn" class="nav-btn nav-btn-primary" style="width:100%;justify-content:center;padding:12px;font-size:15px;margin-top:8px;">
                    🚀 கணக்கு துவங்கி OTP பெறுக
                </button>
                <div style="margin-top:14px;text-align:center;font-size:13px;color:#94a3b8;">
                    ஏற்கனவே கணக்கு உள்ளதா? <a href="javascript:void(0)" onclick="switchAuthTab(true)" style="color:#ef4444;font-weight:600;text-decoration:underline;">உள்நுழைக</a>
                </div>
            </form>

            <!-- 3. OTP Verification Step (Shown after Registration OR for unverified login) -->
            <form id="otpVerifyForm" style="display:none;" onsubmit="handleVerifyOTP(event)">
                <div style="background:#111827;border:1px solid #1e293b;border-radius:8px;padding:12px;margin-bottom:16px;display:flex;justify-content:space-between;align-items:center;">
                    <div>
                        <div style="font-size:11px;color:#94a3b8;">சரிபார்ப்பு மின்னஞ்சல்:</div>
                        <div id="otpEmailTarget" style="font-size:13px;font-weight:600;color:#f8fafc;"></div>
                    </div>
                    <button type="button" onclick="cancelOtpStep()" class="nav-btn nav-btn-outline" style="font-size:11px;padding:4px 8px;">
                        மாற்றுக ✎
                    </button>
                </div>

                <div class="form-row">
                    <label class="form-label" for="otpInput">மின்னஞ்சல் OTP சரிபார்ப்புக் குறியீடு (6-Digit OTP) <span class="req">*</span></label>
                    <input type="text" id="otpInput" class="form-input otp-box" placeholder="••••••" maxlength="6" pattern="[0-9]{6}" inputmode="numeric" autocomplete="one-time-code" required>
                    <p style="color:#94a3b8;font-size:12px;margin-top:6px;text-align:center;">
                        உங்கள் மின்னஞ்சல் இன்பாக்ஸ் அல்லது ஸ்பேம் ஃபோல்டரில் உள்ள 6-இலக்கக் குறியீட்டை உள்ளிடவும்.
                    </p>
                </div>

                <button type="submit" id="verifyOtpBtn" class="nav-btn nav-btn-primary" style="width:100%;justify-content:center;padding:12px;font-size:15px;margin-top:6px;">
                    <span>🚀 சரிபார்த்து உள்நுழைக</span>
                </button>

                <div style="margin-top:16px;text-align:center;">
                    <span id="resendTimerText" style="color:#64748b;font-size:12px;">மீண்டும் OTP அனுப்ப: <span id="resendCountdown">60</span>s</span>
                    <button type="button" id="resendOtpBtn" onclick="handleResendOTP()" style="display:none;background:none;border:none;color:#ef4444;font-size:13px;cursor:pointer;font-weight:600;text-decoration:underline;">
                        மீண்டும் OTP அனுப்புக ↺
                    </button>
                </div>
            </form>

            <!-- 4. Quick Passwordless OTP Request Form -->
            <form id="quickOtpForm" style="display:none;" onsubmit="handleQuickOtpSend(event)">
                <div class="form-row">
                    <label class="form-label" for="quickOtpEmail">மின்னஞ்சல் முகவரி (Email)</label>
                    <input type="email" id="quickOtpEmail" class="form-input" placeholder="name@example.com" required>
                </div>
                <button type="submit" id="quickOtpBtn" class="nav-btn nav-btn-primary" style="width:100%;justify-content:center;padding:12px;font-size:15px;margin-top:6px;">
                    📩 OTP சரிபார்ப்புக் குறியீடு பெறுக
                </button>
                <div style="margin-top:14px;text-align:center;">
                    <button type="button" onclick="switchAuthTab(true)" style="background:none;border:none;color:#94a3b8;font-size:12px;cursor:pointer;text-decoration:underline;">
                        ← கடவுச்சொல் உள்நுழைவிற்கு திரும்புக
                    </button>
                </div>
            </form>

            <!-- Feedback Notice Banner -->
            <div id="authNotice" style="display:none;margin-top:16px;padding:12px;border-radius:8px;font-size:13px;text-align:center;line-height:1.5;"></div>
        </div>
    </main>

    <script>
        let currentEmail = '';
        let timerInterval = null;

        function showNotice(text, isError) {
            const n = document.getElementById('authNotice');
            n.style.display = 'block';
            n.textContent = text;
            if (isError) {
                n.style.background = 'rgba(239, 68, 68, 0.15)';
                n.style.color = '#ef4444';
                n.style.border = '1px solid rgba(239, 68, 68, 0.3)';
            } else {
                n.style.background = 'rgba(34, 197, 94, 0.15)';
                n.style.color = '#22c55e';
                n.style.border = '1px solid rgba(34, 197, 94, 0.3)';
            }
        }

        function clearNotice() {
            const n = document.getElementById('authNotice');
            n.style.display = 'none';
        }

        function switchAuthTab(login) {
            clearNotice();
            document.getElementById('authTabsContainer').style.display = 'flex';
            document.getElementById('loginForm').style.display = login ? 'block' : 'none';
            document.getElementById('registerForm').style.display = login ? 'none' : 'block';
            document.getElementById('otpVerifyForm').style.display = 'none';
            document.getElementById('quickOtpForm').style.display = 'none';

            document.getElementById('tabLogin').className = 'auth-tab-btn ' + (login ? 'active' : '');
            document.getElementById('tabRegister').className = 'auth-tab-btn ' + (!login ? 'active' : '');
        }

        function showQuickOtpLogin() {
            clearNotice();
            document.getElementById('authTabsContainer').style.display = 'none';
            document.getElementById('loginForm').style.display = 'none';
            document.getElementById('registerForm').style.display = 'none';
            document.getElementById('otpVerifyForm').style.display = 'none';
            document.getElementById('quickOtpForm').style.display = 'block';
        }

        function startResendCountdown() {
            let left = 60;
            const textSpan = document.getElementById('resendTimerText');
            const countSpan = document.getElementById('resendCountdown');
            const btn = document.getElementById('resendOtpBtn');
            textSpan.style.display = 'inline';
            btn.style.display = 'none';
            countSpan.textContent = left;

            if (timerInterval) clearInterval(timerInterval);
            timerInterval = setInterval(() => {
                left--;
                countSpan.textContent = left;
                if (left <= 0) {
                    clearInterval(timerInterval);
                    textSpan.style.display = 'none';
                    btn.style.display = 'inline';
                }
            }, 1000);
        }

        function cancelOtpStep() {
            clearInterval(timerInterval);
            switchAuthTab(true);
        }

        // 1. Traditional Registration with Email OTP Verification
        async function handleRegisterSubmit(e) {
            e.preventDefault();
            clearNotice();

            const fullName = document.getElementById('regFullName').value.trim();
            const username = document.getElementById('regUsername').value.trim();
            const email = document.getElementById('regEmail').value.trim();
            const password = document.getElementById('regPassword').value;

            if (!fullName || !username || !email || !password) {
                showNotice('அனைத்து புலங்களையும் நிரப்பவும்.', true);
                return;
            }

            const btn = document.getElementById('regBtn');
            btn.disabled = true;
            btn.innerHTML = '🔄 கணக்கு துவங்கப்படுகிறது...';

            try {
                const res = await fetch('/auth/register', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        fullName: fullName,
                        username: username,
                        email: email,
                        password: password
                    })
                });
                const data = await res.json();
                if (data.success && data.data && data.data.requiresOtp) {
                    currentEmail = email;
                    document.getElementById('otpEmailTarget').textContent = email;
                    document.getElementById('authTabsContainer').style.display = 'none';
                    document.getElementById('registerForm').style.display = 'none';
                    document.getElementById('otpVerifyForm').style.display = 'block';
                    startResendCountdown();
                    showNotice(data.message || 'OTP குறியீடு அனுப்பப்பட்டது! குறியீட்டை உள்ளிட்டு கணக்கை உறுதிப்படுத்தவும்.', false);
                    setTimeout(() => {
                        document.getElementById('otpInput').focus();
                    }, 100);
                } else if (data.success) {
                    // Fallback direct login if requiresOtp was false
                    window.location.href = '/submit';
                } else {
                    showNotice(data.message || 'பதிவு தோல்வியடைந்தது.', true);
                }
            } catch (err) {
                showNotice('இணைப்பு பிழை: ' + err.message, true);
            } finally {
                btn.disabled = false;
                btn.innerHTML = '🚀 கணக்கு துவங்கி OTP பெறுக';
            }
        }

        // 2. Traditional Email + Password Login
        async function handleLoginSubmit(e) {
            e.preventDefault();
            clearNotice();

            const email = document.getElementById('loginEmail').value.trim();
            const password = document.getElementById('loginPassword').value;

            const btn = document.getElementById('loginBtn');
            btn.disabled = true;
            btn.innerHTML = '🔄 உள்நுழைகிறது...';

            try {
                const res = await fetch('/auth/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email: email, password: password })
                });
                const data = await res.json();
                if (data.success) {
                    showNotice('✅ உள்நுழைவு வெற்றிகரமானது!', false);
                    const urlParams = new URLSearchParams(window.location.search);
                    const next = urlParams.get('next') || '/submit';
                    setTimeout(() => {
                        window.location.href = next;
                    }, 500);
                } else if (data.requiresOtp) {
                    // User exists but has not verified email yet -> trigger OTP verification
                    currentEmail = email;
                    document.getElementById('otpEmailTarget').textContent = email;
                    document.getElementById('authTabsContainer').style.display = 'none';
                    document.getElementById('loginForm').style.display = 'none';
                    document.getElementById('otpVerifyForm').style.display = 'block';
                    startResendCountdown();
                    showNotice(data.message || 'மின்னஞ்சல் சரிபார்க்கப்படவில்லை. புதிய OTP அனுப்பப்பட்டுள்ளது.', false);
                    setTimeout(() => {
                        document.getElementById('otpInput').focus();
                    }, 100);
                } else {
                    showNotice(data.message || 'மின்னஞ்சல் அல்லது கடவுச்சொல் தவறானது.', true);
                }
            } catch (err) {
                showNotice('இணைப்பு பிழை: ' + err.message, true);
            } finally {
                btn.disabled = false;
                btn.innerHTML = 'உள்நுழைக';
            }
        }

        // 3. Quick OTP Request
        async function handleQuickOtpSend(e) {
            e.preventDefault();
            clearNotice();
            const email = document.getElementById('quickOtpEmail').value.trim();
            if (!email) return;

            const btn = document.getElementById('quickOtpBtn');
            btn.disabled = true;
            btn.innerHTML = '🔄 அனுப்பப்படுகிறது...';

            try {
                const res = await fetch('/auth/otp/send', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email: email })
                });
                const data = await res.json();
                if (data.success) {
                    currentEmail = email;
                    document.getElementById('otpEmailTarget').textContent = email;
                    document.getElementById('quickOtpForm').style.display = 'none';
                    document.getElementById('otpVerifyForm').style.display = 'block';
                    startResendCountdown();
                    showNotice('OTP அனுப்பப்பட்டுள்ளது. சரிபார்க்கவும்.', false);
                    setTimeout(() => {
                        document.getElementById('otpInput').focus();
                    }, 100);
                } else {
                    showNotice(data.message || 'OTP அனுப்ப முடியவில்லை.', true);
                }
            } catch (err) {
                showNotice('இணைப்பு பிழை: ' + err.message, true);
            } finally {
                btn.disabled = false;
                btn.innerHTML = '📩 OTP சரிபார்ப்புக் குறியீடு பெறுக';
            }
        }

        // 4. Verify OTP & Finalize First-Time Login
        async function handleVerifyOTP(e) {
            e.preventDefault();
            clearNotice();

            const otp = document.getElementById('otpInput').value.trim();
            if (!otp || otp.length < 6) {
                showNotice('6-இலக்க OTP குறியீட்டை முழுமையாக உள்ளிடவும்.', true);
                return;
            }

            const btn = document.getElementById('verifyOtpBtn');
            btn.disabled = true;
            btn.innerHTML = '🔄 சரிபார்க்கப்படுகிறது...';

            try {
                const res = await fetch('/auth/otp/verify', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email: currentEmail, otp: otp })
                });
                const data = await res.json();
                if (data.success) {
                    showNotice('✅ சரிபார்க்கப்பட்டது! உள்நுழைந்துள்ளீர்கள்...', false);
                    const urlParams = new URLSearchParams(window.location.search);
                    const next = urlParams.get('next') || '/submit';
                    setTimeout(() => {
                        window.location.href = next;
                    }, 500);
                } else {
                    showNotice(data.message || 'தவறான அல்லது காலாவதியான OTP.', true);
                    btn.disabled = false;
                    btn.innerHTML = '<span>🚀 சரிபார்த்து உள்நுழைக</span>';
                }
            } catch (err) {
                showNotice('இணைப்பு பிழை: ' + err.message, true);
                btn.disabled = false;
                btn.innerHTML = '<span>🚀 சரிபார்த்து உள்நுழைக</span>';
            }
        }

        async function handleResendOTP() {
            if (!currentEmail) return;
            const btn = document.getElementById('resendOtpBtn');
            btn.disabled = true;
            btn.textContent = 'அனுப்பப்படுகிறது...';
            try {
                const res = await fetch('/auth/otp/send', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email: currentEmail })
                });
                const data = await res.json();
                if (data.success) {
                    showNotice('புதிய OTP குறியீடு அனுப்பப்பட்டது!', false);
                    startResendCountdown();
                } else {
                    showNotice(data.message || 'மீண்டும் அனுப்புவதில் பிழை.', true);
                    btn.disabled = false;
                    btn.textContent = 'மீண்டும் OTP அனுப்புக ↺';
                }
            } catch (err) {
                showNotice('இணைப்பு பிழை: ' + err.message, true);
                btn.disabled = false;
                btn.textContent = 'மீண்டும் OTP அனுப்புக ↺';
            }
        }
    </script>
</body>
</html>`
}

// renderFrontProfilePageHTML renders the user's personal profile editor
func renderFrontProfilePageHTML(userID, username, fullName string) string {
	head := getCommonFrontHead("எனது சுயவிவரம் | TN24 வாசகர் அரங்கம்", "உங்கள் வாசகர் அரங்கம் சுயவிவரம், பெயர் மற்றும் எழுத்தாளர் தகவல்களை மாற்றியமைக்கவும்.", "https://www.tn24.in/user/profile")

	return head + `
    <style>
        .profile-avatar-preview {
            width: 88px;
            height: 88px;
            border-radius: 50%;
            border: 2px solid #ef4444;
            object-fit: cover;
            background: #1e293b;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 32px;
            color: #f8fafc;
        }
        .stat-pill {
            display: inline-flex;
            align-items: center;
            gap: 6px;
            background: rgba(30, 41, 59, 0.8);
            border: 1px solid #1e293b;
            border-radius: 20px;
            padding: 6px 14px;
            font-size: 13px;
            color: #cbd5e1;
        }
    </style>
</head>
<body>
    ` + getCommonFrontHeader("profile", userID, fullName) + `

    <main class="front-container" style="max-width: 680px; margin-top: 30px; margin-bottom: 60px;">
        <!-- Header & Breadcrumbs -->
        <div style="margin-bottom: 24px;">
            <div style="display:flex;align-items:center;gap:8px;font-size:13px;color:#94a3b8;margin-bottom:8px;">
                <a href="/">முகப்பு</a> <span>&rsaquo;</span> <a href="/user-corner">வாசகர் அரங்கம்</a> <span>&rsaquo;</span> <span style="color:#f8fafc;">சுயவிவரம்</span>
            </div>
            <h1 style="font-size: 24px; font-weight: 800; color: #f8fafc; letter-spacing: -0.5px;">
                👤 எனது சுயவிவரம் (Profile Settings)
            </h1>
            <p style="color: #94a3b8; font-size: 14px; margin-top: 4px;">
                வாசகர் அரங்கில் உங்களின் கட்டுரைகளுடன் காட்டப்படும் எழுத்தாளர் விவரங்களை இங்கே புதுப்பிக்கலாம்.
            </p>
        </div>

        <!-- Contributor Status Summary Card -->
        <div class="front-card" style="margin-bottom: 20px; display:flex; gap: 20px; align-items: center; flex-wrap: wrap;">
            <div id="avatarBox" class="profile-avatar-preview">
                👤
            </div>
            <div style="flex: 1; min-width: 220px;">
                <div style="display:flex; align-items:center; gap: 10px; flex-wrap: wrap;">
                    <h2 id="headerDisplayName" style="font-size: 19px; font-weight: 800; color: #f8fafc; margin: 0;"></h2>
                    <span id="headerVerifiedBadge" class="status-badge badge-published" style="font-size: 11px;">✅ சரிபார்க்கப்பட்டது</span>
                </div>
                <div id="headerUsername" style="font-size: 13px; color: #38bdf8; margin-top: 2px;"></div>
                
                <div style="display:flex; gap: 8px; margin-top: 10px; flex-wrap: wrap;">
                    <span class="stat-pill">⭐ <span id="headerLevel">New Contributor</span></span>
                    <span class="stat-pill">🏅 <span id="headerPoints">0</span> புள்ளிகள்</span>
                    <span class="stat-pill">📰 <span id="headerApproved">0</span> கட்டுரைகள்</span>
                </div>
            </div>
        </div>

        <!-- Profile Edit Form Card -->
        <div class="front-card" style="margin-bottom: 20px;">
            <h3 style="font-size: 16px; font-weight: 700; color: #f8fafc; margin-bottom: 16px; border-bottom: 1px solid #1e293b; padding-bottom: 10px;">
                📝 அடிப்படை தகவல்கள் (Basic Details)
            </h3>

            <form id="profileEditForm" onsubmit="handleProfileUpdate(event)">
                <div class="form-row">
                    <label class="form-label" for="profileFullName">முழுப் பெயர் (Display Name) <span class="req">*</span></label>
                    <input type="text" id="profileFullName" class="form-input" required>
                </div>

                <div class="form-grid-2 form-row">
                    <div>
                        <label class="form-label" for="profileUsername">பயனர் பெயர் (Username)</label>
                        <input type="text" id="profileUsername" class="form-input" readonly disabled style="opacity: 0.6; cursor: not-allowed;">
                    </div>
                    <div>
                        <label class="form-label" for="profileEmail">மின்னஞ்சல் (Email Address)</label>
                        <input type="email" id="profileEmail" class="form-input" readonly disabled style="opacity: 0.6; cursor: not-allowed;">
                    </div>
                </div>

                <div class="form-row">
                    <label class="form-label" for="profileBio">சுயவிவரக் குறிப்பு / அறிமுகம் (Author Bio)</label>
                    <textarea id="profileBio" class="form-textarea" rows="3" placeholder="உங்களைப்பற்றிய சிறு அறிமுகம், எழுத்து ஆர்வம் அல்லது ஊடக அனுபவம்..."></textarea>
                    <p style="color:#64748b;font-size:12px;margin-top:4px;">💡 இந்த அறிமுகம் உங்கள் கட்டுரைகளின் கீழ் வாசகர்களுக்குக் காட்டப்படும்.</p>
                </div>

                <div class="form-row">
                    <label class="form-label" for="profileAvatarUrl">சுயவிவரப் படம் (Avatar Image URL)</label>
                    <input type="url" id="profileAvatarUrl" class="form-input" placeholder="https://example.com/my-photo.jpg" oninput="updateAvatarPreview()">
                </div>

                <div style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 24px;">
                    <button type="submit" id="saveProfileBtn" class="nav-btn nav-btn-primary" style="padding: 10px 24px; font-size: 14px;">
                        💾 சுயவிவரத்தை சேமிக்க
                    </button>
                </div>
            </form>
        </div>

        <!-- Password Change Card -->
        <div class="front-card" style="margin-bottom: 24px;">
            <h3 style="font-size: 16px; font-weight: 700; color: #f8fafc; margin-bottom: 16px; border-bottom: 1px solid #1e293b; padding-bottom: 10px;">
                🔑 கடவுச்சொல் மாற்றம் (Change Password)
            </h3>

            <form id="passwordChangeForm" onsubmit="handlePasswordUpdate(event)">
                <div class="form-grid-2 form-row">
                    <div>
                        <label class="form-label" for="newPassword">புதிய கடவுச்சொல் (New Password)</label>
                        <input type="password" id="newPassword" class="form-input" placeholder="குறைந்தது 6 எழுத்துக்கள்" minlength="6" required>
                    </div>
                    <div>
                        <label class="form-label" for="confirmPassword">மீண்டும் புதிய கடவுச்சொல்</label>
                        <input type="password" id="confirmPassword" class="form-input" placeholder="••••••••" minlength="6" required>
                    </div>
                </div>

                <div style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 16px;">
                    <button type="submit" id="savePasswordBtn" class="nav-btn nav-btn-outline" style="padding: 9px 20px; font-size: 14px;">
                        கடவுச்சொல்லை மாற்றுக
                    </button>
                </div>
            </form>
        </div>

        <!-- Quick Links Hub -->
        <div style="display:flex; gap:12px; justify-content:center; flex-wrap:wrap;">
            <a href="/user/my-posts" class="nav-btn nav-btn-outline">📋 எனது கட்டுரைகள்</a>
            <a href="/submit" class="nav-btn nav-btn-primary">✍️ புதிய செய்தி எழுதுக</a>
            <a id="publicAuthorLink" href="/user-corner" target="_blank" class="nav-btn nav-btn-outline">👁️ பொது பக்கத்தைக் காண்க ↗</a>
        </div>

        <!-- Toast Feedback Notice -->
        <div id="profileNotice" style="display:none;position:fixed;bottom:24px;right:24px;padding:14px 20px;border-radius:10px;box-shadow:0 10px 30px rgba(0,0,0,0.5);font-size:14px;z-index:999;"></div>
    </main>

    <footer>
        &copy; 2026 TN24 — Tamil Nadu's Independent Digital News Platform.
    </footer>

    <script>
        function showProfileNotice(msg, isError) {
            const n = document.getElementById('profileNotice');
            n.textContent = msg;
            n.style.display = 'block';
            if (isError) {
                n.style.background = '#ef4444';
                n.style.color = '#fff';
            } else {
                n.style.background = '#22c55e';
                n.style.color = '#fff';
            }
            setTimeout(() => { n.style.display = 'none'; }, 3500);
        }

        function updateAvatarPreview() {
            const url = document.getElementById('profileAvatarUrl').value.trim();
            const box = document.getElementById('avatarBox');
            if (url && (url.startsWith('http://') || url.startsWith('https://'))) {
                box.innerHTML = '<img src="' + url + '" alt="Avatar" style="width:100%;height:100%;border-radius:50%;object-fit:cover;" onerror="this.onerror=null;this.parentElement.innerHTML=\'👤\';">';
            } else {
                box.innerHTML = '👤';
            }
        }

        async function loadProfileData() {
            try {
                const res = await fetch('/auth/me');
                const data = await res.json();
                if (!data.success || !data.data) {
                    window.location.href = '/login?next=/user/profile';
                    return;
                }

                const u = data.data;
                document.getElementById('headerDisplayName').textContent = u.fullName || u.username;
                document.getElementById('headerUsername').textContent = '@' + u.username;
                document.getElementById('profileFullName').value = u.fullName || '';
                document.getElementById('profileUsername').value = u.username || '';
                document.getElementById('profileEmail').value = u.email || '';
                document.getElementById('profileBio').value = u.bio || '';
                document.getElementById('profileAvatarUrl').value = u.avatarUrl || '';

                if (u.level) document.getElementById('headerLevel').textContent = u.level;
                if (u.points !== undefined) document.getElementById('headerPoints').textContent = u.points;
                if (u.approvedCount !== undefined) document.getElementById('headerApproved').textContent = u.approvedCount;

                const publicLink = document.getElementById('publicAuthorLink');
                if (publicLink && u.username) {
                    publicLink.href = '/author/' + encodeURIComponent(u.username);
                }

                updateAvatarPreview();
            } catch (err) {
                showProfileNotice('சுயவிவர தகவல்களை பெற முடியவில்லை.', true);
            }
        }

        async function handleProfileUpdate(e) {
            e.preventDefault();
            const btn = document.getElementById('saveProfileBtn');
            btn.disabled = true;
            btn.innerHTML = '🔄 சேமிக்கப்படுகிறது...';

            const fullName = document.getElementById('profileFullName').value.trim();
            const bio = document.getElementById('profileBio').value.trim();
            const avatarUrl = document.getElementById('profileAvatarUrl').value.trim();

            try {
                const res = await fetch('/auth/profile', {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        fullName: fullName,
                        bio: bio,
                        avatarUrl: avatarUrl
                    })
                });
                const data = await res.json();
                if (data.success) {
                    showProfileNotice('✅ சுயவிவரம் வெற்றிகரமாக சேமிக்கப்பட்டது!', false);
                    document.getElementById('headerDisplayName').textContent = fullName;
                } else {
                    showProfileNotice(data.message || 'சேமிக்க முடியவில்லை.', true);
                }
            } catch (err) {
                showProfileNotice('பிழை: ' + err.message, true);
            } finally {
                btn.disabled = false;
                btn.innerHTML = '💾 சுயவிவரத்தை சேமிக்க';
            }
        }

        async function handlePasswordUpdate(e) {
            e.preventDefault();
            const p1 = document.getElementById('newPassword').value;
            const p2 = document.getElementById('confirmPassword').value;

            if (p1 !== p2) {
                showProfileNotice('கடவுச்சொற்கள் பொருந்தவில்லை.', true);
                return;
            }

            const btn = document.getElementById('savePasswordBtn');
            btn.disabled = true;
            btn.innerHTML = 'மாற்றப்படுகிறது...';

            try {
                const fullName = document.getElementById('profileFullName').value.trim();
                const bio = document.getElementById('profileBio').value.trim();
                const avatarUrl = document.getElementById('profileAvatarUrl').value.trim();

                const res = await fetch('/auth/profile', {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        fullName: fullName,
                        bio: bio,
                        avatarUrl: avatarUrl,
                        newPassword: p1
                    })
                });
                const data = await res.json();
                if (data.success) {
                    showProfileNotice('✅ கடவுச்சொல் வெற்றிகரமாக மாற்றப்பட்டது!', false);
                    document.getElementById('newPassword').value = '';
                    document.getElementById('confirmPassword').value = '';
                } else {
                    showProfileNotice(data.message || 'கடவுச்சொல்லை மாற்ற முடியவில்லை.', true);
                }
            } catch (err) {
                showProfileNotice('பிழை: ' + err.message, true);
            } finally {
                btn.disabled = false;
                btn.innerHTML = 'கடவுச்சொல்லை மாற்றுக';
            }
        }

        document.addEventListener('DOMContentLoaded', loadProfileData);
    </script>
</body>
</html>`
}

// renderFrontMyPostsPageHTML renders the user's personal submission tracker
func renderFrontMyPostsPageHTML(userID, username, fullName string) string {
	head := getCommonFrontHead("எனது பதிவுகள் | TN24 வாசகர் அரங்கம்", "நீங்கள் சமர்ப்பித்த கட்டுரைகள் மற்றும் அவற்றின் நிலை (Under Review, Published).", "https://www.tn24.in/user/my-posts")

	return head + `
</head>
<body>
    ` + getCommonFrontHeader("myposts", userID, fullName) + `

    <main class="front-container">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:24px;flex-wrap:wrap;gap:12px;">
            <div>
                <h1 style="font-size:24px;font-weight:800;color:#f8fafc;">📋 எனது கட்டுரைகள் (My Submissions)</h1>
                <p style="color:#94a3b8;font-size:13px;margin-top:2px;">நீங்கள் சமர்ப்பித்த செய்திகளின் பரிசீலனை நிலை மற்றும் நேரடி இணைப்புகள்</p>
            </div>
            <a href="/submit" class="nav-btn nav-btn-primary">+ புதிய செய்தி எழுதுக</a>
        </div>

        <div id="postsListContainer">
            <div style="text-align:center;padding:40px;color:#94a3b8;">🔄 உங்கள் கட்டுரைகள் ஏற்றப்படுகின்றன...</div>
        </div>
    </main>

    <footer>
        &copy; 2026 TN24 — Tamil Nadu's Independent Digital News Platform.
    </footer>

    <script>
        async function loadMyPosts() {
            const container = document.getElementById('postsListContainer');
            try {
                const res = await fetch('/api/user/my-posts');
                const data = await res.json();
                if (!data.success || !data.data || data.data.length === 0) {
                    container.innerHTML = '<div class="front-card" style="text-align:center;padding:48px 20px;">' +
                        '<div style="font-size:42px;margin-bottom:12px;">✍️</div>' +
                        '<h3 style="color:#f8fafc;font-size:18px;margin-bottom:6px;">இன்னும் செய்திகள் சமர்ப்பிக்கப்படவில்லை</h3>' +
                        '<p style="color:#94a3b8;font-size:14px;margin-bottom:20px;">உங்கள் பகுதி முக்கிய நிகழ்வுகளை உடனே எழுதி வாசகர் அரங்கில் பகிருங்கள்.</p>' +
                        '<a href="/submit" class="nav-btn nav-btn-primary">முதல் செய்தியை எழுதுக →</a>' +
                    '</div>';
                    return;
                }

                let html = '';
                data.data.forEach(p => {
                    let badge = '<span class="status-badge badge-pending">🟡 பரிசீலனையில் (Under Review)</span>';
                    let viewLink = '';
                    if (p.status === 'PUBLISHED') {
                        badge = '<span class="status-badge badge-published">🟢 வெளியிடப்பட்டது (Live)</span>';
                        viewLink = '<a href="/news/' + p.slug + '" target="_blank" class="nav-btn nav-btn-outline" style="font-size:12px;padding:6px 12px;">கட்டுரையைக் காண்க ↗</a>';
                    } else if (p.status === 'REJECTED') {
                        badge = '<span class="status-badge badge-rejected">🔴 நிராகரிக்கப்பட்டது (Rejected)</span>';
                    }

                    const dateStr = new Date(p.createdAt).toLocaleDateString('ta-IN', { year: 'numeric', month: 'short', day: 'numeric' });

                    html += '<div class="front-card" style="margin-bottom:16px;padding:20px;display:flex;gap:16px;flex-wrap:wrap;align-items:center;justify-content:space-between;">' +
                        '<div style="flex:1;min-width:260px;">' +
                            '<div style="display:flex;align-items:center;gap:8px;margin-bottom:6px;">' +
                                badge +
                                '<span style="font-size:12px;color:#64748b;">• ' + (p.district || 'தமிழ்நாடு') + ' • ' + (p.category || 'User Corner') + '</span>' +
                            '</div>' +
                            '<h2 style="font-size:17px;font-weight:700;color:#f8fafc;margin-bottom:6px;">' + p.title + '</h2>' +
                            '<p style="font-size:13px;color:#94a3b8;margin:0;">' + (p.description || '') + '</p>' +
                            '<div style="font-size:12px;color:#64748b;margin-top:6px;">🕒 சமர்ப்பிக்கப்பட்ட தேதி: ' + dateStr + '</div>' +
                        '</div>' +
                        '<div>' + viewLink + '</div>' +
                    '</div>';
                });
                container.innerHTML = html;
            } catch (err) {
                container.innerHTML = '<div style="color:#ef4444;text-align:center;padding:24px;">ஏற்றுவதில் பிழை: ' + err.message + '</div>';
            }
        }

        function handleLogout() {
            fetch('/auth/logout', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: '{}' })
                .then(() => { window.location.href = '/'; });
        }

        loadMyPosts();
    </script>
</body>
</html>`
}

// renderFrontUserCornerHubHTML renders the dedicated User Corner landing page
func renderFrontUserCornerHubHTML() string {
	head := getCommonFrontHead(
		"வாசகர் அரங்கம் – தமிழ்நாடு செய்திகள் & சமூகம் | TN24 User Corner",
		"தமிழ்நாட்டின் அனைத்து மாவட்ட வாசகர்களின் களச் செய்திகள், பார்வைகள் மற்றும் சமூக பதிவுகள். TN24 வாசகர் அரங்கம்.",
		"https://www.tn24.in/user-corner",
	)

	return head + `
</head>
<body>
    ` + getCommonFrontHeader("user-corner", "", "") + `

    <main class="front-container">
        <!-- Hero Hub Banner -->
        <div class="front-card" style="background: linear-gradient(135deg, #0d1421 0%, #151e33 100%); border-color: #2563eb; margin-bottom: 32px; padding: 36px 28px;">
            <div style="display:inline-block;background:rgba(37,99,235,0.2);color:#60a5fa;border:1px solid rgba(37,99,235,0.4);padding:4px 12px;border-radius:20px;font-size:12px;font-weight:700;margin-bottom:12px;">
                ✍️ CITIZEN JOURNALISM & OPINIONS
            </div>
            <h1 style="font-size: 30px; font-weight: 800; color: #fff; line-height: 1.3; margin-bottom: 8px;">
                வாசகர் அரங்கம் | TN24 Community Voices
            </h1>
            <p style="color: #94a3b8; font-size: 15px; max-width: 680px; margin-bottom: 20px;">
                தமிழ்நாட்டின் 38 மாவட்டங்களைச் சேர்ந்த வாசகர்களின் நேரடி களச் செய்திகள், சமூகப் பிரச்சனைகள், தனித்துவக் கட்டுரைகள் மற்றும் மக்கள் குரல். ஆசிரியர் குழுவால் சரிபார்க்கப்பட்டு வெளியிடப்படுகிறது.
            </p>
            <div style="display:flex;gap:12px;flex-wrap:wrap;">
                <a href="/submit" class="nav-btn nav-btn-primary" style="padding:10px 22px;">+ உங்கள் செய்தியை சமர்ப்பிக்கவும்</a>
                <a href="/" class="nav-btn nav-btn-outline" style="padding:10px 20px;">அனைத்து செய்திகள்</a>
            </div>
        </div>

        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px;">
            <h2 style="font-size:20px;font-weight:700;color:#f8fafc;">அண்மைய வாசகர் பதிவுகள்</h2>
            <span style="font-size:13px;color:#94a3b8;">சரிபார்க்கப்பட்ட கட்டுரைகள்</span>
        </div>

        <div id="userCornerFeed" style="display:grid;grid-template-columns:repeat(auto-fill, minmax(280px, 1fr));gap:20px;">
            <div style="grid-column:1/-1;text-align:center;padding:40px;color:#94a3b8;">🔄 கட்டுரைகள் ஏற்றப்படுகின்றன...</div>
        </div>
    </main>

    <footer>
        &copy; 2026 TN24 — Tamil Nadu's Independent Digital News Platform. வாசகர் அரங்கம்.
    </footer>

    <script>
        async function loadUserCornerFeed() {
            const container = document.getElementById('userCornerFeed');
            try {
                const res = await fetch('/api/portal/feed?category=User%20Corner');
                const data = await res.json();
                const articles = (data.centerArticles || []).concat(data.leftFeed || []);
                if (!articles || articles.length === 0) {
                    container.innerHTML = '<div class="front-card" style="grid-column:1/-1;text-align:center;padding:48px 20px;">' +
                        '<div style="font-size:42px;margin-bottom:12px;">✍️</div>' +
                        '<h3 style="color:#f8fafc;font-size:18px;margin-bottom:6px;">முதல் செய்தியை நீங்கள் பதிவிடுங்கள்!</h3>' +
                        '<p style="color:#94a3b8;font-size:14px;margin-bottom:20px;">உங்கள் பகுதி நடப்புகளை நேரடியாக தமிழ்நாடு முழுவதும் கொண்டு செல்லுங்கள்.</p>' +
                        '<a href="/submit" class="nav-btn nav-btn-primary">+ சமர்ப்பிக்கவும்</a>' +
                    '</div>';
                    return;
                }

                let html = '';
                articles.forEach(a => {
                    const thumb = a.thumbnail || '/assets/brand/tn24-icon.svg';
                    html += '<article class="front-card" style="padding:16px;display:flex;flex-direction:column;justify-content:space-between;">' +
                        '<div>' +
                            '<div style="height:160px;border-radius:8px;overflow:hidden;margin-bottom:12px;background:#000;">' +
                                '<img src="' + thumb + '" alt="' + a.title + '" style="width:100%;height:100%;object-fit:cover;">' +
                            '</div>' +
                            '<div style="display:flex;align-items:center;gap:6px;margin-bottom:6px;">' +
                                '<span style="font-size:11px;font-weight:700;color:#60a5fa;background:rgba(37,99,235,0.15);padding:2px 8px;border-radius:4px;">' + (a.district || 'தமிழ்நாடு') + '</span>' +
                            '</div>' +
                            '<h3 style="font-size:16px;font-weight:700;color:#f8fafc;line-height:1.4;margin-bottom:8px;">' +
                                '<a href="/news/' + a.slug + '">' + a.title + '</a>' +
                            '</h3>' +
                            '<p style="font-size:13px;color:#94a3b8;line-height:1.5;margin-bottom:12px;display:-webkit-box;-webkit-line-clamp:3;-webkit-box-orient:vertical;overflow:hidden;">' + (a.description || '') + '</p>' +
                        '</div>' +
                        '<div style="display:flex;justify-content:space-between;align-items:center;border-top:1px solid #1e293b;padding-top:10px;font-size:12px;color:#64748b;">' +
                            '<span>✍️ ' + (a.author || 'வாசகர் நிருபர்') + '</span>' +
                            '<a href="/news/' + a.slug + '" style="color:#ef4444;font-weight:700;">படிக்க →</a>' +
                        '</div>' +
                    '</article>';
                });
                container.innerHTML = html;
            } catch (err) {
                container.innerHTML = '<div style="grid-column:1/-1;color:#ef4444;text-align:center;padding:24px;">ஏற்றுவதில் பிழை: ' + err.message + '</div>';
            }
        }
        loadUserCornerFeed();
    </script>
</body>
</html>`
}

// renderAuthorProfilePageHTML renders Google E-E-A-T compliant contributor author profile
func renderAuthorProfilePageHTML(username, fullName, bio, avatarUrl string, approvedCount, points, trustScore int, articles []ArticleItem) string {
	head := getCommonFrontHead(
		fmt.Sprintf("%s | TN24 வாசகர் நிருபர் சுயவிவரம்", fullName),
		fmt.Sprintf("%s அவர்களின் கட்டுரைகள் மற்றும் கள அறிக்கைகள். TN24 வாசகர் அரங்கம்.", fullName),
		fmt.Sprintf("https://www.tn24.in/author/%s", username),
	)

	articlesHTML := ""
	for _, a := range articles {
		thumb := a.Thumbnail
		if thumb == "" {
			thumb = "/assets/brand/tn24-icon.svg"
		}
		articlesHTML += fmt.Sprintf(`
		<div class="front-card" style="padding:16px;margin-bottom:12px;display:flex;gap:16px;align-items:center;">
			<img src="%s" alt="%s" style="width:100px;height:75px;border-radius:6px;object-fit:cover;background:#000;">
			<div style="flex:1;">
				<div style="font-size:11px;color:#60a5fa;font-weight:700;margin-bottom:4px;">%s</div>
				<h3 style="font-size:15px;font-weight:700;color:#f8fafc;margin-bottom:4px;">
					<a href="/news/%s">%s</a>
				</h3>
				<span style="font-size:12px;color:#64748b;">🕒 %s</span>
			</div>
		</div>`, thumb, html.EscapeString(a.Title), html.EscapeString(a.District), a.Slug, html.EscapeString(a.Title), a.CreatedAt.Format("02 Jan 2006"))
	}
	if articlesHTML == "" {
		articlesHTML = `<div class="front-card" style="text-align:center;padding:32px;color:#94a3b8;">இன்னும் கட்டுரைகள் வெளியிடப்படவில்லை.</div>`
	}

	return head + fmt.Sprintf(`
    <!-- JSON-LD Person Schema for Google E-E-A-T -->
    <script type="application/ld+json">
    {
      "@context": "https://schema.org",
      "@type": "Person",
      "name": "%s",
      "url": "https://www.tn24.in/author/%s",
      "jobTitle": "Community Contributor",
      "worksFor": {
        "@type": "NewsMediaOrganization",
        "name": "TN24 — Tamil Nadu News"
      }
    }
    </script>
</head>
<body>
    <header class="front-header">
        <a href="/" title="TN24 முகப்பு">
            <img src="/assets/brand/tn24-logo.svg" alt="TN24 Logo" class="brand-logo-img">
        </a>
        <a href="/user-corner" class="nav-btn nav-btn-outline">✍️ வாசகர் அரங்கம்</a>
    </header>

    <main class="front-container">
        <!-- Author Profile Card -->
        <div class="front-card" style="display:flex;gap:24px;align-items:center;margin-bottom:32px;flex-wrap:wrap;">
            <div style="width:84px;height:84px;border-radius:50%%;background:linear-gradient(135deg,#ef4444,#2563eb);display:flex;align-items:center;justify-content:center;font-size:36px;color:#fff;font-weight:800;">
                %s
            </div>
            <div style="flex:1;">
                <div style="display:flex;align-items:center;gap:10px;margin-bottom:4px;">
                    <h1 style="font-size:24px;font-weight:800;color:#f8fafc;">%s</h1>
                    <span style="font-size:12px;background:rgba(34,197,94,0.15);color:#22c55e;padding:2px 8px;border-radius:12px;font-weight:700;">✓ சரிபார்க்கப்பட்ட நிருபர்</span>
                </div>
                <div style="font-size:13px;color:#94a3b8;margin-bottom:8px;">@%s • வாசகர் எழுத்தாளர் (Community Contributor)</div>
                <p style="font-size:14px;color:#cbd5e1;line-height:1.5;">%s</p>
                <div style="display:flex;gap:16px;margin-top:12px;font-size:13px;color:#64748b;">
                    <span>📰 வெளியிட்ட கட்டுரைகள்: <strong style="color:#f8fafc;">%d</strong></span>
                    <span>⭐ நன்மதிப்பு புள்ளிகள்: <strong style="color:#f59e0b;">%d</strong></span>
                </div>
            </div>
        </div>

        <h2 style="font-size:18px;font-weight:700;color:#f8fafc;margin-bottom:16px;">இவரது படைப்புகள் (Published Articles)</h2>
        %s
    </main>

    <footer>
        &copy; 2026 TN24 — Tamil Nadu's Independent Digital News Platform.
    </footer>
</body>
</html>`, html.EscapeString(fullName), username, string([]rune(fullName)[:1]), html.EscapeString(fullName), username, html.EscapeString(bio), approvedCount, points, articlesHTML)
}

