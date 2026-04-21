CREATE TABLE conversions (
    id UUID PRIMARY KEY,
    
    user_id UUID NOT NULL,

    from_currency TEXT,
    to_currency TEXT,

    rate NUMERIC,
    amount_in BIGINT,
    amount_out BIGINT,

    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_conversions_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_conversions_user_id ON conversions(user_id);