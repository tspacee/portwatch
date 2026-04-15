package history

import (
	"bytes"
	"log"
	"testing"
	"time"
)

func TestNewCleaner_NilHistory(t *testing.T) {
	_, err := NewCleaner(nil, time.Hour, time.Minute, nil)
	if err != ErrNilHistory {
		t.Fatalf("expected ErrNilHistory, got %v", err)
	}
}

func TestNewCleaner_Valid(t *testing.T) {
	h := New(10)
	c, err := NewCleaner(h, time.Hour, time.Minute, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Cleaner")
	}
}

func TestCleaner_Sweep_RemovesOldEntries(t *testing.T) {
	h := New(20)
	now := time.Now()

	// Add two old entries and one fresh entry.
	h.Add(Entry{Timestamp: now.Add(-2 * time.Hour), Added: []int{80}})
	h.Add(Entry{Timestamp: now.Add(-90 * time.Minute), Added: []int{443}})
	h.Add(Entry{Timestamp: now.Add(-1 * time.Minute), Added: []int{8080}})

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	c, _ := NewCleaner(h, time.Hour, time.Minute, logger)

	removed := c.sweep()
	if removed != 2 {
		t.Fatalf("expected 2 removed, got %d", removed)
	}
	if h.Len() != 1 {
		t.Fatalf("expected 1 remaining entry, got %d", h.Len())
	}
	if h.Entries()[0].Added[0] != 8080 {
		t.Errorf("expected remaining port 8080, got %v", h.Entries()[0].Added)
	}
}

func TestCleaner_Sweep_NothingToRemove(t *testing.T) {
	h := New(10)
	h.Add(Entry{Timestamp: time.Now(), Added: []int{22}})

	c, _ := NewCleaner(h, time.Hour, time.Minute, nil)
	removed := c.sweep()
	if removed != 0 {
		t.Fatalf("expected 0 removed, got %d", removed)
	}
}

func TestCleaner_Run_StopsOnClose(t *testing.T) {
	h := New(10)
	c, _ := NewCleaner(h, time.Hour, 10*time.Millisecond, nil)

	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		c.Run(stop)
		close(done)
	}()

	time.Sleep(25 * time.Millisecond)
	close(stop)

	select {
	case <-done:
		// ok
	case <-time.After(time.Second):
		t.Fatal("Cleaner.Run did not stop after channel close")
	}
}
