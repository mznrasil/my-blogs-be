CREATE TABLE users (
  id VARCHAR(35) PRIMARY KEY,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  profile_image TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  UNIQUE(email)
);
