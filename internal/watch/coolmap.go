package watch

import (
	"errors"
	"sync"
	"time"
)

// CoolMap tracks per-port cooldown periods, preventing repeated actions
// on the same port within a configurable window.
type CoolMap struct {
	mu     sync.Mutex
	window time.Duration
	entries map[int]time.Time
}

// NewCoolMap creates a CoolMap with the given cooldown window.
// Returns an error if window is zero or negative.
func NewCoolMap(window time.Duration) (*CoolMap, error) {
	if window <= 0 {
		return nil, errors.New("coolmap: window must be positive")
	}
	return &CoolMap{
		window:  window,
		entries: make(map[int]time.Time),
	}, nil
}

// Ready returns true if the port is not currently in cooldown.
// Calling Ready with a port that is not cooling down marks it as active.
func (c *CoolMap) Ready(port int) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	if t, ok := c.entries[port]; ok && now.Before(t) {
		return false
	}
	c.entries[port] = now.Add(c.window)
	return true
}

// Reset clears the cooldown for a specific port.
func (c *CoolMap) Reset(port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, port)
}

// Len returns the number of ports currently tracked.
func (c *CoolMap) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.entries)
}

// Purge removes all expired entries from the map.
func (c *CoolMap) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for port, exp := range c.entries {
		if now.After(exp) {
			delete(c.entries, port)
		}
	}
}
