package transaction

import (
	"context"
	"database/sql"
)

type TransactionRecord struct {
	ID       string
	UserID   string
	Type     string
	Status   string
	Currency string
	Amount   int64
	RefID    string
}

type TransactionWriter interface {
	Create(ctx context.Context, tx *sql.Tx, t TransactionRecord) error
}