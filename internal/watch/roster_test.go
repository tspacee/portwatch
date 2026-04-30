package watch

import (
	"testing"
)

func TestNewRoster_Empty(t *testing.T) {
	r := NewRoster()
	if r.Len() != 0 {
		t.Fatalf("expected 0, got %d", r.Len())
	}
}

func TestRoster_Enroll_Valid(t *testing.T) {
	r := NewRoster()
	if err := r.Enroll(80); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.IsActive(80) {
		t.Fatal("expected port 80 to be active")
	}
}

func TestRoster_Enroll_InvalidPort(t *testing.T) {
	r := NewRoster()
	if err := r.Enroll(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := r.Enroll(65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestRoster_Enroll_Idempotent(t *testing.T) {
	r := NewRoster()
	_ = r.Enroll(443)
	_ = r.Enroll(443)
	if r.Len() != 1 {
		t.Fatalf("expected 1, got %d", r.Len())
	}
}

func TestRoster_Withdraw_RemovesPort(t *testing.T) {
	r := NewRoster()
	_ = r.Enroll(22)
	r.Withdraw(22)
	if r.IsActive(22) {
		t.Fatal("expected port 22 to be inactive after withdraw")
	}
}

func TestRoster_Withdraw_NonExistent_NoOp(t *testing.T) {
	r := NewRoster()
	r.Withdraw(9999) // should not panic
	if r.Len() != 0 {
		t.Fatalf("expected 0, got %d", r.Len())
	}
}

func TestRoster_Members_ReturnsSorted(t *testing.T) {
	r := NewRoster()
	ports := []int{443, 80, 22, 8080}
	for _, p := range ports {
		_ = r.Enroll(p)
	}
	members := r.Members()
	expected := []int{22, 80, 443, 8080}
	if len(members) != len(expected) {
		t.Fatalf("expected %d members, got %d", len(expected), len(members))
	}
	for i, p := range expected {
		if members[i] != p {
			t.Errorf("index %d: expected %d, got %d", i, p, members[i])
		}
	}
}

func TestRoster_Members_ReturnsCopy(t *testing.T) {
	r := NewRoster()
	_ = r.Enroll(80)
	m1 := r.Members()
	m1[0] = 9999
	m2 := r.Members()
	if m2[0] == 9999 {
		t.Fatal("Members should return a copy, not a reference")
	}
}
