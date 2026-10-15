-- Migration 000006: Add SEO slug to content table for clean URL routing (/news/<slug>)

-- 1. Add slug column to content table
ALTER TABLE content ADD COLUMN IF NOT EXISTS slug VARCHAR(300);

-- 2. Create unique index for fast O(1) slug resolution in web portal routing
CREATE UNIQUE INDEX IF NOT EXISTS idx_content_slug ON content(slug) WHERE slug IS NOT NULL;

-- 3. Backfill existing records with clean deterministic fallback slug
UPDATE content
SET slug = CONCAT('news-', SUBSTRING(id::text, 1, 8))
WHERE slug IS NULL;
