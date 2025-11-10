package main

import (
	"hexa/m/v2/horizon/ports/inbound"
)

type TransferCommandHTTP struct {
	FromAccount    string `json:"from_account"`
	ToAccount      string `json:"to_account"`
	AmountCents    int64  `json:"amount_cents"`
	IdempotencyKey string `json:"idempotency_key"`
}

func (dto *TransferCommandHTTP) ToCommand() inbound.Command {
	return TransferCommand{
		fromAccount:    dto.FromAccount,
		toAccount:      dto.ToAccount,
		amountCents:    dto.AmountCents,
		idempotencyKey: dto.IdempotencyKey,
	}
}

// TransferCommand defines the external API payload for /transfer from any transport.
type TransferCommand struct {
	fromAccount    string
	toAccount      string
	amountCents    int64
	idempotencyKey string
}

// FromAccount returns the source account ID.
func (t TransferCommand) FromAccount() string {
	return t.fromAccount
}

// ToAccount returns the destination account ID.
func (t TransferCommand) ToAccount() string {
	return t.toAccount
}

// AmountCents returns the transfer amount in cents.
func (t TransferCommand) AmountCents() int64 {
	return t.amountCents
}

// IdempotencyKey returns the idempotency key for the transfer.
func (t TransferCommand) IdempotencyKey() string {
	return t.idempotencyKey
}

type TransferResult struct {
	transactionID string
	status        inbound.ResultStatus
	message       string
}

func (r TransferResult) TransactionID() string        { return r.transactionID }
func (r TransferResult) Message() string              { return r.message }
func (r TransferResult) Status() inbound.ResultStatus { return r.status }

func (r TransferResult) Encode(s inbound.Sink) {
	s.Write(r.status.String(), TransferResponse{
		TransactionID: r.transactionID,
		Status:        r.status.String(),
		Message:       r.message,
	})
}

type TransferResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}
