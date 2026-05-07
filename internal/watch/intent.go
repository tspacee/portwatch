package watch

import (
	"errors"
	"sync"
)

// Intent records the declared purpose or expected state for a port.
// It allows operators to annotate ports with a reason they should be open,
// making unexpected ports easier to identify during evaluation.
type Intent struct {
	mu      sync.RWMutex
	entries map[int]string
}

// NewIntent returns an initialised Intent registry.
func NewIntent() *Intent {
	return &Intent{
		entries: make(map[int]string),
	}
}

// Declare associates a port with a stated intent.
// Returns an error if the port is out of range or the intent string is empty.
func (i *Intent) Declare(port int, intent string) error {
	if port < 1 || port > 65535 {
		return errors.New("intent: port out of range")
	}
	if intent == "" {
		return errors.New("intent: intent string must not be empty")
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	i.entries[port] = intent
	return nil
}

// Lookup returns the declared intent for a port and whether one exists.
func (i *Intent) Lookup(port int) (string, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	v, ok := i.entries[port]
	return v, ok
}

// Revoke removes a declared intent for a port.
func (i *Intent) Revoke(port int) {
	i.mu.Lock()
	defer i.mu.Unlock()
	delete(i.entries, port)
}

// Len returns the number of declared intents.
func (i *Intent) Len() int {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return len(i.entries)
}
