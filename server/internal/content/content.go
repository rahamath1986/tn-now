package content

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"tn-now/server/internal/db"
)

type ContentHandler struct {
	q    *db.Queries
	conn *pgxpool.Pool
}

func NewContentHandler(q *db.Queries, conn *pgxpool.Pool) *ContentHandler {
	return &ContentHandler{q: q, conn: conn}
}

type PaginatedContentResponse struct {
	Items          []ContentDTO `json:"items"`
	NextCursorTime string       `json:"nextCursorTime,omitempty"`
	NextCursorID   string       `json:"nextCursorId,omitempty"`
	HasMore        bool         `json:"hasMore"`
}

func (h *ContentHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/content", h.HandleListContent)
}

func (h *ContentHandler) HandleListContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"success":false,"data":null,"message":"Method not allowed","errors":[]}`))
		return
	}

	categoryIDStr := r.URL.Query().Get("category_id")
	if categoryIDStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"success":false,"data":null,"message":"category_id parameter is required","errors":[]}`))
		return
	}

	ctx := r.Context()

	// Parse Limit
	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if limit > 50 {
		limit = 50
	}

	// Parse Cursor ID
	cursorIDStr := r.URL.Query().Get("cursor_id")
	if cursorIDStr == "" {
		cursorIDStr = "ffffffff-ffff-ffff-ffff-ffffffffffff"
	}

	// Parse Cursor Time
	var dbCursorTime pgtype.Timestamptz
	cursorTimeStr := r.URL.Query().Get("cursor_time")
	if cursorTimeStr != "" {
		t, err := time.Parse(time.RFC3339, cursorTimeStr)
		if err == nil {
			dbCursorTime = pgtype.Timestamptz{Time: t, Valid: true}
		} else {
			dbCursorTime = pgtype.Timestamptz{Time: time.Now().Add(24 * time.Hour), Valid: true}
		}
	} else {
		dbCursorTime = pgtype.Timestamptz{Time: time.Now().Add(24 * time.Hour), Valid: true}
	}

	dbCategoryID := stringToUUID(categoryIDStr)
	dbCursorID := stringToUUID(cursorIDStr)

	// Fetch database contents
	dbItems, err := h.q.ListContentByCategoryPaginated(ctx, db.ListContentByCategoryPaginatedParams{
		CategoryID:  dbCategoryID,
		PublishedAt: dbCursorTime,
		ID:          dbCursorID,
		Limit:       int32(limit),
	})
	if err != nil {
		slog.Error("Failed to fetch paginated content by category", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	items := make([]ContentDTO, len(dbItems))
	for i, item := range dbItems {
		var pubAt *time.Time
		if item.PublishedAt.Valid {
			t := item.PublishedAt.Time
			pubAt = &t
		}
		var desc string
		if item.Description.Valid {
			desc = item.Description.String
		}

		items[i] = ContentDTO{
			ID:          uuidToString(item.ID),
			Title:       item.Title,
			Description: desc,
			ContentType: item.ContentType,
			CategoryID:  uuidToString(item.CategoryID),
			DistrictID:  uuidToString(item.DistrictID),
			PublishedAt: pubAt,
			CreatedAt:   item.CreatedAt.Time,
		}
	}

	hasMore := len(items) == limit
	var nextCursorTime string
	var nextCursorID string

	if hasMore && len(dbItems) > 0 {
		lastItem := dbItems[len(dbItems)-1]
		if lastItem.PublishedAt.Valid {
			nextCursorTime = lastItem.PublishedAt.Time.Format(time.RFC3339)
		}
		nextCursorID = uuidToString(lastItem.ID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": PaginatedContentResponse{
			Items:          items,
			NextCursorTime: nextCursorTime,
			NextCursorID:   nextCursorID,
			HasMore:        hasMore,
		},
		"message": nil,
		"errors":  []interface{}{},
	})
}

func stringToUUID(s string) pgtype.UUID {
	var u pgtype.UUID
	_ = u.Scan(s)
	return u
}
