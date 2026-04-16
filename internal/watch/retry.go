package watch

import (
	"errors"
	"time"
)

// ErrMaxRetries is returned when all retry attempts are exhausted.
var ErrMaxRetries = errors.New("max retries exceeded")

// RetryPolicy defines how retries are attempted.
type RetryPolicy struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64
}

// Retry executes fn up to MaxAttempts times, backing off between attempts.
// It returns nil on first success, or ErrMaxRetries if all attempts fail.
func (p RetryPolicy) Retry(fn func() error) error {
	if p.MaxAttempts <= 0 {
		return ErrMaxRetries
	}
	delay := p.Delay
	for attempt := 0; attempt < p.MaxAttempts; attempt++ {
		if err := fn(); err == nil {
			return nil
		}
		if attempt < p.MaxAttempts-1 && delay > 0 {
			time.Sleep(delay)
			if p.Multiplier > 1 {
				delay = time.Duration(float64(delay) * p.Multiplier)
			}
		}
	}
	return ErrMaxRetries
}

// NewRetryPolicy returns a RetryPolicy with sensible defaults.
func NewRetryPolicy(maxAttempts int, delay time.Duration, multiplier float64) (RetryPolicy, error) {
	if maxAttempts <= 0 {
		return RetryPolicy{}, errors.New("maxAttempts must be greater than zero")
	}
	if multiplier != 0 && multiplier < 1 {
		return RetryPolicy{}, errors.New("multiplier must be >= 1 or zero")
	}
	return RetryPolicy{
		MaxAttempts: maxAttempts,
		Delay:       delay,
		Multiplier:  multiplier,
	}, nil
}
