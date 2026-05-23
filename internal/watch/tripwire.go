package watch

import (
	"errors"
	"sync"
)

// Tripwire monitors a set of ports and fires a callback the first time
// any of them disappears from an observed port set. Once tripped, it
// remains in the tripped state until explicitly reset.
type Tripwire struct {
	mu      sync.Mutex
	ports   map[int]struct{}
	tripped bool
	cb      func(port int)
}

// NewTripwire creates a Tripwire that watches the given ports.
// cb is invoked once per trip event with the port that caused the trip.
func NewTripwire(ports []int, cb func(port int)) (*Tripwire, error) {
	if len(ports) == 0 {
		return nil, errors.New("tripwire: at least one port required")
	}
	if cb == nil {
		return nil, errors.New("tripwire: callback must not be nil")
	}
	set := make(map[int]struct{}, len(ports))
	for _, p := range ports {
		if p < 1 || p > 65535 {
			return nil, errors.New("tripwire: port out of range")
		}
		set[p] = struct{}{}
	}
	return &Tripwire{ports: set, cb: cb}, nil
}

// Observe checks the provided open ports against the watched set.
// If any watched port is absent and the tripwire has not already fired,
// the callback is invoked and the tripwire is marked as tripped.
func (t *Tripwire) Observe(open []int) {
	present := make(map[int]struct{}, len(open))
	for _, p := range open {
		present[p] = struct{}{}
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.tripped {
		return
	}
	for p := range t.ports {
		if _, ok := present[p]; !ok {
			t.tripped = true
			t.cb(p)
			return
		}
	}
}

// Tripped reports whether the tripwire has been triggered.
func (t *Tripwire) Tripped() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.tripped
}

// Reset clears the tripped state so the tripwire can fire again.
func (t *Tripwire) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.tripped = false
}
