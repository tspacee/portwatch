package watch

import (
	"errors"
	"sync"
	"time"
)

// ErrInvalidTTL is returned when a zero or negative TTL is provided.
var ErrInvalidTTL = errors.New("ttl must be greater than zero")

// Expiry tracks per-port expiration timestamps.
type Expiry struct {
	mu      sync.Mutex
	ttl     time.Duration
	entries map[int]time.Time
}

// NewExpiry creates a new Expiry tracker with the given TTL.
func NewExpiry(ttl time.Duration) (*Expiry, error) {
	if ttl <= 0 {
		return nil, ErrInvalidTTL
	}
	return &Expiry{
		ttl:     ttl,
		entries: make(map[int]time.Time),
	}, nil
}

// Set records the current time for the given port.
func (e *Expiry) Set(port int) error {
	if port < 1 || port > 65535 {
		return ErrInvalidPort
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.entries[port] = time.Now()
	return nil
}

// Expired reports whether the given port's entry has expired.
// Ports never set are considered expired.
func (e *Expiry) Expired(port int) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	t, ok := e.entries[port]
	if !ok {
		return true
	}
	return time.Since(t) >= e.ttl
}

// Delete removes the entry for the given port.
func (e *Expiry) Delete(port int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.entries, port)
}

// Len returns the number of tracked ports.
func (e *Expiry) Len() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.entries)
}
