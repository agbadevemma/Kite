package deposit

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByIdempotencyKey(ctx context.Context, key string) (*Deposit, error) {
	query := `
		SELECT id, user_id, currency, amount, idempotency_key, created_at
		FROM deposits
		WHERE idempotency_key = $1
	`

	var d Deposit
	err := r.db.QueryRowContext(ctx, query, key).Scan(
		&d.ID,
		&d.UserID,
		&d.Currency,
		&d.Amount,
		&d.IdempotencyKey,
		&d.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *Repository) Create(ctx context.Context, tx *sql.Tx, d Deposit) error {
	query := `
		INSERT INTO deposits (id, user_id, currency, amount, idempotency_key)
		VALUES ($1,$2,$3,$4,$5)
	`

	_, err := tx.ExecContext(ctx, query,
		d.ID,
		d.UserID,
		d.Currency,
		d.Amount,
		d.IdempotencyKey,
	)

	return err
}