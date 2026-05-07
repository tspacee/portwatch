package watch

import (
	"errors"
	"sync"
	"time"
)

// Trench tracks ports that have been continuously absent for a configurable
// duration. Once a port exceeds the absence window it is considered entrenched.
type Trench struct {
	mu      sync.Mutex
	window  time.Duration
	absent  map[int]time.Time
}

// NewTrench creates a Trench with the given absence window.
func NewTrench(window time.Duration) (*Trench, error) {
	if window <= 0 {
		return nil, errors.New("trench: window must be positive")
	}
	return &Trench{
		window: window,
		absent: make(map[int]time.Time),
	}, nil
}

// Observe records the current set of open ports. Ports not present in the
// provided slice are tracked as absent from this moment onward (if not already
// tracked). Ports that reappear are removed from the absence map.
func (t *Trench) Observe(ports []int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	present := make(map[int]struct{}, len(ports))
	for _, p := range ports {
		present[p] = struct{}{}
	}

	// Remove ports that have reappeared.
	for p := range t.absent {
		if _, ok := present[p]; ok {
			delete(t.absent, p)
		}
	}
}

// MarkAbsent records port p as absent starting now, if it is not already
// tracked. Returns an error for invalid port numbers.
func (t *Trench) MarkAbsent(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("trench: port out of range")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.absent[port]; !exists {
		t.absent[port] = time.Now()
	}
	return nil
}

// Entrenched returns true if port has been absent for at least the configured
// window duration.
func (t *Trench) Entrenched(port int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	since, ok := t.absent[port]
	if !ok {
		return false
	}
	return time.Since(since) >= t.window
}

// Forget removes a port from absence tracking entirely.
func (t *Trench) Forget(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.absent, port)
}

// Len returns the number of ports currently tracked as absent.
func (t *Trench) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.absent)
}
