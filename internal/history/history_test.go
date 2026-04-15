package history

import (
	"testing"
	"time"
)

func makeEntry(violations ...string) Entry {
	return Entry{
		Timestamp:  time.Now(),
		OpenPorts:  []int{80, 443},
		Violations: violations,
	}
}

func TestNew_DefaultMaxSize(t *testing.T) {
	h := New(0)
	if h.maxSize != 100 {
		t.Fatalf("expected default maxSize 100, got %d", h.maxSize)
	}
}

func TestAdd_And_Len(t *testing.T) {
	h := New(10)
	h.Add(makeEntry())
	h.Add(makeEntry())
	if h.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", h.Len())
	}
}

func TestAdd_Evicts_Oldest(t *testing.T) {
	h := New(3)
	for i := 0; i < 5; i++ {
		h.Add(makeEntry())
	}
	if h.Len() != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", h.Len())
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	h := New(10)
	h.Add(makeEntry("port 9999 unexpected"))

	entries := h.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	// Mutate the copy — original must be unaffected.
	entries[0].Violations = nil
	if len(h.Entries()[0].Violations) == 0 {
		t.Fatal("mutating returned slice affected internal state")
	}
}

func TestViolationCount_SumsAcrossEntries(t *testing.T) {
	h := New(10)
	h.Add(makeEntry("v1", "v2"))
	h.Add(makeEntry())
	h.Add(makeEntry("v3"))

	if got := h.ViolationCount(); got != 3 {
		t.Fatalf("expected violation count 3, got %d", got)
	}
}

func TestViolationCount_EmptyHistory(t *testing.T) {
	h := New(10)
	if got := h.ViolationCount(); got != 0 {
		t.Fatalf("expected 0 violations on empty history, got %d", got)
	}
}
