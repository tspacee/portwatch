package watch

import (
	"errors"
	"sync"
	"time"
)

// Poll tracks per-port polling intervals, allowing each port to have an
// independent cadence rather than a single global scan frequency.
type Poll struct {
	mu       sync.Mutex
	default_ time.Duration
	intervals map[int]time.Duration
	lastPoll  map[int]time.Time
}

// NewPoll creates a Poll with the given default interval.
// Returns an error if the default interval is zero or negative.
func NewPoll(defaultInterval time.Duration) (*Poll, error) {
	if defaultInterval <= 0 {
		return nil, errors.New("poll: default interval must be positive")
	}
	return &Poll{
		default_:  defaultInterval,
		intervals: make(map[int]time.Duration),
		lastPoll:  make(map[int]time.Time),
	}, nil
}

// Set assigns a custom polling interval for the given port.
// Returns an error if the port is out of range or the interval is non-positive.
func (p *Poll) Set(port int, interval time.Duration) error {
	if port < 1 || port > 65535 {
		return errors.New("poll: port out of range")
	}
	if interval <= 0 {
		return errors.New("poll: interval must be positive")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.intervals[port] = interval
	return nil
}

// Due reports whether the given port is due for a poll based on its interval.
// A port that has never been polled is always due.
func (p *Poll) Due(port int) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	interval, ok := p.intervals[port]
	if !ok {
		interval = p.default_
	}
	last, seen := p.lastPoll[port]
	if !seen {
		return true
	}
	return time.Since(last) >= interval
}

// Mark records the current time as the last poll time for the given port.
func (p *Poll) Mark(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("poll: port out of range")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.lastPoll[port] = time.Now()
	return nil
}

// Interval returns the effective polling interval for the given port.
func (p *Poll) Interval(port int) time.Duration {
	p.mu.Lock()
	defer p.mu.Unlock()
	if d, ok := p.intervals[port]; ok {
		return d
	}
	return p.default_
}
