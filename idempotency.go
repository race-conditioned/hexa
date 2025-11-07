package main

import (
	"fmt"
	"hexa/m/v2/horizon/ports/inbound"
)

// store IdempotencyI[Res],
// counterMetrics CounterMetrics,

type IdempotentHandler = inbound.UnaryHandler[AppCtxI, IdempotentCommand, inbound.Result]

// Idempotency is a middleware that provides idempotency support for commands implementing IdempotentCommand.
func Idempotency(next IdempotentHandler) IdempotentHandler {
	return func(ctx AppCtxI, meta inbound.RequestMeta, cmd IdempotentCommand) (inbound.Result, error) {
		fmt.Println("idempotency")
		// WARN: can store be nil?
		if cached, ok := ctx.Idempotency().Get(cmd.IdempotencyKey()); ok {
			ctx.Metrics().IncIdempotentHit()
			return cached, nil
		}

		res, err := next(ctx, meta, cmd)
		// WARN: haven't checked if idempotencykey can be empty
		if err == nil {
			ctx.Idempotency().Store(cmd.IdempotencyKey(), res)
		}

		return res, err
	}
}

// // Idempotency is a middleware that provides idempotency support for commands implementing IdempotentCommand.
// func Idempotency[Com inbound.IdempotentCommand, Res inbound.Result](
// 	store IdempotencyI[Res],
// 	counterMetrics outbound.CounterMetrics,
// ) inbound.UnaryMiddleware[Com, Res] {
// 	return func(next inbound.UnaryHandler[Com, Res]) inbound.UnaryHandler[Com, Res] {
// 		return func(ctx context.Context, meta inbound.RequestMeta, cmd Com) (Res, error) {
// 			if store != nil {
// 				if cached, ok := store.Get(cmd.IdempotencyKey()); ok {
// 					if counterMetrics != nil {
// 						counterMetrics.IncIdempotentHit()
// 					}
// 					return cached, nil
// 				}
// 			}
// 			res, err := next(ctx, meta, cmd)
// 			if err == nil && store != nil && cmd.IdempotencyKey() != "" {
// 				store.Store(cmd.IdempotencyKey(), res)
// 			}
// 			return res, err
// 		}
// 	}
// }
