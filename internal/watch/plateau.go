package watch

import (
	"errors"
	"sync"
	"time"
)

// Plateau detects when a port's open/closed state has remained unchanged
// for a minimum duration, indicating a stable condition.
type Plateau struct {
	mu       sync.Mutex
	minStable time.Duration
	first    map[int]time.Time
	last     map[int]bool
}

// NewPlateau creates a Plateau tracker requiring the given minimum stable duration.
// Returns an error if minStable is zero or negative.
func NewPlateau(minStable time.Duration) (*Plateau, error) {
	if minStable <= 0 {
		return nil, errors.New("plateau: minStable must be positive")
	}
	return &Plateau{
		minStable: minStable,
		first:     make(map[int]time.Time),
		last:      make(map[int]bool),
	}, nil
}

// Observe records the current state of a port. Returns an error for invalid ports.
func (p *Plateau) Observe(port int, open bool) error {
	if port < 1 || port > 65535 {
		return errors.New("plateau: port out of range")
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	prev, seen := p.last[port]
	if !seen || prev != open {
		p.first[port] = time.Now()
		p.last[port] = open
	}
	return nil
}

// Stable reports whether the given port has held its current state for at
// least the configured minStable duration. Returns false for unseen ports.
func (p *Plateau) Stable(port int) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	t, ok := p.first[port]
	if !ok {
		return false
	}
	return time.Since(t) >= p.minStable
}

// Reset clears all tracking state for all ports.
func (p *Plateau) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.first = make(map[int]time.Time)
	p.last = make(map[int]bool)
}

// Len returns the number of ports currently being tracked.
func (p *Plateau) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.last)
}
