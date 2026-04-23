package deposit

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/himarnoel/kite/internal/db"
	"github.com/himarnoel/kite/internal/ledger"
	"github.com/himarnoel/kite/internal/money"
	 "github.com/himarnoel/kite/internal/transactions"
)

type Service struct {
	db       *sql.DB
	repo     *Repository
	ledger   ledger.LedgerWriter
	txWriter transaction.TransactionWriter
}

func NewService(db *sql.DB, r *Repository, l ledger.LedgerWriter, t transaction.TransactionWriter) *Service {
	return &Service{
		db:       db,
		repo:     r,
		ledger:   l,
		txWriter: t,
	}
}

type Request struct {
	UserID         string
	Currency       string
	Amount         float64
	IdempotencyKey string
}

func (s *Service) Create(ctx context.Context, req Request) error {

	amount := money.ToSmallestUnit(
		req.Amount,
		req.Currency,
	)
	
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	existing, err := s.repo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	if err != nil {
		return err
	}
	if existing != nil {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	depositID := uuid.NewString()

	err = s.repo.Create(ctx, tx, Deposit{
		ID:             depositID,
		UserID:         req.UserID,
		Currency:       req.Currency,
		Amount:         amount,
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		if db.IsUniqueViolation(err) {
			return nil
		}
		return err
	}

	err = s.ledger.InsertEntries(ctx, tx, []ledger.LedgerEntry{
		{
			ID:        uuid.NewString(),
			UserID:    req.UserID,
			Currency:  req.Currency,
			Amount:    amount,
			Type:      "credit",
			RefType:   "deposit",
			RefID:     depositID,
			CreatedAt: time.Now(),
		},
	})
	if err != nil {
		return err
	}

	err = s.txWriter.Create(ctx, tx, transaction.TransactionRecord{
		ID:       uuid.NewString(),
		UserID:   req.UserID,
		Type:     "deposit",
		Status:   "completed",
		Currency: req.Currency,
		Amount:   amount,
		RefID:    depositID,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}
