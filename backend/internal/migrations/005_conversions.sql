CREATE TABLE conversions (
    id UUID PRIMARY KEY,
    
    user_id UUID NOT NULL,

    from_currency TEXT NOT NULL,
    to_currency TEXT NOT NULL,

    rate NUMERIC NOT NULL,
    amount_in BIGINT NOT NULL,
    amount_out BIGINT NOT NULL,
    fee BIGINT NOT NULL,

    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_conversions_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_conversions_user_id ON conversions(user_id);