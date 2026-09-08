package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password123@localhost:5432/tnnow_dev?sslmode=disable"
		logger.Info("DATABASE_URL not set, using default local development database", slog.String("db", dbURL))
	}

	logger.Info("Connecting to database for migrations...", slog.String("db", dbURL))

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		logger.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Ensure schema_migrations table exists
	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		logger.Error("Failed to create schema_migrations table", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Backfill historical migrations if existing tables are already present
	var tableExists bool
	_ = conn.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users')").Scan(&tableExists)
	if tableExists {
		_, _ = conn.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ('000001_init_schema.up.sql') ON CONFLICT DO NOTHING")
	}
	var sessionsExists bool
	_ = conn.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'user_sessions')").Scan(&sessionsExists)
	if sessionsExists {
		_, _ = conn.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ('000002_add_sessions.up.sql') ON CONFLICT DO NOTHING")
	}

	// Scan migrations directory
	migrationDir := "migrations"
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		migrationDir = filepath.Join("..", "migrations")
		files, err = os.ReadDir(migrationDir)
		if err != nil {
			logger.Error("Failed to read migrations directory", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}

	var upMigrations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			upMigrations = append(upMigrations, file.Name())
		}
	}
	sort.Strings(upMigrations)

	for _, filename := range upMigrations {
		var alreadyApplied bool
		err := conn.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = $1)", filename).Scan(&alreadyApplied)
		if err == nil && alreadyApplied {
			logger.Info("Migration already applied, skipping", slog.String("file", filename))
			continue
		}

		filePath := filepath.Join(migrationDir, filename)
		logger.Info("Applying schema migration file", slog.String("file", filePath))

		sqlBytes, err := os.ReadFile(filePath)
		if err != nil {
			logger.Error("Failed to read migration file", slog.String("file", filePath), slog.String("error", err.Error()))
			os.Exit(1)
		}

		tx, err := conn.Begin(ctx)
		if err != nil {
			logger.Error("Failed to start transaction", slog.String("error", err.Error()))
			os.Exit(1)
		}

		if _, err = tx.Exec(ctx, string(sqlBytes)); err != nil {
			_ = tx.Rollback(ctx)
			logger.Error("Failed to execute migration SQL", slog.String("file", filePath), slog.String("error", err.Error()))
			os.Exit(1)
		}

		if _, err = tx.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", filename); err != nil {
			_ = tx.Rollback(ctx)
			logger.Error("Failed to record migration version", slog.String("file", filename), slog.String("error", err.Error()))
			os.Exit(1)
		}

		if err = tx.Commit(ctx); err != nil {
			logger.Error("Failed to commit migration", slog.String("file", filename), slog.String("error", err.Error()))
			os.Exit(1)
		}

		logger.Info("Successfully applied migration", slog.String("file", filename))
	}

	logger.Info("All database migrations completed successfully.")
}
