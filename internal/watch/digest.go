package watch

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
)

// Digest computes and caches a fingerprint of a port set,
// allowing fast detection of changes between scans.
type Digest struct {
	mu   sync.Mutex
	last string
}

// NewDigest returns a new Digest instance.
func NewDigest() *Digest {
	return &Digest{}
}

// Compute returns the SHA-256 hex digest of the sorted port list.
func (d *Digest) Compute(ports []int) string {
	sorted := make([]int, len(ports))
	copy(sorted, ports)
	sort.Ints(sorted)

	h := sha256.New()
	for _, p := range sorted {
		fmt.Fprintf(h, "%d:", p)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Changed returns true if the digest of ports differs from the last stored digest.
// It updates the stored digest on every call.
func (d *Digest) Changed(ports []int) bool {
	current := d.Compute(ports)
	d.mu.Lock()
	defer d.mu.Unlock()
	if current == d.last {
		return false
	}
	d.last = current
	return true
}

// Last returns the most recently stored digest string.
func (d *Digest) Last() string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.last
}

// Reset clears the stored digest.
func (d *Digest) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.last = ""
}
