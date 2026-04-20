CREATE TABLE deposits (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    amount BIGINT NOT NULL,
    currency TEXT NOT NULL,
    idempotency_key TEXT UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);