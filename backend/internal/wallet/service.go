package wallet

import (
	"context"

	"github.com/himarnoel/kite/internal/ledger"
)

var currencies = []string{"USD", "NGN", "EUR", "GBP", "KES"}

type Service struct {
	ledgerRepo *ledger.Repository
}

func NewService(l *ledger.Repository) *Service {
	return &Service{ledgerRepo: l}
}

type Balance struct {
	Currency string `json:"currency"`
	Amount   int64  `json:"amount"`
}

func (s *Service) GetBalances(ctx context.Context, userID string) ([]Balance, error) {

	var result []Balance

	for _, currency := range currencies {
		bal, err := s.ledgerRepo.GetBalance(ctx, userID, currency)
		if err != nil {
			return nil, err
		}

		result = append(result, Balance{
			Currency: currency,
			Amount:   bal,
		})
	}

	return result, nil
}