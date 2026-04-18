package watch

import (
	"errors"
	"math/rand"
	"time"
)

// Jitter adds randomised delay to scan intervals to avoid thundering-herd
// problems when multiple portwatch instances run in parallel.
type Jitter struct {
	base   time.Duration
	factor float64
	rng    *rand.Rand
}

// NewJitter returns a Jitter that spreads intervals by up to factor (0–1) of
// the base duration. E.g. base=10s, factor=0.2 yields delays in [10s, 12s).
func NewJitter(base time.Duration, factor float64) (*Jitter, error) {
	if base <= 0 {
		return nil, errors.New("jitter: base duration must be positive")
	}
	if factor < 0 || factor > 1 {
		return nil, errors.New("jitter: factor must be between 0 and 1")
	}
	return &Jitter{
		base:   base,
		factor: factor,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

// Next returns the next interval with jitter applied.
func (j *Jitter) Next() time.Duration {
	if j.factor == 0 {
		return j.base
	}
	max := float64(j.base) * j.factor
	offset := time.Duration(j.rng.Float64() * max)
	return j.base + offset
}

// Base returns the configured base duration.
func (j *Jitter) Base() time.Duration {
	return j.base
}

// Factor returns the configured jitter factor.
func (j *Jitter) Factor() float64 {
	return j.factor
}
