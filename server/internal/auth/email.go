package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"tn-now/server/internal/config"
)

// GenerateNumericOTP generates a secure 6-digit numeric OTP string (e.g. "849201")
func GenerateNumericOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

// GetEffectiveSMTPConfig retrieves active email settings from database system_settings table, falling back to env config
func GetEffectiveSMTPConfig(ctx context.Context, conn *pgxpool.Pool, baseCfg *config.Config) *config.Config {
	if baseCfg == nil {
		baseCfg = &config.Config{
			SMTPPort:        "587",
			SMTPFrom:        "noreply@tn24.in",
			EmailProvider:   "resend",
			ResendFromEmail: "TN24 News <onboarding@resend.dev>",
		}
	}

	effective := &config.Config{
		SMTPHost:        baseCfg.SMTPHost,
		SMTPPort:        baseCfg.SMTPPort,
		SMTPUser:        baseCfg.SMTPUser,
		SMTPPass:        baseCfg.SMTPPass,
		SMTPFrom:        baseCfg.SMTPFrom,
		EmailProvider:   baseCfg.EmailProvider,
		ResendAPIKey:    baseCfg.ResendAPIKey,
		ResendFromEmail: baseCfg.ResendFromEmail,
	}

	if conn == nil {
		return effective
	}

	var host, port, user, pass, from, provider, resendKey, resendFrom string
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'smtp_host'").Scan(&host)
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'smtp_port'").Scan(&port)
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'smtp_user'").Scan(&user)
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'smtp_pass'").Scan(&pass)
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'smtp_from'").Scan(&from)
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'email_provider'").Scan(&provider)
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'resend_api_key'").Scan(&resendKey)
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'resend_from_email'").Scan(&resendFrom)

	if strings.TrimSpace(host) != "" {
		effective.SMTPHost = strings.TrimSpace(host)
	}
	if strings.TrimSpace(port) != "" {
		effective.SMTPPort = strings.TrimSpace(port)
	}
	if strings.TrimSpace(user) != "" {
		effective.SMTPUser = strings.TrimSpace(user)
	}
	if strings.TrimSpace(pass) != "" {
		effective.SMTPPass = strings.TrimSpace(pass)
	}
	if strings.TrimSpace(from) != "" {
		effective.SMTPFrom = strings.TrimSpace(from)
	}
	if strings.TrimSpace(provider) != "" {
		effective.EmailProvider = strings.TrimSpace(provider)
	}
	if strings.TrimSpace(resendKey) != "" {
		effective.ResendAPIKey = strings.TrimSpace(resendKey)
	}
	if strings.TrimSpace(resendFrom) != "" {
		effective.ResendFromEmail = strings.TrimSpace(resendFrom)
	}

	return effective
}

// SendEmailViaResend dispatches an OTP verification email using Resend's REST API over HTTPS (Port 443)
func SendEmailViaResend(ctx context.Context, apiKey, fromEmail, toEmail, otp string) error {
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return fmt.Errorf("resend API key is not configured")
	}

	from := strings.TrimSpace(fromEmail)
	if from == "" {
		from = "TN24 News <onboarding@resend.dev>"
	}

	subject := fmt.Sprintf("TN24 சரிபார்ப்புக் குறியீடு (OTP): %s", otp)
	htmlContent := renderEmailOTPBody(from, toEmail, subject, otp)

	payload := map[string]interface{}{
		"from":    from,
		"to":      []string{toEmail},
		"subject": subject,
		"html":    htmlContent,
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed encoding email payload: %w", err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "POST", "https://api.resend.com/emails", bytes.NewReader(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed creating HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "TN24-News/1.0")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Resend API HTTP request failed", slog.String("error", err.Error()), slog.String("to", toEmail))
		return fmt.Errorf("failed connecting to Resend API: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp struct {
			Message string `json:"message"`
			Name    string `json:"name"`
		}
		_ = json.Unmarshal(respBody, &errResp)
		msg := errResp.Message
		if msg == "" {
			msg = string(respBody)
		}
		slog.Error("Resend API error response", slog.Int("status", resp.StatusCode), slog.String("message", msg))
		return fmt.Errorf("Resend error (HTTP %d): %s", resp.StatusCode, msg)
	}

	slog.Info("Successfully dispatched OTP email via Resend HTTPS API", slog.String("to", toEmail), slog.Int("status", resp.StatusCode))
	return nil
}

// loginAuth implements smtp.Auth for AUTH LOGIN mechanism (widely used by Exchange, Office365, Brevo, SendGrid, and cPanel)
type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username: username, password: password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		serverPrompt := strings.ToLower(string(fromServer))
		if strings.Contains(serverPrompt, "username") || strings.Contains(serverPrompt, "user") || strings.Contains(serverPrompt, "vxnlyw5hbwu") {
			return []byte(a.username), nil
		}
		if strings.Contains(serverPrompt, "password") || strings.Contains(serverPrompt, "pass") || strings.Contains(serverPrompt, "ugfzc3dvcmq") {
			return []byte(a.password), nil
		}
		// Default to password if username already sent in Start
		return []byte(a.password), nil
	}
	return nil, nil
}

// renderEmailOTPBody builds the HTML email template for the OTP
func renderEmailOTPBody(from, toEmail, subject, otp string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<style>
  body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background: #080f1c; color: #f8fafc; margin: 0; padding: 24px; }
  .box { max-width: 480px; margin: 0 auto; background: #0d1421; border: 1px solid #1e293b; border-radius: 12px; padding: 32px 24px; text-align: center; }
  .logo { font-size: 24px; font-weight: 800; color: #ef4444; margin-bottom: 12px; letter-spacing: -0.5px; }
  .otp-code { font-size: 36px; font-weight: 900; letter-spacing: 6px; color: #38bdf8; background: #111827; border: 2px dashed #1e293b; border-radius: 8px; padding: 14px 20px; display: inline-block; margin: 20px 0; }
  .note { color: #94a3b8; font-size: 13px; line-height: 1.6; }
  .footer { color: #64748b; font-size: 11px; margin-top: 24px; border-top: 1px solid #1e293b; padding-top: 16px; }
</style>
</head>
<body>
  <div class="box">
    <div class="logo">TN24 — தமிழ்நாடு செய்திகள்</div>
    <h2 style="font-size: 18px; color: #f8fafc; margin: 0 0 8px;">உள்நுழைவு சரிபார்ப்புக் குறியீடு (OTP)</h2>
    <p class="note">வாசகர் அரங்கம் மற்றும் TN24 போர்ட்டலில் உங்கள் கணக்கை சரிபார்க்க கீழ்க்கண்ட 6-இலக்க ஒருமுறை கடவுச்சொல்லை (OTP) பயன்படுத்தவும்:</p>
    
    <div class="otp-code">%s</div>
    
    <p class="note">இக்குறியீடு <strong>10 நிமிடங்கள்</strong> மட்டுமே செல்லுபடியாகும். இக்குறியீட்டை யாருடனும் பகிர வேண்டாம்.</p>
    
    <div class="footer">
      இந்த கோரிக்கையை நீங்கள் விடுக்கவில்லை எனில், இச்செய்தியைப் புறக்கணிக்கலாம்.<br>
      &copy; 2026 TN24 News Portal (www.tn24.in).
    </div>
  </div>
</body>
</html>`, otp)
}

// sendEmailViaSMTP performs standard SMTP TCP dispatch
func sendEmailViaSMTP(ctx context.Context, cfg *config.Config, toEmail, otp string) error {
	subject := fmt.Sprintf("TN24 சரிபார்ப்புக் குறியீடு (OTP): %s", otp)
	from := strings.TrimSpace(cfg.SMTPFrom)
	if from == "" {
		from = strings.TrimSpace(cfg.SMTPUser)
	}

	htmlContent := renderEmailOTPBody(from, toEmail, subject, otp)
	body := fmt.Sprintf("From: TN24 News <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		from, toEmail, subject, htmlContent)

	host := strings.TrimSpace(cfg.SMTPHost)
	port := strings.TrimSpace(cfg.SMTPPort)

	if h, p, err := net.SplitHostPort(host); err == nil {
		host = h
		if port == "" {
			port = p
		}
	}
	if port == "" {
		port = "587"
	}
	addr := net.JoinHostPort(host, port)

	timeoutCtx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	dialer := &net.Dialer{
		Timeout: 15 * time.Second,
	}

	rawConn, err := dialer.DialContext(timeoutCtx, "tcp", addr)
	if err != nil {
		rawConn, err = dialer.DialContext(timeoutCtx, "tcp4", addr)
	}
	if err != nil {
		slog.Error("Failed to connect to SMTP server", slog.String("addr", addr), slog.String("error", err.Error()))
		return fmt.Errorf("failed to connect to SMTP server %s: %w", addr, err)
	}
	defer rawConn.Close()

	var client *smtp.Client

	if port == "465" {
		tlsConfig := &tls.Config{
			ServerName: host,
			MinVersion: tls.VersionTLS12,
		}
		tlsConn := tls.Client(rawConn, tlsConfig)
		if err := tlsConn.HandshakeContext(timeoutCtx); err != nil {
			slog.Error("TLS handshake failed with SMTP server", slog.String("host", host), slog.String("error", err.Error()))
			return fmt.Errorf("TLS handshake failed with %s: %w", host, err)
		}
		c, err := smtp.NewClient(tlsConn, host)
		if err != nil {
			return fmt.Errorf("failed creating SMTP SSL client: %w", err)
		}
		client = c
	} else {
		c, err := smtp.NewClient(rawConn, host)
		if err != nil {
			return fmt.Errorf("failed creating SMTP client: %w", err)
		}
		client = c

		if ok, _ := client.Extension("STARTTLS"); ok {
			tlsConfig := &tls.Config{
				ServerName: host,
				MinVersion: tls.VersionTLS12,
			}
			if err := client.StartTLS(tlsConfig); err != nil {
				client.Close()
				slog.Error("STARTTLS negotiation failed", slog.String("host", host), slog.String("error", err.Error()))
				return fmt.Errorf("STARTTLS encryption failed with %s: %w", host, err)
			}
		}
	}
	defer client.Quit()

	if cfg.SMTPUser != "" && cfg.SMTPPass != "" {
		hasAuth, authExt := client.Extension("AUTH")
		authUpper := strings.ToUpper(authExt)

		var auth smtp.Auth
		if hasAuth && strings.Contains(authUpper, "LOGIN") && !strings.Contains(authUpper, "PLAIN") {
			auth = LoginAuth(cfg.SMTPUser, cfg.SMTPPass)
		} else {
			auth = smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, host)
		}

		if err := client.Auth(auth); err != nil {
			if strings.Contains(authUpper, "LOGIN") {
				if errLogin := client.Auth(LoginAuth(cfg.SMTPUser, cfg.SMTPPass)); errLogin == nil {
					err = nil
				}
			}
			if err != nil {
				slog.Error("SMTP authentication failed", slog.String("user", cfg.SMTPUser), slog.String("error", err.Error()))
				return fmt.Errorf("SMTP authentication failed for user %s: %w", cfg.SMTPUser, err)
			}
		}
	}

	if err := client.Mail(from); err != nil {
		slog.Error("SMTP sender rejected", slog.String("from", from), slog.String("error", err.Error()))
		return fmt.Errorf("SMTP sender (%s) rejected: %w", from, err)
	}

	if err := client.Rcpt(toEmail); err != nil {
		slog.Error("SMTP recipient rejected", slog.String("to", toEmail), slog.String("error", err.Error()))
		return fmt.Errorf("SMTP recipient (%s) rejected: %w", toEmail, err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP data initialization failed: %w", err)
	}

	if _, err := w.Write([]byte(body)); err != nil {
		return fmt.Errorf("failed writing email message: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed finalizing email delivery: %w", err)
	}

	slog.Info("Successfully dispatched OTP email via SMTP", slog.String("to", toEmail), slog.String("host", host))
	return nil
}

// SendOTPEmailWithContext sends a formatted email containing the 6-digit OTP code using Resend (HTTPS) or SMTP (TCP)
func SendOTPEmailWithContext(ctx context.Context, cfg *config.Config, toEmail, otp string) error {
	toEmail = strings.TrimSpace(toEmail)
	if toEmail == "" {
		return fmt.Errorf("recipient email address is empty")
	}

	// 1. If provider is "resend" (or if Resend API key is available and provider is not forced to smtp), use Resend HTTPS API
	if cfg != nil && (cfg.EmailProvider == "resend" || cfg.ResendAPIKey != "") && cfg.EmailProvider != "smtp" {
		if strings.TrimSpace(cfg.ResendAPIKey) != "" {
			return SendEmailViaResend(ctx, cfg.ResendAPIKey, cfg.ResendFromEmail, toEmail, otp)
		}
	}

	// 2. If SMTP is configured, use standard SMTP
	if cfg != nil && strings.TrimSpace(cfg.SMTPHost) != "" && strings.TrimSpace(cfg.SMTPUser) != "" {
		return sendEmailViaSMTP(ctx, cfg, toEmail, otp)
	}

	// 3. Graceful dev console log fallback
	slog.Info("✉️ [EMAIL OTP DISPATCH - DEV LOG]",
		slog.String("to", toEmail),
		slog.String("otp_code", otp),
		slog.String("note", "Configure Resend API key or SMTP in Admin Console (or .env) to dispatch real emails"),
	)
	return nil
}

// SendOTPEmail sends a formatted email containing the 6-digit OTP code
func SendOTPEmail(cfg *config.Config, toEmail, otp string) error {
	return SendOTPEmailWithContext(context.Background(), cfg, toEmail, otp)
}
