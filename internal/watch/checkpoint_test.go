package watch

import (
	"testing"
	"time"
)

func TestNewCheckpoint_InitialState(t *testing.T) {
	cp := NewCheckpoint()
	if cp.Seq() != 0 {
		t.Fatalf("expected seq 0, got %d", cp.Seq())
	}
	if !cp.Timestamp().IsZero() {
		t.Fatal("expected zero timestamp")
	}
	if len(cp.Ports()) != 0 {
		t.Fatal("expected empty ports")
	}
}

func TestCheckpoint_Commit_Valid(t *testing.T) {
	cp := NewCheckpoint()
	now := time.Now()
	ports := []int{80, 443, 8080}

	if err := cp.Commit(1, now, ports); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cp.Seq() != 1 {
		t.Fatalf("expected seq 1, got %d", cp.Seq())
	}
	if !cp.Timestamp().Equal(now) {
		t.Fatal("timestamp mismatch")
	}
}

func TestCheckpoint_Commit_ZeroSeq_ReturnsError(t *testing.T) {
	cp := NewCheckpoint()
	err := cp.Commit(0, time.Now(), []int{80})
	if err == nil {
		t.Fatal("expected error for seq=0")
	}
}

func TestCheckpoint_Ports_ReturnsCopy(t *testing.T) {
	cp := NewCheckpoint()
	original := []int{22, 80}
	_ = cp.Commit(1, time.Now(), original)

	got := cp.Ports()
	got[0] = 9999

	if cp.Ports()[0] == 9999 {
		t.Fatal("Ports should return a copy, not the internal slice")
	}
}

func TestCheckpoint_Commit_EmptyPorts(t *testing.T) {
	cp := NewCheckpoint()
	if err := cp.Commit(2, time.Now(), []int{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cp.Ports()) != 0 {
		t.Fatal("expected empty ports slice")
	}
}

func TestCheckpoint_Reset_ClearsState(t *testing.T) {
	cp := NewCheckpoint()
	_ = cp.Commit(5, time.Now(), []int{80, 443})
	cp.Reset()

	if cp.Seq() != 0 {
		t.Fatalf("expected seq 0 after reset, got %d", cp.Seq())
	}
	if !cp.Timestamp().IsZero() {
		t.Fatal("expected zero timestamp after reset")
	}
	if len(cp.Ports()) != 0 {
		t.Fatal("expected empty ports after reset")
	}
}

func TestCheckpoint_Commit_OverwritesPrevious(t *testing.T) {
	cp := NewCheckpoint()
	_ = cp.Commit(1, time.Now(), []int{80})
	later := time.Now().Add(time.Second)
	_ = cp.Commit(2, later, []int{443, 8080})

	if cp.Seq() != 2 {
		t.Fatalf("expected seq 2, got %d", cp.Seq())
	}
	if len(cp.Ports()) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(cp.Ports()))
	}
}
