package transaction

import "context"

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) List(ctx context.Context, userID string, page, limit int) ([]Transaction, error) {

	offset := (page - 1) * limit

	return s.repo.GetTransactions(ctx, userID, limit, offset)
}