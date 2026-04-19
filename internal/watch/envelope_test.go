package watch

import (
	"testing"
	"time"
)

func TestNewEnvelopeBuilder_EmptySource(t *testing.T) {
	_, err := NewEnvelopeBuilder("")
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestNewEnvelopeBuilder_Valid(t *testing.T) {
	b, err := NewEnvelopeBuilder("host1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil builder")
	}
}

func TestEnvelopeBuilder_Wrap_IncrementsSeq(t *testing.T) {
	b, _ := NewEnvelopeBuilder("host1")
	e1 := b.Wrap([]int{80})
	e2 := b.Wrap([]int{443})
	if e1.Seq != 1 {
		t.Errorf("expected seq 1, got %d", e1.Seq)
	}
	if e2.Seq != 2 {
		t.Errorf("expected seq 2, got %d", e2.Seq)
	}
}

func TestEnvelopeBuilder_Wrap_SetsSource(t *testing.T) {
	b, _ := NewEnvelopeBuilder("scanner-a")
	env := b.Wrap(nil)
	if env.Source != "scanner-a" {
		t.Errorf("expected source scanner-a, got %s", env.Source)
	}
}

func TestEnvelopeBuilder_Wrap_CopiesPorts(t *testing.T) {
	b, _ := NewEnvelopeBuilder("host1")
	orig := []int{22, 80}
	env := b.Wrap(orig)
	orig[0] = 9999
	if env.Ports[0] == 9999 {
		t.Error("envelope should not share slice with original")
	}
}

func TestEnvelopeBuilder_Wrap_SetsScannedAt(t *testing.T) {
	b, _ := NewEnvelopeBuilder("host1")
	before := time.Now()
	env := b.Wrap([]int{80})
	after := time.Now()
	if env.ScannedAt.Before(before) || env.ScannedAt.After(after) {
		t.Error("ScannedAt outside expected range")
	}
}

func TestEnvelopeBuilder_CurrentSeq_ReflectsState(t *testing.T) {
	b, _ := NewEnvelopeBuilder("host1")
	if b.CurrentSeq() != 0 {
		t.Errorf("expected 0 before any wrap, got %d", b.CurrentSeq())
	}
	b.Wrap([]int{})
	if b.CurrentSeq() != 1 {
		t.Errorf("expected 1 after one wrap, got %d", b.CurrentSeq())
	}
}
