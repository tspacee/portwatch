package watch

import (
	"errors"
	"sync"
)

// Pivot tracks the last known "pivot point" for a port — the scan index at
// which its open/closed state last changed. This is useful for understanding
// how long a port has been in its current state.
type Pivot struct {
	mu      sync.RWMutex
	entries map[int]uint64 // port -> scan index of last state change
	states  map[int]bool   // port -> last known state (true = open)
}

// NewPivot returns an initialised Pivot.
func NewPivot() *Pivot {
	return &Pivot{
		entries: make(map[int]uint64),
		states:  make(map[int]bool),
	}
}

// Observe records the current state of a port at the given scan index.
// If the state has changed since the last observation the pivot index is
// updated. Returns an error for invalid ports.
func (p *Pivot) Observe(port int, open bool, scanIndex uint64) error {
	if port < 1 || port > 65535 {
		return errors.New("pivot: port out of range")
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	prev, seen := p.states[port]
	if !seen || prev != open {
		p.entries[port] = scanIndex
		p.states[port] = open
	}
	return nil
}

// Since returns the scan index at which the port last changed state.
// The second return value is false if the port has never been observed.
func (p *Pivot) Since(port int) (uint64, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	idx, ok := p.entries[port]
	return idx, ok
}

// Reset removes all recorded pivot data for a port.
func (p *Pivot) Reset(port int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.entries, port)
	delete(p.states, port)
}

// Len returns the number of ports currently tracked.
func (p *Pivot) Len() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.entries)
}
