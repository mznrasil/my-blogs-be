CREATE TABLE subscriptions (
    id VARCHAR(36) PRIMARY KEY,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    user_id VARCHAR(35) UNIQUE NOT NULL,
    plan_id INTEGER NOT NULL,
    payment_id VARCHAR(36),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT subscriptions_users_id_fk
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON UPDATE CASCADE,
    CONSTRAINT subscriptions_plans_id_fk
        FOREIGN KEY (plan_id)
        REFERENCES plans(id)
        ON UPDATE CASCADE,
    CONSTRAINT subscriptions_payments_id_fk
        FOREIGN KEY (payment_id)
        REFERENCES payments(id)
        ON UPDATE CASCADE
);
