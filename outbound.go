package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Dispatcher submits jobs to the worker pool.
type Dispatcher interface {
	Submit(ctx context.Context, cmd TransferCommand) TransferResult
	QueueDepth() int64
	ActiveWorkers() int64
}

// Dispatcher: immediately succeed
type immediateDispatcher struct{}

// Submit: immediately return success
func (d *immediateDispatcher) Submit(_ context.Context, _ TransferCommand) TransferResult {
	fmt.Println("submitting dispatch")
	return TransferResult{"success", "ok"}
}
func (d *immediateDispatcher) QueueDepth() int64    { return 0 }
func (d *immediateDispatcher) ActiveWorkers() int64 { return 0 }

// Idempotency defines caching for idempotency keys.
type IdempotencyI[T any] interface {
	Get(key string) (T, bool)
	Store(key string, res T)
}

// A simple in-memory idempotency store for demos/tests.
type memIdempotency[T any] struct {
	mu sync.RWMutex
	m  map[string]T
}

func newMemIdempotency[T any]() *memIdempotency[T] {
	return &memIdempotency[T]{m: make(map[string]T)}
}

func (s *memIdempotency[T]) Get(key string) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[key]
	return v, ok
}

func (s *memIdempotency[T]) Store(key string, res T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = res
}

// Limiter defines rate limiting behavior.
type LimiterI interface {
	Allow(clientID string) bool
}

type dummyMetrics struct{}

func (dummyMetrics) Snapshot() MetricsSnapshot {
	return MetricsSnapshot{
		RequestsTotal: 100,
		SuccessRate:   99.5,
		AvgLatencyMs:  150.0,
		ActiveWorkers: 5,
		QueueDepth:    10,
	}
}
func (dummyMetrics) IncRequest()       { fmt.Println("metric ++ requests") }
func (dummyMetrics) IncSuccess()       { fmt.Println("metric ++ success") }
func (dummyMetrics) IncRateLimited()   { fmt.Println("metric ++ rate_limited") }
func (dummyMetrics) IncTimeout()       { fmt.Println("metric ++ timeout") }
func (dummyMetrics) IncIdempotentHit() { fmt.Println("metric ++ idempotent_hit") }
func (dummyMetrics) ObserveLatency(d time.Duration) {
	fmt.Println("metric observe latency", d.Milliseconds())
}

type dummyLimiter struct{}

func (dummyLimiter) Allow(clientID string) bool {
	fmt.Println("rate-limit check ok")
	return true
}
