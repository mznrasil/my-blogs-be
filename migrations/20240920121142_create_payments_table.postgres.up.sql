CREATE TYPE payment_status AS ENUM ('Initiated', 'Pending', 'Completed', 'Expired', 'Refunded', 'User canceled', 'Partially Refunded');
CREATE TABLE payments (
    id VARCHAR(36) PRIMARY KEY,
    pidx TEXT UNIQUE NOT NULL,
    status payment_status NOT NULL,
    transaction_id TEXT,
    amount DECIMAL(8, 2),
    mobile CHAR(10),
    total_amount DECIMAL(8, 2),
    plan_id INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT payments_plans_id_fk
        FOREIGN KEY (plan_id)
        REFERENCES plans(id)
        ON UPDATE CASCADE
);
