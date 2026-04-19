package watch

import (
	"errors"
	"sync"
	"time"
)

// Clock tracks scan timing metadata: when scans started, ended, and elapsed duration.
type Clock struct {
	mu       sync.Mutex
	started  time.Time
	finished time.Time
	elapsed  time.Duration
	count    int64
}

// NewClock returns an initialized Clock.
func NewClock() *Clock {
	return &Clock{}
}

// Start records the beginning of a scan cycle.
func (c *Clock) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.started.IsZero() && c.finished.IsZero() {
		return errors.New("clock: scan already in progress")
	}
	c.started = time.Now()
	c.finished = time.Time{}
	return nil
}

// Stop records the end of a scan cycle and computes elapsed time.
func (c *Clock) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.started.IsZero() {
		return errors.New("clock: no scan in progress")
	}
	c.finished = time.Now()
	c.elapsed = c.finished.Sub(c.started)
	c.count++
	return nil
}

// Elapsed returns the duration of the most recent completed scan.
func (c *Clock) Elapsed() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.elapsed
}

// LastScan returns the time the most recent scan finished.
func (c *Clock) LastScan() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.finished
}

// Count returns the total number of completed scans.
func (c *Clock) Count() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}
