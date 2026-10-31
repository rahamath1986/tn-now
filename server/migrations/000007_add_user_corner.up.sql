-- Migration 000007: Add User Corner category and performance indexes for user submissions

-- 1. Insert User Corner category
INSERT INTO categories (name, slug)
VALUES ('User Corner', 'user-corner')
ON CONFLICT (name) DO NOTHING;

-- 2. Index for filtering user submissions by status (moderation queue & portal feeds)
CREATE INDEX IF NOT EXISTS idx_content_source_status ON content(source_type, status, created_at DESC);

-- 3. Index for user submissions lookup by author (user profile & my-posts dashboard)
CREATE INDEX IF NOT EXISTS idx_content_author_created ON content(author_user_id, created_at DESC);

