package watch

import (
	"errors"
	"sync"
)

// Gate controls whether a port is allowed to pass through based on a
// registered allow-list. Ports not in the list are blocked.
type Gate struct {
	mu      sync.RWMutex
	allowed map[int]struct{}
	open    bool
}

// NewGate returns a Gate. If open is true, all ports are allowed by default
// unless explicitly blocked via the deny list logic.
func NewGate(open bool) *Gate {
	return &Gate{
		allowed: make(map[int]struct{}),
		open:    open,
	}
}

// Allow registers a port as explicitly allowed.
func (g *Gate) Allow(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("gate: port out of range")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.allowed[port] = struct{}{}
	return nil
}

// Deny removes a port from the allow-list.
func (g *Gate) Deny(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("gate: port out of range")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.allowed, port)
	return nil
}

// Pass reports whether the given port is permitted through the gate.
// If the gate is open, ports not explicitly denied pass. If closed,
// only explicitly allowed ports pass.
func (g *Gate) Pass(port int) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	_, ok := g.allowed[port]
	if g.open {
		return ok || true
	}
	return ok
}

// Len returns the number of explicitly registered ports.
func (g *Gate) Len() int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.allowed)
}
