package watch

import (
	"testing"
)

func TestNewParity_Empty(t *testing.T) {
	p := NewParity()
	if p == nil {
		t.Fatal("expected non-nil Parity")
	}
	if p.Mismatches() != 0 {
		t.Errorf("expected 0 mismatches, got %d", p.Mismatches())
	}
}

func TestParity_SetReference_InvalidPort(t *testing.T) {
	p := NewParity()
	if err := p.SetReference([]int{0}); err == nil {
		t.Error("expected error for port 0")
	}
	if err := p.SetReference([]int{65536}); err == nil {
		t.Error("expected error for port 65536")
	}
}

func TestParity_SetReference_Valid(t *testing.T) {
	p := NewParity()
	if err := p.SetReference([]int{80, 443}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParity_InSync_EmptyBothSides(t *testing.T) {
	p := NewParity()
	p.Compare([]int{})
	if !p.InSync() {
		t.Error("expected in-sync when both sides are empty")
	}
}

func TestParity_InSync_MatchingPorts(t *testing.T) {
	p := NewParity()
	_ = p.SetReference([]int{80, 443})
	p.Compare([]int{443, 80})
	if !p.InSync() {
		t.Errorf("expected in-sync, got %d mismatches", p.Mismatches())
	}
}

func TestParity_Mismatches_ExtraInCurrent(t *testing.T) {
	p := NewParity()
	_ = p.SetReference([]int{80})
	p.Compare([]int{80, 8080})
	if p.Mismatches() != 1 {
		t.Errorf("expected 1 mismatch, got %d", p.Mismatches())
	}
}

func TestParity_Mismatches_MissingFromCurrent(t *testing.T) {
	p := NewParity()
	_ = p.SetReference([]int{80, 443})
	p.Compare([]int{80})
	if p.Mismatches() != 1 {
		t.Errorf("expected 1 mismatch, got %d", p.Mismatches())
	}
}

func TestParity_Mismatches_BothDiverge(t *testing.T) {
	p := NewParity()
	_ = p.SetReference([]int{80, 443})
	p.Compare([]int{22, 8080})
	if p.Mismatches() != 4 {
		t.Errorf("expected 4 mismatches, got %d", p.Mismatches())
	}
}

func TestParity_Compare_UpdatesOnSubsequentCall(t *testing.T) {
	p := NewParity()
	_ = p.SetReference([]int{80})
	p.Compare([]int{22})
	if p.InSync() {
		t.Error("expected out-of-sync")
	}
	p.Compare([]int{80})
	if !p.InSync() {
		t.Error("expected in-sync after corrective compare")
	}
}
