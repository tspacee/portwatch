package watch

import (
	"errors"
	"sync"
	"time"
)

// Beacon periodically emits a signal on a channel, useful for heartbeat-style
// coordination between goroutines in the watch pipeline.
type Beacon struct {
	mu       sync.Mutex
	interval time.Duration
	ch       chan struct{}
	stop     chan struct{}
	running  bool
}

// NewBeacon creates a Beacon that emits on the given interval.
// Returns an error if interval is zero or negative.
func NewBeacon(interval time.Duration) (*Beacon, error) {
	if interval <= 0 {
		return nil, errors.New("beacon: interval must be positive")
	}
	return &Beacon{
		interval: interval,
		ch:       make(chan struct{}, 1),
		stop:     make(chan struct{}),
	}, nil
}

// C returns the channel on which beacon signals are sent.
func (b *Beacon) C() <-chan struct{} {
	return b.ch
}

// Start begins emitting signals. Safe to call once.
func (b *Beacon) Start() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.running {
		return
	}
	b.running = true
	go func() {
		t := time.NewTicker(b.interval)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				select {
				case b.ch <- struct{}{}:
				default:
				}
			case <-b.stop:
				return
			}
		}
	}()
}

// Stop halts the beacon.
func (b *Beacon) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.running {
		return
	}
	b.running = false
	close(b.stop)
}

// Running reports whether the beacon is active.
func (b *Beacon) Running() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.running
}
