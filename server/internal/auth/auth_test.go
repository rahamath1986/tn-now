package auth

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"tn-now/server/internal/config"
	"tn-now/server/internal/db"
)

func TestPasswordHashing(t *testing.T) {
	password := "mypassword123"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		t.Error("Hashed password comparison failed")
	}

	err = bcrypt.CompareHashAndPassword(hashed, []byte("wrongpassword"))
	if err == nil {
		t.Error("Expected error for mismatched password, got nil")
	}
}

func TestJWTTokenLifecycle(t *testing.T) {
	secret := "mysecretkeyforauthsignings!!!"
	issuer := "test-issuer"
	audience := "test-audience"
	userID := "b9683c3e-8fbf-4a39-ae7f-fb18dcf92a9c"
	role := "CONTRIBUTOR"

	token, err := GenerateAccessToken(userID, role, secret, issuer, audience)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := VerifyAccessToken(token, secret, issuer, audience)
	if err != nil {
		t.Fatalf("Failed to verify token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected UserID %s, got %s", userID, claims.UserID)
	}

	if claims.Role != role {
		t.Errorf("Expected Role %s, got %s", role, claims.Role)
	}

	// Verify mismatched secret fails
	_, err = VerifyAccessToken(token, "wrong-secret-key-goes-here", issuer, audience)
	if err == nil {
		t.Error("Expected verification failure with wrong secret, got nil")
	}
}

func TestRegistrationAndLoginFlow(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping auth database integration tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Skipf("Failed to connect to database at %s, skipping test: %v", dbURL, err)
	}
	defer pool.Close()

	// Apply migrations
	migrationDir := filepath.Join("..", "..", "migrations")
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		t.Fatalf("Failed to read migrations directory: %v", err)
	}

	var upMigrations []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			upMigrations = append(upMigrations, file.Name())
		}
	}
	sort.Strings(upMigrations)

	tx, err := pool.Begin(context.Background())
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	for _, filename := range upMigrations {
		filePath := filepath.Join(migrationDir, filename)
		sqlBytes, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read migration file %s: %v", filePath, err)
		}
		_, err = tx.Exec(context.Background(), string(sqlBytes))
		if err != nil {
			t.Fatalf("Failed to execute migration %s: %v", filePath, err)
		}
	}

	queries := db.New(tx)
	cfg := config.Load()
	authService := NewService(queries, pool, cfg)

	// Since authService.Register creates its own transaction using s.conn, we must commit our schema setup transaction before running Service tests so they run on connection cleanly!
	err = tx.Commit(context.Background())
	if err != nil {
		t.Fatalf("Failed to commit migrations transaction: %v", err)
	}

	// Clean up after test runs
	defer func() {
		pool.Exec(context.Background(), "TRUNCATE users CASCADE;")
	}()

	// 1. Test Register
	username := "auth_test_user"
	email := "auth_test@tnnow.in"
	password := "securepassword123"

	resp, err := authService.Register(context.Background(), username, email, password)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if resp.Username != username {
		t.Errorf("Expected username %s, got %s", username, resp.Username)
	}

	if resp.Email != email {
		t.Errorf("Expected email %s, got %s", email, resp.Email)
	}

	// 2. Test Login
	loginResp, err := authService.Login(context.Background(), email, password, "mobile-test-device")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if loginResp.UserID != resp.UserID {
		t.Errorf("Expected matching user IDs on login, got %s vs %s", resp.UserID, loginResp.UserID)
	}

	// 3. Test Token Rotation (Refresh)
	refreshResp, err := authService.Refresh(context.Background(), loginResp.RefreshToken, "mobile-test-device")
	if err != nil {
		t.Fatalf("Refresh failed: %v", err)
	}

	if refreshResp.AccessToken == "" || refreshResp.RefreshToken == "" {
		t.Error("Expected rotated access and refresh tokens, got empty strings")
	}

	// 4. Test Logout
	err = authService.Logout(context.Background(), refreshResp.RefreshToken)
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}

	// Verify refreshed token is now revoked
	_, err = authService.Refresh(context.Background(), refreshResp.RefreshToken, "mobile-test-device")
	if err == nil {
		t.Error("Expected refresh call using logged out token to fail, got nil error")
	}
}
