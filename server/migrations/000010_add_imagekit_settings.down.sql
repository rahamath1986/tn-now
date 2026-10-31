DELETE FROM system_settings WHERE key IN (
    'imagekit_url_endpoint',
    'imagekit_public_key',
    'imagekit_private_key',
    'imagekit_folder',
    'imagekit_auto_optimize'
);

