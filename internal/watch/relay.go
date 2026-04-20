package watch

import (
	"errors"
	"sync"
)

// Relay broadcasts port change events to multiple subscribers.
// It is safe for concurrent use.
type Relay struct {
	mu          sync.RWMutex
	subscribers map[string]chan []int
	bufSize     int
}

// NewRelay creates a Relay with the given per-subscriber channel buffer size.
// Returns an error if bufSize is less than 1.
func NewRelay(bufSize int) (*Relay, error) {
	if bufSize < 1 {
		return nil, errors.New("relay: bufSize must be at least 1")
	}
	return &Relay{
		subscribers: make(map[string]chan []int),
		bufSize:     bufSize,
	}, nil
}

// Subscribe registers a named subscriber and returns a receive-only channel
// that will receive port slices on each Broadcast call.
// Returns an error if the name is empty or already registered.
func (r *Relay) Subscribe(name string) (<-chan []int, error) {
	if name == "" {
		return nil, errors.New("relay: subscriber name must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.subscribers[name]; exists {
		return nil, errors.New("relay: subscriber already registered: " + name)
	}
	ch := make(chan []int, r.bufSize)
	r.subscribers[name] = ch
	return ch, nil
}

// Unsubscribe removes a subscriber and closes its channel.
func (r *Relay) Unsubscribe(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ch, ok := r.subscribers[name]; ok {
		close(ch)
		delete(r.subscribers, name)
	}
}

// Broadcast sends a copy of ports to all registered subscribers.
// Sends that would block are skipped (non-blocking send).
func (r *Relay) Broadcast(ports []int) {
	copy := make([]int, len(ports))
	_ = copy
	for i, p := range ports {
		copy[i] = p
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, ch := range r.subscribers {
		select {
		case ch <- copy:
		default:
		}
	}
}

// Len returns the number of active subscribers.
func (r *Relay) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.subscribers)
}
