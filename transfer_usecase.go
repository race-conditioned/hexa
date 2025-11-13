package main

import (
	"errors"
	"fmt"

	"github.com/race-conditioned/hexa/apperr"
)

// TransferUsecase handles transfer requests.
type TransferUsecase struct {
	dispatcher Dispatcher
	metrics    Metrics
}

// NewTransferUsecase creates a new TransferUsecase.
func NewTransferUsecase(d Dispatcher, m Metrics) *TransferUsecase {
	return &TransferUsecase{dispatcher: d, metrics: m}
}

// SubmitTransfer is a usecase that validates and submits a transfer command.
func (s *TransferUsecase) SubmitTransfer(ctx AppCtxI, cmd TransferCommand) (TransferResult, error) {
	fmt.Printf("Processing transfer from %s to %s of %d\n", cmd.FromAccount(), cmd.ToAccount(), cmd.AmountCents())
	if err := validate(cmd); err != nil {
		return TransferResult{}, apperr.Invalid(err.Error())
	}

	// Delegate to worker pool via outbound port (transport-agnostic).
	return s.dispatcher.Submit(ctx, cmd), nil
}

// validate checks the transfer command for required fields.
func validate(cmd TransferCommand) error {
	if cmd.FromAccount() == "" || cmd.ToAccount() == "" {
		return errors.New("missing account IDs")
	}
	if cmd.AmountCents() <= 0 {
		return errors.New("amount must be positive")
	}
	if cmd.IdempotencyKey() == "" {
		return errors.New("missing idempotency key")
	}
	return nil
}
