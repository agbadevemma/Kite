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
	"github.com/himarnoel/kite/internal/money"
	transaction "github.com/himarnoel/kite/internal/transactions"
)

type Service struct {
	repo              *Repository
	ledgerWriter      ledger.LedgerWriter
	transactionWriter transaction.TransactionWriter
	db                *sql.DB
}

func NewService(
	r *Repository,
	lw ledger.LedgerWriter,
	tw transaction.TransactionWriter,
	db *sql.DB,
) *Service {
	return &Service{
		repo:              r,
		ledgerWriter:      lw,
		transactionWriter: tw,
		db:                db,
	}
}

func (s *Service) Quote(
	ctx context.Context,
	userID,
	from,
	to string,
	amount float64,
) (*Quote, error) {

	if from == to {
		return nil, errors.New("invalid conversion pair")
	}

	amountInSmallestUnit := money.ToSmallestUnit(
		amount,
		from,
	)

	balance, err := s.ledgerWriter.GetBalance(
		ctx,
		userID,
		from,
	)
	if err != nil {
		return nil, err
	}


	log.Println(
		"Checking conversion balance:",
		balance,
		"required:",
		amountInSmallestUnit,
	)
	
	if balance < amountInSmallestUnit {
	return nil, errors.New("insufficient balance")
}


	rate := GetRate(from, to)

	spread := 0.01
	finalRate := rate * (1 + spread)

	amountMajor := float64(amountInSmallestUnit) / 100
	convertedMajor := amountMajor * finalRate
	feeMajor := convertedMajor * 0.01
	amountOutMajor := convertedMajor - feeMajor

	amountOut := int64(math.Round(amountOutMajor * 100))
	fee := int64(math.Round(feeMajor * 100))


	quote := &Quote{
		ID:           uuid.NewString(),
		UserID:       userID,
		FromCurrency: from,
		ToCurrency:   to,

		Rate: finalRate,

		// stored in smallest units
		AmountIn: amountInSmallestUnit,

		// stored in smallest units
		AmountOut: amountOut,

		// stored in smallest units
		Fee: fee,

		ExpiresAt: time.Now().Add(60 * time.Second),
		CreatedAt: time.Now(),
	}

	err = s.repo.CreateQuote(ctx, *quote)
	if err != nil {
		return nil, err
	}

	return quote, nil
}

func (s *Service) Execute(
	ctx context.Context,
	quoteID string,
) error {

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
	if err != nil {
		return err
	}

	// user pays amount_in only
	// fee already deducted from amount_out
	requiredBalance := quote.AmountIn

	if balance < requiredBalance {
		return errors.New("insufficient balance")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	txRecord := transaction.TransactionRecord{
		ID:     uuid.NewString(),
		UserID: quote.UserID,
		Type:   "conversion",
		Status: "completed",

		Currency: quote.FromCurrency,

		Amount: quote.AmountIn,

		RefID: quote.ID,
	}

	err = s.transactionWriter.Create(
		ctx,
		tx,
		txRecord,
	)
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
			ID:       uuid.NewString(),
			UserID:   quote.UserID,
			Currency: quote.ToCurrency,

			Amount: quote.AmountOut,

			Type: "credit",

			RefType: "conversion",
			RefID:   quote.ID,

			CreatedAt: time.Now(),
		},
	}

	err = s.ledgerWriter.InsertEntries(
		ctx,
		tx,
		entries,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
