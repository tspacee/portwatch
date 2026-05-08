package watch

import (
	"errors"
	"sync"
)

// Forecast predicts whether a port is likely to appear in the next scan
// based on its historical presence ratio.
type Forecast struct {
	mu      sync.Mutex
	history map[int][]bool
	window  int
}

// NewForecast creates a Forecast with the given rolling window size.
func NewForecast(window int) (*Forecast, error) {
	if window < 1 {
		return nil, errors.New("forecast: window must be at least 1")
	}
	return &Forecast{
		history: make(map[int][]bool),
		window:  window,
	}, nil
}

// Observe records whether a port was seen (present=true) in a scan.
func (f *Forecast) Observe(port int, present bool) error {
	if port < 1 || port > 65535 {
		return errors.New("forecast: port out of range")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	h := append(f.history[port], present)
	if len(h) > f.window {
		h = h[len(h)-f.window:]
	}
	f.history[port] = h
	return nil
}

// Likelihood returns the fraction of recent scans in which the port was present.
// Returns 0 if the port has never been observed.
func (f *Forecast) Likelihood(port int) float64 {
	f.mu.Lock()
	defer f.mu.Unlock()
	h := f.history[port]
	if len(h) == 0 {
		return 0
	}
	var seen int
	for _, v := range h {
		if v {
			seen++
		}
	}
	return float64(seen) / float64(len(h))
}

// Predict returns true if the port's likelihood meets or exceeds the given threshold.
func (f *Forecast) Predict(port int, threshold float64) bool {
	return f.Likelihood(port) >= threshold
}

// Reset clears all observed history for a port.
func (f *Forecast) Reset(port int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.history, port)
}
