package content

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"tn-now/server/internal/db"
)

type HomeHandler struct {
	q    *db.Queries
	conn *pgxpool.Pool
}

func NewHomeHandler(q *db.Queries, conn *pgxpool.Pool) *HomeHandler {
	return &HomeHandler{q: q, conn: conn}
}

type CategoryDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ContentDTO struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	ContentType string     `json:"contentType"`
	CategoryID  string     `json:"categoryId"`
	DistrictID  string     `json:"districtId"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type EventDTO struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description,omitempty"`
	EventName     string    `json:"eventName"`
	EventDate     time.Time `json:"eventDate"`
	StartTime     string    `json:"startTime,omitempty"`
	EndTime       string    `json:"endTime,omitempty"`
	OrganizerName string    `json:"organizerName"`
}

type HomeFeedResponse struct {
	Categories []CategoryDTO `json:"categories"`
	Trending   []ContentDTO  `json:"trending"`
	Latest     []ContentDTO  `json:"latest"`
	Events     []EventDTO    `json:"events"`
}

func (h *HomeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/home", h.HandleGetHome)
}

func (h *HomeHandler) HandleGetHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"success":false,"data":null,"message":"Method not allowed","errors":[]}`))
		return
	}

	ctx := r.Context()

	// 1. Fetch categories
	dbCategories, err := h.q.ListCategories(ctx)
	if err != nil {
		slog.Error("Failed to fetch categories for home feed", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	categories := make([]CategoryDTO, len(dbCategories))
	for i, c := range dbCategories {
		categories[i] = CategoryDTO{
			ID:   uuidToString(c.ID),
			Name: c.Name,
			Slug: c.Slug,
		}
	}

	// 2. Fetch latest content (limit 10)
	dbLatest, err := h.q.GetLatestContent(ctx, 10)
	if err != nil {
		slog.Error("Failed to fetch latest content for home feed", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	latest := make([]ContentDTO, len(dbLatest))
	for i, item := range dbLatest {
		var pubAt *time.Time
		if item.PublishedAt.Valid {
			t := item.PublishedAt.Time
			pubAt = &t
		}
		var desc string
		if item.Description.Valid {
			desc = item.Description.String
		}

		latest[i] = ContentDTO{
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

	// 3. Fetch upcoming events
	today := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	dbEvents, err := h.q.GetUpcomingEvents(ctx, db.GetUpcomingEventsParams{
		EventDate: today,
		Limit:     5,
	})
	if err != nil {
		slog.Error("Failed to fetch upcoming events for home feed", slog.String("error", err.Error()))
		writeError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	events := make([]EventDTO, len(dbEvents))
	for i, ev := range dbEvents {
		var desc string
		if ev.Description.Valid {
			desc = ev.Description.String
		}
		var st, et string
		if ev.StartTime.Valid {
			st = ev.StartTime.String
		}
		if ev.EndTime.Valid {
			et = ev.EndTime.String
		}

		events[i] = EventDTO{
			ID:            uuidToString(ev.ID),
			Title:         ev.Title,
			Description:   desc,
			EventName:     ev.EventName,
			EventDate:     ev.EventDate.Time,
			StartTime:     st,
			EndTime:       et,
			OrganizerName: ev.OrganizerName,
		}
	}

	// 4. Populate trending placeholder (simply uses the latest list for now)
	trending := latest

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": HomeFeedResponse{
			Categories: categories,
			Trending:   trending,
			Latest:     latest,
			Events:     events,
		},
		"message": nil,
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
	_, _ = w.Write([]byte(`{"success":false,"data":null,"message":"` + message + `","errors":[]}`))
}
