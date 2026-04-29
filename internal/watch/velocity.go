package watch

import (
	"errors"
	"sync"
	"time"
)

// Velocity tracks the rate of change for observed ports over a sliding window.
// It measures how frequently a port's state changes per unit time.
type Velocity struct {
	mu      sync.Mutex
	window  time.Duration
	events  map[int][]time.Time
}

// NewVelocity creates a Velocity tracker with the given sliding window duration.
func NewVelocity(window time.Duration) (*Velocity, error) {
	if window <= 0 {
		return nil, errors.New("velocity: window must be positive")
	}
	return &Velocity{
		window: window,
		events: make(map[int][]time.Time),
	}, nil
}

// Record registers a state-change event for the given port.
func (v *Velocity) Record(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("velocity: port out of range")
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	now := time.Now()
	v.evict(port, now)
	v.events[port] = append(v.events[port], now)
	return nil
}

// Rate returns the number of recorded events for a port within the current window.
func (v *Velocity) Rate(port int) int {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.evict(port, time.Now())
	return len(v.events[port])
}

// Reset clears all recorded events for a port.
func (v *Velocity) Reset(port int) {
	v.mu.Lock()
	defer v.mu.Unlock()
	delete(v.events, port)
}

// evict removes events outside the current window. Must be called with mu held.
func (v *Velocity) evict(port int, now time.Time) {
	cutoff := now.Add(-v.window)
	evts := v.events[port]
	i := 0
	for i < len(evts) && evts[i].Before(cutoff) {
		i++
	}
	v.events[port] = evts[i:]
}
