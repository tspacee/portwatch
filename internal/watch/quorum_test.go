package watch

import (
	"testing"
)

func TestNewQuorum_InvalidMinVotes(t *testing.T) {
	_, err := NewQuorum(0)
	if err == nil {
		t.Fatal("expected error for minVotes=0")
	}
}

func TestNewQuorum_Valid(t *testing.T) {
	q, err := NewQuorum(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q == nil {
		t.Fatal("expected non-nil Quorum")
	}
}

func TestQuorum_Vote_InvalidPort(t *testing.T) {
	q, _ := NewQuorum(1)
	_, err := q.Vote(0, "scanner-a")
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestQuorum_Vote_EmptySource(t *testing.T) {
	q, _ := NewQuorum(1)
	_, err := q.Vote(80, "")
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestQuorum_Vote_SingleSource_BelowQuorum(t *testing.T) {
	q, _ := NewQuorum(2)
	ok, err := q.Vote(80, "scanner-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected quorum not reached with one vote")
	}
}

func TestQuorum_Vote_ReachesQuorum(t *testing.T) {
	q, _ := NewQuorum(2)
	q.Vote(443, "scanner-a")
	ok, err := q.Vote(443, "scanner-b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected quorum reached with two votes")
	}
}

func TestQuorum_Vote_DuplicateSource_DoesNotDoubleCount(t *testing.T) {
	q, _ := NewQuorum(2)
	q.Vote(8080, "scanner-a")
	ok, _ := q.Vote(8080, "scanner-a")
	if ok {
		t.Fatal("duplicate source should not satisfy quorum")
	}
	if q.VoteCount(8080) != 1 {
		t.Fatalf("expected vote count 1, got %d", q.VoteCount(8080))
	}
}

func TestQuorum_Confirmed_False(t *testing.T) {
	q, _ := NewQuorum(3)
	q.Vote(22, "a")
	if q.Confirmed(22) {
		t.Fatal("expected port not confirmed")
	}
}

func TestQuorum_Reset_ClearsVotes(t *testing.T) {
	q, _ := NewQuorum(1)
	q.Vote(9000, "scanner-a")
	q.Reset(9000)
	if q.Confirmed(9000) {
		t.Fatal("expected port not confirmed after reset")
	}
	if q.VoteCount(9000) != 0 {
		t.Fatalf("expected vote count 0 after reset, got %d", q.VoteCount(9000))
	}
}
