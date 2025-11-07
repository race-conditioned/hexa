// Package intake is the standard package for defining ingress policies
package intake

// Spec is the universal ingress policy *shape*.
// It defines slots that transports interpret.
// The functions are optional; nil = no behaviour.
// Library can provide small helpers to fill them.
type Spec struct {
	// --- Wire-level guards ---
	MaxBodyBytes       int64
	MaxInFlight        int
	EnableRecover      bool   // recover panics in decode path
	EnableReqID        bool   // assign/extract request id at adapter
	EnableReqLog       bool   // basic access log at adapter
	EnableRateLimiting bool   // enable rate limiting
	IPRateLimiter      func() // optional cheap IP throttle

	// --- Behavioural slots ---
	OnBodyTooLarge    func()
	OnTooManyInFlight func()
	OnRecover         func()
	OnRequestID       func()
	OnLog             func()
	OnRateLimit       func()

	// --- Middleware chain extension ---
	Middleware []func() // user-supplied extras

	// --- App-level hints ---
	DefaultTimeoutMs int
}
