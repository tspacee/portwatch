package watch

import (
	"errors"
	"sync"
)

// Norm tracks the normal (expected) set of open ports and detects deviations.
// It maintains a learned baseline and reports whether an observed set of ports
// conforms to the established norm.
type Norm struct {
	mu       sync.RWMutex
	baseline map[int]struct{}
	frozen   bool
}

// NewNorm creates a new Norm with an empty baseline.
func NewNorm() *Norm {
	return &Norm{
		baseline: make(map[int]struct{}),
	}
}

// Learn adds a port to the baseline. Returns an error if the norm is frozen
// or the port is out of range.
func (n *Norm) Learn(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("norm: port out of range")
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.frozen {
		return errors.New("norm: baseline is frozen")
	}
	n.baseline[port] = struct{}{}
	return nil
}

// Freeze locks the baseline, preventing further learning.
func (n *Norm) Freeze() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.frozen = true
}

// IsFrozen reports whether the baseline has been frozen.
func (n *Norm) IsFrozen() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.frozen
}

// Conforms reports whether all ports in the observed slice are within the
// baseline. Returns false if any port is unexpected.
func (n *Norm) Conforms(ports []int) bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	for _, p := range ports {
		if _, ok := n.baseline[p]; !ok {
			return false
		}
	}
	return true
}

// Deviations returns the ports in observed that are not in the baseline.
func (n *Norm) Deviations(ports []int) []int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	var out []int
	for _, p := range ports {
		if _, ok := n.baseline[p]; !ok {
			out = append(out, p)
		}
	}
	return out
}

// Size returns the number of ports in the baseline.
func (n *Norm) Size() int {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return len(n.baseline)
}
