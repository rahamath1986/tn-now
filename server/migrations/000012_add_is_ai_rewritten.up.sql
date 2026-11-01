-- Migration 000012: Ensure is_ai_rewritten, source_type, and indexes on content table
ALTER TABLE content ADD COLUMN IF NOT EXISTS is_viral BOOLEAN DEFAULT FALSE;
ALTER TABLE content ADD COLUMN IF NOT EXISTS language VARCHAR(20) DEFAULT 'ta';
ALTER TABLE content ADD COLUMN IF NOT EXISTS is_ai_rewritten BOOLEAN DEFAULT FALSE;
ALTER TABLE content ADD COLUMN IF NOT EXISTS source_type VARCHAR(50) DEFAULT 'SCRAPER';
ALTER TABLE content ADD COLUMN IF NOT EXISTS slug VARCHAR(300);

CREATE INDEX IF NOT EXISTS idx_content_is_ai_rewritten ON content(is_ai_rewritten);
CREATE INDEX IF NOT EXISTS idx_content_status_created ON content(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_content_created_at ON content(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_content_source_url ON content(source_url);
