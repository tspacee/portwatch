package watch

import (
	"errors"
	"sync"
	"time"
)

// Horizon tracks how long each port has been continuously open.
// It records the first time a port was observed and computes its age.
type Horizon struct {
	mu      sync.Mutex
	first   map[int]time.Time
	cutoff  time.Duration
}

// NewHorizon creates a Horizon with the given age cutoff.
// Ports older than cutoff are considered "beyond the horizon".
func NewHorizon(cutoff time.Duration) (*Horizon, error) {
	if cutoff <= 0 {
		return nil, errors.New("horizon: cutoff must be positive")
	}
	return &Horizon{
		first:  make(map[int]time.Time),
		cutoff: cutoff,
	}, nil
}

// Observe records the first-seen time for a port if not already tracked.
func (h *Horizon) Observe(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("horizon: port out of range")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.first[port]; !ok {
		h.first[port] = time.Now()
	}
	return nil
}

// Age returns how long the port has been observed. Returns zero if never seen.
func (h *Horizon) Age(port int) time.Duration {
	h.mu.Lock()
	defer h.mu.Unlock()
	t, ok := h.first[port]
	if !ok {
		return 0
	}
	return time.Since(t)
}

// Beyond reports whether the port has been open longer than the cutoff.
func (h *Horizon) Beyond(port int) bool {
	return h.Age(port) >= h.cutoff
}

// Forget removes tracking for a port.
func (h *Horizon) Forget(port int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.first, port)
}

// Len returns the number of tracked ports.
func (h *Horizon) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.first)
}
