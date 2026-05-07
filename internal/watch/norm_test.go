package watch

import (
	"testing"
)

func TestNewNorm_Empty(t *testing.T) {
	n := NewNorm()
	if n.Size() != 0 {
		t.Fatalf("expected size 0, got %d", n.Size())
	}
	if n.IsFrozen() {
		t.Fatal("expected norm to be unfrozen initially")
	}
}

func TestNorm_Learn_Valid(t *testing.T) {
	n := NewNorm()
	if err := n.Learn(80); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n.Size() != 1 {
		t.Fatalf("expected size 1, got %d", n.Size())
	}
}

func TestNorm_Learn_InvalidPort(t *testing.T) {
	n := NewNorm()
	if err := n.Learn(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := n.Learn(65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestNorm_Learn_AfterFreeze_ReturnsError(t *testing.T) {
	n := NewNorm()
	_ = n.Learn(443)
	n.Freeze()
	if err := n.Learn(8080); err == nil {
		t.Fatal("expected error after freeze")
	}
	if n.Size() != 1 {
		t.Fatalf("expected size to remain 1, got %d", n.Size())
	}
}

func TestNorm_Conforms_AllInBaseline(t *testing.T) {
	n := NewNorm()
	_ = n.Learn(80)
	_ = n.Learn(443)
	if !n.Conforms([]int{80, 443}) {
		t.Fatal("expected ports to conform")
	}
}

func TestNorm_Conforms_WithDeviation(t *testing.T) {
	n := NewNorm()
	_ = n.Learn(80)
	if n.Conforms([]int{80, 9090}) {
		t.Fatal("expected non-conformance due to port 9090")
	}
}

func TestNorm_Deviations_ReturnsUnexpected(t *testing.T) {
	n := NewNorm()
	_ = n.Learn(80)
	_ = n.Learn(443)
	devs := n.Deviations([]int{80, 443, 9090, 3000})
	if len(devs) != 2 {
		t.Fatalf("expected 2 deviations, got %d", len(devs))
	}
}

func TestNorm_Deviations_EmptyWhenAllConform(t *testing.T) {
	n := NewNorm()
	_ = n.Learn(22)
	devs := n.Deviations([]int{22})
	if len(devs) != 0 {
		t.Fatalf("expected no deviations, got %d", len(devs))
	}
}

func TestNorm_Freeze_SetsFlag(t *testing.T) {
	n := NewNorm()
	n.Freeze()
	if !n.IsFrozen() {
		t.Fatal("expected norm to be frozen")
	}
}
