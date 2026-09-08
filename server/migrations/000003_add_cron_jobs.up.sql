CREATE TABLE IF NOT EXISTS cron_jobs (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    schedule_interval VARCHAR(32) NOT NULL, -- e.g., '1m', '5m', '15m', '1h', '6h', '24h'
    job_type VARCHAR(64) NOT NULL,          -- e.g., 'DEAD_LINK_CHECKER', 'AUTO_MODERATION', 'GRIEVANCE_SLA', 'REPUTATION_UPDATE', 'CUSTOM'
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_run_at TIMESTAMPTZ,
    next_run_at TIMESTAMPTZ,
    run_count INT NOT NULL DEFAULT 0,
    failure_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS cron_job_logs (
    id BIGSERIAL PRIMARY KEY,
    job_id VARCHAR(64) NOT NULL REFERENCES cron_jobs(id) ON DELETE CASCADE,
    status VARCHAR(32) NOT NULL, -- 'SUCCESS', 'FAILURE', 'RUNNING'
    duration_ms INT NOT NULL DEFAULT 0,
    message TEXT NOT NULL,
    executed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cron_job_logs_job_id ON cron_job_logs(job_id);
CREATE INDEX IF NOT EXISTS idx_cron_job_logs_executed_at ON cron_job_logs(executed_at DESC);

-- Seed core platform background tasks
INSERT INTO cron_jobs (id, name, description, schedule_interval, job_type, is_active)
VALUES
    ('dead_link_checker', 'Dead Video Link Scanner', 'Inspects video links from YouTube, Instagram, and X for broken URLs or removed media.', '15m', 'DEAD_LINK_CHECKER', true),
    ('auto_moderation_batch', 'Auto-Moderation Toxicity Evaluator', 'Scans incoming pending submissions against toxic keywords and spam patterns.', '5m', 'AUTO_MODERATION', true),
    ('grievance_sla_monitor', 'IT Rules 2021 Grievance SLA Monitor', 'Tracks 24-hour receipt acknowledgment and 15-day resolution deadlines for legal compliance.', '1h', 'GRIEVANCE_SLA', true),
    ('reputation_recalculator', 'Contributor Trust Score Recalculator', 'Calculates contributor trust points, updates progression tiers, and evaluates badges.', '6h', 'REPUTATION_UPDATE', true)
ON CONFLICT (id) DO NOTHING;
