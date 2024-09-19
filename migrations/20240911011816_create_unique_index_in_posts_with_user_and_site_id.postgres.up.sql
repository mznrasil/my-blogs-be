-- Drop the old unique constraint on slug (if exists)
ALTER TABLE posts DROP CONSTRAINT IF EXISTS posts_slug_key;

-- Add the new constraint on user_id, site_id and slug
ALTER TABLE posts ADD CONSTRAINT posts_user_id_site_id_slug_key UNIQUE (user_id, site_id, slug);
