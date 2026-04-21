package transaction

import "context"

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) List(ctx context.Context, userID string, page, limit int) ([]Transaction,int,  error) {

	offset := (page - 1) * limit

	total, err:=s.repo.CountByUser(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	txs, err := s.repo.GetTransactions(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return txs, total, nil
}