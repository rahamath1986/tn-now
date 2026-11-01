-- 000011_add_resend_email_settings.up.sql
-- Seed default Resend.com HTTP email provider settings into system_settings table

INSERT INTO system_settings (key, value, updated_at)
VALUES 
    ('email_provider', 'resend', NOW()),
    ('resend_api_key', '', NOW()),
    ('resend_from_email', 'TN24 News <onboarding@resend.dev>', NOW())
ON CONFLICT (key) DO NOTHING;

