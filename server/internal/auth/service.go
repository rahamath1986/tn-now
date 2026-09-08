package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"tn-now/server/internal/config"
	"tn-now/server/internal/db"
)

type Service struct {
	q    *db.Queries
	conn *pgxpool.Pool
	cfg  *config.Config
}

type AuthResponse struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	UserID       string    `json:"userId"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
}

func NewService(q *db.Queries, conn *pgxpool.Pool, cfg *config.Config) *Service {
	return &Service{q: q, conn: conn, cfg: cfg}
}

// Register creates a new user, profile, contributor record and issues tokens in a transaction
func (s *Service) Register(ctx context.Context, username, email, password string) (*AuthResponse, error) {
	// Hash password
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	passwordHash := string(hashedBytes)

	// Begin transaction
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	// Create user
	user, err := qtx.CreateUser(ctx, db.CreateUserParams{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         "CONTRIBUTOR", // default signup role
	})
	if err != nil {
		return nil, err
	}

	// Create profile with default username as name
	_, err = qtx.CreateUserProfile(ctx, db.CreateUserProfileParams{
		UserID:   user.ID,
		FullName: username,
	})
	if err != nil {
		return nil, err
	}

	// Create contributor profile
	_, err = qtx.CreateContributorProfile(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// Generate access & refresh tokens
	userIDStr := user.ID.String()
	accessToken, err := GenerateAccessToken(userIDStr, user.Role, s.cfg.JWTSecret, s.cfg.JWTIssuer, s.cfg.JWTAudience)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Save session record (7 days expiry)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	_, err = s.q.CreateSession(ctx, db.CreateSessionParams{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       userIDStr,
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
	}, nil
}

// Login validates user credentials and returns new session tokens
func (s *Service) Login(ctx context.Context, email, password, deviceInfo string) (*AuthResponse, error) {
	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	userIDStr := user.ID.String()
	accessToken, err := GenerateAccessToken(userIDStr, user.Role, s.cfg.JWTSecret, s.cfg.JWTIssuer, s.cfg.JWTAudience)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Save session
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	_, err = s.q.CreateSession(ctx, db.CreateSessionParams{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		DeviceInfo:   pgtype.Text{String: deviceInfo, Valid: deviceInfo != ""},
		ExpiresAt:    pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       userIDStr,
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
	}, nil
}

// Refresh rotates session refresh token and issues a new access token
func (s *Service) Refresh(ctx context.Context, oldRefreshToken, deviceInfo string) (*AuthResponse, error) {
	session, err := s.q.GetSessionByToken(ctx, oldRefreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("invalid session")
		}
		return nil, err
	}

	if session.IsRevoked {
		return nil, errors.New("session has been revoked")
	}

	if !session.ExpiresAt.Valid || session.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("session has expired")
	}

	// Retrieve user profile information
	user, err := s.q.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	// Revoke old session
	_, err = s.q.RevokeSession(ctx, oldRefreshToken)
	if err != nil {
		return nil, err
	}

	userIDStr := user.ID.String()
	accessToken, err := GenerateAccessToken(userIDStr, user.Role, s.cfg.JWTSecret, s.cfg.JWTIssuer, s.cfg.JWTAudience)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Create new rotated session
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	_, err = s.q.CreateSession(ctx, db.CreateSessionParams{
		UserID:       user.ID,
		RefreshToken: newRefreshToken,
		DeviceInfo:   pgtype.Text{String: deviceInfo, Valid: deviceInfo != ""},
		ExpiresAt:    pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		UserID:       userIDStr,
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
	}, nil
}

// Logout invalidates the active refresh token session
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	_, err := s.q.RevokeSession(ctx, refreshToken)
	return err
}
