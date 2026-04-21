CREATE TABLE transactions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL,              -- deposit | conversion | payout
    status TEXT NOT NULL,            -- pending | processing | completed | failed

    currency TEXT NOT NULL,
    amount BIGINT NOT NULL,          -- smallest unit (kobo, cents)

    ref_id TEXT,                     -- links to deposit / conversion / payout

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);
CREATE INDEX idx_transactions_type ON transactions(type);