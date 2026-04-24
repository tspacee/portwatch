package watch

import (
	"errors"
	"sync"
	"time"
)

// HoldDown suppresses repeated alerts for a port within a configurable
// quiet period. Once a port has triggered, further triggers are ignored
// until the hold-down window expires.
type HoldDown struct {
	mu     sync.Mutex
	window time.Duration
	held   map[int]time.Time
}

// NewHoldDown creates a HoldDown with the given quiet window.
// Returns an error if window is zero or negative.
func NewHoldDown(window time.Duration) (*HoldDown, error) {
	if window <= 0 {
		return nil, errors.New("holddown: window must be positive")
	}
	return &HoldDown{
		window: window,
		held:   make(map[int]time.Time),
	}, nil
}

// Suppressed reports whether the given port is currently held down.
// If it is not held down, it arms the hold-down and returns false.
// If it is held down and the window has not expired, it returns true.
// If the window has expired, it re-arms and returns false.
func (h *HoldDown) Suppressed(port int) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now()
	if until, ok := h.held[port]; ok {
		if now.Before(until) {
			return true
		}
	}
	h.held[port] = now.Add(h.window)
	return false
}

// Release removes the hold-down for a port immediately, allowing the
// next trigger to pass through without waiting for the window to expire.
func (h *HoldDown) Release(port int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.held, port)
}

// Len returns the number of ports currently held down.
func (h *HoldDown) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	count := 0
	now := time.Now()
	for _, until := range h.held {
		if now.Before(until) {
			count++
		}
	}
	return count
}
