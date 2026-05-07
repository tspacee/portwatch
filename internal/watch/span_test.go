package watch

import (
	"testing"
	"time"
)

func TestNewSpan_Empty(t *testing.T) {
	s := NewSpan()
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
}

func TestSpan_Observe_Valid(t *testing.T) {
	s := NewSpan()
	if err := s.Observe(8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 1 {
		t.Fatalf("expected 1, got %d", s.Len())
	}
}

func TestSpan_Observe_InvalidPort(t *testing.T) {
	s := NewSpan()
	if err := s.Observe(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := s.Observe(65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestSpan_Observe_Idempotent(t *testing.T) {
	s := NewSpan()
	_ = s.Observe(443)
	first, _ := s.Duration(443)
	time.Sleep(2 * time.Millisecond)
	_ = s.Observe(443)
	second, _ := s.Duration(443)
	if second < first {
		t.Fatal("re-observing a port should not reset its start time")
	}
	if s.Len() != 1 {
		t.Fatalf("expected 1 tracked port, got %d", s.Len())
	}
}

func TestSpan_Duration_NeverObserved(t *testing.T) {
	s := NewSpan()
	d, ok := s.Duration(22)
	if ok {
		t.Fatal("expected ok=false for unseen port")
	}
	if d != 0 {
		t.Fatalf("expected 0 duration, got %v", d)
	}
}

func TestSpan_Duration_Positive(t *testing.T) {
	s := NewSpan()
	_ = s.Observe(9000)
	time.Sleep(5 * time.Millisecond)
	d, ok := s.Duration(9000)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if d < time.Millisecond {
		t.Fatalf("expected positive duration, got %v", d)
	}
}

func TestSpan_Reset_RemovesPort(t *testing.T) {
	s := NewSpan()
	_ = s.Observe(3306)
	if err := s.Reset(3306); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
	_, ok := s.Duration(3306)
	if ok {
		t.Fatal("expected port to be absent after reset")
	}
}

func TestSpan_Reset_InvalidPort(t *testing.T) {
	s := NewSpan()
	if err := s.Reset(0); err == nil {
		t.Fatal("expected error for port 0")
	}
}
