package main

import (
	"hexa/m/v2/symphony"
)

// PolicyStage represents different stages where policies can be applied.

// Policy represents a middleware policy at a specific stage.
type Policy struct {
	Stage symphony.PolicyStage
	MW    any
}

// DefaultPolicyOrder defines the order in which policies are applied.
// Idempotency policy is fired before rate limit policy as a design decision.
// Retries are expected (client backoff + jitter).
// When a duplicate key is used, the cached result is returned immediately, even if the client has hit other limits.
// That reduces client churn and lowers total work across the system.
// Idempotency lookup is fast & cheap (in-memory) and resilient.
// Under a thundering herd of retries, this short-circuits work earlier than rate limit does.
var DefaultPolicyOrder = []symphony.PolicyStage{
	"idempotency",
	"rate_limit",
	"timeout",
	"latency",
}

//
// // Idempotency is a middleware that provides idempotency support for commands implementing IdempotentCommand.
// func Idempotency[Com inbound.IdempotentCommand, Res inbound.Result](
// 	store IdempotencyI[Res],
// 	counterMetrics CounterMetrics,
// ) inbound.UnaryMiddleware {
// 	return func(next inbound.UnaryHandler) inbound.UnaryHandler {
// 		return func(ctx context.Context, meta inbound.RequestMeta, cmd inbound.Command) (inbound.Result, error) {
// 			/*
// 				if store != nil {
// 					if cached, ok := store.Get(cmd.IdempotencyKey()); ok {
// 						if counterMetrics != nil {
// 							counterMetrics.IncIdempotentHit()
// 						}
// 						return cached, nil
// 					}
// 				}
// 			*/
// 			res, err := next(ctx, meta, cmd)
// 			/*
// 				if err == nil && store != nil && cmd.IdempotencyKey() != "" {
// 					store.Store(cmd.IdempotencyKey(), res)
// 				}
// 			*/
// 			return res, err
// 		}
// 	}
// }
