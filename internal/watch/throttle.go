package watch

import (
	"sync"
	"time"
)

// Throttle limits how frequently an action can be triggered within a
// rolling time window. It is safe for concurrent use.
type Throttle struct {
	mu       sync.Mutex
	window   time.Duration
	maxCount int
	events   []time.Time
}

// NewThrottle creates a Throttle that allows at most maxCount events
// within the given window duration.
// Returns an error if maxCount < 1 or window <= 0.
func NewThrottle(window time.Duration, maxCount int) (*Throttle, error) {
	if window <= 0 {
		return nil, ErrInvalidWindow
	}
	if maxCount < 1 {
		return nil, ErrInvalidMaxCount
	}
	return &Throttle{
		window:   window,
		maxCount: maxCount,
		events:   make([]time.Time, 0, maxCount),
	}, nil
}

// Allow reports whether a new event should be permitted right now.
// If allowed, the event timestamp is recorded internally.
func (t *Throttle) Allow() bool {
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()

	// Evict events outside the window.
	cutoff := now.Add(-t.window)
	valid := t.events[:0]
	for _, ts := range t.events {
		if ts.After(cutoff) {
			valid = append(valid, ts)
		}
	}
	t.events = valid

	if len(t.events) >= t.maxCount {
		return false
	}
	t.events = append(t.events, now)
	return true
}

// Reset clears all recorded events.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = t.events[:0]
}

// Count returns the number of events currently within the window.
func (t *Throttle) Count() int {
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()
	cutoff := now.Add(-t.window)
	count := 0
	for _, ts := range t.events {
		if ts.After(cutoff) {
			count++
		}
	}
	return count
}
