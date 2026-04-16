package watch

import (
	"errors"
	"sync"
	"time"
)

// HealthStatus represents the current health of the watcher.
type HealthStatus struct {
	Healthy     bool
	LastSuccess time.Time
	LastError   error
	ConsecFails int
}

// HealthChecker tracks watcher health based on scan outcomes.
type HealthChecker struct {
	mu          sync.RWMutex
	status      HealthStatus
	maxFails    int
}

// NewHealthChecker creates a HealthChecker that marks unhealthy after maxFails consecutive failures.
func NewHealthChecker(maxFails int) (*HealthChecker, error) {
	if maxFails < 1 {
		return nil, errors.New("maxFails must be at least 1")
	}
	return &HealthChecker{
		maxFails: maxFails,
		status:   HealthStatus{Healthy: true},
	}, nil
}

// RecordSuccess marks the last scan as successful.
func (h *HealthChecker) RecordSuccess() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.status.LastSuccess = time.Now()
	h.status.ConsecFails = 0
	h.status.LastError = nil
	h.status.Healthy = true
}

// RecordFailure records a scan failure and may mark the checker unhealthy.
func (h *HealthChecker) RecordFailure(err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.status.LastError = err
	h.status.ConsecFails++
	if h.status.ConsecFails >= h.maxFails {
		h.status.Healthy = false
	}
}

// Status returns a snapshot of the current health status.
func (h *HealthChecker) Status() HealthStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.status
}
