package watch

import (
	"testing"
)

func TestNewMask_Empty(t *testing.T) {
	m := NewMask()
	if m.Len() != 0 {
		t.Fatalf("expected 0, got %d", m.Len())
	}
}

func TestMask_Add_Valid(t *testing.T) {
	m := NewMask()
	if err := m.Add(8080, "maintenance"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Len() != 1 {
		t.Fatalf("expected 1, got %d", m.Len())
	}
}

func TestMask_Add_InvalidPort(t *testing.T) {
	m := NewMask()
	if err := m.Add(0, "reason"); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := m.Add(65536, "reason"); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestMask_Add_EmptyReason(t *testing.T) {
	m := NewMask()
	if err := m.Add(443, ""); err == nil {
		t.Fatal("expected error for empty reason")
	}
}

func TestMask_Masked_True(t *testing.T) {
	m := NewMask()
	_ = m.Add(22, "ssh suppressed")
	if !m.Masked(22) {
		t.Fatal("expected port 22 to be masked")
	}
}

func TestMask_Masked_False(t *testing.T) {
	m := NewMask()
	if m.Masked(9090) {
		t.Fatal("expected port 9090 to not be masked")
	}
}

func TestMask_Reason_ReturnsValue(t *testing.T) {
	m := NewMask()
	_ = m.Add(3306, "db maintenance")
	if got := m.Reason(3306); got != "db maintenance" {
		t.Fatalf("expected 'db maintenance', got %q", got)
	}
}

func TestMask_Reason_Missing(t *testing.T) {
	m := NewMask()
	if got := m.Reason(1234); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestMask_Remove_ClearsPort(t *testing.T) {
	m := NewMask()
	_ = m.Add(80, "planned")
	m.Remove(80)
	if m.Masked(80) {
		t.Fatal("expected port 80 to be unmasked after Remove")
	}
}

func TestMask_Clear_RemovesAll(t *testing.T) {
	m := NewMask()
	_ = m.Add(80, "a")
	_ = m.Add(443, "b")
	m.Clear()
	if m.Len() != 0 {
		t.Fatalf("expected 0 after Clear, got %d", m.Len())
	}
}
