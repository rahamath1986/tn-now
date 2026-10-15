-- Migration 000006 rollback: Remove SEO slug from content table
DROP INDEX IF EXISTS idx_content_slug;
ALTER TABLE content DROP COLUMN IF EXISTS slug;
