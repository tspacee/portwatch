package watch

import (
	"testing"
	"time"
)

func TestNewWindow_InvalidSize(t *testing.T) {
	_, err := NewWindow(0)
	if err == nil {
		t.Fatal("expected error for zero size")
	}
	_, err = NewWindow(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative size")
	}
}

func TestNewWindow_Valid(t *testing.T) {
	w, err := NewWindow(time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil window")
	}
}

func TestWindow_Count_Empty(t *testing.T) {
	w, _ := NewWindow(time.Second)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestWindow_Record_IncrementsCount(t *testing.T) {
	w, _ := NewWindow(time.Second)
	w.Record()
	w.Record()
	if got := w.Count(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestWindow_Count_EvictsExpired(t *testing.T) {
	w, _ := NewWindow(50 * time.Millisecond)
	w.Record()
	w.Record()
	time.Sleep(80 * time.Millisecond)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after expiry, got %d", got)
	}
}

func TestWindow_Reset_ClearsEvents(t *testing.T) {
	w, _ := NewWindow(time.Second)
	w.Record()
	w.Record()
	w.Reset()
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestWindow_Record_AfterExpiry_CountsOnlyRecent(t *testing.T) {
	w, _ := NewWindow(60 * time.Millisecond)
	w.Record()
	time.Sleep(80 * time.Millisecond)
	w.Record()
	if got := w.Count(); got != 1 {
		t.Fatalf("expected 1 recent event, got %d", got)
	}
}
