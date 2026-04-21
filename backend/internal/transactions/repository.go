package transaction

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

// Unified query (simplified)
func (r *Repository) GetTransactions(ctx context.Context, userID string, limit, offset int) ([]Transaction, error) {

	query := `
		SELECT id, user_id, type, status, currency, amount, ref_id, created_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []Transaction

	for rows.Next() {
		var t Transaction

		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Type,
			&t.Status,
			&t.Currency,
			&t.Amount,
			&t.RefID,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		txs = append(txs, t)
	}

	return txs, nil
}

func (r *Repository) Create(
	ctx context.Context,
	tx *sql.Tx,
	t TransactionRecord,
) error {

	query := `
		INSERT INTO transactions
		(id, user_id, type, status, currency, amount, ref_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`

	_, err := tx.ExecContext(ctx, query,
		t.ID,
		t.UserID,
		t.Type,
		t.Status,
		t.Currency,
		t.Amount,
		t.RefID,
	)

	return err
}
