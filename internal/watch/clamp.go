package watch

import (
	"fmt"
	"sync"
)

// Clamp constrains port scan counts to a configured min/max range.
// Values below the minimum are raised; values above the maximum are lowered.
type Clamp struct {
	mu  sync.RWMutex
	min int
	max int
	counts map[int]int
}

// NewClamp creates a Clamp with the given inclusive [min, max] bounds.
func NewClamp(min, max int) (*Clamp, error) {
	if min < 0 {
		return nil, fmt.Errorf("clamp: min must be non-negative, got %d", min)
	}
	if max < min {
		return nil, fmt.Errorf("clamp: max (%d) must be >= min (%d)", max, min)
	}
	return &Clamp{
		min:    min,
		max:    max,
		counts: make(map[int]int),
	}, nil
}

// Set records a raw count for the given port and clamps it.
func (c *Clamp) Set(port, count int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("clamp: invalid port %d", port)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if count < c.min {
		count = c.min
	}
	if count > c.max {
		count = c.max
	}
	c.counts[port] = count
	return nil
}

// Get returns the clamped count for the given port.
// Returns 0 and false if the port has not been set.
func (c *Clamp) Get(port int) (int, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.counts[port]
	return v, ok
}

// Reset clears all stored counts.
func (c *Clamp) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts = make(map[int]int)
}

// Len returns the number of ports currently tracked.
func (c *Clamp) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.counts)
}
