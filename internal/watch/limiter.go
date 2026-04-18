package watch

import (
	"errors"
	"sync"
	"time"
)

// ScanLimiter enforces a minimum interval between scans and a maximum
// number of scans per rolling window.
type ScanLimiter struct {
	mu          sync.Mutex
	minInterval time.Duration
	window      time.Duration
	maxPerWindow int
	last        time.Time
	times       []time.Time
}

// NewScanLimiter creates a ScanLimiter.
// minInterval is the minimum time between any two scans.
// window and maxPerWindow define a rolling rate limit.
func NewScanLimiter(minInterval, window time.Duration, maxPerWindow int) (*ScanLimiter, error) {
	if minInterval <= 0 {
		return nil, errors.New("minInterval must be positive")
	}
	if window <= 0 {
		return nil, errors.New("window must be positive")
	}
	if maxPerWindow < 1 {
		return nil, errors.New("maxPerWindow must be at least 1")
	}
	return &ScanLimiter{
		minInterval:  minInterval,
		window:       window,
		maxPerWindow: maxPerWindow,
	}, nil
}

// Allow returns true if a scan is permitted at time now.
func (l *ScanLimiter) Allow(now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.last.IsZero() && now.Sub(l.last) < l.minInterval {
		return false
	}

	cutoff := now.Add(-l.window)
	filtered := l.times[:0]
	for _, t := range l.times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	l.times = filtered

	if len(l.times) >= l.maxPerWindow {
		return false
	}

	l.last = now
	l.times = append(l.times, now)
	return true
}

// Reset clears the limiter state.
func (l *ScanLimiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = time.Time{}
	l.times = nil
}
