package watch

import (
	"fmt"
	"sync"
	"time"
)

// RateLimiter controls how frequently scan cycles can be triggered,
// preventing runaway scanning under high alert conditions.
type RateLimiter struct {
	mu          sync.Mutex
	minInterval time.Duration
	lastAllowed time.Time
	skipped     int
}

// ErrRateLimited is returned when a scan is rejected due to rate limiting.
type ErrRateLimited struct {
	NextAllowed time.Time
}

func (e *ErrRateLimited) Error() string {
	return fmt.Sprintf("rate limited: next scan allowed at %s", e.NextAllowed.Format(time.RFC3339))
}

// NewRateLimiter creates a RateLimiter that enforces a minimum interval between scans.
// minInterval must be positive.
func NewRateLimiter(minInterval time.Duration) (*RateLimiter, error) {
	if minInterval <= 0 {
		return nil, fmt.Errorf("minInterval must be positive, got %s", minInterval)
	}
	return &RateLimiter{
		minInterval: minInterval,
	}, nil
}

// Allow returns nil if a scan may proceed, or ErrRateLimited if the minimum
// interval has not yet elapsed since the last allowed scan.
func (r *RateLimiter) Allow() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	next := r.lastAllowed.Add(r.minInterval)

	if !r.lastAllowed.IsZero() && now.Before(next) {
		r.skipped++
		return &ErrRateLimited{NextAllowed: next}
	}

	r.lastAllowed = now
	return nil
}

// Skipped returns the total number of scans that have been rate-limited.
func (r *RateLimiter) Skipped() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.skipped
}

// Reset clears the last-allowed timestamp, permitting the next scan immediately.
func (r *RateLimiter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastAllowed = time.Time{}
}
