CREATE TABLE plans (
    id SERIAL PRIMARY KEY,
    plan_name TEXT UNIQUE NOT NULL,
    amount DECIMAL(8, 2) NOT NULL,
    interval VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
