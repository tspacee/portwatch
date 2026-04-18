package watch

import (
	"testing"
)

func TestNewRegistry_Empty(t *testing.T) {
	r := NewRegistry()
	if r.Len() != 0 {
		t.Fatalf("expected 0, got %d", r.Len())
	}
}

func TestRegistry_Track_Valid(t *testing.T) {
	r := NewRegistry()
	if err := r.Track(443, "tcp"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Len() != 1 {
		t.Fatalf("expected 1, got %d", r.Len())
	}
}

func TestRegistry_Track_InvalidPort(t *testing.T) {
	r := NewRegistry()
	if err := r.Track(0, "tcp"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := r.Track(65536, "tcp"); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestRegistry_Track_IncrementsSeenCount(t *testing.T) {
	r := NewRegistry()
	_ = r.Track(80, "tcp")
	_ = r.Track(80, "tcp")
	m, ok := r.Get(80)
	if !ok {
		t.Fatal("expected entry")
	}
	if m.SeenCount != 2 {
		t.Fatalf("expected SeenCount=2, got %d", m.SeenCount)
	}
}

func TestRegistry_Get_Missing(t *testing.T) {
	r := NewRegistry()
	_, ok := r.Get(9999)
	if ok {
		t.Fatal("expected missing")
	}
}

func TestRegistry_Remove(t *testing.T) {
	r := NewRegistry()
	_ = r.Track(22, "tcp")
	r.Remove(22)
	if r.Len() != 0 {
		t.Fatalf("expected 0 after remove, got %d", r.Len())
	}
}

func TestRegistry_Get_ReturnsCopy(t *testing.T) {
	r := NewRegistry()
	_ = r.Track(8080, "tcp")
	m, _ := r.Get(8080)
	m.SeenCount = 999
	orig, _ := r.Get(8080)
	if orig.SeenCount == 999 {
		t.Fatal("Get should return a copy, not a reference")
	}
}
