package main

import "hexa/m/v2/horizon/ports/inbound"

// Idempotent is an optional Command Capability
type Idempotent interface {
	IdempotencyKey() string
}

// IdempotentCommand is a Command that supports Idempotency
type IdempotentCommand interface {
	inbound.Command
	Idempotent
}
