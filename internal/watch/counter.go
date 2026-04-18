package watch

import (
	"errors"
	"sync"
)

// Counter tracks how many times each port has been seen open across scans.
// It is safe for concurrent use.
type Counter struct {
	mu     sync.Mutex
	counts map[int]int
}

// NewCounter returns an initialized Counter.
func NewCounter() *Counter {
	return &Counter{
		counts: make(map[int]int),
	}
}

// Increment adds one to the count for the given port.
// Returns an error if the port is out of range.
func (c *Counter) Increment(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("counter: port out of range")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[port]++
	return nil
}

// Get returns the current count for the given port.
func (c *Counter) Get(port int) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counts[port]
}

// Reset clears the count for the given port.
func (c *Counter) Reset(port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.counts, port)
}

// Snapshot returns a copy of all current counts.
func (c *Counter) Snapshot() map[int]int {
	c.mu.Lock()
	defer c.mu.Unlock()
	copy := make(map[int]int, len(c.counts))
	for k, v := range c.counts {
		copy[k] = v
	}
	return copy
}
