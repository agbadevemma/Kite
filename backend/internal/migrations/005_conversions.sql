CREATE TABLE conversions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    from_currency TEXT,
    to_currency TEXT,
    rate NUMERIC,
    amount_in BIGINT,
    amount_out BIGINT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);