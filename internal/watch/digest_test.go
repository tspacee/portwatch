package watch

import (
	"testing"
)

func TestNewDigest_InitialLastEmpty(t *testing.T) {
	d := NewDigest()
	if d.Last() != "" {
		t.Fatalf("expected empty last digest, got %q", d.Last())
	}
}

func TestDigest_Compute_Deterministic(t *testing.T) {
	d := NewDigest()
	ports := []int{443, 80, 22}
	a := d.Compute(ports)
	b := d.Compute(ports)
	if a != b {
		t.Fatalf("expected deterministic digest, got %q and %q", a, b)
	}
}

func TestDigest_Compute_OrderIndependent(t *testing.T) {
	d := NewDigest()
	a := d.Compute([]int{80, 443, 22})
	b := d.Compute([]int{22, 80, 443})
	if a != b {
		t.Fatalf("expected order-independent digest, got %q and %q", a, b)
	}
}

func TestDigest_Changed_FirstCall_ReturnsTrue(t *testing.T) {
	d := NewDigest()
	if !d.Changed([]int{80, 443}) {
		t.Fatal("expected Changed to return true on first call")
	}
}

func TestDigest_Changed_SamePorts_ReturnsFalse(t *testing.T) {
	d := NewDigest()
	ports := []int{80, 443}
	d.Changed(ports)
	if d.Changed(ports) {
		t.Fatal("expected Changed to return false for identical ports")
	}
}

func TestDigest_Changed_DifferentPorts_ReturnsTrue(t *testing.T) {
	d := NewDigest()
	d.Changed([]int{80})
	if !d.Changed([]int{80, 443}) {
		t.Fatal("expected Changed to return true when ports differ")
	}
}

func TestDigest_Reset_ClearsLast(t *testing.T) {
	d := NewDigest()
	d.Changed([]int{80, 443})
	if d.Last() == "" {
		t.Fatal("expected non-empty last digest after Changed")
	}
	d.Reset()
	if d.Last() != "" {
		t.Fatalf("expected empty last digest after Reset, got %q", d.Last())
	}
}

func TestDigest_Last_ReflectsLatestCompute(t *testing.T) {
	d := NewDigest()
	ports := []int{22, 80}
	d.Changed(ports)
	expected := d.Compute(ports)
	if d.Last() != expected {
		t.Fatalf("expected Last()=%q, got %q", expected, d.Last())
	}
}
