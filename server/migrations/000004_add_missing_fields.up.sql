-- Migration 000004: Add missing columns and seed default cron jobs

-- 1. Ensure is_viral and language exist on content table
ALTER TABLE content ADD COLUMN IF NOT EXISTS is_viral BOOLEAN DEFAULT FALSE;
ALTER TABLE content ADD COLUMN IF NOT EXISTS language VARCHAR(20) DEFAULT 'ta';

-- 2. Performance indexes
CREATE INDEX IF NOT EXISTS idx_content_is_viral ON content(is_viral);
CREATE INDEX IF NOT EXISTS idx_content_language ON content(language);
CREATE INDEX IF NOT EXISTS idx_content_status_created ON content(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_content_source_url ON content(source_url);

-- 3. Ensure source_urls exists on cron_jobs
ALTER TABLE cron_jobs ADD COLUMN IF NOT EXISTS source_urls TEXT[] DEFAULT '{}';

-- 4. Seed all platform background cron jobs into cron_jobs table
INSERT INTO cron_jobs (id, name, description, schedule_interval, job_type, is_active, source_urls)
VALUES
    ('tn_live_news_cron', 'TN Live News & Video Scraper', 'Scrapes Tamil Nadu text stories, video links, and images across regional sources and stages them for review.', '30m', 'INGESTION', true, '{"https://www.thehindu.com/news/national/tamil-nadu/feeder/default.rss"}'),
    ('content_retention_cleanup', '24-Hour Content Retention Auto-Purge', 'Automatically purges content in any status (Pending, Published, Rejected) older than the configured retention policy (default: 24 hours).', '15m', 'MAINTENANCE', true, '{}'),
    ('content_deduplication_cleanup', 'Automatic Duplicate Post Purge (Retain Latest)', 'Automatically identifies duplicate news posts across providers, merges metadata/photos into the newest version, and purges older duplicate records.', '15m', 'MAINTENANCE', true, '{}')
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    schedule_interval = EXCLUDED.schedule_interval,
    job_type = EXCLUDED.job_type;

-- 5. Ensure system_settings table exists
CREATE TABLE IF NOT EXISTS system_settings (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
INSERT INTO system_settings (key, value) VALUES ('content_retention_hours', '24') ON CONFLICT (key) DO NOTHING;
INSERT INTO system_settings (key, value) VALUES ('auto_cleanup_enabled', 'true') ON CONFLICT (key) DO NOTHING;
INSERT INTO system_settings (key, value) VALUES ('default_language', 'all') ON CONFLICT (key) DO NOTHING;
INSERT INTO system_settings (key, value) VALUES ('enabled_languages', 'ta,en,ta-en') ON CONFLICT (key) DO NOTHING;
INSERT INTO system_settings (key, value) VALUES ('scraper_language_policy', 'bilingual') ON CONFLICT (key) DO NOTHING;
