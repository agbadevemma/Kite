CREATE TABLE ledger_entries (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    currency TEXT NOT NULL,
    amount BIGINT NOT NULL,
    type TEXT NOT NULL, -- debit / credit
    ref_type TEXT,
    ref_id TEXT,
    created_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_ledger_entries_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_ledger_user_currency ON ledger_entries(user_id, currency);