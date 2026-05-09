package watch

import (
	"fmt"
	"sync"
)

// Peak tracks the maximum number of simultaneously open ports observed
// across scan cycles. It provides a high-water mark for port activity.
type Peak struct {
	mu    sync.Mutex
	value int
	set   bool
}

// NewPeak returns a new Peak tracker.
func NewPeak() *Peak {
	return &Peak{}
}

// Observe records the current port count and updates the peak if the
// provided count exceeds the previously recorded maximum.
func (p *Peak) Observe(count int) error {
	if count < 0 {
		return fmt.Errorf("peak: count must be non-negative, got %d", count)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.set || count > p.value {
		p.value = count
		p.set = true
	}
	return nil
}

// Value returns the highest port count observed. The second return value
// is false if no observation has been recorded yet.
func (p *Peak) Value() (int, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.value, p.set
}

// Reset clears the recorded peak, allowing tracking to restart.
func (p *Peak) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.value = 0
	p.set = false
}
