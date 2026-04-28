package watch

import (
	"errors"
	"sync"
	"time"
)

// Flap detects ports that rapidly alternate between open and closed states.
// A port is considered flapping when it changes state more than Threshold
// times within the configured Window.
type Flap struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	events    map[int][]time.Time
}

// NewFlap creates a Flap detector with the given window and threshold.
// threshold must be >= 2; window must be > 0.
func NewFlap(window time.Duration, threshold int) (*Flap, error) {
	if window <= 0 {
		return nil, errors.New("flap: window must be positive")
	}
	if threshold < 2 {
		return nil, errors.New("flap: threshold must be at least 2")
	}
	return &Flap{
		window:    window,
		threshold: threshold,
		events:    make(map[int][]time.Time),
	}, nil
}

// Record registers a state-change event for the given port at now.
func (f *Flap) Record(port int, now time.Time) error {
	if port < 1 || port > 65535 {
		return errors.New("flap: port out of range")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	cutoff := now.Add(-f.window)
	filtered := f.evicts(f.events[port], cutoff)
	filtered = append(filtered, now)
	f.events[port] = filtered
	return nil
}

// IsFlapping reports whether the port has exceeded the change threshold
// within the configured window relative to now.
func (f *Flap) IsFlapping(port int, now time.Time) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	cutoff := now.Add(-f.window)
	valid := f.evicts(f.events[port], cutoff)
	f.events[port] = valid
	return len(valid) >= f.threshold
}

// Reset clears all recorded events for the given port.
func (f *Flap) Reset(port int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.events, port)
}

func (f *Flap) evicts(ts []time.Time, cutoff time.Time) []time.Time {
	out := ts[:0]
	for _, t := range ts {
		if t.After(cutoff) {
			out = append(out, t)
		}
	}
	return out
}
