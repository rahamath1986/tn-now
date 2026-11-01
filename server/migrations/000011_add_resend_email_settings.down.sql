-- 000011_add_resend_email_settings.down.sql
DELETE FROM system_settings 
WHERE key IN ('email_provider', 'resend_api_key', 'resend_from_email');

