package watch

import (
	"testing"
)

func TestNewVeto_Empty(t *testing.T) {
	v := NewVeto()
	if v.Len() != 0 {
		t.Fatalf("expected 0 vetoed ports, got %d", v.Len())
	}
}

func TestVeto_Block_Valid(t *testing.T) {
	v := NewVeto()
	if err := v.Block(8080, "testing"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Len() != 1 {
		t.Fatalf("expected 1 vetoed port, got %d", v.Len())
	}
}

func TestVeto_Block_InvalidPort(t *testing.T) {
	v := NewVeto()
	if err := v.Block(0, "reason"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := v.Block(65536, "reason"); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestVeto_Block_EmptyReason(t *testing.T) {
	v := NewVeto()
	if err := v.Block(443, ""); err == nil {
		t.Fatal("expected error for empty reason")
	}
}

func TestVeto_IsVetoed_True(t *testing.T) {
	v := NewVeto()
	_ = v.Block(22, "ssh blocked")
	if !v.IsVetoed(22) {
		t.Fatal("expected port 22 to be vetoed")
	}
}

func TestVeto_IsVetoed_False(t *testing.T) {
	v := NewVeto()
	if v.IsVetoed(22) {
		t.Fatal("expected port 22 to not be vetoed")
	}
}

func TestVeto_Reason_ReturnsValue(t *testing.T) {
	v := NewVeto()
	_ = v.Block(3306, "db port")
	if got := v.Reason(3306); got != "db port" {
		t.Fatalf("expected 'db port', got %q", got)
	}
}

func TestVeto_Reason_Missing_ReturnsEmpty(t *testing.T) {
	v := NewVeto()
	if got := v.Reason(9999); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestVeto_Lift_RemovesPort(t *testing.T) {
	v := NewVeto()
	_ = v.Block(80, "http")
	if !v.Lift(80) {
		t.Fatal("expected Lift to return true")
	}
	if v.IsVetoed(80) {
		t.Fatal("expected port 80 to no longer be vetoed")
	}
}

func TestVeto_Lift_NotPresent_ReturnsFalse(t *testing.T) {
	v := NewVeto()
	if v.Lift(443) {
		t.Fatal("expected Lift to return false for non-vetoed port")
	}
}
