package transaction

import "time"

type Transaction struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Type      string    `json:"type"`   // deposit | conversion | payout
	Status    string    `json:"status"` // pending | completed | failed
	Currency  string    `json:"currency"`
	Amount    int64     `json:"amount"`
	RefID     string    `json:"ref_id"`
	CreatedAt time.Time `json:"created_at"`
}