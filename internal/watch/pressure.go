package watch

import (
	"errors"
	"sync"
	"time"
)

// Pressure tracks scan queue depth over time and reports whether the system
// is under load. It uses a sliding window to count pending scans.
type Pressure struct {
	mu       sync.Mutex
	window   time.Duration
	threshold int
	events   []time.Time
}

// NewPressure creates a Pressure monitor with the given sliding window and
// alert threshold. Returns an error if either value is invalid.
func NewPressure(window time.Duration, threshold int) (*Pressure, error) {
	if window <= 0 {
		return nil, errors.New("pressure: window must be positive")
	}
	if threshold < 1 {
		return nil, errors.New("pressure: threshold must be at least 1")
	}
	return &Pressure{
		window:    window,
		threshold: threshold,
	}, nil
}

// Record registers a new scan event at the current time.
func (p *Pressure) Record() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.evict()
	p.events = append(p.events, time.Now())
}

// Count returns the number of events within the current window.
func (p *Pressure) Count() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.evict()
	return len(p.events)
}

// High returns true when the event count meets or exceeds the threshold.
func (p *Pressure) High() bool {
	return p.Count() >= p.threshold
}

// evict removes events outside the sliding window. Must be called with mu held.
func (p *Pressure) evict() {
	cutoff := time.Now().Add(-p.window)
	i := 0
	for i < len(p.events) && p.events[i].Before(cutoff) {
		i++
	}
	p.events = p.events[i:]
}
