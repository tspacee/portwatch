package watch

import (
	"errors"
	"sync"
	"time"
)

// Latch holds the first occurrence of a port event and prevents re-triggering
// until explicitly reset. Useful for one-shot alerting on port changes.
type Latch struct {
	mu      sync.Mutex
	latched map[int]time.Time
}

// NewLatch returns an initialized Latch.
func NewLatch() *Latch {
	return &Latch{
		latched: make(map[int]time.Time),
	}
}

// Arm records the first time a port is seen. Returns true if the port was
// newly latched, false if it was already held.
func (l *Latch) Arm(port int) (bool, error) {
	if port < 1 || port > 65535 {
		return false, errors.New("latch: port out of range")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.latched[port]; ok {
		return false, nil
	}
	l.latched[port] = time.Now()
	return true, nil
}

// Reset clears the latch for a port so it can be armed again.
func (l *Latch) Reset(port int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.latched, port)
}

// IsArmed reports whether the given port is currently latched.
func (l *Latch) IsArmed(port int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	_, ok := l.latched[port]
	return ok
}

// ArmedAt returns the time a port was latched and whether it was found.
func (l *Latch) ArmedAt(port int) (time.Time, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	t, ok := l.latched[port]
	return t, ok
}

// Snapshot returns a copy of all currently latched ports and their times.
func (l *Latch) Snapshot() map[int]time.Time {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make(map[int]time.Time, len(l.latched))
	for k, v := range l.latched {
		out[k] = v
	}
	return out
}
