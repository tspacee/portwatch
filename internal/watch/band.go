package watch

import (
	"fmt"
	"sync"
)

// Band groups ports into named frequency bands and tracks how often each
// band is active during a scan cycle. Useful for classifying port activity
// by well-known service tiers (e.g., "system", "registered", "dynamic").
type Band struct {
	mu    sync.RWMutex
	bands map[string][2]int // name -> [min, max]
	hits  map[string]int
}

// NewBand returns an initialised Band with no ranges defined.
func NewBand() *Band {
	return &Band{
		bands: make(map[string][2]int),
		hits:  make(map[string]int),
	}
}

// Define registers a named band covering ports in [min, max] inclusive.
// Returns an error if the name is empty or the range is invalid.
func (b *Band) Define(name string, min, max int) error {
	if name == "" {
		return fmt.Errorf("band name must not be empty")
	}
	if min < 1 || max > 65535 || min > max {
		return fmt.Errorf("band %q: invalid port range [%d, %d]", name, min, max)
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.bands[name] = [2]int{min, max}
	return nil
}

// Classify records which bands the given ports fall into and increments
// their hit counters. Ports that match no band are silently ignored.
func (b *Band) Classify(ports []int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, p := range ports {
		for name, r := range b.bands {
			if p >= r[0] && p <= r[1] {
				b.hits[name]++
			}
		}
	}
}

// Hits returns the current hit count for the named band.
// Returns 0 if the band is unknown.
func (b *Band) Hits(name string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.hits[name]
}

// Reset clears all hit counters without removing band definitions.
func (b *Band) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for k := range b.hits {
		delete(b.hits, k)
	}
}

// Snapshot returns a copy of the current hit map.
func (b *Band) Snapshot() map[string]int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make(map[string]int, len(b.hits))
	for k, v := range b.hits {
		out[k] = v
	}
	return out
}
