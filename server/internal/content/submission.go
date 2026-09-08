package content

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"tn-now/server/internal/db"
	"tn-now/server/internal/video"
)

type SubmissionHandler struct {
	q    *db.Queries
	conn *pgxpool.Pool
}

func NewSubmissionHandler(q *db.Queries, conn *pgxpool.Pool) *SubmissionHandler {
	return &SubmissionHandler{q: q, conn: conn}
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

	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Body) == "" || req.CategoryID == "" || req.DistrictID == "" || req.AuthorUserID == "" {
		writeError(w, http.StatusUnprocessableEntity, "Title, Body, CategoryID, DistrictID, and AuthorUserID are required")
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
		ContentType:  "TEXT_STORY",
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

	var photoURL pgtype.Text
	if req.SinglePhotoURL != "" {
		photoURL = pgtype.Text{String: req.SinglePhotoURL, Valid: true}
	}

	headline := req.Headline
	if headline == "" {
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

	writeSuccess(w, http.StatusCreated, map[string]string{"contentId": uuidToString(content.ID)}, "Story submitted for moderation")
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
