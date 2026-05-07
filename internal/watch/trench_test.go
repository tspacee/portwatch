package watch

import (
	"testing"
	"time"
)

func TestNewTrench_InvalidWindow(t *testing.T) {
	_, err := NewTrench(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewTrench_Valid(t *testing.T) {
	tr, err := NewTrench(time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil Trench")
	}
}

func TestTrench_MarkAbsent_InvalidPort(t *testing.T) {
	tr, _ := NewTrench(time.Second)
	if err := tr.MarkAbsent(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := tr.MarkAbsent(65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestTrench_MarkAbsent_Valid(t *testing.T) {
	tr, _ := NewTrench(time.Second)
	if err := tr.MarkAbsent(8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Len() != 1 {
		t.Fatalf("expected len 1, got %d", tr.Len())
	}
}

func TestTrench_MarkAbsent_Idempotent(t *testing.T) {
	tr, _ := NewTrench(time.Second)
	_ = tr.MarkAbsent(9090)
	_ = tr.MarkAbsent(9090)
	if tr.Len() != 1 {
		t.Fatalf("expected len 1 after duplicate mark, got %d", tr.Len())
	}
}

func TestTrench_Entrenched_NotYet(t *testing.T) {
	tr, _ := NewTrench(time.Hour)
	_ = tr.MarkAbsent(443)
	if tr.Entrenched(443) {
		t.Fatal("port should not be entrenched before window elapses")
	}
}

func TestTrench_Entrenched_AfterWindow(t *testing.T) {
	tr, _ := NewTrench(time.Millisecond)
	_ = tr.MarkAbsent(22)
	time.Sleep(5 * time.Millisecond)
	if !tr.Entrenched(22) {
		t.Fatal("port should be entrenched after window elapses")
	}
}

func TestTrench_Observe_ClearsReappearedPort(t *testing.T) {
	tr, _ := NewTrench(time.Hour)
	_ = tr.MarkAbsent(3000)
	tr.Observe([]int{3000})
	if tr.Len() != 0 {
		t.Fatalf("expected len 0 after port reappeared, got %d", tr.Len())
	}
}

func TestTrench_Forget_RemovesPort(t *testing.T) {
	tr, _ := NewTrench(time.Hour)
	_ = tr.MarkAbsent(5000)
	tr.Forget(5000)
	if tr.Len() != 0 {
		t.Fatalf("expected len 0 after forget, got %d", tr.Len())
	}
}
