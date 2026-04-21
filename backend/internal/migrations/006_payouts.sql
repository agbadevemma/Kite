CREATE TABLE payouts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    amount BIGINT NOT NULL,
    currency TEXT NOT NULL,
    status TEXT NOT NULL,
    account_number TEXT,
    bank_code TEXT,
    account_name TEXT,
    created_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_payouts_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_payouts_user_id ON payouts(user_id);
CREATE INDEX idx_payouts_status ON payouts(status);