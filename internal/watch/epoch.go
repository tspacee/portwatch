package watch

import (
	"errors"
	"sync"
	"time"
)

// Epoch tracks a monotonically incrementing scan generation counter.
// Each time a new scan cycle begins, the epoch is advanced. Components
// can use the epoch number to correlate observations across a single scan.
type Epoch struct {
	mu      sync.RWMutex
	current uint64
	started time.Time
	times   []time.Time
}

// NewEpoch returns a new Epoch starting at generation zero.
func NewEpoch() *Epoch {
	return &Epoch{
		times: make([]time.Time, 0),
	}
}

// Advance increments the epoch counter and records the timestamp.
// Returns the new epoch number.
func (e *Epoch) Advance() uint64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.current++
	now := time.Now()
	if e.current == 1 {
		e.started = now
	}
	e.times = append(e.times, now)
	return e.current
}

// Current returns the current epoch number.
func (e *Epoch) Current() uint64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.current
}

// Since returns the time elapsed since the given epoch was recorded.
// Returns an error if the epoch number is out of range.
func (e *Epoch) Since(epoch uint64) (time.Duration, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if epoch == 0 || epoch > uint64(len(e.times)) {
		return 0, errors.New("epoch out of range")
	}
	return time.Since(e.times[epoch-1]), nil
}

// Reset clears all epoch history and resets the counter to zero.
func (e *Epoch) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.current = 0
	e.started = time.Time{}
	e.times = make([]time.Time, 0)
}

// Age returns the duration since the first epoch was recorded.
// Returns zero if no epochs have been advanced.
func (e *Epoch) Age() time.Duration {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.started.IsZero() {
		return 0
	}
	return time.Since(e.started)
}
