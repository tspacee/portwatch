package watch

import (
	"errors"
	"sync"
)

// Watermark tracks the high-water mark (maximum observed count) for each port
// across scan cycles. It is useful for detecting ports that have exceeded a
// historically significant threshold.
type Watermark struct {
	mu     sync.RWMutex
	marks  map[int]int
}

// NewWatermark returns an initialised Watermark.
func NewWatermark() *Watermark {
	return &Watermark{
		marks: make(map[int]int),
	}
}

// Record updates the high-water mark for port if value exceeds the current
// mark. Returns an error if port is out of the valid 1–65535 range.
func (w *Watermark) Record(port, value int) error {
	if port < 1 || port > 65535 {
		return errors.New("watermark: port out of range")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if value > w.marks[port] {
		w.marks[port] = value
	}
	return nil
}

// Peak returns the highest recorded value for port and whether an entry
// exists.
func (w *Watermark) Peak(port int) (int, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	v, ok := w.marks[port]
	return v, ok
}

// Reset clears the high-water mark for port.
func (w *Watermark) Reset(port int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.marks, port)
}

// Snapshot returns a copy of all current high-water marks keyed by port.
func (w *Watermark) Snapshot() map[int]int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make(map[int]int, len(w.marks))
	for k, v := range w.marks {
		out[k] = v
	}
	return out
}
