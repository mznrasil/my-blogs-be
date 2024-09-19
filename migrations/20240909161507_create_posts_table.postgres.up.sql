CREATE TABLE posts (
  id VARCHAR(36) PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  article_content JSONB,
  small_description VARCHAR(255),
  image TEXT,
  slug VARCHAR(255),
  created_at TIMESTAMP, 
  updated_at TIMESTAMP, 
  user_id VARCHAR(35),
  site_id VARCHAR(36),
  UNIQUE(slug),
  CONSTRAINT posts_users_id_fk
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT posts_sites_id_fk
    FOREIGN KEY (site_id)
    REFERENCES sites(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
