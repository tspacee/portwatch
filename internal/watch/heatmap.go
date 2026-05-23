package watch

import (
	"errors"
	"sync"
)

// HeatMap tracks how frequently each port has been observed across scans.
// Ports observed more often accumulate a higher heat value.
type HeatMap struct {
	mu     sync.RWMutex
	heat   map[int]int
	total  int
}

// NewHeatMap returns an initialised HeatMap.
func NewHeatMap() *HeatMap {
	return &HeatMap{
		heat: make(map[int]int),
	}
}

// Record increments the heat value for the given port.
func (h *HeatMap) Record(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("heatmap: port out of range")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.heat[port]++
	h.total++
	return nil
}

// Heat returns the raw observation count for a port.
func (h *HeatMap) Heat(port int) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.heat[port]
}

// Hottest returns the port with the highest heat value.
// Returns 0 if no ports have been recorded.
func (h *HeatMap) Hottest() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	var best, bestVal int
	for port, val := range h.heat {
		if val > bestVal {
			best = port
			bestVal = val
		}
	}
	return best
}

// Snapshot returns a copy of the current heat map.
func (h *HeatMap) Snapshot() map[int]int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make(map[int]int, len(h.heat))
	for k, v := range h.heat {
		out[k] = v
	}
	return out
}

// Reset clears all heat values.
func (h *HeatMap) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.heat = make(map[int]int)
	h.total = 0
}
