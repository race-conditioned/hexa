package main

import (
	"fmt"
	"hexa/m/v2/horizon/ports/inbound"
)

// CountRequests is a middleware that counts incoming requests using the provided metrics.
func CountRequests(next AppHandler) AppHandler {
	return func(ctx AppCtxI, meta inbound.RequestMeta, cmd inbound.Command) (inbound.Result, error) {
		fmt.Println("Counting request...")
		ctx.Metrics().IncRequest()
		return next(ctx, meta, cmd)
	}
}

// CountSuccess is a middleware that counts incoming requests using the provided metrics.
func CountSuccess(next AppHandler) AppHandler {
	return func(ctx AppCtxI, meta inbound.RequestMeta, cmd inbound.Command) (inbound.Result, error) {
		r, err := next(ctx, meta, cmd)
		if err == nil {
			fmt.Println("Counting success...")
			ctx.Metrics().IncSuccess()
		}
		return r, err
	}
}
