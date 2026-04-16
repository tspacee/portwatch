package watch

import (
	"sync"
	"time"
)

// EventType classifies a port change event.
type EventType string

const (
	EventPortOpened EventType = "opened"
	EventPortClosed EventType = "closed"
)

// PortEvent represents a single detected port change.
type PortEvent struct {
	Type      EventType
	Port      int
	Protocol  string
	DetectedAt time.Time
}

// EventLog stores a bounded in-memory log of recent port events.
type EventLog struct {
	mu      sync.RWMutex
	events  []PortEvent
	maxSize int
}

// NewEventLog creates an EventLog with the given capacity.
// Returns an error if maxSize < 1.
func NewEventLog(maxSize int) (*EventLog, error) {
	if maxSize < 1 {
		return nil, ErrInvalidEventLogSize
	}
	return &EventLog{maxSize: maxSize}, nil
}

// Add appends an event, evicting the oldest if at capacity.
func (l *EventLog) Add(e PortEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.events) >= l.maxSize {
		l.events = l.events[1:]
	}
	l.events = append(l.events, e)
}

// Entries returns a copy of all stored events.
func (l *EventLog) Entries() []PortEvent {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]PortEvent, len(l.events))
	copy(out, l.events)
	return out
}

// Len returns the current number of stored events.
func (l *EventLog) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.events)
}
