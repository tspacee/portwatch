package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func TestNew_CopiesPorts(t *testing.T) {
	ports := []int{80, 443, 8080}
	s := snapshot.New(ports)
	ports[0] = 9999 // mutate original
	if s.Ports[0] == 9999 {
		t.Error("expected snapshot to be independent of original slice")
	}
	if s.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := snapshot.New([]int{22, 80, 443})
	if err := orig.Save(path); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if len(loaded.Ports) != len(orig.Ports) {
		t.Errorf("port count mismatch: got %d want %d", len(loaded.Ports), len(orig.Ports))
	}
}

func TestLoad_MissingFile_ReturnsNil(t *testing.T) {
	s, err := snapshot.Load("/nonexistent/path/snap.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s != nil {
		t.Error("expected nil snapshot for missing file")
	}
}

func TestLoad_InvalidJSON_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)

	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestDiff_AddedAndRemoved(t *testing.T) {
	prev := snapshot.New([]int{80, 443, 22})
	curr := snapshot.New([]int{80, 8080})

	added, removed := snapshot.Diff(prev, curr)

	if len(added) != 1 || added[0] != 8080 {
		t.Errorf("expected added=[8080], got %v", added)
	}
	if len(removed) != 2 {
		t.Errorf("expected 2 removed ports, got %v", removed)
	}
}

func TestDiff_NilPrev_AllAdded(t *testing.T) {
	curr := snapshot.New([]int{80, 443})
	added, removed := snapshot.Diff(nil, curr)

	if len(added) != 2 {
		t.Errorf("expected 2 added ports, got %v", added)
	}
	if len(removed) != 0 {
		t.Errorf("expected no removed ports, got %v", removed)
	}
}
