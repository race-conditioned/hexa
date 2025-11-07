package main

import "time"

// Metrics aggregates all metric capabilities.
// Existing code can still depend on this for convenience.
type Metrics interface {
	CounterMetrics
	LatencyMetrics
	SnapshotMetrics
}

// CounterMetrics defines event counters.
type CounterMetrics interface {
	IncRequest()
	IncSuccess()
	IncRateLimited()
	IncTimeout()
	IncIdempotentHit()
}

// LatencyMetrics defines latency observation.
type LatencyMetrics interface {
	ObserveLatency(d time.Duration)
}

// SnapshotMetrics defines exportable snapshotting.
type SnapshotMetrics interface {
	Snapshot() MetricsSnapshot
}

// MetricsSnapshot defines JSON output for /metrics.
type MetricsSnapshot struct {
	RequestsTotal int64   `json:"requests_total"`
	SuccessRate   float64 `json:"success_rate"`
	AvgLatencyMs  float64 `json:"avg_latency_ms"`
	ActiveWorkers int64   `json:"active_workers"`
	QueueDepth    int64   `json:"queue_depth"`
}
