package config

import (
	"os"
	"strings"
)

type Config struct {
	Port                string
	AppEnv              string
	DatabaseURL         string
	RedisURL            string
	JWTSecret           string
	JWTIssuer           string
	JWTAudience         string
	AllowedVideoDomains []string
	CORSAllowedOrigins  []string
	GoogleClientID      string
	GoogleClientSecret  string
	GoogleRedirectURL   string
	GeminiAPIKey        string
	AdminUsername       string
	AdminPassword       string
	SMTPHost            string
	SMTPPort            string
	SMTPUser            string
	SMTPPass            string
	SMTPFrom            string
	EmailProvider       string
	ResendAPIKey        string
	ResendFromEmail     string
	ImageKitURLEndpoint string
	ImageKitPublicKey   string
	ImageKitPrivateKey  string
	ImageKitFolder      string
}

func loadEnvFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if os.Getenv(key) == "" {
				_ = os.Setenv(key, val)
			}
		}
	}
}

func Load() *Config {
	loadEnvFile(".env")
	loadEnvFile("server/.env")

	port := getEnv("PORT", "8080")
	appEnv := getEnv("APP_ENV", "development")
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:password123@localhost:5432/tnnow_dev?sslmode=disable")
	redisURL := getEnv("REDIS_URL", "redis://localhost:6379")
	jwtSecret := getEnv("JWT_SECRET", "tnnow-default-secret-key-32-chars-long!")
	jwtIssuer := getEnv("JWT_ISSUER", "tn-now")
	jwtAudience := getEnv("JWT_AUDIENCE", "tn-now-client")

	allowedVideoDomainsStr := getEnv("ALLOWED_VIDEO_DOMAINS", "*.youtube.com,youtu.be,*.instagram.com,instagram.com,*.facebook.com,facebook.com")
	allowedVideoDomains := splitAndTrim(allowedVideoDomainsStr)

	corsOriginsStr := getEnv("CORS_ALLOWED_ORIGINS", "*")
	corsOrigins := splitAndTrim(corsOriginsStr)

	googleClientID := getEnv("GOOGLE_CLIENT_ID", "")
	googleClientSecret := getEnv("GOOGLE_CLIENT_SECRET", "")
	googleRedirectURL := getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/admin/auth/google/callback")
	geminiAPIKey := getEnv("GEMINI_API_KEY", "")

	adminUsername := getEnv("ADMIN_USERNAME", "admin")
	adminPassword := getEnv("ADMIN_PASSWORD", "tn24admin")

	smtpHost := getEnv("SMTP_HOST", "")
	smtpPort := getEnv("SMTP_PORT", "587")
	smtpUser := getEnv("SMTP_USER", "")
	smtpPass := getEnv("SMTP_PASS", "")
	smtpFrom := getEnv("SMTP_FROM", "noreply@tn24.in")

	emailProvider := getEnv("EMAIL_PROVIDER", "resend")
	resendAPIKey := getEnv("RESEND_API_KEY", "")
	resendFromEmail := getEnv("RESEND_FROM_EMAIL", "TN24 News <onboarding@resend.dev>")

	ikEndpoint := getEnv("IMAGEKIT_URL_ENDPOINT", "")
	ikPublic := getEnv("IMAGEKIT_PUBLIC_KEY", "")
	ikPrivate := getEnv("IMAGEKIT_PRIVATE_KEY", "")
	ikFolder := getEnv("IMAGEKIT_FOLDER", "/articles")

	return &Config{
		Port:                port,
		AppEnv:              appEnv,
		DatabaseURL:         dbURL,
		RedisURL:            redisURL,
		JWTSecret:           jwtSecret,
		JWTIssuer:           jwtIssuer,
		JWTAudience:         jwtAudience,
		AllowedVideoDomains: allowedVideoDomains,
		CORSAllowedOrigins:  corsOrigins,
		GoogleClientID:      googleClientID,
		GoogleClientSecret:  googleClientSecret,
		GoogleRedirectURL:   googleRedirectURL,
		GeminiAPIKey:        geminiAPIKey,
		AdminUsername:       adminUsername,
		AdminPassword:       adminPassword,
		SMTPHost:            smtpHost,
		SMTPPort:            smtpPort,
		SMTPUser:            smtpUser,
		SMTPPass:            smtpPass,
		SMTPFrom:            smtpFrom,
		EmailProvider:       emailProvider,
		ResendAPIKey:        resendAPIKey,
		ResendFromEmail:     resendFromEmail,
		ImageKitURLEndpoint: ikEndpoint,
		ImageKitPublicKey:   ikPublic,
		ImageKitPrivateKey:  ikPrivate,
		ImageKitFolder:      ikFolder,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func splitAndTrim(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
