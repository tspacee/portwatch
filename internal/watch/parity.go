package watch

import (
	"errors"
	"sync"
)

// Parity tracks whether the current set of open ports matches a known-good
// reference set, reporting the degree of divergence as a count of mismatches.
type Parity struct {
	mu        sync.RWMutex
	reference map[int]struct{}
	mismatches int
}

// NewParity returns a Parity instance with an empty reference set.
func NewParity() *Parity {
	return &Parity{
		reference: make(map[int]struct{}),
	}
}

// SetReference replaces the reference port set used for comparison.
// Ports must be in the range [1, 65535].
func (p *Parity) SetReference(ports []int) error {
	for _, port := range ports {
		if port < 1 || port > 65535 {
			return errors.New("parity: port out of range")
		}
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.reference = make(map[int]struct{}, len(ports))
	for _, port := range ports {
		p.reference[port] = struct{}{}
	}
	return nil
}

// Compare evaluates the provided ports against the reference set and stores
// the mismatch count (ports only in one of the two sets).
func (p *Parity) Compare(ports []int) {
	current := make(map[int]struct{}, len(ports))
	for _, port := range ports {
		current[port] = struct{}{}
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	count := 0
	for port := range p.reference {
		if _, ok := current[port]; !ok {
			count++
		}
	}
	for port := range current {
		if _, ok := p.reference[port]; !ok {
			count++
		}
	}
	p.mismatches = count
}

// Mismatches returns the number of ports that differ between the last
// compared set and the reference set.
func (p *Parity) Mismatches() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.mismatches
}

// InSync reports whether the last compared port set exactly matches the
// reference set (zero mismatches).
func (p *Parity) InSync() bool {
	return p.Mismatches() == 0
}
