package watch

import (
	"errors"
	"sync"
	"time"
)

// Anchor records the first time a port was observed open and treats that
// timestamp as an immutable reference point for downstream logic.
type Anchor struct {
	mu      sync.RWMutex
	records map[int]time.Time
}

// NewAnchor returns an initialised Anchor.
func NewAnchor() *Anchor {
	return &Anchor{
		records: make(map[int]time.Time),
	}
}

// Pin records the anchor time for port if it has not been pinned before.
// Returns an error for ports outside the valid range.
func (a *Anchor) Pin(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("anchor: port out of range")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, exists := a.records[port]; !exists {
		a.records[port] = time.Now()
	}
	return nil
}

// Since returns how long ago the port was first pinned.
// If the port has never been pinned the second return value is false.
func (a *Anchor) Since(port int) (time.Duration, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	t, ok := a.records[port]
	if !ok {
		return 0, false
	}
	return time.Since(t), true
}

// Release removes the anchor record for port so it can be re-pinned.
func (a *Anchor) Release(port int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.records, port)
}

// Len returns the number of currently anchored ports.
func (a *Anchor) Len() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.records)
}
