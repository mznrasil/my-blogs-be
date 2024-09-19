CREATE TABLE sites (
  id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(35) NOT NULL,
  description VARCHAR(150),
  subdirectory VARCHAR(40) NOT NULL,
  image_url TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  user_id VARCHAR(35),
  UNIQUE(subdirectory),
  CONSTRAINT sites_users_id_fk 
    FOREIGN KEY (user_id) 
      REFERENCES users(id)
      ON UPDATE CASCADE
      ON DELETE CASCADE
);
