package conversion

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/himarnoel/kite/internal/ledger"
	 "github.com/himarnoel/kite/internal/transactions"
)

type Service struct {
	repo              *Repository
	ledgerWriter      ledger.LedgerWriter
	transactionWriter transaction.TransactionWriter
	db                *sql.DB
}

func NewService(r *Repository, lw ledger.LedgerWriter, tw transaction.TransactionWriter, db *sql.DB) *Service {
	return &Service{
		repo:              r,
		ledgerWriter:      lw,
		transactionWriter: tw,
		db:                db,
	}
}

func (s *Service) Quote(ctx context.Context, userID, from, to string, amount int64) (*Quote, error) {

	if from == to {
		return nil, errors.New("invalid conversion pair")
	}

	balance, err := s.ledgerWriter.GetBalance(
		ctx,
		userID,
		from,
	)
	log.Println("Retrieved balance:", balance)

	if err != nil {
		return nil, err
	}

	log.Println("Checking if balance is sufficient for conversion:", amount)

	if balance < amount {
		return nil, errors.New("insufficient balance")
	}
	rate := GetRate(from, to)

	spread := 0.01
	finalRate := rate * (1 + spread)

	out := float64(amount) * finalRate
	fee := int64(math.Round(out * 0.01))

	quote := &Quote{
		ID:           uuid.NewString(),
		UserID:       userID,
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         finalRate,
		AmountIn:     amount,
		AmountOut:    int64(out) - fee,
		Fee:          fee,
		ExpiresAt:    time.Now().Add(60 * time.Second),
		CreatedAt:    time.Now(),
	}

	err = s.repo.CreateQuote(ctx, *quote)

	if err != nil {
		return nil, err
	}

	return quote, nil
}

func (s *Service) Execute(ctx context.Context, quoteID string) error {

	quote, err := s.repo.GetQuote(ctx, quoteID)
	if err != nil {
		return err
	}

	if time.Now().After(quote.ExpiresAt) {
		return errors.New("quote expired")
	}

	balance, err := s.ledgerWriter.GetBalance(
		ctx,
		quote.UserID,
		quote.FromCurrency,
	)
	log.Println("Retrieved balance:", balance)

	if err != nil {
		return err
	}
	log.Println("Checking if balance is sufficient for conversion:", quote.AmountIn)

	totalAmount := quote.AmountIn + quote.Fee
	if balance < totalAmount {
		return errors.New("insufficient balance")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	txRecord := transaction.TransactionRecord{
		ID:       uuid.NewString(),
		UserID:   quote.UserID,
		Type:     "conversion",
		Status:   "completed",
		Currency: quote.FromCurrency,
		Amount:   totalAmount,
		RefID:    quote.ID,
	}

	err = s.transactionWriter.Create(ctx, tx, txRecord)
	if err != nil {
		return err
	}

	entries := []ledger.LedgerEntry{
		{
			ID:        uuid.NewString(),
			UserID:    quote.UserID,
			Currency:  quote.FromCurrency,
			Amount:    quote.AmountIn,
			Type:      "debit",
			RefType:   "conversion",
			RefID:     quote.ID,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.NewString(),
			UserID:    quote.UserID,
			Currency:  quote.ToCurrency,
			Amount:    quote.AmountOut,
			Type:      "credit",
			RefType:   "conversion",
			RefID:     quote.ID,
			CreatedAt: time.Now(),
		},
	}

	err = s.ledgerWriter.InsertEntries(ctx, tx, entries)
	if err != nil {
		return err
	}

	return tx.Commit()
}
