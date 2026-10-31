-- Migration 000009: Add user email verification status and system SMTP settings

-- 1. Add is_verified column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_verified BOOLEAN NOT NULL DEFAULT FALSE;

-- Grandfather existing accounts as verified
UPDATE users SET is_verified = TRUE WHERE is_verified = FALSE;

-- Index for quick lookup on verification status
CREATE INDEX IF NOT EXISTS idx_users_email_verified ON users(email, is_verified);

-- 2. Seed default SMTP configuration keys into system_settings if not already present
INSERT INTO system_settings (key, value) VALUES
    ('smtp_host', ''),
    ('smtp_port', '587'),
    ('smtp_user', ''),
    ('smtp_pass', ''),
    ('smtp_from', 'noreply@tn24.in')
ON CONFLICT (key) DO NOTHING;

