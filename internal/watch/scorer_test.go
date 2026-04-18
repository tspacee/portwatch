package watch

import (
	"testing"
)

func TestNewScorer_NegativeDefault(t *testing.T) {
	_, err := NewScorer(-1)
	if err == nil {
		t.Fatal("expected error for negative default score")
	}
}

func TestNewScorer_Valid(t *testing.T) {
	s, err := NewScorer(1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil scorer")
	}
}

func TestScorer_Score_ReturnsDefault(t *testing.T) {
	s, _ := NewScorer(2.5)
	if got := s.Score(80); got != 2.5 {
		t.Errorf("expected 2.5, got %f", got)
	}
}

func TestScorer_SetWeight_Valid(t *testing.T) {
	s, _ := NewScorer(1.0)
	if err := s.SetWeight(443, 9.0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.Score(443); got != 9.0 {
		t.Errorf("expected 9.0, got %f", got)
	}
}

func TestScorer_SetWeight_InvalidPort(t *testing.T) {
	s, _ := NewScorer(1.0)
	if err := s.SetWeight(0, 5.0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := s.SetWeight(65536, 5.0); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestScorer_SetWeight_NegativeWeight(t *testing.T) {
	s, _ := NewScorer(1.0)
	if err := s.SetWeight(80, -1.0); err == nil {
		t.Fatal("expected error for negative weight")
	}
}

func TestScorer_ScoreAll(t *testing.T) {
	s, _ := NewScorer(1.0)
	_ = s.SetWeight(22, 8.0)
	_ = s.SetWeight(80, 3.0)

	result := s.ScoreAll([]int{22, 80, 443})
	if result[22] != 8.0 {
		t.Errorf("expected 8.0 for port 22, got %f", result[22])
	}
	if result[80] != 3.0 {
		t.Errorf("expected 3.0 for port 80, got %f", result[80])
	}
	if result[443] != 1.0 {
		t.Errorf("expected default 1.0 for port 443, got %f", result[443])
	}
}

func TestScorer_ScoreAll_EmptyPorts(t *testing.T) {
	s, _ := NewScorer(1.0)
	result := s.ScoreAll([]int{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
