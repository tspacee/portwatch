package watch

import (
	"errors"
	"sync"
	"time"
)

// DedupCache suppresses repeated events within a TTL window.
type DedupCache struct {
	mu    sync.Mutex
	entries map[string]time.Time
	ttl   time.Duration
}

var ErrInvalidTTL = errors.New("dedup cache: TTL must be greater than zero")

// NewDedupCache creates a DedupCache with the given TTL.
func NewDedupCache(ttl time.Duration) (*DedupCache, error) {
	if ttl <= 0 {
		return nil, ErrInvalidTTL
	}
	return &DedupCache{
		entries: make(map[string]time.Time),
		ttl:     ttl,
	}, nil
}

// Seen returns true if the key was seen within the TTL window.
// If not seen (or expired), it records the key and returns false.
func (d *DedupCache) Seen(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	if exp, ok := d.entries[key]; ok && now.Before(exp) {
		return true
	}
	d.entries[key] = now.Add(d.ttl)
	return false
}

// Evict removes expired entries from the cache.
func (d *DedupCache) Evict() {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	for k, exp := range d.entries {
		if now.After(exp) {
			delete(d.entries, k)
		}
	}
}

// Len returns the number of active (non-expired) entries.
func (d *DedupCache) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	count := 0
	for _, exp := range d.entries {
		if now.Before(exp) {
			count++
		}
	}
	return count
}
