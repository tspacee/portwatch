package watch

import (
	"fmt"
	"sync"
	"time"
)

// TrendDirection indicates whether a port's activity is increasing, stable, or decreasing.
type TrendDirection int

const (
	// TrendStable indicates no significant change in activity.
	TrendStable TrendDirection = iota
	// TrendRising indicates increasing port activity.
	TrendRising
	// TrendFalling indicates decreasing port activity.
	TrendFalling
)

// String returns a human-readable label for the trend direction.
func (d TrendDirection) String() string {
	switch d {
	case TrendRising:
		return "rising"
	case TrendFalling:
		return "falling"
	default:
		return "stable"
	}
}

// trendBucket holds a count and the time it was recorded.
type trendBucket struct {
	at    time.Time
	count int
}

// Trend tracks per-port activity over two consecutive windows to determine
// whether port visibility is rising, falling, or stable.
type Trend struct {
	mu       sync.Mutex
	window   time.Duration
	threshold float64
	buckets  map[int][2]trendBucket // [0]=previous window, [1]=current window
	now      func() time.Time
}

// NewTrend creates a Trend tracker. window is the duration of each measurement
// period. threshold is the fractional change (0.0–1.0) required to declare a
// rising or falling trend; e.g. 0.25 means a 25% change.
func NewTrend(window time.Duration, threshold float64) (*Trend, error) {
	if window <= 0 {
		return nil, fmt.Errorf("trend: window must be positive, got %s", window)
	}
	if threshold < 0 || threshold > 1 {
		return nil, fmt.Errorf("trend: threshold must be between 0.0 and 1.0, got %.2f", threshold)
	}
	return &Trend{
		window:    window,
		threshold: threshold,
		buckets:   make(map[int][2]trendBucket),
		now:       time.Now,
	}, nil
}

// Record increments the activity count for the given port in the current window.
// If the current window has expired, the previous window is replaced and a new
// current window begins.
func (t *Trend) Record(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("trend: port %d out of range", port)
	}
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	pair := t.buckets[port]
	current := pair[1]

	if current.at.IsZero() || now.Sub(current.at) >= t.window {
		// Rotate: current becomes previous, start a fresh current bucket.
		pair[0] = pair[1]
		pair[1] = trendBucket{at: now, count: 1}
	} else {
		pair[1].count++
	}
	t.buckets[port] = pair
	return nil
}

// Direction returns the TrendDirection for the given port by comparing the
// count in the previous window to the count in the current window.
func (t *Trend) Direction(port int) TrendDirection {
	t.mu.Lock()
	defer t.mu.Unlock()

	pair, ok := t.buckets[port]
	if !ok {
		return TrendStable
	}
	prev := float64(pair[0].count)
	curr := float64(pair[1].count)

	if prev == 0 && curr == 0 {
		return TrendStable
	}
	if prev == 0 {
		return TrendRising
	}
	change := (curr - prev) / prev
	if change >= t.threshold {
		return TrendRising
	}
	if change <= -t.threshold {
		return TrendFalling
	}
	return TrendStable
}

// Reset clears all recorded data for the given port.
func (t *Trend) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.buckets, port)
}
