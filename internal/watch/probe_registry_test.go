package watch

import (
	"errors"
	"testing"
	"time"
)

func TestNewProbeRegistry_InvalidTimeout(t *testing.T) {
	_, err := NewProbeRegistry(0)
	if err == nil {
		t.Fatal("expected error for zero timeout")
	}
}

func TestNewProbeRegistry_Valid(t *testing.T) {
	r, err := NewProbeRegistry(time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Len() != 0 {
		t.Fatalf("expected empty registry")
	}
}

func TestProbeRegistry_Record_Valid(t *testing.T) {
	r, _ := NewProbeRegistry(time.Second)
	if err := r.Record(8080, true, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Len() != 1 {
		t.Fatalf("expected 1 record, got %d", r.Len())
	}
}

func TestProbeRegistry_Record_InvalidPort(t *testing.T) {
	r, _ := NewProbeRegistry(time.Second)
	if err := r.Record(0, true, nil); err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestProbeRegistry_Get_Present(t *testing.T) {
	r, _ := NewProbeRegistry(time.Second)
	sentinel := errors.New("probe failed")
	_ = r.Record(443, false, sentinel)
	rec := r.Get(443)
	if rec == nil {
		t.Fatal("expected record, got nil")
	}
	if rec.Open {
		t.Fatal("expected closed port")
	}
	if rec.Err != sentinel {
		t.Fatalf("expected sentinel error, got %v", rec.Err)
	}
}

func TestProbeRegistry_Get_Missing(t *testing.T) {
	r, _ := NewProbeRegistry(time.Second)
	if r.Get(9999) != nil {
		t.Fatal("expected nil for missing port")
	}
}

func TestProbeRegistry_Get_ReturnsCopy(t *testing.T) {
	r, _ := NewProbeRegistry(time.Second)
	_ = r.Record(80, true, nil)
	rec := r.Get(80)
	rec.Open = false
	if got := r.Get(80); !got.Open {
		t.Fatal("mutation of returned record should not affect registry")
	}
}

func TestProbeRegistry_Clear_ResetsLen(t *testing.T) {
	r, _ := NewProbeRegistry(time.Second)
	_ = r.Record(80, true, nil)
	_ = r.Record(443, true, nil)
	r.Clear()
	if r.Len() != 0 {
		t.Fatalf("expected 0 after clear, got %d", r.Len())
	}
}
