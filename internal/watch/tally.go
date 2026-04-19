package watch

import (
	"errors"
	"sync"
)

// Tally tracks cumulative counts per port across scan cycles.
type Tally struct {
	mu     sync.Mutex
	counts map[int]int
}

// NewTally returns an initialised Tally.
func NewTally() *Tally {
	return &Tally{counts: make(map[int]int)}
}

// Inc increments the count for the given port by one.
func (t *Tally) Inc(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("tally: port out of range")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.counts[port]++
	return nil
}

// Get returns the current count for the given port.
func (t *Tally) Get(port int) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.counts[port]
}

// Reset clears the count for the given port.
func (t *Tally) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.counts, port)
}

// Snapshot returns a copy of all current counts.
func (t *Tally) Snapshot() map[int]int {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make(map[int]int, len(t.counts))
	for k, v := range t.counts {
		out[k] = v
	}
	return out
}

// Top returns the port with the highest count, or 0 if empty.
func (t *Tally) Top() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	var topPort, topCount int
	for port, count := range t.counts {
		if count > topCount {
			topCount = count
			topPort = port
		}
	}
	return topPort
}
