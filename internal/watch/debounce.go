package watch

import (
	"sync"
	"time"
)

// Debounce delays execution of a function until after a quiet period has elapsed.
// If the trigger is called again before the delay expires, the timer resets.
type Debounce struct {
	delay  time.Duration
	mu     sync.Mutex
	timer  *time.Timer
	closed chan struct{}
}

// NewDebounce creates a Debounce with the given delay.
// Returns an error if delay is zero or negative.
func NewDebounce(delay time.Duration) (*Debounce, error) {
	if delay <= 0 {
		return nil, ErrInvalidWindow
	}
	return &Debounce{
		delay:  delay,
		closed: make(chan struct{}),
	}, nil
}

// Trigger schedules fn to be called after the debounce delay.
// If Trigger is called again before the delay elapses, the previous call is cancelled.
func (d *Debounce) Trigger(fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	select {
	case <-d.closed:
		return
	default:
	}

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.delay, func() {
		select {
		case <-d.closed:
		default:
			fn()
		}
	})
}

// Stop cancels any pending debounced call.
func (d *Debounce) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()
	select {
	case <-d.closed:
	default:
		close(d.closed)
	}
	if d.timer != nil {
		d.timer.Stop()
	}
}
