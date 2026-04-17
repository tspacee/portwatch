package watch

import (
	"errors"
	"sync"
	"time"
)

// Cooldown suppresses repeated triggers within a fixed duration window.
// Once triggered, further calls to Ready return false until the window elapses.
type Cooldown struct {
	mu       sync.Mutex
	window   time.Duration
	lastFire time.Time
}

// ErrInvalidCooldownWindow is returned when the window is non-positive.
var ErrInvalidCooldownWindow = errors.New("cooldown: window must be greater than zero")

// NewCooldown creates a Cooldown with the given suppression window.
func NewCooldown(window time.Duration) (*Cooldown, error) {
	if window <= 0 {
		return nil, ErrInvalidCooldownWindow
	}
	return &Cooldown{window: window}, nil
}

// Ready returns true if enough time has passed since the last fire.
// If ready, it records the current time as the last fire time.
func (c *Cooldown) Ready() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	if c.lastFire.IsZero() || now.Sub(c.lastFire) >= c.window {
		c.lastFire = now
		return true
	}
	return false
}

// Reset clears the last fire time, making the cooldown immediately ready.
func (c *Cooldown) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastFire = time.Time{}
}

// Remaining returns how much time is left in the current cooldown window.
// Returns zero if the cooldown has already elapsed or was never triggered.
func (c *Cooldown) Remaining() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lastFire.IsZero() {
		return 0
	}
	elapsed := time.Since(c.lastFire)
	if elapsed >= c.window {
		return 0
	}
	return c.window - elapsed
}
