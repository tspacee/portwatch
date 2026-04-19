package watch

import (
	"testing"
	"time"
)

func TestNewLatch_Empty(t *testing.T) {
	l := NewLatch()
	if len(l.Snapshot()) != 0 {
		t.Fatal("expected empty latch")
	}
}

func TestLatch_Arm_Valid(t *testing.T) {
	l := NewLatch()
	ok, err := l.Arm(8080)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected true on first arm")
	}
}

func TestLatch_Arm_AlreadyLatched(t *testing.T) {
	l := NewLatch()
	l.Arm(8080)
	ok, err := l.Arm(8080)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected false on second arm")
	}
}

func TestLatch_Arm_InvalidPort(t *testing.T) {
	l := NewLatch()
	_, err := l.Arm(0)
	if err == nil {
		t.Fatal("expected error for port 0")
	}
	_, err = l.Arm(65536)
	if err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestLatch_IsArmed_True(t *testing.T) {
	l := NewLatch()
	l.Arm(443)
	if !l.IsArmed(443) {
		t.Fatal("expected port to be armed")
	}
}

func TestLatch_IsArmed_False(t *testing.T) {
	l := NewLatch()
	if l.IsArmed(443) {
		t.Fatal("expected port to not be armed")
	}
}

func TestLatch_Reset_ClearsPort(t *testing.T) {
	l := NewLatch()
	l.Arm(9090)
	l.Reset(9090)
	if l.IsArmed(9090) {
		t.Fatal("expected port to be cleared after reset")
	}
}

func TestLatch_ArmedAt_ReturnsTime(t *testing.T) {
	l := NewLatch()
	before := time.Now()
	l.Arm(3000)
	after := time.Now()
	t2, ok := l.ArmedAt(3000)
	if !ok {
		t.Fatal("expected armed time to exist")
	}
	if t2.Before(before) || t2.After(after) {
		t.Fatal("armed time outside expected range")
	}
}

func TestLatch_Snapshot_ReturnsCopy(t *testing.T) {
	l := NewLatch()
	l.Arm(80)
	l.Arm(443)
	s := l.Snapshot()
	if len(s) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(s))
	}
	delete(s, 80)
	if !l.IsArmed(80) {
		t.Fatal("snapshot mutation affected latch state")
	}
}
