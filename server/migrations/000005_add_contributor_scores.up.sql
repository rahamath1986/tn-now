-- Migration 000005: Add contributor_scores table
-- This table is required by the reputation_recalculator cron job and admin contributor leaderboard.

CREATE TABLE IF NOT EXISTS contributor_scores (
    user_id     UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    level       VARCHAR(50)  NOT NULL DEFAULT 'Newcomer',
    points      INTEGER      NOT NULL DEFAULT 0,
    trust_score INTEGER      NOT NULL DEFAULT 0,
    approved_count INTEGER   NOT NULL DEFAULT 0,
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contributor_scores_points ON contributor_scores(points DESC);
CREATE INDEX IF NOT EXISTS idx_contributor_scores_trust  ON contributor_scores(trust_score DESC);

-- View profiles for backward compatibility with queries joining profiles
CREATE OR REPLACE VIEW profiles AS
SELECT 
    u.id AS user_id,
    u.username,
    COALESCE(up.full_name, u.username) AS display_name,
    COALESCE(up.bio, '') AS bio,
    COALESCE(up.avatar_url, '') AS avatar_url,
    (u.role = 'ADMIN') AS is_verified,
    COALESCE(cp.is_founding_contributor, false) AS is_founding_contributor
FROM users u
LEFT JOIN user_profiles up ON u.id = up.user_id
LEFT JOIN contributor_profiles cp ON u.id = cp.user_id;

