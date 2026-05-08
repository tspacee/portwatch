package watch

import (
	"testing"
	"time"
)

func TestNewAnchor_Empty(t *testing.T) {
	a := NewAnchor()
	if a.Len() != 0 {
		t.Fatalf("expected 0, got %d", a.Len())
	}
}

func TestAnchor_Pin_Valid(t *testing.T) {
	a := NewAnchor()
	if err := a.Pin(80); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Len() != 1 {
		t.Fatalf("expected 1, got %d", a.Len())
	}
}

func TestAnchor_Pin_InvalidPort(t *testing.T) {
	a := NewAnchor()
	if err := a.Pin(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := a.Pin(65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestAnchor_Pin_Idempotent(t *testing.T) {
	a := NewAnchor()
	_ = a.Pin(443)
	time.Sleep(5 * time.Millisecond)
	_ = a.Pin(443) // second pin should not overwrite
	d, ok := a.Since(443)
	if !ok {
		t.Fatal("expected port to be anchored")
	}
	if d < 5*time.Millisecond {
		t.Fatalf("expected duration >= 5ms, got %v", d)
	}
	if a.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", a.Len())
	}
}

func TestAnchor_Since_NeverPinned(t *testing.T) {
	a := NewAnchor()
	_, ok := a.Since(8080)
	if ok {
		t.Fatal("expected false for unseen port")
	}
}

func TestAnchor_Since_ReturnsPositiveDuration(t *testing.T) {
	a := NewAnchor()
	_ = a.Pin(22)
	time.Sleep(2 * time.Millisecond)
	d, ok := a.Since(22)
	if !ok {
		t.Fatal("expected port to be anchored")
	}
	if d <= 0 {
		t.Fatalf("expected positive duration, got %v", d)
	}
}

func TestAnchor_Release_RemovesEntry(t *testing.T) {
	a := NewAnchor()
	_ = a.Pin(3306)
	a.Release(3306)
	if a.Len() != 0 {
		t.Fatalf("expected 0 after release, got %d", a.Len())
	}
	_, ok := a.Since(3306)
	if ok {
		t.Fatal("expected false after release")
	}
}
