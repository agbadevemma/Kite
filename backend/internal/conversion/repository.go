package conversion

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

func (r *Repository) CreateQuote(ctx context.Context, q Quote) error {
	query := `
	INSERT INTO  conversions
	(id, user_id, from_currency, to_currency, rate, amount_in, amount_out, fee, expires_at, created_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`

	_, err := r.db.ExecContext(ctx, query,
		q.ID, q.UserID, q.FromCurrency, q.ToCurrency,
		q.Rate, q.AmountIn, q.AmountOut, q.Fee,
		q.ExpiresAt, q.CreatedAt,
	)

	return err
}

func (r *Repository) GetQuote(ctx context.Context, id string) (*Quote, error) {
	query := `
	SELECT id, user_id, from_currency, to_currency, rate, amount_in, amount_out, fee, expires_at, created_at
	FROM conversion_quotes WHERE id = $1
	`

	var q Quote

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&q.ID,
		&q.UserID,
		&q.FromCurrency,
		&q.ToCurrency,
		&q.Rate,
		&q.AmountIn,
		&q.AmountOut,
		&q.Fee,
		&q.ExpiresAt,
		&q.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &q, nil
}
