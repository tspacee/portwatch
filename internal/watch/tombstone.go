package watch

import (
	"fmt"
	"sync"
	"time"
)

// Tombstone tracks ports that have been permanently closed and should no
// longer trigger alerts. Entries expire after a configurable TTL.
type Tombstone struct {
	mu      sync.RWMutex
	entries map[int]time.Time
	ttl     time.Duration
}

// NewTombstone creates a Tombstone with the given TTL. TTL must be positive.
func NewTombstone(ttl time.Duration) (*Tombstone, error) {
	if ttl <= 0 {
		return nil, fmt.Errorf("tombstone: ttl must be positive, got %s", ttl)
	}
	return &Tombstone{
		entries: make(map[int]time.Time),
		ttl:     ttl,
	}, nil
}

// Bury marks a port as tombstoned at the current time.
func (t *Tombstone) Bury(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("tombstone: invalid port %d", port)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[port] = time.Now()
	return nil
}

// IsBuried reports whether the port is currently tombstoned (buried and not expired).
func (t *Tombstone) IsBuried(port int) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	at, ok := t.entries[port]
	if !ok {
		return false
	}
	return time.Since(at) < t.ttl
}

// Unbury removes a port from the tombstone, allowing it to trigger alerts again.
func (t *Tombstone) Unbury(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, port)
}

// Sweep removes all expired tombstone entries.
func (t *Tombstone) Sweep() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	removed := 0
	for port, at := range t.entries {
		if time.Since(at) >= t.ttl {
			delete(t.entries, port)
			removed++
		}
	}
	return removed
}

// Len returns the number of active (non-expired) tombstone entries.
func (t *Tombstone) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	count := 0
	for _, at := range t.entries {
		if time.Since(at) < t.ttl {
			count++
		}
	}
	return count
}
