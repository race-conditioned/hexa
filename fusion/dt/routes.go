package dt

import (
	"hexa/m/v2/horizon"
	"hexa/m/v2/horizon/ports/inbound"
)

// Route describes how an HTTP path maps to a gateway handler.
type Route[c inbound.Ctx] struct {
	HandlerKey horizon.HandlerKey
	Path       string     // optional; fallback to "/"+handlerKey
	NewPayload func() any // <— just “give me a pointer to something”
}

// JSONRoute creates a route that decodes JSON into T and sends it to the gateway handler.
func JSONRoute[c inbound.Ctx, T any](key horizon.HandlerKey) Route[c] {
	return Route[c]{
		HandlerKey: key,
		NewPayload: func() any { return new(T) },
	}
}

// JSONRoutePath is the same but lets you override the path (/v1/...).
func JSONRoutePath[c inbound.Ctx, T any](key horizon.HandlerKey, path string) Route[c] {
	return Route[c]{
		HandlerKey: key,
		Path:       path,
		NewPayload: func() any { return new(T) },
	}
}
