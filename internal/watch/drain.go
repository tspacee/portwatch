package watch

import (
	"errors"
	"sync"
	"time"
)

// Drain holds ports that are pending removal and releases them only after a
// quiet period has elapsed. This prevents transient port closures from
// triggering alerts before they are confirmed stable.
type Drain struct {
	mu      sync.Mutex
	window  time.Duration
	pending map[int]time.Time
}

// NewDrain creates a Drain with the given quiet-period window.
// Returns an error if window is zero or negative.
func NewDrain(window time.Duration) (*Drain, error) {
	if window <= 0 {
		return nil, errors.New("drain: window must be positive")
	}
	return &Drain{
		window:  window,
		pending: make(map[int]time.Time),
	}, nil
}

// Stage marks a port as pending removal, recording the current time.
// If the port is already staged its timestamp is not updated.
func (d *Drain) Stage(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("drain: port out of range")
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, ok := d.pending[port]; !ok {
		d.pending[port] = time.Now()
	}
	return nil
}

// Unstage removes a port from the pending set (e.g. it re-opened).
func (d *Drain) Unstage(port int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.pending, port)
}

// Drained returns ports whose quiet period has fully elapsed and removes them
// from the pending set.
func (d *Drain) Drained() []int {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := time.Now()
	var out []int
	for port, staged := range d.pending {
		if now.Sub(staged) >= d.window {
			out = append(out, port)
			delete(d.pending, port)
		}
	}
	return out
}

// Len returns the number of ports currently staged.
func (d *Drain) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.pending)
}
