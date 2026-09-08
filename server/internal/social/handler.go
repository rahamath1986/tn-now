package social

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"tn-now/server/internal/db"
)

type Handler struct {
	q    *db.Queries
	conn *pgxpool.Pool
}

func NewHandler(q *db.Queries, conn *pgxpool.Pool) *Handler {
	return &Handler{q: q, conn: conn}
}

type likeRequest struct {
	ContentID string `json:"contentId"`
	UserID    string `json:"userId"`
	Action    string `json:"action"` // 'like' | 'unlike'
}

type bookmarkRequest struct {
	ContentID string `json:"contentId"`
	UserID    string `json:"userId"`
	Action    string `json:"action"` // 'bookmark' | 'unbookmark'
}

type commentRequest struct {
	ContentID string `json:"contentId"`
	UserID    string `json:"userId"`
	Body      string `json:"body"`
}

type CommentDTO struct {
	ID        string    `json:"id"`
	ContentID string    `json:"contentId"`
	UserID    string    `json:"userId"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/social/like", h.HandleLike)
	mux.HandleFunc("/social/bookmark", h.HandleBookmark)
	mux.HandleFunc("/social/comment", h.HandleComment)
	mux.HandleFunc("/social/comments", h.HandleGetComments)
}

func (h *Handler) HandleLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req likeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.ContentID == "" || req.UserID == "" {
		writeError(w, http.StatusUnprocessableEntity, "contentId and userId are required")
		return
	}

	ctx := r.Context()
	cID := stringToUUID(req.ContentID)
	uID := stringToUUID(req.UserID)

	if req.Action == "unlike" {
		err := h.q.DeleteLike(ctx, db.DeleteLikeParams{ContentID: cID, UserID: uID})
		if err != nil {
			slog.Error("Failed to remove like", slog.String("error", err.Error()))
			writeError(w, http.StatusInternalServerError, "Failed to remove like")
			return
		}
		writeSuccess(w, http.StatusOK, map[string]string{"contentId": req.ContentID, "liked": "false"}, "Like removed")
		return
	}

	_, err := h.q.CreateLike(ctx, db.CreateLikeParams{ContentID: cID, UserID: uID})
	if err != nil {
		slog.Error("Failed to create like", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Failed to add like")
		return
	}

	writeSuccess(w, http.StatusOK, map[string]string{"contentId": req.ContentID, "liked": "true"}, "Like added")
}

func (h *Handler) HandleBookmark(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req bookmarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.ContentID == "" || req.UserID == "" {
		writeError(w, http.StatusUnprocessableEntity, "contentId and userId are required")
		return
	}

	ctx := r.Context()
	cID := stringToUUID(req.ContentID)
	uID := stringToUUID(req.UserID)

	if req.Action == "unbookmark" {
		err := h.q.DeleteBookmark(ctx, db.DeleteBookmarkParams{ContentID: cID, UserID: uID})
		if err != nil {
			slog.Error("Failed to remove bookmark", slog.String("error", err.Error()))
			writeError(w, http.StatusInternalServerError, "Failed to remove bookmark")
			return
		}
		writeSuccess(w, http.StatusOK, map[string]string{"contentId": req.ContentID, "bookmarked": "false"}, "Bookmark removed")
		return
	}

	_, err := h.q.CreateBookmark(ctx, db.CreateBookmarkParams{ContentID: cID, UserID: uID})
	if err != nil {
		slog.Error("Failed to add bookmark", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Failed to add bookmark")
		return
	}

	writeSuccess(w, http.StatusOK, map[string]string{"contentId": req.ContentID, "bookmarked": "true"}, "Bookmark added")
}

func (h *Handler) HandleComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req commentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.ContentID == "" || req.UserID == "" || strings.TrimSpace(req.Body) == "" {
		writeError(w, http.StatusUnprocessableEntity, "contentId, userId, and non-empty body are required")
		return
	}

	ctx := r.Context()
	cID := stringToUUID(req.ContentID)
	uID := stringToUUID(req.UserID)

	comment, err := h.q.CreateComment(ctx, db.CreateCommentParams{
		ContentID: cID,
		UserID:    uID,
		Body:      req.Body,
	})
	if err != nil {
		slog.Error("Failed to create comment", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Failed to add comment")
		return
	}

	dto := CommentDTO{
		ID:        uuidToString(comment.ID),
		ContentID: uuidToString(comment.ContentID),
		UserID:    uuidToString(comment.UserID),
		Body:      comment.Body,
		CreatedAt: comment.CreatedAt.Time,
	}

	writeSuccess(w, http.StatusCreated, dto, "Comment posted successfully")
}

func (h *Handler) HandleGetComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	contentID := r.URL.Query().Get("content_id")
	if contentID == "" {
		writeError(w, http.StatusBadRequest, "content_id parameter is required")
		return
	}

	ctx := r.Context()
	cID := stringToUUID(contentID)

	dbComments, err := h.q.ListCommentsByContentID(ctx, cID)
	if err != nil {
		slog.Error("Failed to fetch comments", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	comments := make([]CommentDTO, len(dbComments))
	for i, c := range dbComments {
		comments[i] = CommentDTO{
			ID:        uuidToString(c.ID),
			ContentID: uuidToString(c.ContentID),
			UserID:    uuidToString(c.UserID),
			Body:      c.Body,
			CreatedAt: c.CreatedAt.Time,
		}
	}

	writeSuccess(w, http.StatusOK, comments, "Comments retrieved")
}

func stringToUUID(s string) pgtype.UUID {
	var u pgtype.UUID
	_ = u.Scan(s)
	return u
}

func uuidToString(uuid pgtype.UUID) string {
	if !uuid.Valid {
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid.Bytes[0:4], uuid.Bytes[4:6], uuid.Bytes[6:8], uuid.Bytes[8:10], uuid.Bytes[10:])
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

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"data":    nil,
		"message": message,
		"errors":  []interface{}{},
	})
}
