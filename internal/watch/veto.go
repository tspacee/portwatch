package watch

import (
	"errors"
	"sync"
)

// Veto tracks ports that have been explicitly blocked from alerting.
// A vetoed port is suppressed regardless of other rule evaluations.
type Veto struct {
	mu     sync.RWMutex
	vetoed map[int]string
}

// NewVeto returns an empty Veto registry.
func NewVeto() *Veto {
	return &Veto{
		vetoed: make(map[int]string),
	}
}

// Block adds a port to the veto list with a reason.
// Returns an error if the port is out of range or the reason is empty.
func (v *Veto) Block(port int, reason string) error {
	if port < 1 || port > 65535 {
		return errors.New("veto: port out of range")
	}
	if reason == "" {
		return errors.New("veto: reason must not be empty")
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	v.vetoed[port] = reason
	return nil
}

// Lift removes a port from the veto list.
// Returns false if the port was not vetoed.
func (v *Veto) Lift(port int) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	_, ok := v.vetoed[port]
	if ok {
		delete(v.vetoed, port)
	}
	return ok
}

// IsVetoed reports whether the given port is currently blocked.
func (v *Veto) IsVetoed(port int) bool {
	v.mu.RLock()
	defer v.mu.RUnlock()
	_, ok := v.vetoed[port]
	return ok
}

// Reason returns the veto reason for the given port.
// Returns an empty string if the port is not vetoed.
func (v *Veto) Reason(port int) string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.vetoed[port]
}

// Len returns the number of currently vetoed ports.
func (v *Veto) Len() int {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return len(v.vetoed)
}
