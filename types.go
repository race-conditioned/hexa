package main

import "hexa/m/v2/horizon/ports/inbound"

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
	status  inbound.ResultStatus
	message string
}

func (r TransferResult) Message() string              { return r.message }
func (r TransferResult) Status() inbound.ResultStatus { return r.status }
