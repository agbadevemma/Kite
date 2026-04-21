package conversion

import "time"

type Quote struct {
	ID           string
	UserID       string
	FromCurrency string
	ToCurrency   string
	Rate         float64
	AmountIn     int64
	AmountOut    int64
	Fee          int64
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

type ExecuteRequest struct {
	QuoteID string `json:"quote_id" binding:"required"`
}