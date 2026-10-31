package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserID       string `json:"userId"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Role         string `json:"role"`
}

type RegisterResponse struct {
	RequiresOTP bool   `json:"requiresOtp"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Message     string `json:"message"`
}

type UpdateProfileParams struct {
	FullName    string `json:"fullName"`
	Bio         string `json:"bio"`
	AvatarURL   string `json:"avatarUrl"`
	NewPassword string `json:"newPassword,omitempty"`
}

func NewService(q *db.Queries, conn *pgxpool.Pool, cfg *config.Config) *Service {
	return &Service{q: q, conn: conn, cfg: cfg}
}

// Register creates a new user, profile, contributor record with is_verified=false, and dispatches an Email OTP
func (s *Service) Register(ctx context.Context, username, email, password, fullName string) (*RegisterResponse, error) {
	username = strings.TrimSpace(username)
	email = strings.ToLower(strings.TrimSpace(email))
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		fullName = username
	}

	// Check if username already exists
	var count int
	_ = s.conn.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if count > 0 {
		return nil, errors.New("இந்த பயனர் பெயர் ஏற்கனவே பயன்பாட்டில் உள்ளது (Username is already taken)")
	}

	// Check if email already exists
	_ = s.conn.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if count > 0 {
		return nil, errors.New("இந்த மின்னஞ்சல் முகவரியில் ஏற்கனவே கணக்கு உள்ளது (Email already registered)")
	}

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

	// Explicitly mark user as not verified until OTP confirmation
	_, _ = tx.Exec(ctx, "UPDATE users SET is_verified = FALSE WHERE id = $1", user.ID)

	// Create profile
	_, err = qtx.CreateUserProfile(ctx, db.CreateUserProfileParams{
		UserID:   user.ID,
		FullName: fullName,
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

	// Generate 6-digit numeric OTP and send verification email
	otp, err := GenerateNumericOTP()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(10 * time.Minute)
	_, err = s.conn.Exec(ctx, `
		INSERT INTO email_otps (email, otp_code, attempts, expires_at, created_at)
		VALUES ($1, $2, 0, $3, NOW())
	`, email, otp, expiresAt)
	if err != nil {
		return nil, err
	}

	// Dispatch OTP email using effective SMTP settings (DB settings with env fallback)
	smtpCfg := GetEffectiveSMTPConfig(ctx, s.conn, s.cfg)
	_ = SendOTPEmail(smtpCfg, email, otp)

	return &RegisterResponse{
		RequiresOTP: true,
		Email:       email,
		Username:    username,
		Message:     "பதிவு வெற்றிகரமானது. உங்கள் மின்னஞ்சலுக்கு அனுப்பப்பட்ட OTP சரிபார்ப்புக் குறியீட்டை உள்ளிடவும்.",
	}, nil
}

// Login validates user credentials and returns new session tokens
func (s *Service) Login(ctx context.Context, email, password, deviceInfo string) (*AuthResponse, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := s.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("மின்னஞ்சல் அல்லது கடவுச்சொல் தவறானது (Invalid email or password)")
		}
		return nil, err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("மின்னஞ்சல் அல்லது கடவுச்சொல் தவறானது (Invalid email or password)")
	}

	// Check if user is verified
	var isVerified bool
	_ = s.conn.QueryRow(ctx, "SELECT is_verified FROM users WHERE id = $1", user.ID).Scan(&isVerified)
	if !isVerified {
		// Send a fresh OTP to email
		_ = s.SendEmailOTP(ctx, user.Email)
		return nil, errors.New("EMAIL_NOT_VERIFIED")
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

type UserProfileResponse struct {
	UserID        string `json:"userId"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	FullName      string `json:"fullName"`
	Bio           string `json:"bio"`
	AvatarURL     string `json:"avatarUrl"`
	IsVerified    bool   `json:"isVerified"`
	Points        int    `json:"points"`
	Level         string `json:"level"`
	ApprovedCount int    `json:"approvedCount"`
}

// GetMe retrieves user account and profile information by user ID
func (s *Service) GetMe(ctx context.Context, userID string) (*UserProfileResponse, error) {
	var uid pgtype.UUID
	if err := uid.Scan(userID); err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := s.q.GetUserByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	profile, _ := s.q.GetUserProfile(ctx, uid)

	fullName := profile.FullName
	if strings.TrimSpace(fullName) == "" {
		fullName = user.Username
	}

	var isVerified bool
	_ = s.conn.QueryRow(ctx, "SELECT is_verified FROM users WHERE id = $1", uid).Scan(&isVerified)

	var points, approvedCount int
	var level string
	_ = s.conn.QueryRow(ctx, "SELECT points, level, approved_count FROM contributor_profiles WHERE user_id = $1", uid).Scan(&points, &level, &approvedCount)
	if level == "" {
		level = "New Contributor"
	}

	return &UserProfileResponse{
		UserID:        user.ID.String(),
		Username:      user.Username,
		Email:         user.Email,
		Role:          user.Role,
		FullName:      fullName,
		Bio:           profile.Bio.String,
		AvatarURL:     profile.AvatarUrl.String,
		IsVerified:    isVerified,
		Points:        points,
		Level:         level,
		ApprovedCount: approvedCount,
	}, nil
}

// UpdateUserProfile updates user's display name, bio, avatar, and optional password
func (s *Service) UpdateUserProfile(ctx context.Context, userIDStr string, params UpdateProfileParams) (*UserProfileResponse, error) {
	var uid pgtype.UUID
	if err := uid.Scan(userIDStr); err != nil {
		return nil, errors.New("invalid user id")
	}

	params.FullName = strings.TrimSpace(params.FullName)
	params.Bio = strings.TrimSpace(params.Bio)
	params.AvatarURL = strings.TrimSpace(params.AvatarURL)

	if params.FullName == "" {
		return nil, errors.New("முழுப் பெயர் (Full Name) அவசியம்")
	}

	// Update or insert user_profiles
	_, err := s.conn.Exec(ctx, `
		INSERT INTO user_profiles (user_id, full_name, bio, avatar_url, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			full_name = EXCLUDED.full_name,
			bio = EXCLUDED.bio,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = NOW();
	`, uid, params.FullName, params.Bio, params.AvatarURL)
	if err != nil {
		return nil, err
	}

	// If new password provided, update users table
	if strings.TrimSpace(params.NewPassword) != "" {
		if len(params.NewPassword) < 6 {
			return nil, errors.New("கடவுச்சொல் குறைந்தது 6 எழுத்துக்கள் இருக்க வேண்டும் (Password must be >= 6 chars)")
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(params.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		_, err = s.conn.Exec(ctx, "UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2", string(hashed), uid)
		if err != nil {
			return nil, err
		}
	}

	return s.GetMe(ctx, userIDStr)
}

// VerifyToken verifies access token claims
func (s *Service) VerifyToken(tokenString string) (*JWTClaims, error) {
	return VerifyAccessToken(tokenString, s.cfg.JWTSecret, s.cfg.JWTIssuer, s.cfg.JWTAudience)
}

// SendEmailOTP generates a 6-digit OTP, stores it with 10-minute expiry, and dispatches it via email
func (s *Service) SendEmailOTP(ctx context.Context, email string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return errors.New("செல்லுபடியாகும் மின்னஞ்சல் முகவரியை உள்ளிடவும் (Valid email required)")
	}

	// Rate limiting: prevent spamming OTP requests within 60 seconds
	var recentCount int
	_ = s.conn.QueryRow(ctx, "SELECT COUNT(*) FROM email_otps WHERE email = $1 AND created_at > NOW() - INTERVAL '60 seconds'", email).Scan(&recentCount)
	if recentCount > 0 {
		return errors.New("தயவுசெய்து 1 நிமிடம் காத்திருந்து புதிய OTP கோரவும் (Please wait 60s before requesting another OTP)")
	}

	otp, err := GenerateNumericOTP()
	if err != nil {
		return err
	}

	// Clear previous unconsumed OTPs for this email
	_, _ = s.conn.Exec(ctx, "DELETE FROM email_otps WHERE email = $1", email)

	// Save new OTP with 10 minute expiry
	_, err = s.conn.Exec(ctx, `
		INSERT INTO email_otps (email, otp_code, expires_at)
		VALUES ($1, $2, NOW() + INTERVAL '10 minutes')
	`, email, otp)
	if err != nil {
		return err
	}

	// Dispatch email using effective SMTP configuration (DB system_settings with env fallback)
	smtpCfg := GetEffectiveSMTPConfig(ctx, s.conn, s.cfg)
	return SendOTPEmail(smtpCfg, email, otp)
}

// VerifyEmailOTP checks the 6-digit OTP, logs in the user, or creates a new contributor account automatically
func (s *Service) VerifyEmailOTP(ctx context.Context, email, otp, deviceInfo string) (*AuthResponse, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	otp = strings.TrimSpace(otp)

	if email == "" || otp == "" {
		return nil, errors.New("மின்னஞ்சல் மற்றும் OTP இரண்டும் தேவை (Email and OTP are required)")
	}

	var otpID pgtype.UUID
	var storedOTP string
	var attempts int
	var expiresAt time.Time

	err := s.conn.QueryRow(ctx, `
		SELECT id, otp_code, attempts, expires_at
		FROM email_otps
		WHERE email = $1 AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`, email).Scan(&otpID, &storedOTP, &attempts, &expiresAt)

	if err != nil {
		return nil, errors.New("OTP தவறானது அல்லது காலாவதியானது (OTP is invalid or has expired)")
	}

	if attempts >= 5 {
		_, _ = s.conn.Exec(ctx, "DELETE FROM email_otps WHERE email = $1", email)
		return nil, errors.New("அதிக தவறான முயற்சிகள். புதிய OTP கோரவும் (Too many failed attempts. Request a new OTP)")
	}

	if storedOTP != otp {
		_, _ = s.conn.Exec(ctx, "UPDATE email_otps SET attempts = attempts + 1 WHERE id = $1", otpID)
		return nil, errors.New("தவறான OTP குறியீடு (Invalid OTP code)")
	}

	// Consume OTP so it cannot be reused
	_, _ = s.conn.Exec(ctx, "DELETE FROM email_otps WHERE email = $1", email)

	// Look up user by email
	user, err := s.q.GetUserByEmail(ctx, email)
	var finalUser db.User

	if err != nil {
		// New user -> Auto-register in transaction
		baseUsername := strings.Split(email, "@")[0]
		baseUsername = strings.ReplaceAll(baseUsername, ".", "_")
		baseUsername = strings.ReplaceAll(baseUsername, "+", "_")
		if len(baseUsername) > 20 {
			baseUsername = baseUsername[:20]
		}

		// Ensure unique username
		username := baseUsername
		var count int
		_ = s.conn.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
		if count > 0 {
			randSuffix, _ := GenerateNumericOTP()
			username = fmt.Sprintf("%s_%s", baseUsername, randSuffix[:4])
		}

		// Generate random secure password placeholder
		randomPass, _ := GenerateRefreshToken()
		hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(randomPass), bcrypt.DefaultCost)

		tx, txErr := s.conn.Begin(ctx)
		if txErr != nil {
			return nil, txErr
		}
		defer tx.Rollback(ctx)

		qtx := s.q.WithTx(tx)

		newUser, createErr := qtx.CreateUser(ctx, db.CreateUserParams{
			Username:     username,
			Email:        email,
			PasswordHash: string(hashedBytes),
			Role:         "CONTRIBUTOR",
		})
		if createErr != nil {
			return nil, createErr
		}

		_, _ = qtx.CreateUserProfile(ctx, db.CreateUserProfileParams{
			UserID:   newUser.ID,
			FullName: username,
		})

		_, _ = qtx.CreateContributorProfile(ctx, newUser.ID)

		if commitErr := tx.Commit(ctx); commitErr != nil {
			return nil, commitErr
		}

		finalUser = db.User{
			ID:        newUser.ID,
			Username:  newUser.Username,
			Email:     newUser.Email,
			Role:      newUser.Role,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
		}
	} else {
		finalUser = user
	}

	// Mark user account as verified upon successful OTP check
	_, _ = s.conn.Exec(ctx, "UPDATE users SET is_verified = TRUE WHERE id = $1", finalUser.ID)

	// Issue JWT tokens and create session
	userIDStr := finalUser.ID.String()
	accessToken, err := GenerateAccessToken(userIDStr, finalUser.Role, s.cfg.JWTSecret, s.cfg.JWTIssuer, s.cfg.JWTAudience)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	sessionExpiry := time.Now().Add(7 * 24 * time.Hour)
	_, err = s.q.CreateSession(ctx, db.CreateSessionParams{
		UserID:       finalUser.ID,
		RefreshToken: refreshToken,
		DeviceInfo:   pgtype.Text{String: deviceInfo, Valid: deviceInfo != ""},
		ExpiresAt:    pgtype.Timestamptz{Time: sessionExpiry, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       userIDStr,
		Username:     finalUser.Username,
		Email:        finalUser.Email,
		Role:         finalUser.Role,
	}, nil
}

