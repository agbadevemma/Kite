package deposit

import "time"

type Deposit struct {
	ID             string
	UserID         string
	Currency       string
	Amount         int64
	IdempotencyKey string
	CreatedAt      time.Time
}