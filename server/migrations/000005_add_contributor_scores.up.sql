-- Migration 000005: Add contributor_scores table
-- This table is required by the reputation_recalculator cron job and admin contributor leaderboard.

CREATE TABLE IF NOT EXISTS contributor_scores (
    user_id     UUID PRIMARY KEY REFERENCES profiles(user_id) ON DELETE CASCADE,
    level       VARCHAR(50)  NOT NULL DEFAULT 'Newcomer',
    points      INTEGER      NOT NULL DEFAULT 0,
    trust_score INTEGER      NOT NULL DEFAULT 0,
    approved_count INTEGER   NOT NULL DEFAULT 0,
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_contributor_scores_points ON contributor_scores(points DESC);
CREATE INDEX IF NOT EXISTS idx_contributor_scores_trust  ON contributor_scores(trust_score DESC);
