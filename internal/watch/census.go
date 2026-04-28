package watch

import (
	"errors"
	"sync"
)

// Census tracks how many times each port has been observed open across scans.
// It provides a cumulative frequency map useful for identifying persistent
// versus transient open ports.
type Census struct {
	mu     sync.RWMutex
	counts map[int]int
	scans  int
}

// NewCensus creates an empty Census ready for recording observations.
func NewCensus() *Census {
	return &Census{
		counts: make(map[int]int),
	}
}

// Record increments the observation count for each port in the provided list
// and increments the total scan count by one.
func (c *Census) Record(ports []int) error {
	for _, p := range ports {
		if p < 1 || p > 65535 {
			return errors.New("census: port out of range")
		}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scans++
	for _, p := range ports {
		c.counts[p]++
	}
	return nil
}

// Frequency returns the number of scans in which the given port was observed open.
func (c *Census) Frequency(port int) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.counts[port]
}

// Scans returns the total number of scans recorded.
func (c *Census) Scans() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.scans
}

// Snapshot returns a copy of the current frequency map.
func (c *Census) Snapshot() map[int]int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[int]int, len(c.counts))
	for k, v := range c.counts {
		out[k] = v
	}
	return out
}

// Reset clears all recorded observations and scan counts.
func (c *Census) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts = make(map[int]int)
	c.scans = 0
}
