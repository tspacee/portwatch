package watch

import (
	"errors"
	"sync"
	"time"
)

// Burst tracks the number of events occurring within a short window and
// determines whether the rate constitutes a burst above a given threshold.
type Burst struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	events    []time.Time
}

// NewBurst creates a Burst detector with the given time window and event threshold.
// Returns an error if window is zero or threshold is less than 1.
func NewBurst(window time.Duration, threshold int) (*Burst, error) {
	if window <= 0 {
		return nil, errors.New("burst: window must be greater than zero")
	}
	if threshold < 1 {
		return nil, errors.New("burst: threshold must be at least 1")
	}
	return &Burst{
		window:    window,
		threshold: threshold,
		events:    make([]time.Time, 0),
	}, nil
}

// Record registers a new event at the current time.
func (b *Burst) Record() {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	b.evict(now)
	b.events = append(b.events, now)
}

// IsBursting returns true if the number of events within the window meets
// or exceeds the configured threshold.
func (b *Burst) IsBursting() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.evict(time.Now())
	return len(b.events) >= b.threshold
}

// Count returns the number of events currently within the active window.
func (b *Burst) Count() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.evict(time.Now())
	return len(b.events)
}

// Reset clears all recorded events.
func (b *Burst) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.events = b.events[:0]
}

func (b *Burst) evict(now time.Time) {
	cutoff := now.Add(-b.window)
	i := 0
	for i < len(b.events) && b.events[i].Before(cutoff) {
		i++
	}
	b.events = b.events[i:]
}
