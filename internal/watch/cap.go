package watch

import (
	"errors"
	"sync"
)

// Cap enforces an upper bound on the number of concurrently tracked ports.
// Once the cap is reached, new ports are rejected until existing ones are released.
type Cap struct {
	mu      sync.Mutex
	max     int
	active  map[int]struct{}
}

// NewCap creates a Cap with the given maximum concurrent port limit.
func NewCap(max int) (*Cap, error) {
	if max < 1 {
		return nil, errors.New("cap: max must be at least 1")
	}
	return &Cap{
		max:    max,
		active: make(map[int]struct{}),
	}, nil
}

// Acquire attempts to track port. Returns true if within cap, false if cap exceeded.
func (c *Cap) Acquire(port int) (bool, error) {
	if port < 1 || port > 65535 {
		return false, errors.New("cap: port out of range")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.active[port]; ok {
		return true, nil
	}
	if len(c.active) >= c.max {
		return false, nil
	}
	c.active[port] = struct{}{}
	return true, nil
}

// Release removes port from the active set.
func (c *Cap) Release(port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.active, port)
}

// Len returns the number of currently active ports.
func (c *Cap) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.active)
}

// Full returns true when the active set has reached the maximum.
func (c *Cap) Full() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.active) >= c.max
}
