package watch

import (
	"errors"
	"sync"
	"time"
)

// ErrCircuitOpen is returned when the circuit breaker is open.
var ErrCircuitOpen = errors.New("circuit breaker is open")

// CircuitState represents the state of the circuit breaker.
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

// CircuitBreaker trips after a threshold of consecutive failures and
// recovers after a configurable timeout.
type CircuitBreaker struct {
	mu           sync.Mutex
	state        CircuitState
	failures     int
	threshold    int
	recoveryTime time.Duration
	openedAt     time.Time
}

// NewCircuitBreaker creates a CircuitBreaker with the given failure threshold
// and recovery timeout. Both values must be positive.
func NewCircuitBreaker(threshold int, recovery time.Duration) (*CircuitBreaker, error) {
	if threshold <= 0 {
		return nil, errors.New("threshold must be greater than zero")
	}
	if recovery <= 0 {
		return nil, errors.New("recovery duration must be greater than zero")
	}
	return &CircuitBreaker{
		threshold:    threshold,
		recoveryTime: recovery,
	}, nil
}

// Allow returns nil if the call is permitted, or ErrCircuitOpen if not.
func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	switch cb.state {
	case StateOpen:
		if time.Since(cb.openedAt) >= cb.recoveryTime {
			cb.state = StateHalfOpen
			return nil
		}
		return ErrCircuitOpen
	}
	return nil
}

// RecordSuccess resets the breaker on a successful call.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	cb.state = StateClosed
}

// RecordFailure increments the failure count and trips the breaker if needed.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	if cb.failures >= cb.threshold {
		cb.state = StateOpen
		cb.openedAt = time.Now()
	}
}

// State returns the current circuit state.
func (cb *CircuitBreaker) State() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}
