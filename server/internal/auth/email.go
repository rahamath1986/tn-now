package auth

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log/slog"
	"math/big"
	"net/smtp"
	"strings"

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

// GetEffectiveSMTPConfig retrieves active SMTP settings from database system_settings table, falling back to env config
func GetEffectiveSMTPConfig(ctx context.Context, conn *pgxpool.Pool, baseCfg *config.Config) *config.Config {
	if baseCfg == nil {
		baseCfg = &config.Config{
			SMTPPort: "587",
			SMTPFrom: "noreply@tn24.in",
		}
	}

	effective := &config.Config{
		SMTPHost: baseCfg.SMTPHost,
		SMTPPort: baseCfg.SMTPPort,
		SMTPUser: baseCfg.SMTPUser,
		SMTPPass: baseCfg.SMTPPass,
		SMTPFrom: baseCfg.SMTPFrom,
	}

	if conn == nil {
		return effective
	}

	var host, port, user, pass, from string
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'smtp_host'").Scan(&host)
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'smtp_port'").Scan(&port)
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'smtp_user'").Scan(&user)
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'smtp_pass'").Scan(&pass)
	_ = conn.QueryRow(ctx, "SELECT value FROM system_settings WHERE key = 'smtp_from'").Scan(&from)

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

	return effective
}

// SendOTPEmail sends a formatted email containing the 6-digit OTP code
func SendOTPEmail(cfg *config.Config, toEmail, otp string) error {
	toEmail = strings.TrimSpace(toEmail)

	// If SMTP is not configured, gracefully log the OTP to server console (ideal for dev/testing)
	if cfg == nil || cfg.SMTPHost == "" || cfg.SMTPUser == "" {
		slog.Info("✉️ [EMAIL OTP DISPATCH - DEV LOG]",
			slog.String("to", toEmail),
			slog.String("otp_code", otp),
			slog.String("note", "Configure SMTP credentials in Admin Console (or .env) to dispatch real emails"),
		)
		return nil
	}

	subject := fmt.Sprintf("TN24 சரிபார்ப்புக் குறியீடு (OTP): %s", otp)
	from := cfg.SMTPFrom
	if from == "" {
		from = cfg.SMTPUser
	}

	body := fmt.Sprintf(`From: TN24 News <%s>
To: %s
Subject: %s
MIME-Version: 1.0
Content-Type: text/html; charset=UTF-8

<!DOCTYPE html>
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
</html>
`, from, toEmail, subject, otp)

	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)

	// SSL/TLS connection for port 465
	if cfg.SMTPPort == "465" {
		tlsconfig := &tls.Config{
			ServerName: cfg.SMTPHost,
		}
		conn, err := tls.Dial("tcp", addr, tlsconfig)
		if err != nil {
			slog.Error("Failed to dial SMTP SSL (465)", slog.String("error", err.Error()), slog.String("to", toEmail))
			return err
		}
		defer conn.Close()

		c, err := smtp.NewClient(conn, cfg.SMTPHost)
		if err != nil {
			return err
		}
		defer c.Quit()

		if err = c.Auth(auth); err != nil {
			return err
		}
		if err = c.Mail(from); err != nil {
			return err
		}
		if err = c.Rcpt(toEmail); err != nil {
			return err
		}
		w, err := c.Data()
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(body))
		if err != nil {
			return err
		}
		return w.Close()
	}

	// Standard STARTTLS connection for port 587
	err := smtp.SendMail(addr, auth, from, []string{toEmail}, []byte(body))
	if err != nil {
		slog.Error("Failed to send OTP email via SMTP", slog.String("error", err.Error()), slog.String("to", toEmail))
		return err
	}

	slog.Info("Successfully dispatched OTP email via SMTP", slog.String("to", toEmail))
	return nil
}

