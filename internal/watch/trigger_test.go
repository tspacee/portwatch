package watch

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestNewTrigger_InvalidWindow(t *testing.T) {
	_, err := NewTrigger(0, 2, func() {})
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewTrigger_InvalidThreshold(t *testing.T) {
	_, err := NewTrigger(time.Second, 0, func() {})
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNewTrigger_NilCallback(t *testing.T) {
	_, err := NewTrigger(time.Second, 1, nil)
	if err == nil {
		t.Fatal("expected error for nil callback")
	}
}

func TestNewTrigger_Valid(t *testing.T) {
	tr, err := NewTrigger(time.Second, 3, func() {})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil trigger")
	}
}

func TestTrigger_Count_Empty(t *testing.T) {
	tr, _ := NewTrigger(time.Second, 3, func() {})
	if tr.Count() != 0 {
		t.Fatal("expected count 0")
	}
}

func TestTrigger_Record_BelowThreshold(t *testing.T) {
	var fired int32
	tr, _ := NewTrigger(time.Second, 3, func() { atomic.AddInt32(&fired, 1) })
	tr.Record()
	tr.Record()
	time.Sleep(20 * time.Millisecond)
	if atomic.LoadInt32(&fired) != 0 {
		t.Fatal("callback should not have fired")
	}
	if tr.Count() != 2 {
		t.Fatalf("expected count 2, got %d", tr.Count())
	}
}

func TestTrigger_Record_FiresAtThreshold(t *testing.T) {
	var fired int32
	tr, _ := NewTrigger(time.Second, 3, func() { atomic.AddInt32(&fired, 1) })
	tr.Record()
	tr.Record()
	tr.Record()
	time.Sleep(50 * time.Millisecond)
	if atomic.LoadInt32(&fired) != 1 {
		t.Fatalf("expected callback fired once, got %d", atomic.LoadInt32(&fired))
	}
	if tr.Count() != 0 {
		t.Fatal("expected count reset after firing")
	}
}

func TestTrigger_Reset_ClearsEvents(t *testing.T) {
	tr, _ := NewTrigger(time.Second, 3, func() {})
	tr.Record()
	tr.Record()
	tr.Reset()
	if tr.Count() != 0 {
		t.Fatal("expected count 0 after reset")
	}
}
