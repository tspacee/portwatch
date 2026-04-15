package watch

import (
	"testing"
)

func TestNewRangeFilter_Valid(t *testing.T) {
	f, err := NewRangeFilter(1024, 9000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Min != 1024 || f.Max != 9000 {
		t.Errorf("unexpected bounds: %d-%d", f.Min, f.Max)
	}
}

func TestNewRangeFilter_InvalidBounds(t *testing.T) {
	if _, err := NewRangeFilter(0, 100); err == nil {
		t.Error("expected error for min < 1")
	}
	if _, err := NewRangeFilter(100, 70000); err == nil {
		t.Error("expected error for max > 65535")
	}
	if _, err := NewRangeFilter(9000, 1024); err == nil {
		t.Error("expected error for min > max")
	}
}

func TestRangeFilter_Allow(t *testing.T) {
	f, _ := NewRangeFilter(1000, 2000)
	if !f.Allow(1000) || !f.Allow(1500) || !f.Allow(2000) {
		t.Error("expected ports within range to be allowed")
	}
	if f.Allow(999) || f.Allow(2001) {
		t.Error("expected ports outside range to be blocked")
	}
}

func TestExcludeFilter_Allow(t *testing.T) {
	f := NewExcludeFilter([]int{22, 80, 443})
	if f.Allow(22) || f.Allow(80) || f.Allow(443) {
		t.Error("expected excluded ports to be blocked")
	}
	if !f.Allow(8080) {
		t.Error("expected non-excluded port to be allowed")
	}
}

func TestExcludeFilter_EmptyList(t *testing.T) {
	f := NewExcludeFilter(nil)
	if !f.Allow(22) {
		t.Error("expected all ports allowed when exclusion list is empty")
	}
}

func TestNewChainFilter_NilFilter(t *testing.T) {
	_, err := NewChainFilter(nil)
	if err == nil {
		t.Error("expected error for nil filter in chain")
	}
}

func TestChainFilter_AllAllow(t *testing.T) {
	rf, _ := NewRangeFilter(1000, 9000)
	ef := NewExcludeFilter([]int{22})
	chain, err := NewChainFilter(rf, ef)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !chain.Allow(8080) {
		t.Error("expected 8080 to be allowed")
	}
}

func TestChainFilter_OneBlocks(t *testing.T) {
	rf, _ := NewRangeFilter(1000, 9000)
	ef := NewExcludeFilter([]int{8080})
	chain, _ := NewChainFilter(rf, ef)
	if chain.Allow(8080) {
		t.Error("expected 8080 to be blocked by exclude filter")
	}
	if chain.Allow(80) {
		t.Error("expected 80 to be blocked by range filter")
	}
}
