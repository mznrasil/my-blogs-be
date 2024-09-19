-- Drop the new unique constraint on user_id, site_id and slug
ALTER TABLE posts DROP CONSTRAINT IF EXISTS posts_user_id_site_id_slug_key;

-- Re-add the old unique constraint on slug
ALTER TABLE posts ADD CONSTRAINT posts_slug_key UNIQUE (slug);
