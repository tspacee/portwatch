package watch

import (
	"errors"
	"sync"
)

// NoteMap stores arbitrary string notes keyed by port number.
// It is safe for concurrent use.
type NoteMap struct {
	mu    sync.RWMutex
	notes map[int]string
}

// ErrInvalidNotePort is returned when a port number is out of range.
var ErrInvalidNotePort = errors.New("notemap: port must be between 1 and 65535")

// ErrEmptyNote is returned when an empty note string is provided.
var ErrEmptyNote = errors.New("notemap: note must not be empty")

// NewNoteMap creates an empty NoteMap.
func NewNoteMap() *NoteMap {
	return &NoteMap{notes: make(map[int]string)}
}

// Set associates a note with a port. Returns an error for invalid inputs.
func (n *NoteMap) Set(port int, note string) error {
	if port < 1 || port > 65535 {
		return ErrInvalidNotePort
	}
	if note == "" {
		return ErrEmptyNote
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.notes[port] = note
	return nil
}

// Get retrieves the note for a port. Returns empty string if not found.
func (n *NoteMap) Get(port int) string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.notes[port]
}

// Delete removes the note for a port.
func (n *NoteMap) Delete(port int) {
	n.mu.Lock()
	defer n.mu.Unlock()
	delete(n.notes, port)
}

// Snapshot returns a copy of all current notes.
func (n *NoteMap) Snapshot() map[int]string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	copy := make(map[int]string, len(n.notes))
	for k, v := range n.notes {
		copy[k] = v
	}
	return copy
}
