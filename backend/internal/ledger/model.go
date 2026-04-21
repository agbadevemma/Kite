package ledger


import (
	"context"
	"database/sql"
	"time"
)


type LedgerEntry struct {
	ID        string
	UserID    string
	Currency  string
	Amount    int64 // always in smallest unit
	Type      string // "debit" | "credit"
	RefType   string
	RefID     string
	CreatedAt time.Time
}

type LedgerWriter interface {
	InsertEntries(ctx context.Context, tx *sql.Tx, entries []LedgerEntry) error
}