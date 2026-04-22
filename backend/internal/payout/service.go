
package payout

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/himarnoel/kite/internal/ledger"
	"github.com/himarnoel/kite/internal/transactions"
)

type Service struct {
	repo         *Repository
	ledgerWriter ledger.LedgerWriter
	txWriter     transaction.TransactionWriter
	db           *sql.DB
}

func NewService(
	r *Repository,
	lw ledger.LedgerWriter,
	tw transaction.TransactionWriter,
	db *sql.DB,
) *Service {
	return &Service{
		repo:         r,
		ledgerWriter: lw,
		txWriter:     tw,
		db:           db,
	}
}

func (s *Service) Create(
	ctx context.Context,
	userID string,
	req CreatePayoutDTO	,
) error {

	balance, err := s.ledgerWriter.GetBalance(
		ctx,
		userID,
		req.Currency,
	)
	if err != nil {
		return err
	}

	if balance < req.Amount {
		return errors.New("insufficient balance")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	payoutID := uuid.NewString()

	err = s.repo.Create(ctx, tx, Payout{
		ID:             payoutID,
		UserID:         userID,
		Currency:       req.Currency,
		Amount:         req.Amount,
		Destination:    req.Destination,
		Status:         "completed",
		IdempotencyKey: req.IdempotencyKey,
		CreatedAt:      time.Now(),
	})
	if err != nil {
		return err
	}

	err = s.txWriter.Create(ctx, tx, transaction.TransactionRecord{
		ID:       uuid.NewString(),
		UserID:   userID,
		Type:     "payout",
		Status:   "completed",
		Currency: req.Currency,
		Amount:   req.Amount,
		RefID:    payoutID,
	})
	if err != nil {
		return err
	}

	err = s.ledgerWriter.InsertEntries(
		ctx,
		tx,
		[]ledger.LedgerEntry{
			{
				ID:        uuid.NewString(),
				UserID:    userID,
				Currency:  req.Currency,
				Amount:    req.Amount,
				Type:      "debit",
				RefType:   "payout",
				RefID:     payoutID,
				CreatedAt: time.Now(),
			},
		},
	)

	if err != nil {
		return err
	}

	return tx.Commit()
}