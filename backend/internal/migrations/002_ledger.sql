CREATE TABLE ledger_entries (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    currency TEXT NOT NULL,
    amount BIGINT NOT NULL,
    type TEXT NOT NULL, -- debit / credit
    ref_type TEXT,
    ref_id TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_ledger_user_currency ON ledger_entries(user_id, currency);