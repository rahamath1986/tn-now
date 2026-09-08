package moderation

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
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

type decideRequest struct {
	ContentID       string `json:"contentId"`
	Action          string `json:"action"` // 'APPROVE' | 'REJECT' | 'QUARANTINE'
	ModeratorUserID string `json:"moderatorUserId"`
	Reason          string `json:"reason,omitempty"`
}

type QueueItemDTO struct {
	ID               string    `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description,omitempty"`
	ContentType      string    `json:"contentType"`
	CategoryID       string    `json:"categoryId"`
	DistrictID       string    `json:"districtId"`
	AuthorUserID     string    `json:"authorUserId"`
	Status           string    `json:"status"`
	ModerationStatus string    `json:"moderationStatus"`
	CreatedAt        time.Time `json:"createdAt"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/moderation/queue", h.HandleGetQueue)
	mux.HandleFunc("/moderation/decide", h.HandleDecideContent)
}

func (h *Handler) HandleGetQueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	limit := 20
	if lStr := r.URL.Query().Get("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}

	ctx := r.Context()
	dbItems, err := h.q.ListPendingModerationQueue(ctx, int32(limit))
	if err != nil {
		slog.Error("Failed to fetch moderation queue", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	items := make([]QueueItemDTO, len(dbItems))
	for i, item := range dbItems {
		var desc string
		if item.Description.Valid {
			desc = item.Description.String
		}

		items[i] = QueueItemDTO{
			ID:               uuidToString(item.ID),
			Title:            item.Title,
			Description:      desc,
			ContentType:      item.ContentType,
			CategoryID:       uuidToString(item.CategoryID),
			DistrictID:       uuidToString(item.DistrictID),
			AuthorUserID:     uuidToString(item.AuthorUserID),
			Status:           item.Status,
			ModerationStatus: item.ModerationStatus,
			CreatedAt:        item.CreatedAt.Time,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    items,
		"message": "Moderation queue retrieved",
		"errors":  []interface{}{},
	})
}

func (h *Handler) HandleDecideContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req decideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.ContentID == "" || req.Action == "" || req.ModeratorUserID == "" {
		writeError(w, http.StatusUnprocessableEntity, "ContentID, Action, and ModeratorUserID are required")
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

	var newStatus string
	var newModStatus string
	var pubAt pgtype.Timestamptz

	switch req.Action {
	case "APPROVE":
		newStatus = "PUBLISHED"
		newModStatus = "SAFE"
		pubAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	case "REJECT":
		newStatus = "REJECTED"
		newModStatus = "REJECTED"
	case "QUARANTINE":
		newStatus = "QUARANTINED"
		newModStatus = "QUARANTINE"
	default:
		writeError(w, http.StatusUnprocessableEntity, "Invalid action. Allowed: APPROVE, REJECT, QUARANTINE")
		return
	}

	contentID := stringToUUID(req.ContentID)
	modUserID := stringToUUID(req.ModeratorUserID)

	// Update Content master record
	_, err = qtx.UpdateContentStatus(ctx, db.UpdateContentStatusParams{
		ID:               contentID,
		Status:           newStatus,
		ModerationStatus: newModStatus,
		PublishedAt:      pubAt,
	})
	if err != nil {
		slog.Error("Failed to update content moderation status", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Failed to update content status")
		return
	}

	// Create Moderation Result Record
	_, err = qtx.CreateModerationResult(ctx, db.CreateModerationResultParams{
		ContentID:         contentID,
		AutoModerated:     false,
		TextToxicityScore: 0.0,
		ImageSafetyScore:  0.0,
		FlagsRaised:       []string{},
		Decision:          newModStatus,
	})
	if err != nil {
		slog.Error("Failed to insert moderation result", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Failed to record moderation result")
		return
	}

	// Create Audit Log Entry
	details := fmt.Sprintf("Action: %s, Reason: %s", req.Action, req.Reason)
	_, err = qtx.CreateAuditLog(ctx, db.CreateAuditLogParams{
		ActorUserID:  modUserID,
		Action:       "MODERATION_" + req.Action,
		TargetEntity: "content",
		TargetID:     contentID,
		Details:      pgtype.Text{String: details, Valid: true},
	})
	if err != nil {
		slog.Error("Failed to create audit log", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Failed to write audit trail")
		return
	}

	if err := tx.Commit(ctx); err != nil {
		writeError(w, http.StatusInternalServerError, "Transaction commit failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    map[string]string{"contentId": req.ContentID, "status": newStatus},
		"message": "Moderation decision applied successfully",
		"errors":  []interface{}{},
	})
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
