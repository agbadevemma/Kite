package conversion

import "time"

type Quote struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	FromCurrency string	`json:"from_currency"`
	ToCurrency   string `json:"to_currency"`
	Rate         float64 `json:"rate"`
	AmountIn     int64 `json:"amount_in"`
	AmountOut    int64 `json:"amount_out"`
	Fee          int64 `json:"fee"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`	
}

