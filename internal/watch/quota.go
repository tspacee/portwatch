package watch

import (
	"errors"
	"sync"
	"time"
)

// Quota enforces a maximum number of scan cycles within a rolling time window.
// Once the quota is exhausted, Allow returns false until the window resets.
type Quota struct {
	mu       sync.Mutex
	window   time.Duration
	max      int
	timestamps []time.Time
}

// NewQuota creates a Quota with the given rolling window and max allowed scans.
func NewQuota(window time.Duration, max int) (*Quota, error) {
	if window <= 0 {
		return nil, errors.New("quota: window must be positive")
	}
	if max <= 0 {
		return nil, errors.New("quota: max must be positive")
	}
	return &Quota{window: window, max: max}, nil
}

// Allow returns true if a scan is permitted under the current quota.
func (q *Quota) Allow() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	now := time.Now()
	q.evict(now)
	if len(q.timestamps) >= q.max {
		return false
	}
	q.timestamps = append(q.timestamps, now)
	return true
}

// Remaining returns the number of scans still allowed in the current window.
func (q *Quota) Remaining() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.evict(time.Now())
	r := q.max - len(q.timestamps)
	if r < 0 {
		return 0
	}
	return r
}

// Reset clears all recorded timestamps.
func (q *Quota) Reset() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.timestamps = q.timestamps[:0]
}

func (q *Quota) evict(now time.Time) {
	cutoff := now.Add(-q.window)
	i := 0
	for i < len(q.timestamps) && q.timestamps[i].Before(cutoff) {
		i++
	}
	q.timestamps = q.timestamps[i:]
}
