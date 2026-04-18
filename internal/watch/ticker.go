package watch

import (
	"errors"
	"time"
)

// Ticker emits ticks at a configurable interval with optional jitter.
// It wraps time.Ticker and adds jitter support via a Jitter instance.
type Ticker struct {
	base    time.Duration
	jitter  *Jitter
	stop    chan struct{}
	C       <-chan time.Time
	internal chan time.Time
}

// NewTicker creates a Ticker with the given base interval and optional jitter factor.
// factor must be in [0, 1]. Use 0 for no jitter.
func NewTicker(base time.Duration, factor float64) (*Ticker, error) {
	if base <= 0 {
		return nil, errors.New("ticker: base interval must be positive")
	}
	j, err := NewJitter(base, factor)
	if err != nil {
		return nil, err
	}
	ch := make(chan time.Time, 1)
	t := &Ticker{
		base:     base,
		jitter:   j,
		stop:     make(chan struct{}),
		C:        ch,
		internal: ch,
	}
	go t.run()
	return t, nil
}

func (t *Ticker) run() {
	for {
		d := t.jitter.Next()
		select {
		case <-time.After(d):
			select {
			case t.internal <- time.Now():
			default:
			}
		case <-t.stop:
			return
		}
	}
}

// Stop halts the ticker. Safe to call multiple times.
func (t *Ticker) Stop() {
	select {
	case <-t.stop:
	default:
		close(t.stop)
	}
}
