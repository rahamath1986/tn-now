DROP INDEX IF EXISTS idx_content_author_created;
DROP INDEX IF EXISTS idx_content_source_status;
DELETE FROM categories WHERE slug = 'user-corner';

