-- Migration 000010: Add ImageKit CDN Configuration Keys to system_settings

INSERT INTO system_settings (key, value) VALUES
    ('imagekit_url_endpoint', ''),
    ('imagekit_public_key', ''),
    ('imagekit_private_key', ''),
    ('imagekit_folder', '/articles'),
    ('imagekit_auto_optimize', 'true')
ON CONFLICT (key) DO NOTHING;

