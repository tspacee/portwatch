package watch

import (
	"testing"
)

func TestNewFence_EmptyPorts(t *testing.T) {
	_, err := NewFence([]int{})
	if err == nil {
		t.Fatal("expected error for empty port list")
	}
}

func TestNewFence_InvalidPort(t *testing.T) {
	_, err := NewFence([]int{0})
	if err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestNewFence_Valid(t *testing.T) {
	f, err := NewFence([]int{80, 443})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil fence")
	}
}

func TestFence_Allow_Permitted(t *testing.T) {
	f, _ := NewFence([]int{8080})
	if !f.Allow(8080) {
		t.Error("expected port 8080 to be allowed")
	}
}

func TestFence_Allow_Blocked(t *testing.T) {
	f, _ := NewFence([]int{8080})
	if f.Allow(9090) {
		t.Error("expected port 9090 to be blocked")
	}
}

func TestFence_Add_Valid(t *testing.T) {
	f, _ := NewFence([]int{80})
	if err := f.Add(443); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Allow(443) {
		t.Error("expected port 443 to be allowed after Add")
	}
}

func TestFence_Add_InvalidPort(t *testing.T) {
	f, _ := NewFence([]int{80})
	if err := f.Add(70000); err == nil {
		t.Fatal("expected error for out-of-range port")
	}
}

func TestFence_Remove_RemovesPort(t *testing.T) {
	f, _ := NewFence([]int{80, 443})
	f.Remove(80)
	if f.Allow(80) {
		t.Error("expected port 80 to be removed")
	}
}

func TestFence_Snapshot_ReturnsCopy(t *testing.T) {
	f, _ := NewFence([]int{80, 443})
	snap := f.Snapshot()
	if len(snap) != 2 {
		t.Errorf("expected 2 ports, got %d", len(snap))
	}
	snap[0] = 9999
	if f.Allow(9999) {
		t.Error("snapshot mutation should not affect fence")
	}
}
