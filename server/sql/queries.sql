-- --- USERS & AUTH ---

-- name: CreateUser :one
INSERT INTO users (username, email, password_hash, role)
VALUES ($1, $2, $3, $4)
RETURNING id, username, email, role, created_at, updated_at;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash, role, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, username, email, password_hash, role, created_at, updated_at
FROM users
WHERE id = $1;


-- --- USER PROFILES ---

-- name: CreateUserProfile :one
INSERT INTO user_profiles (user_id, full_name, avatar_url, bio)
VALUES ($1, $2, $3, $4)
RETURNING user_id, full_name, avatar_url, bio, created_at, updated_at;

-- name: GetUserProfile :one
SELECT user_id, full_name, avatar_url, bio, created_at, updated_at
FROM user_profiles
WHERE user_id = $1;


-- --- CONTRIBUTOR PROFILES ---

-- name: CreateContributorProfile :one
INSERT INTO contributor_profiles (user_id)
VALUES ($1)
RETURNING user_id, level, points, trust_score;

-- name: GetContributorProfile :one
SELECT user_id, level, points, trust_score, approved_count, rejected_count, report_count, violation_count, featured_count, is_founding_contributor, founding_badge_granted_at, created_at, updated_at
FROM contributor_profiles
WHERE user_id = $1;

-- name: UpdateContributorReputation :one
UPDATE contributor_profiles
SET points = $2, trust_score = $3, level = $4, updated_at = NOW()
WHERE user_id = $1
RETURNING user_id, level, points, trust_score;

-- name: UpdateFoundingStatus :one
UPDATE contributor_profiles
SET is_founding_contributor = $2, founding_badge_granted_at = $3, updated_at = NOW()
WHERE user_id = $1
RETURNING user_id, is_founding_contributor, founding_badge_granted_at;


-- --- GEOGRAPHIC & TAXONOMY ---

-- name: ListDistricts :many
SELECT id, name, created_at
FROM districts
ORDER BY name ASC;

-- name: ListCategories :many
SELECT id, name, slug, created_at
FROM categories
ORDER BY name ASC;


-- --- CONTENT CREATION & LISTING ---

-- name: CreateContent :one
INSERT INTO content (title, description, content_type, category_id, district_id, location_id, source_type, source_url, author_user_id, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, title, description, content_type, category_id, district_id, location_id, source_type, source_url, author_user_id, status, created_at;

-- name: GetContentByID :one
SELECT id, title, description, content_type, category_id, district_id, location_id, source_type, source_url, author_user_id, status, moderation_status, verification_status, published_at, created_at, updated_at
FROM content
WHERE id = $1;

-- name: UpdateContentStatus :one
UPDATE content
SET status = $2, moderation_status = $3, published_at = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, status, moderation_status, published_at;

-- name: ListContentPaginated :many
SELECT id, title, description, content_type, category_id, district_id, location_id, source_type, source_url, author_user_id, status, published_at, created_at
FROM content
WHERE status = 'PUBLISHED'
  AND (published_at < $1 OR (published_at = $1 AND id < $2))
ORDER BY published_at DESC, id DESC
LIMIT $3;


-- --- VIDEO LINKS SPECIFIC ---

-- name: CreateVideoLink :one
INSERT INTO video_links (content_id, platform, external_video_id, canonical_url, embed_html, thumbnail_url, duration_seconds)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING content_id, platform, external_video_id, canonical_url, embed_html, thumbnail_url, duration_seconds, link_status;

-- name: GetVideoLinkByContentID :one
SELECT content_id, platform, external_video_id, canonical_url, embed_html, thumbnail_url, duration_seconds, link_status, last_checked_at, created_at
FROM video_links
WHERE content_id = $1;


-- --- PHOTOS SPECIFIC ---

-- name: CreatePhoto :one
INSERT INTO photos (content_id, photo_count, photo_urls, caption_en, caption_ta)
VALUES ($1, $2, $3, $4, $5)
RETURNING content_id, photo_count, photo_urls;


-- --- STORIES SPECIFIC ---

-- name: CreateStory :one
INSERT INTO stories (content_id, headline, body, single_photo_url)
VALUES ($1, $2, $3, $4)
RETURNING content_id, headline, body, single_photo_url;

-- name: GetStoryByContentID :one
SELECT content_id, headline, body, single_photo_url
FROM stories
WHERE content_id = $1;


-- --- EVENTS SPECIFIC ---

-- name: CreateEvent :one
INSERT INTO events (content_id, event_name, event_date, start_time, end_time, organizer_name, organizer_contact)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING content_id, event_name, event_date, start_time, end_time, organizer_name, organizer_contact;

-- name: GetEventByContentID :one
SELECT content_id, event_name, event_date, start_time, end_time, organizer_name, organizer_contact, interested_count
FROM events
WHERE content_id = $1;


-- --- GRIEVANCES (COMPLIANCE) ---

-- name: CreateGrievance :one
INSERT INTO grievances (complainant_user_id, contact_email, category, description, content_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, complainant_user_id, contact_email, category, description, content_id, status, created_at;

-- name: GetGrievanceByID :one
SELECT id, complainant_user_id, contact_email, category, description, content_id, status, resolution_notes, resolved_at, created_at, updated_at
FROM grievances
WHERE id = $1;

-- name: UpdateGrievanceStatus :one
UPDATE grievances
SET status = $2, resolution_notes = $3, resolved_at = $4, updated_at = NOW()
WHERE id = $1
RETURNING id, status, resolution_notes, resolved_at;


-- --- USER SESSIONS ---

-- name: CreateSession :one
INSERT INTO user_sessions (user_id, refresh_token, device_info, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, refresh_token, expires_at, created_at;

-- name: GetSessionByToken :one
SELECT id, user_id, refresh_token, device_info, is_revoked, expires_at, created_at, updated_at
FROM user_sessions
WHERE refresh_token = $1;

-- name: RevokeSession :one
UPDATE user_sessions
SET is_revoked = TRUE, updated_at = NOW()
WHERE refresh_token = $1
RETURNING id, is_revoked, updated_at;

-- name: DeleteUserSessions :exec
DELETE FROM user_sessions
WHERE user_id = $1;


-- --- DEVICES ---

-- name: CreateDevice :one
INSERT INTO devices (user_id, device_token, device_type)
VALUES ($1, $2, $3)
ON CONFLICT (device_token) DO UPDATE SET updated_at = NOW()
RETURNING id, device_token, device_type;

-- name: DeleteDevice :exec
DELETE FROM devices
WHERE device_token = $1;


-- --- HOME FEED ---

-- name: GetLatestContent :many
SELECT id, title, description, content_type, category_id, district_id, location_id, source_type, source_url, author_user_id, status, published_at, created_at, updated_at
FROM content
WHERE status = 'PUBLISHED'
ORDER BY published_at DESC, id DESC
LIMIT $1;

-- name: GetUpcomingEvents :many
SELECT c.id, c.title, c.description, c.category_id, c.district_id, c.author_user_id, c.status, e.event_name, e.event_date, e.start_time, e.end_time, e.organizer_name
FROM content c
JOIN events e ON c.id = e.content_id
WHERE c.status = 'PUBLISHED' AND e.event_date >= $1
ORDER BY e.event_date ASC, c.id DESC
LIMIT $2;


-- --- CATEGORY FEED ---

-- name: ListContentByCategoryPaginated :many
SELECT id, title, description, content_type, category_id, district_id, location_id, source_type, source_url, author_user_id, status, published_at, created_at, updated_at
FROM content
WHERE status = 'PUBLISHED'
  AND category_id = $1
  AND (published_at < $2 OR (published_at = $2 AND id < $3))
ORDER BY published_at DESC, id DESC
LIMIT $4;


-- --- MODERATION & AUDIT ---

-- name: ListPendingModerationQueue :many
SELECT id, title, description, content_type, category_id, district_id, location_id, source_type, source_url, author_user_id, status, moderation_status, verification_status, published_at, created_at, updated_at
FROM content
WHERE status = 'PENDING' OR moderation_status = 'QUARANTINE'
ORDER BY created_at ASC
LIMIT $1;

-- name: CreateModerationResult :one
INSERT INTO moderation_results (content_id, auto_moderated, text_toxicity_score, image_safety_score, flags_raised, decision)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (content_id) DO UPDATE SET
  auto_moderated = EXCLUDED.auto_moderated,
  text_toxicity_score = EXCLUDED.text_toxicity_score,
  image_safety_score = EXCLUDED.image_safety_score,
  flags_raised = EXCLUDED.flags_raised,
  decision = EXCLUDED.decision,
  processed_at = NOW()
RETURNING content_id, decision, processed_at;

-- name: CreateAuditLog :one
INSERT INTO audit_logs (actor_user_id, action, target_entity, target_id, details)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, action, created_at;


-- --- SOCIAL & ENGAGEMENT ---

-- name: CreateLike :one
INSERT INTO likes (content_id, user_id)
VALUES ($1, $2)
ON CONFLICT (content_id, user_id) DO NOTHING
RETURNING content_id, user_id, created_at;

-- name: DeleteLike :exec
DELETE FROM likes
WHERE content_id = $1 AND user_id = $2;

-- name: GetLikeCount :one
SELECT COUNT(*) FROM likes WHERE content_id = $1;

-- name: CreateBookmark :one
INSERT INTO bookmarks (content_id, user_id)
VALUES ($1, $2)
ON CONFLICT (content_id, user_id) DO NOTHING
RETURNING content_id, user_id, created_at;

-- name: DeleteBookmark :exec
DELETE FROM bookmarks
WHERE content_id = $1 AND user_id = $2;

-- name: CreateFollow :one
INSERT INTO follows (follower_user_id, followed_user_id)
VALUES ($1, $2)
ON CONFLICT (follower_user_id, followed_user_id) DO NOTHING
RETURNING follower_user_id, followed_user_id, created_at;

-- name: DeleteFollow :exec
DELETE FROM follows
WHERE follower_user_id = $1 AND followed_user_id = $2;

-- name: CreateComment :one
INSERT INTO comments (content_id, user_id, body)
VALUES ($1, $2, $3)
RETURNING id, content_id, user_id, body, created_at;

-- name: ListCommentsByContentID :many
SELECT id, content_id, user_id, body, created_at, updated_at
FROM comments
WHERE content_id = $1
ORDER BY created_at ASC;


-- --- SEARCH & DISCOVERY ---

-- name: SearchContent :many
SELECT id, title, description, content_type, category_id, district_id, location_id, source_type, source_url, author_user_id, status, published_at, created_at, updated_at
FROM content
WHERE status = 'PUBLISHED'
  AND (
    title ILIKE '%' || $1 || '%' OR COALESCE(description, '') ILIKE '%' || $1 || '%'
  )
ORDER BY published_at DESC, id DESC
LIMIT $2;

