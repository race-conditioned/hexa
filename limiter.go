package main

// Limiter defines rate limiting behavior.
type Limiter interface {
	Allow(clientID string) bool
}
