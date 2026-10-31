package content

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"tn-now/server/internal/config"
	"tn-now/server/internal/db"
	"tn-now/server/internal/media"
	"tn-now/server/internal/scraper"
	"tn-now/server/internal/video"
)

type SubmissionHandler struct {
	q    *db.Queries
	conn *pgxpool.Pool
	cfg  *config.Config
}

func NewSubmissionHandler(q *db.Queries, conn *pgxpool.Pool, cfg ...*config.Config) *SubmissionHandler {
	var c *config.Config
	if len(cfg) > 0 {
		c = cfg[0]
	}
	return &SubmissionHandler{q: q, conn: conn, cfg: c}
}

type submitVideoRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	CategoryID   string `json:"categoryId"`
	DistrictID   string `json:"districtId"`
	LocationID   string `json:"locationId,omitempty"`
	AuthorUserID string `json:"authorUserId"`
	URL          string `json:"url"`
}

type submitPhotoRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	CategoryID   string   `json:"categoryId"`
	DistrictID   string   `json:"districtId"`
	AuthorUserID string   `json:"authorUserId"`
	PhotoURLs    []string `json:"photoUrls"`
	CaptionEn    string   `json:"captionEn,omitempty"`
	CaptionTa    string   `json:"captionTa,omitempty"`
}

type submitStoryRequest struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	CategoryID     string `json:"categoryId"`
	DistrictID     string `json:"districtId"`
	AuthorUserID   string `json:"authorUserId"`
	Headline       string `json:"headline"`
	Body           string `json:"body"`
	SinglePhotoURL string `json:"singlePhotoUrl,omitempty"`
}

type submitEventRequest struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	CategoryID       string `json:"categoryId"`
	DistrictID       string `json:"districtId"`
	AuthorUserID     string `json:"authorUserId"`
	EventName        string `json:"eventName"`
	EventDate        string `json:"eventDate"`
	StartTime        string `json:"startTime,omitempty"`
	EndTime          string `json:"endTime,omitempty"`
	OrganizerName    string `json:"organizerName"`
	OrganizerContact string `json:"organizerContact,omitempty"`
}

func (h *SubmissionHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/content/submit/video", h.HandleSubmitVideo)
	mux.HandleFunc("/content/submit/photo", h.HandleSubmitPhoto)
	mux.HandleFunc("/content/submit/story", h.HandleSubmitStory)
	mux.HandleFunc("/content/submit/event", h.HandleSubmitEvent)
	mux.HandleFunc("/content/upload/image", h.HandleUploadImage)
	mux.HandleFunc("/api/upload/image", h.HandleUploadImage)
}

func (h *SubmissionHandler) HandleUploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 10 MB maximum upload memory limit
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid form data or file exceeds 10MB limit")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		file, header, err = r.FormFile("file")
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "Missing image file in form data")
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil || len(fileBytes) == 0 {
		writeError(w, http.StatusBadRequest, "Failed to read image bytes")
		return
	}

	// Validate content type
	contentType := http.DetectContentType(fileBytes)
	if !strings.HasPrefix(contentType, "image/") {
		writeError(w, http.StatusBadRequest, "Only image files (JPEG, PNG, WebP) are allowed")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	ikCfg := media.GetEffectiveImageKitConfig(ctx, h.conn, h.cfg)
	res, err := media.UploadToImageKit(ctx, ikCfg, fileBytes, header.Filename)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Image upload failed: " + err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"url":          res.URL,
		"thumbnailUrl": res.ThumbnailURL,
		"width":        res.Width,
		"height":       res.Height,
		"size":         res.Size,
		"message":      "Image uploaded successfully to ImageKit CDN",
	})
}

func (h *SubmissionHandler) HandleSubmitVideo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req submitVideoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.URL) == "" || req.CategoryID == "" || req.DistrictID == "" || req.AuthorUserID == "" {
		writeError(w, http.StatusUnprocessableEntity, "Title, URL, CategoryID, DistrictID, and AuthorUserID are required")
		return
	}

	// Parse video URL
	videoMeta, err := video.ParseVideoURL(req.URL)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	ctx := r.Context()
	tx, err := h.conn.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to start database transaction")
		return
	}
	defer tx.Rollback(ctx)

	qtx := h.q.WithTx(tx)

	// Insert Content master record
	var locID pgtype.UUID
	if req.LocationID != "" {
		locID = stringToUUID(req.LocationID)
	}

	var desc pgtype.Text
	if req.Description != "" {
		desc = pgtype.Text{String: req.Description, Valid: true}
	}

	content, err := qtx.CreateContent(ctx, db.CreateContentParams{
		Title:        req.Title,
		Description:  desc,
		ContentType:  "VIDEO_LINK",
		CategoryID:   stringToUUID(req.CategoryID),
		DistrictID:   stringToUUID(req.DistrictID),
		LocationID:   locID,
		SourceType:   "USER",
		SourceUrl:    pgtype.Text{String: req.URL, Valid: true},
		AuthorUserID: stringToUUID(req.AuthorUserID),
		Status:       "PENDING",
	})
	if err != nil {
		slog.Error("Failed to insert video content", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Failed to create submission record")
		return
	}

	// Insert Video Link extension record
	var embedHTML pgtype.Text
	if videoMeta.EmbedHTML != "" {
		embedHTML = pgtype.Text{String: videoMeta.EmbedHTML, Valid: true}
	}
	var thumbURL pgtype.Text
	if videoMeta.ThumbnailURL != "" {
		thumbURL = pgtype.Text{String: videoMeta.ThumbnailURL, Valid: true}
	}

	_, err = qtx.CreateVideoLink(ctx, db.CreateVideoLinkParams{
		ContentID:       content.ID,
		Platform:        videoMeta.Platform,
		ExternalVideoID: videoMeta.ExternalVideoID,
		CanonicalUrl:    videoMeta.CanonicalURL,
		EmbedHtml:       embedHTML,
		ThumbnailUrl:    thumbURL,
	})
	if err != nil {
		slog.Error("Failed to insert video link detail", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Failed to create video details")
		return
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, "Transaction commit failed")
		return
	}

	writeSuccess(w, http.StatusCreated, map[string]string{"contentId": uuidToString(content.ID)}, "Video submission received and queued for moderation")
}

func (h *SubmissionHandler) HandleSubmitPhoto(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req submitPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if strings.TrimSpace(req.Title) == "" || len(req.PhotoURLs) == 0 || req.CategoryID == "" || req.DistrictID == "" || req.AuthorUserID == "" {
		writeError(w, http.StatusUnprocessableEntity, "Title, PhotoURLs, CategoryID, DistrictID, and AuthorUserID are required")
		return
	}

	ctx := r.Context()
	tx, err := h.conn.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Transaction error")
		return
	}
	defer tx.Rollback(ctx)

	qtx := h.q.WithTx(tx)

	var desc pgtype.Text
	if req.Description != "" {
		desc = pgtype.Text{String: req.Description, Valid: true}
	}

	content, err := qtx.CreateContent(ctx, db.CreateContentParams{
		Title:        req.Title,
		Description:  desc,
		ContentType:  "PHOTO",
		CategoryID:   stringToUUID(req.CategoryID),
		DistrictID:   stringToUUID(req.DistrictID),
		SourceType:   "USER",
		AuthorUserID: stringToUUID(req.AuthorUserID),
		Status:       "PENDING",
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var capEn, capTa pgtype.Text
	if req.CaptionEn != "" {
		capEn = pgtype.Text{String: req.CaptionEn, Valid: true}
	}
	if req.CaptionTa != "" {
		capTa = pgtype.Text{String: req.CaptionTa, Valid: true}
	}

	_, err = qtx.CreatePhoto(ctx, db.CreatePhotoParams{
		ContentID:  content.ID,
		PhotoCount: int32(len(req.PhotoURLs)),
		PhotoUrls:  req.PhotoURLs,
		CaptionEn:  capEn,
		CaptionTa:  capTa,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, http.StatusCreated, map[string]string{"contentId": uuidToString(content.ID)}, "Photo post submitted for moderation")
}

func (h *SubmissionHandler) resolveDistrictID(ctx context.Context, input string) pgtype.UUID {
	var uid pgtype.UUID
	if err := uid.Scan(input); err == nil && uid.Valid {
		return uid
	}
	input = strings.TrimSpace(input)
	if input == "" {
		input = "Tamil Nadu"
	}
	var res pgtype.UUID
	err := h.conn.QueryRow(ctx, "SELECT id FROM districts WHERE name ILIKE $1 OR name ILIKE $2 LIMIT 1", input, "%"+input+"%").Scan(&res)
	if err == nil && res.Valid {
		return res
	}
	_ = h.conn.QueryRow(ctx, "SELECT id FROM districts ORDER BY name ASC LIMIT 1").Scan(&res)
	return res
}

func (h *SubmissionHandler) resolveCategoryID(ctx context.Context, input string) pgtype.UUID {
	var uid pgtype.UUID
	if err := uid.Scan(input); err == nil && uid.Valid {
		return uid
	}
	input = strings.TrimSpace(input)
	if input == "" {
		input = "user-corner"
	}
	var res pgtype.UUID
	err := h.conn.QueryRow(ctx, "SELECT id FROM categories WHERE slug = $1 OR name ILIKE $1 OR name ILIKE $2 LIMIT 1", strings.ToLower(input), "%"+input+"%").Scan(&res)
	if err == nil && res.Valid {
		return res
	}
	_ = h.conn.QueryRow(ctx, "SELECT id FROM categories ORDER BY name ASC LIMIT 1").Scan(&res)
	return res
}

func (h *SubmissionHandler) HandleSubmitStory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req submitStoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Body = strings.TrimSpace(req.Body)
	if req.Title == "" || req.Body == "" {
		writeError(w, http.StatusUnprocessableEntity, "Title and Body are required")
		return
	}

	ctx := r.Context()

	// Resolve Author User ID if not provided in payload
	if req.AuthorUserID == "" {
		// Attempt to extract from token cookie or Authorization header
		token := ""
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else if cookie, err := r.Cookie("tn24_user_token"); err == nil && cookie.Value != "" {
			token = cookie.Value
		}

		if token != "" {
			// Find active session for user
			var uid pgtype.UUID
			_ = h.conn.QueryRow(ctx, "SELECT user_id FROM sessions WHERE is_revoked = false AND (expires_at IS NULL OR expires_at > NOW()) ORDER BY created_at DESC LIMIT 1").Scan(&uid)
			if uid.Valid {
				req.AuthorUserID = uuidToString(uid)
			}
		}
	}

	if req.AuthorUserID == "" {
		// Fallback to first registered user/admin if available
		var fallbackID pgtype.UUID
		_ = h.conn.QueryRow(ctx, "SELECT id FROM users ORDER BY created_at ASC LIMIT 1").Scan(&fallbackID)
		if fallbackID.Valid {
			req.AuthorUserID = uuidToString(fallbackID)
		} else {
			writeError(w, http.StatusUnauthorized, "Author login required to submit articles")
			return
		}
	}

	districtUUID := h.resolveDistrictID(ctx, req.DistrictID)
	categoryUUID := h.resolveCategoryID(ctx, req.CategoryID)

	tx, err := h.conn.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Transaction error")
		return
	}
	defer tx.Rollback(ctx)

	qtx := h.q.WithTx(tx)

	var desc pgtype.Text
	if strings.TrimSpace(req.Description) != "" {
		desc = pgtype.Text{String: strings.TrimSpace(req.Description), Valid: true}
	} else {
		// Auto-generate description excerpt from body (strip HTML tags)
		plainBody := scraper.CleanHTML(req.Body)
		plainRunes := []rune(plainBody)
		if len(plainRunes) > 180 {
			plainBody = string(plainRunes[:180]) + "..."
		}
		desc = pgtype.Text{String: plainBody, Valid: true}
	}

	content, err := qtx.CreateContent(ctx, db.CreateContentParams{
		Title:        req.Title,
		Description:  desc,
		ContentType:  "TEXT_STORY",
		CategoryID:   categoryUUID,
		DistrictID:   districtUUID,
		SourceType:   "USER",
		AuthorUserID: stringToUUID(req.AuthorUserID),
		Status:       "PENDING",
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Generate and persist SEO slug immediately
	contentSlug := scraper.GenerateArticleSlug(req.Title, uuidToString(content.ID))
	_, _ = tx.Exec(ctx, "UPDATE content SET slug = $1 WHERE id = $2", contentSlug, content.ID)

	var photoURL pgtype.Text
	if strings.TrimSpace(req.SinglePhotoURL) != "" {
		photoURL = pgtype.Text{String: strings.TrimSpace(req.SinglePhotoURL), Valid: true}
	}

	headline := req.Headline
	if strings.TrimSpace(headline) == "" {
		headline = req.Title
	}

	_, err = qtx.CreateStory(ctx, db.CreateStoryParams{
		ContentID:      content.ID,
		Headline:       headline,
		Body:           req.Body,
		SinglePhotoUrl: photoURL,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, http.StatusCreated, map[string]interface{}{
		"contentId": uuidToString(content.ID),
		"slug":      contentSlug,
		"status":    "PENDING",
		"message":   "உங்கள் பதிவு வெற்றிகரமாக பெறப்பட்டது! ஆசிரியர் குழுவின் பரிசீலனையில் உள்ளது (Under Review).",
	}, "Article submitted successfully for editorial review")
}

func (h *SubmissionHandler) HandleSubmitEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req submitEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.EventName) == "" || req.EventDate == "" || req.CategoryID == "" || req.DistrictID == "" || req.AuthorUserID == "" || req.OrganizerName == "" {
		writeError(w, http.StatusUnprocessableEntity, "Title, EventName, EventDate, CategoryID, DistrictID, AuthorUserID, and OrganizerName are required")
		return
	}

	eventDate, err := time.Parse(time.RFC3339, req.EventDate)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "Invalid EventDate RFC3339 format")
		return
	}

	ctx := r.Context()
	tx, err := h.conn.Begin(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Transaction error")
		return
	}
	defer tx.Rollback(ctx)

	qtx := h.q.WithTx(tx)

	var desc pgtype.Text
	if req.Description != "" {
		desc = pgtype.Text{String: req.Description, Valid: true}
	}

	content, err := qtx.CreateContent(ctx, db.CreateContentParams{
		Title:        req.Title,
		Description:  desc,
		ContentType:  "EVENT",
		CategoryID:   stringToUUID(req.CategoryID),
		DistrictID:   stringToUUID(req.DistrictID),
		SourceType:   "USER",
		AuthorUserID: stringToUUID(req.AuthorUserID),
		Status:       "PENDING",
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var st, et, contact pgtype.Text
	if req.StartTime != "" {
		st = pgtype.Text{String: req.StartTime, Valid: true}
	}
	if req.EndTime != "" {
		et = pgtype.Text{String: req.EndTime, Valid: true}
	}
	if req.OrganizerContact != "" {
		contact = pgtype.Text{String: req.OrganizerContact, Valid: true}
	}

	_, err = qtx.CreateEvent(ctx, db.CreateEventParams{
		ContentID:        content.ID,
		EventName:        req.EventName,
		EventDate:        pgtype.Timestamptz{Time: eventDate, Valid: true},
		StartTime:        st,
		EndTime:          et,
		OrganizerName:    req.OrganizerName,
		OrganizerContact: contact,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccess(w, http.StatusCreated, map[string]string{"contentId": uuidToString(content.ID)}, "Event submitted for moderation")
}

func writeSuccess(w http.ResponseWriter, statusCode int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
		"message": message,
		"errors":  []interface{}{},
	})
}
