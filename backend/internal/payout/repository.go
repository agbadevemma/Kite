

package payout

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

func (r *Repository) Create(
	ctx context.Context,
	tx *sql.Tx,
	p Payout,
) error {

	query := `
	INSERT INTO payouts
	(
		id,
		user_id,
		currency,
		amount,
		destination,
		status,
		idempotency_key,
		created_at
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		p.ID,
		p.UserID,
		p.Currency,
		p.Amount,
		p.Destination,
		p.Status,
		p.IdempotencyKey,
		p.CreatedAt,
	)

	return err
}