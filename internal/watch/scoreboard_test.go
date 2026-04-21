package watch

import (
	"testing"
)

func TestNewScoreboard_Empty(t *testing.T) {
	sb := NewScoreboard()
	if sb.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", sb.Len())
	}
}

func TestScoreboard_Record_Valid(t *testing.T) {
	sb := NewScoreboard()
	if err := sb.Record(8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sb.Score(8080) != 1 {
		t.Fatalf("expected score 1, got %d", sb.Score(8080))
	}
}

func TestScoreboard_Record_InvalidPort(t *testing.T) {
	sb := NewScoreboard()
	if err := sb.Record(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := sb.Record(65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestScoreboard_Record_Accumulates(t *testing.T) {
	sb := NewScoreboard()
	for i := 0; i < 5; i++ {
		_ = sb.Record(443)
	}
	if sb.Score(443) != 5 {
		t.Fatalf("expected score 5, got %d", sb.Score(443))
	}
}

func TestScoreboard_Score_Missing(t *testing.T) {
	sb := NewScoreboard()
	if sb.Score(9999) != 0 {
		t.Fatal("expected 0 for untracked port")
	}
}

func TestScoreboard_Reset_ClearsCount(t *testing.T) {
	sb := NewScoreboard()
	_ = sb.Record(22)
	_ = sb.Record(22)
	sb.Reset(22)
	if sb.Score(22) != 0 {
		t.Fatalf("expected 0 after reset, got %d", sb.Score(22))
	}
	if sb.Len() != 0 {
		t.Fatalf("expected len 0 after reset, got %d", sb.Len())
	}
}

func TestScoreboard_Snapshot_ReturnsCopy(t *testing.T) {
	sb := NewScoreboard()
	_ = sb.Record(80)
	_ = sb.Record(443)
	snap := sb.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
	// mutating the snapshot must not affect the scoreboard
	snap[80] = 999
	if sb.Score(80) != 1 {
		t.Fatal("snapshot mutation affected scoreboard")
	}
}

func TestScoreboard_Len_TracksEntries(t *testing.T) {
	sb := NewScoreboard()
	_ = sb.Record(1)
	_ = sb.Record(2)
	_ = sb.Record(1) // duplicate port, same entry
	if sb.Len() != 2 {
		t.Fatalf("expected 2 unique ports, got %d", sb.Len())
	}
}
