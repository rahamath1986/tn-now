-- Migration 000012: Down
DROP INDEX IF EXISTS idx_content_is_ai_rewritten;
ALTER TABLE content DROP COLUMN IF EXISTS is_ai_rewritten;

