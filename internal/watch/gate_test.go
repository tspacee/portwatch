package watch

import (
	"testing"
)

func TestNewGate_ClosedByDefault(t *testing.T) {
	g := NewGate(false)
	if g == nil {
		t.Fatal("expected non-nil gate")
	}
	if g.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", g.Len())
	}
}

func TestNewGate_OpenMode_AllowsAll(t *testing.T) {
	g := NewGate(true)
	if !g.Pass(8080) {
		t.Error("expected open gate to pass all ports")
	}
}

func TestGate_Allow_Valid(t *testing.T) {
	g := NewGate(false)
	if err := g.Allow(443); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !g.Pass(443) {
		t.Error("expected port 443 to pass")
	}
}

func TestGate_Allow_InvalidPort(t *testing.T) {
	g := NewGate(false)
	if err := g.Allow(0); err == nil {
		t.Error("expected error for port 0")
	}
	if err := g.Allow(65536); err == nil {
		t.Error("expected error for port 65536")
	}
}

func TestGate_Pass_BlockedWhenClosed(t *testing.T) {
	g := NewGate(false)
	if g.Pass(80) {
		t.Error("expected port 80 to be blocked on closed gate")
	}
}

func TestGate_Deny_RemovesPort(t *testing.T) {
	g := NewGate(false)
	_ = g.Allow(80)
	if err := g.Deny(80); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if g.Pass(80) {
		t.Error("expected port 80 to be denied after Deny()")
	}
}

func TestGate_Deny_InvalidPort(t *testing.T) {
	g := NewGate(false)
	if err := g.Deny(0); err == nil {
		t.Error("expected error for invalid port")
	}
}

func TestGate_Len_Tracks(t *testing.T) {
	g := NewGate(false)
	_ = g.Allow(80)
	_ = g.Allow(443)
	if g.Len() != 2 {
		t.Errorf("expected 2, got %d", g.Len())
	}
	_ = g.Deny(80)
	if g.Len() != 1 {
		t.Errorf("expected 1 after deny, got %d", g.Len())
	}
}
