package watch

import (
	"testing"
	"time"
)

func TestNewGrace_InvalidWindow(t *testing.T) {
	_, err := NewGrace(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
	_, err = NewGrace(-1 * time.Second)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestNewGrace_Valid(t *testing.T) {
	g, err := NewGrace(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.Len() != 0 {
		t.Fatal("expected empty grace tracker")
	}
}

func TestGrace_Observe_TracksPort(t *testing.T) {
	g, _ := NewGrace(100 * time.Millisecond)
	g.Observe(8080)
	if g.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", g.Len())
	}
}

func TestGrace_Observe_Idempotent(t *testing.T) {
	g, _ := NewGrace(100 * time.Millisecond)
	g.Observe(8080)
	g.Observe(8080)
	if g.Len() != 1 {
		t.Fatalf("expected 1 entry after duplicate observe, got %d", g.Len())
	}
}

func TestGrace_Settled_NotYet(t *testing.T) {
	g, _ := NewGrace(500 * time.Millisecond)
	g.Observe(9090)
	if g.Settled(9090) {
		t.Fatal("expected port to not be settled yet")
	}
}

func TestGrace_Settled_AfterWindow(t *testing.T) {
	g, _ := NewGrace(20 * time.Millisecond)
	g.Observe(9090)
	time.Sleep(30 * time.Millisecond)
	if !g.Settled(9090) {
		t.Fatal("expected port to be settled after window elapsed")
	}
}

func TestGrace_Settled_UnknownPort(t *testing.T) {
	g, _ := NewGrace(100 * time.Millisecond)
	if g.Settled(1234) {
		t.Fatal("expected false for unobserved port")
	}
}

func TestGrace_Clear_RemovesEntry(t *testing.T) {
	g, _ := NewGrace(100 * time.Millisecond)
	g.Observe(7070)
	g.Clear(7070)
	if g.Len() != 0 {
		t.Fatal("expected empty tracker after clear")
	}
	if g.Settled(7070) {
		t.Fatal("expected false for cleared port")
	}
}

func TestGrace_Clear_AllowsReobserve(t *testing.T) {
	g, _ := NewGrace(20 * time.Millisecond)
	g.Observe(7070)
	time.Sleep(30 * time.Millisecond)
	g.Clear(7070)
	g.Observe(7070)
	// freshly observed — should not be settled yet
	if g.Settled(7070) {
		t.Fatal("expected port to not be settled after re-observe")
	}
}
