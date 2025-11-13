package main

import (
	"fmt"

	"github.com/race-conditioned/hexa/apperr"
	"github.com/race-conditioned/hexa/horizon/ports/inbound"
)

// RateLimit is a middleware that enforces rate limiting based on the provided limiter.
func RateLimit(next AppHandler) AppHandler {
	return func(ctx AppCtxI, meta inbound.RequestMeta, cmd inbound.Command) (inbound.Result, error) {
		fmt.Println("Applying rate limit...")
		var zero inbound.Result
		if !ctx.Limiter().Allow(meta.ClientID) {
			ctx.Metrics().IncRateLimited()
			return zero, apperr.RateLimited("rate limit exceeded")
		}
		return next(ctx, meta, cmd)
	}
}
