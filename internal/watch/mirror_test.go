package watch

import (
	"testing"
)

func TestNewMirror_InitialState(t *testing.T) {
	m := NewMirror()
	if m.Len() != 0 {
		t.Fatalf("expected len 0, got %d", m.Len())
	}
	if m.Seq() != 0 {
		t.Fatalf("expected seq 0, got %d", m.Seq())
	}
	if m.Ports() != nil {
		t.Fatal("expected nil ports on empty mirror")
	}
}

func TestMirror_Update_Valid(t *testing.T) {
	m := NewMirror()
	ports := []int{80, 443, 8080}
	if err := m.Update(1, ports); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Seq() != 1 {
		t.Fatalf("expected seq 1, got %d", m.Seq())
	}
	if m.Len() != 3 {
		t.Fatalf("expected len 3, got %d", m.Len())
	}
}

func TestMirror_Update_ReturnsCopy(t *testing.T) {
	m := NewMirror()
	ports := []int{22, 80}
	_ = m.Update(1, ports)

	// mutate original slice
	ports[0] = 9999

	got := m.Ports()
	if got[0] == 9999 {
		t.Fatal("mirror should store a copy, not a reference")
	}
}

func TestMirror_Ports_ReturnsCopy(t *testing.T) {
	m := NewMirror()
	_ = m.Update(1, []int{22, 80})

	got := m.Ports()
	got[0] = 9999

	got2 := m.Ports()
	if got2[0] == 9999 {
		t.Fatal("Ports() should return a copy, not a reference")
	}
}

func TestMirror_Update_StaleSeq_ReturnsError(t *testing.T) {
	m := NewMirror()
	_ = m.Update(5, []int{80})

	if err := m.Update(3, []int{443}); err == nil {
		t.Fatal("expected error for non-monotonic sequence number")
	}
}

func TestMirror_Update_SameSeq_ReturnsError(t *testing.T) {
	m := NewMirror()
	_ = m.Update(2, []int{80})

	if err := m.Update(2, []int{443}); err == nil {
		t.Fatal("expected error for duplicate sequence number")
	}
}

func TestMirror_Update_EmptyPorts(t *testing.T) {
	m := NewMirror()
	if err := m.Update(1, []int{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Len() != 0 {
		t.Fatalf("expected len 0, got %d", m.Len())
	}
}
