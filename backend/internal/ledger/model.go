package ledger

import "time"

type Entry struct {
	ID        string
	UserID    string
	Currency  string
	Amount    int64 // always in smallest unit
	Type      string // "debit" | "credit"
	RefType   string
	RefID     string
	CreatedAt time.Time
}