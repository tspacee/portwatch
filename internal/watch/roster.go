package watch

import (
	"errors"
	"sync"
)

// Roster tracks which ports are currently considered "active" members of a
// monitored set. Ports can be enrolled or withdrawn, and the current active
// set can be queried at any time.
type Roster struct {
	mu      sync.RWMutex
	active  map[int]struct{}
}

// NewRoster returns an empty Roster.
func NewRoster() *Roster {
	return &Roster{
		active: make(map[int]struct{}),
	}
}

// Enroll adds a port to the active roster. Returns an error if the port is
// out of the valid range [1, 65535].
func (r *Roster) Enroll(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("roster: port out of range")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.active[port] = struct{}{}
	return nil
}

// Withdraw removes a port from the active roster. No-op if the port was not
// enrolled.
func (r *Roster) Withdraw(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.active, port)
}

// IsActive reports whether the given port is currently enrolled.
func (r *Roster) IsActive(port int) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.active[port]
	return ok
}

// Members returns a sorted snapshot of all currently enrolled ports.
func (r *Roster) Members() []int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]int, 0, len(r.active))
	for p := range r.active {
		out = append(out, p)
	}
	sortPorts(out)
	return out
}

// Len returns the number of enrolled ports.
func (r *Roster) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.active)
}

// sortPorts is a simple insertion sort for small slices of ports.
func sortPorts(ports []int) {
	for i := 1; i < len(ports); i++ {
		for j := i; j > 0 && ports[j] < ports[j-1]; j-- {
			ports[j], ports[j-1] = ports[j-1], ports[j]
		}
	}
}
