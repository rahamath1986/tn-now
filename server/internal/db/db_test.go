package db

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

func TestDatabaseMigrationAndQueries(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping database integration tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		t.Skipf("Failed to connect to database at %s, skipping test: %v", dbURL, err)
	}
	defer conn.Close(context.Background())

	// Run migration up SQL to verify the schema compiles and runs
	migrationPath := filepath.Join("..", "..", "migrations", "000001_init_schema.up.sql")
	sqlBytes, err := os.ReadFile(migrationPath)
	if err != nil {
		t.Fatalf("Failed to read migration SQL file from path %s: %v", migrationPath, err)
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), string(sqlBytes))
	if err != nil {
		t.Fatalf("Failed to execute schema migration SQL: %v", err)
	}

	// Verify pre-seeded districts using sqlc generated code
	queries := New(tx)
	districts, err := queries.ListDistricts(context.Background())
	if err != nil {
		t.Fatalf("Failed to query districts using generated query accessor: %v", err)
	}

	if len(districts) == 0 {
		t.Error("Expected seeded districts, got 0")
	}

	// Verify pre-seeded categories
	categories, err := queries.ListCategories(context.Background())
	if err != nil {
		t.Fatalf("Failed to query categories: %v", err)
	}

	if len(categories) == 0 {
		t.Error("Expected seeded categories, got 0")
	}

	// Test user creation and retrieval
	user, err := queries.CreateUser(context.Background(), CreateUserParams{
		Username:     "test_dev_user",
		Email:        "test_dev@tnnow.in",
		PasswordHash: "$2a$12$securepasswordhashhere",
		Role:         "CONTRIBUTOR",
	})
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	retrievedUser, err := queries.GetUserByID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("Failed to fetch user by ID: %v", err)
	}

	if retrievedUser.Username != "test_dev_user" {
		t.Errorf("Expected username test_dev_user, got %s", retrievedUser.Username)
	}
}
