package watch

import (
	"errors"
	"sync"
	"time"
)

// Spike detects sudden bursts of port activity within a short observation
// window. It is useful for identifying scan events or rapid service churn.
type Spike struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	events    []time.Time
}

// NewSpike creates a Spike detector with the given observation window and
// activation threshold. Returns an error if either parameter is invalid.
func NewSpike(window time.Duration, threshold int) (*Spike, error) {
	if window <= 0 {
		return nil, errors.New("spike: window must be positive")
	}
	if threshold < 1 {
		return nil, errors.New("spike: threshold must be at least 1")
	}
	return &Spike{
		window:    window,
		threshold: threshold,
		events:    make([]time.Time, 0),
	}, nil
}

// Record registers a port-change event at the current time.
func (s *Spike) Record() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prune(time.Now())
	s.events = append(s.events, time.Now())
}

// Triggered reports whether the number of events within the observation
// window meets or exceeds the configured threshold.
func (s *Spike) Triggered() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prune(time.Now())
	return len(s.events) >= s.threshold
}

// Count returns the number of events currently within the observation window.
func (s *Spike) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prune(time.Now())
	return len(s.events)
}

// Reset clears all recorded events.
func (s *Spike) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = s.events[:0]
}

// prune removes events older than the observation window. Must be called
// with s.mu held.
func (s *Spike) prune(now time.Time) {
	cutoff := now.Add(-s.window)
	i := 0
	for i < len(s.events) && s.events[i].Before(cutoff) {
		i++
	}
	s.events = s.events[i:]
}
