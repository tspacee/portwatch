package watch

import (
	"errors"
	"sync"
	"time"
)

// Canary tracks a set of ports that must remain open. If any canary port
// is found closed, the canary is considered tripped.
type Canary struct {
	mu      sync.Mutex
	ports   map[int]struct{}
	tripped map[int]time.Time
}

// NewCanary creates a Canary with no watched ports.
func NewCanary() *Canary {
	return &Canary{
		ports:   make(map[int]struct{}),
		tripped: make(map[int]time.Time),
	}
}

// Watch registers a port as a canary port.
func (c *Canary) Watch(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("canary: port out of range")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ports[port] = struct{}{}
	return nil
}

// Observe checks the given open ports against watched ports.
// Any watched port absent from open is recorded as tripped.
func (c *Canary) Observe(open []int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	openSet := make(map[int]struct{}, len(open))
	for _, p := range open {
		openSet[p] = struct{}{}
	}
	now := time.Now()
	for p := range c.ports {
		if _, ok := openSet[p]; !ok {
			if _, already := c.tripped[p]; !already {
				c.tripped[p] = now
			}
		} else {
			delete(c.tripped, p)
		}
	}
}

// Tripped returns all ports currently tripped along with the time they
// were first detected as closed.
func (c *Canary) Tripped() map[int]time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[int]time.Time, len(c.tripped))
	for k, v := range c.tripped {
		out[k] = v
	}
	return out
}

// Len returns the number of watched ports.
func (c *Canary) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.ports)
}
