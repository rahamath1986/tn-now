package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Warn("DATABASE_URL environment variable is not set. Worker running in idle stand-by mode.")
	}

	logger.Info("Starting TN NOW Dead Link Detection & Moderation Worker...")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Initial run check
	runCheck(ctx, logger, dbURL, httpClient)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Worker shutdown signal received, stopping cleanly...")
			return
		case <-ticker.C:
			runCheck(ctx, logger, dbURL, httpClient)
		}
	}
}

func runCheck(ctx context.Context, logger *slog.Logger, dbURL string, httpClient *http.Client) {
	if dbURL == "" {
		logger.Info("Worker tick ping: standing by for database connection...")
		return
	}

	connCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(connCtx, dbURL)
	if err != nil {
		logger.Error("Worker failed to connect to database", slog.String("error", err.Error()))
		return
	}
	defer conn.Close(context.Background())

	logger.Info("Worker scanning video_links for link verification checks...")

	// Fetch up to 10 video links needing check
	rows, err := conn.Query(ctx, "SELECT content_id, canonical_url, link_status FROM video_links WHERE link_status = 'UNRESOLVED' OR last_checked_at < NOW() - INTERVAL '24 hours' LIMIT 10")
	if err != nil {
		logger.Error("Failed to query video_links", slog.String("error", err.Error()))
		return
	}
	defer rows.Close()

	type item struct {
		contentID string
		url       string
		status    string
	}
	var items []item
	for rows.Next() {
		var contentID, rawURL, status string
		if err := rows.Scan(&contentID, &rawURL, &status); err == nil {
			items = append(items, item{contentID: contentID, url: rawURL, status: status})
		}
	}

	for _, it := range items {
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, it.url, nil)
		if err != nil {
			continue
		}

		resp, err := httpClient.Do(req)
		newStatus := "ALIVE"
		if err != nil || (resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode >= 400)) {
			newStatus = "DEAD"
		}
		if resp != nil {
			_ = resp.Body.Close()
		}

		_, _ = conn.Exec(ctx, "UPDATE video_links SET link_status = $1, last_checked_at = NOW() WHERE content_id = $2", newStatus, it.contentID)
		logger.Info("Updated video link status check", slog.String("content_id", it.contentID), slog.String("status", newStatus))
	}
}
