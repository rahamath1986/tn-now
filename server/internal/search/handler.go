package search

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"tn-now/server/internal/db"
)

type Handler struct {
	q *db.Queries
}

func NewHandler(q *db.Queries) *Handler {
	return &Handler{q: q}
}

type SearchResultDTO struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	ContentType string `json:"contentType"`
	CategoryID  string `json:"categoryId"`
	DistrictID  string `json:"districtId"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/search", h.HandleSearch)
}

func (h *Handler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		writeError(w, http.StatusBadRequest, "q parameter is required")
		return
	}

	limit := 20
	if lStr := r.URL.Query().Get("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}

	ctx := r.Context()
	dbResults, err := h.q.SearchContent(ctx, db.SearchContentParams{
		Column1: pgtype.Text{String: query, Valid: true},
		Limit:   int32(limit),
	})
	if err != nil {
		slog.Error("Search query execution failed", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Search query failed")
		return
	}

	results := make([]SearchResultDTO, len(dbResults))
	for i, item := range dbResults {
		var desc string
		if item.Description.Valid {
			desc = item.Description.String
		}

		results[i] = SearchResultDTO{
			ID:          uuidToString(item.ID),
			Title:       item.Title,
			Description: desc,
			ContentType: item.ContentType,
			CategoryID:  uuidToString(item.CategoryID),
			DistrictID:  uuidToString(item.DistrictID),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    results,
		"message": "Search results retrieved",
		"errors":  []interface{}{},
	})
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
