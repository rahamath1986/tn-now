package reputation

import (
	"encoding/json"
	"log/slog"
	"net/http"
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

type grantFoundingRequest struct {
	UserID string `json:"userId"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/contributor/profile", h.HandleGetProfile)
	mux.HandleFunc("/contributor/founding", h.HandleGrantFounding)
}

func (h *Handler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id parameter is required")
		return
	}

	ctx := r.Context()
	stats, err := GetContributorStats(ctx, h.q, userID)
	if err != nil {
		slog.Error("Failed to fetch contributor profile stats", slog.String("error", err.Error()))
		writeError(w, http.StatusNotFound, "Contributor profile not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    stats,
		"message": "Contributor profile retrieved",
		"errors":  []interface{}{},
	})
}

func (h *Handler) HandleGrantFounding(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req grantFoundingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.UserID == "" {
		writeError(w, http.StatusUnprocessableEntity, "userId is required")
		return
	}

	ctx := r.Context()
	var uID pgtype.UUID
	_ = uID.Scan(req.UserID)

	now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	res, err := h.q.UpdateFoundingStatus(ctx, db.UpdateFoundingStatusParams{
		UserID:                 uID,
		IsFoundingContributor:  true,
		FoundingBadgeGrantedAt: now,
	})
	if err != nil {
		slog.Error("Failed to update founding status", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Failed to grant founding badge")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    map[string]interface{}{"userId": req.UserID, "isFoundingContributor": res.IsFoundingContributor},
		"message": "Founding Contributor 2026 badge granted successfully",
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
