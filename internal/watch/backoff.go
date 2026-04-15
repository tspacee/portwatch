package watch

import (
	"math"
	"time"
)

// Backoff implements exponential backoff for retry delays.
type Backoff struct {
	baseDelay time.Duration
	maxDelay  time.Duration
	multiplier float64
	attempt   int
}

// NewBackoff creates a Backoff with the given base and max delay.
// Returns an error if baseDelay <= 0 or maxDelay < baseDelay.
func NewBackoff(base, max time.Duration) (*Backoff, error) {
	if base <= 0 {
		return nil, ErrInvalidWindow
	}
	if max < base {
		return nil, ErrInvalidWindow
	}
	return &Backoff{
		baseDelay:  base,
		maxDelay:   max,
		multiplier: 2.0,
	}, nil
}

// Next returns the delay for the current attempt and advances the counter.
func (b *Backoff) Next() time.Duration {
	delay := float64(b.baseDelay) * math.Pow(b.multiplier, float64(b.attempt))
	if delay > float64(b.maxDelay) {
		delay = float64(b.maxDelay)
	}
	b.attempt++
	return time.Duration(delay)
}

// Reset resets the attempt counter back to zero.
func (b *Backoff) Reset() {
	b.attempt = 0
}

// Attempt returns the current attempt count.
func (b *Backoff) Attempt() int {
	return b.attempt
}
