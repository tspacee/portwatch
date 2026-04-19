package watch

import (
	"errors"
	"sync"
	"time"
)

// Trigger fires a callback when a condition is met a minimum number of times
// within a sliding time window.
type Trigger struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	events    []time.Time
	callback  func()
}

// NewTrigger creates a Trigger that fires callback after threshold events within window.
func NewTrigger(window time.Duration, threshold int, callback func()) (*Trigger, error) {
	if window <= 0 {
		return nil, errors.New("trigger: window must be positive")
	}
	if threshold < 1 {
		return nil, errors.New("trigger: threshold must be at least 1")
	}
	if callback == nil {
		return nil, errors.New("trigger: callback must not be nil")
	}
	return &Trigger{
		window:    window,
		threshold: threshold,
		callback:  callback,
	}, nil
}

// Record registers an event occurrence and fires the callback if the threshold is met.
func (t *Trigger) Record() {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-t.window)

	filtered := t.events[:0]
	for _, e := range t.events {
		if e.After(cutoff) {
			filtered = append(filtered, e)
		}
	}
	t.events = append(filtered, now)

	if len(t.events) >= t.threshold {
		t.events = nil
		go t.callback()
	}
}

// Count returns the number of events within the current window.
func (t *Trigger) Count() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-t.window)
	count := 0
	for _, e := range t.events {
		if e.After(cutoff) {
			count++
		}
	}
	return count
}

// Reset clears all recorded events.
func (t *Trigger) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = nil
}
