package watch

import (
	"errors"
	"sync"
	"time"
)

// Grace tracks a per-port grace period, suppressing alerts for newly opened
// ports until the period has elapsed. This prevents noise from transient
// connections that close on their own shortly after opening.
type Grace struct {
	mu      sync.Mutex
	window  time.Duration
	entries map[int]time.Time
}

// NewGrace returns a Grace with the given suppression window.
// Returns an error if window is zero or negative.
func NewGrace(window time.Duration) (*Grace, error) {
	if window <= 0 {
		return nil, errors.New("grace: window must be positive")
	}
	return &Grace{
		window:  window,
		entries: make(map[int]time.Time),
	}, nil
}

// Observe records the first time a port was seen. Subsequent calls for the
// same port are ignored until the entry is cleared via Clear.
func (g *Grace) Observe(port int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.entries[port]; !ok {
		g.entries[port] = time.Now()
	}
}

// Settled reports whether the grace period for port has elapsed.
// Returns false if the port has never been observed.
func (g *Grace) Settled(port int) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	t, ok := g.entries[port]
	if !ok {
		return false
	}
	return time.Since(t) >= g.window
}

// Clear removes the grace entry for port, allowing it to be re-observed.
func (g *Grace) Clear(port int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.entries, port)
}

// Len returns the number of ports currently tracked.
func (g *Grace) Len() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return len(g.entries)
}
