-- Enable uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- USERS TABLE
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'USER',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- USER PROFILES TABLE
CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    full_name VARCHAR(150) NOT NULL,
    avatar_url TEXT,
    bio TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- CONTRIBUTOR PROFILES TABLE
CREATE TABLE contributor_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    level VARCHAR(50) NOT NULL DEFAULT 'New Contributor',
    points INT NOT NULL DEFAULT 0,
    trust_score INT NOT NULL DEFAULT 10,
    approved_count INT NOT NULL DEFAULT 0,
    rejected_count INT NOT NULL DEFAULT 0,
    report_count INT NOT NULL DEFAULT 0,
    violation_count INT NOT NULL DEFAULT 0,
    featured_count INT NOT NULL DEFAULT 0,
    is_founding_contributor BOOLEAN NOT NULL DEFAULT FALSE,
    founding_badge_granted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- DISTRICTS TABLE
CREATE TABLE districts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- LOCATIONS TABLE
CREATE TABLE locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) NOT NULL,
    district_id UUID NOT NULL REFERENCES districts(id) ON DELETE CASCADE,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- CATEGORIES TABLE
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- MAIN CONTENT TABLE
CREATE TABLE content (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    content_type VARCHAR(50) NOT NULL, -- 'VIDEO_LINK' | 'PHOTO' | 'TEXT_STORY' | 'EVENT'
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    district_id UUID NOT NULL REFERENCES districts(id) ON DELETE RESTRICT,
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,
    source_type VARCHAR(50) NOT NULL DEFAULT 'USER',
    source_url TEXT,
    author_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- 'PENDING' | 'PUBLISHED' | 'QUARANTINED' | 'REJECTED'
    moderation_status VARCHAR(50) NOT NULL DEFAULT 'UNMODERATED', -- 'UNMODERATED' | 'SAFE' | 'QUARANTINE' | 'REJECTED'
    verification_status VARCHAR(50) NOT NULL DEFAULT 'UNVERIFIED', -- 'UNVERIFIED' | 'SOURCE_IDENTIFIED' | 'SOURCE_CONFIRMED' | 'VERIFIED'
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- VIDEO LINKS TABLE (Extends Content)
CREATE TABLE video_links (
    content_id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,
    platform VARCHAR(50) NOT NULL, -- 'youtube' | 'instagram' | 'facebook' | 'other'
    external_video_id VARCHAR(150) UNIQUE NOT NULL,
    canonical_url TEXT NOT NULL,
    embed_html TEXT,
    thumbnail_url TEXT,
    duration_seconds INT,
    link_status VARCHAR(50) NOT NULL DEFAULT 'UNRESOLVED', -- 'ALIVE' | 'DEAD' | 'REMOVED_BY_SOURCE' | 'UNRESOLVED'
    last_checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- PHOTOS TABLE (Extends Content)
CREATE TABLE photos (
    content_id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,
    photo_count INT NOT NULL DEFAULT 1,
    photo_urls TEXT[] NOT NULL,
    caption_en TEXT,
    caption_ta TEXT
);

-- STORIES TABLE (Extends Content)
CREATE TABLE stories (
    content_id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,
    headline VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    single_photo_url TEXT
);

-- EVENTS TABLE (Extends Content)
CREATE TABLE events (
    content_id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,
    event_name VARCHAR(255) NOT NULL,
    event_date TIMESTAMPTZ NOT NULL,
    start_time VARCHAR(50),
    end_time VARCHAR(50),
    organizer_name VARCHAR(255) NOT NULL,
    organizer_contact VARCHAR(100),
    interested_count INT NOT NULL DEFAULT 0
);

-- COMMENTS TABLE
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- LIKES TABLE (Composite Primary Key)
CREATE TABLE likes (
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (content_id, user_id)
);

-- BOOKMARKS TABLE (Composite Primary Key)
CREATE TABLE bookmarks (
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (content_id, user_id)
);

-- FOLLOWS TABLE (Composite Primary Key)
CREATE TABLE follows (
    follower_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    followed_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (follower_user_id, followed_user_id),
    CONSTRAINT cannot_follow_self CHECK (follower_user_id != followed_user_id)
);

-- REPORTS TABLE
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    reporter_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reason VARCHAR(100) NOT NULL, -- 'Adult/sexual', 'Violence', 'Copyright', 'Spam', 'Harassment', 'Misleading', 'Illegal', 'Dead link', 'Other'
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- 'PENDING' | 'REVIEWED' | 'RESOLVED'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- MODERATION RESULTS TABLE
CREATE TABLE moderation_results (
    content_id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,
    auto_moderated BOOLEAN NOT NULL DEFAULT FALSE,
    text_toxicity_score DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    image_safety_score DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    flags_raised TEXT[] NOT NULL DEFAULT '{}',
    decision VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- 'SAFE' | 'REJECT' | 'QUARANTINE'
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- GRIEVANCES TABLE (Indian IT Rules 2021 Compliant)
CREATE TABLE grievances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    complainant_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    contact_email VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    content_id UUID REFERENCES content(id) ON DELETE SET NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'ACKNOWLEDGED', -- 'ACKNOWLEDGED' | 'UNDER_INVESTIGATION' | 'RESOLVED'
    resolution_notes TEXT,
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- AUDIT LOGS TABLE
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    target_entity VARCHAR(100) NOT NULL,
    target_id UUID,
    details TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- --- INDEXING ---
CREATE INDEX idx_content_status_published ON content (status, published_at DESC);
CREATE INDEX idx_content_category ON content (category_id, status);
CREATE INDEX idx_content_district ON content (district_id, status);
CREATE INDEX idx_content_author ON content (author_user_id);
CREATE INDEX idx_video_links_status ON video_links (link_status);
CREATE INDEX idx_video_links_platform_ext ON video_links (platform, external_video_id);
CREATE INDEX idx_events_date ON events (event_date DESC);
CREATE INDEX idx_reports_status ON reports (status);
CREATE INDEX idx_grievances_status_created ON grievances (status, created_at DESC);

-- --- DEVELOPER SEED DATA ---
INSERT INTO districts (name) VALUES
    ('Madurai'),
    ('Chennai'),
    ('Coimbatore'),
    ('Trichy'),
    ('Salem'),
    ('Tirunelveli'),
    ('Thanjavur'),
    ('Tiruppur'),
    ('Vellore'),
    ('Tamil Nadu')
ON CONFLICT (name) DO NOTHING;

INSERT INTO categories (name, slug) VALUES
    ('Viral Videos', 'viral'),
    ('News', 'news'),
    ('Politics', 'politics'),
    ('Cinema', 'cinema'),
    ('Sports', 'sports'),
    ('Events', 'events'),
    ('Local Happenings', 'local'),
    ('Public Updates', 'public-updates'),
    ('Business', 'business'),
    ('Education', 'education'),
    ('Entertainment', 'entertainment')
ON CONFLICT (name) DO NOTHING;
