CREATE TABLE subscriptions (
  stripe_subscription_id TEXT PRIMARY KEY NOT NULL,
  interval VARCHAR(100) NOT NULL,
  status VARCHAR(100) NOT NULL,
  plan_id TEXT NOT NULL,
  current_period_start INTEGER NOT NULL,
  current_period_end INTEGER NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  user_id VARCHAR(35),
  CONSTRAINT subscriptions_users_id_fk
    FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE,
  UNIQUE(user_id)
);
