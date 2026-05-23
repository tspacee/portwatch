package watch

import (
	"testing"
)

func TestNewNomad_InvalidMaxMove(t *testing.T) {
	_, err := NewNomad(0)
	if err == nil {
		t.Fatal("expected error for maxMove=0")
	}
}

func TestNewNomad_Valid(t *testing.T) {
	n, err := NewNomad(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil Nomad")
	}
}

func TestNomad_MoveCount_NeverObserved(t *testing.T) {
	n, _ := NewNomad(1)
	if got := n.MoveCount(80); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestNomad_Observe_FirstAppearance_NoMove(t *testing.T) {
	n, _ := NewNomad(1)
	n.Observe([]int{80})
	if got := n.MoveCount(80); got != 0 {
		t.Errorf("first appearance should not count as move, got %d", got)
	}
}

func TestNomad_Observe_Reappearance_IncrementsMoveCount(t *testing.T) {
	n, _ := NewNomad(3)
	n.Observe([]int{80})  // seen
	n.Observe([]int{})    // absent
	n.Observe([]int{80})  // reappears -> move count = 1
	if got := n.MoveCount(80); got != 1 {
		t.Errorf("expected move count 1, got %d", got)
	}
}

func TestNomad_IsNomadic_BelowThreshold(t *testing.T) {
	n, _ := NewNomad(2)
	n.Observe([]int{443})
	n.Observe([]int{})
	n.Observe([]int{443}) // move count = 1, threshold = 2
	if n.IsNomadic(443) {
		t.Error("port should not be nomadic below threshold")
	}
}

func TestNomad_IsNomadic_ExceedsThreshold(t *testing.T) {
	n, _ := NewNomad(1)
	for i := 0; i < 3; i++ {
		n.Observe([]int{8080})
		n.Observe([]int{})
	}
	n.Observe([]int{8080}) // move count = 3 > maxMove=1
	if !n.IsNomadic(8080) {
		t.Error("port should be nomadic after exceeding threshold")
	}
}

func TestNomad_Reset_ClearsState(t *testing.T) {
	n, _ := NewNomad(1)
	n.Observe([]int{22})
	n.Observe([]int{})
	n.Observe([]int{22})
	n.Reset()
	if got := n.MoveCount(22); got != 0 {
		t.Errorf("expected 0 after reset, got %d", got)
	}
}
