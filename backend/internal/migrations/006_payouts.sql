CREATE TABLE payouts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    amount BIGINT NOT NULL,
    currency TEXT NOT NULL,
    status TEXT NOT NULL,
    account_number TEXT,
    bank_code TEXT,
    account_name TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);