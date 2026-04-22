package payout

import "time"

type Payout struct {
	ID             string `json:"id"`
	UserID         string `json:"userId"`
	Currency       string `json:"currency"`
	Amount         int64  `json:"amount"`
	Destination    string `json:"destination"`
	Status         string `json:"status"`
	IdempotencyKey string `json:"idempotencyKey"`
	CreatedAt      time.Time `json:"createdAt"`
}