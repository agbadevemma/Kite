package ledger

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

// Insert entries (used in transactions)
func (r *Repository) InsertEntries(ctx context.Context, tx *sql.Tx, entries []LedgerEntry) error {

	query := `
		INSERT INTO ledger_entries (id, user_id, currency, amount, type, ref_type, ref_id, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`

	for _, e := range entries {
		_, err := tx.ExecContext(ctx, query,
			e.ID,
			e.UserID,
			e.Currency,
			e.Amount,
			e.Type,
			e.RefType,
			e.RefID,
			e.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get balance (derived from ledger)
func (r *Repository) GetBalance(ctx context.Context, userID, currency string) (int64, error) {

	query := `
		SELECT COALESCE(SUM(
			CASE 
				WHEN type = 'credit' THEN amount
				ELSE -amount
			END
		), 0)
		FROM ledger_entries
		WHERE user_id = $1 AND currency = $2
	`

	var balance int64
	err := r.db.QueryRowContext(ctx, query, userID, currency).Scan(&balance)
	return balance, err
}