package watch

import (
	"errors"
	"sync"
)

// Baseline tracks the expected set of open ports established during a learning
// period. Once frozen, it can be used to detect deviations from the norm.
type Baseline struct {
	mu     sync.RWMutex
	ports  map[int]struct{}
	frozen bool
}

// NewBaseline returns an empty, unfrozen Baseline.
func NewBaseline() *Baseline {
	return &Baseline{
		ports: make(map[int]struct{}),
	}
}

// Learn adds a port to the baseline. Returns an error if the baseline is
// already frozen or if the port number is out of range.
func (b *Baseline) Learn(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("baseline: port out of range")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.frozen {
		return errors.New("baseline: already frozen")
	}
	b.ports[port] = struct{}{}
	return nil
}

// Freeze locks the baseline, preventing further learning.
func (b *Baseline) Freeze() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.frozen = true
}

// IsFrozen reports whether the baseline has been frozen.
func (b *Baseline) IsFrozen() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.frozen
}

// Contains reports whether port is part of the established baseline.
func (b *Baseline) Contains(port int) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.ports[port]
	return ok
}

// Snapshot returns a copy of the current baseline port set.
func (b *Baseline) Snapshot() []int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]int, 0, len(b.ports))
	for p := range b.ports {
		out = append(out, p)
	}
	return out
}

// Len returns the number of ports in the baseline.
func (b *Baseline) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.ports)
}
