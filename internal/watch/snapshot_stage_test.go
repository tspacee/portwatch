package watch

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func buildManager(t *testing.T) *snapshot.Manager {
	t.Helper()
	dir := t.TempDir()
	m, err := snapshot.NewManager(filepath.Join(dir, "snap.json"))
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	return m
}

func TestNewSnapshotStage_NilManager(t *testing.T) {
	_, err := NewSnapshotStage(nil)
	if err == nil {
		t.Fatal("expected error for nil manager")
	}
}

func TestNewSnapshotStage_Valid(t *testing.T) {
	m := buildManager(t)
	s, err := NewSnapshotStage(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil stage")
	}
}

func TestSnapshotStage_Run_FirstScan_ReturnsAllAdded(t *testing.T) {
	m := buildManager(t)
	s, _ := NewSnapshotStage(m)

	ports := []int{80, 443, 8080}
	added, err := s.Run(context.Background(), ports)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(added) != len(ports) {
		t.Errorf("expected %d added, got %d", len(ports), len(added))
	}
}

func TestSnapshotStage_Run_NoChanges_ReturnsEmpty(t *testing.T) {
	m := buildManager(t)
	s, _ := NewSnapshotStage(m)
	ports := []int{80, 443}

	s.Run(context.Background(), ports) //nolint:errcheck first scan
	added, err := s.Run(context.Background(), ports)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(added) != 0 {
		t.Errorf("expected 0 added, got %d", len(added))
	}
}

func TestSnapshotStage_Run_ContextCancelled(t *testing.T) {
	m := buildManager(t)
	s, _ := NewSnapshotStage(m)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := s.Run(ctx, []int{80})
	if err == nil {
		t.Fatal("expected context error")
	}
}

func TestSnapshotStage_Run_SaveError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	m, _ := snapshot.NewManager(path)
	s, _ := NewSnapshotStage(m)

	// Make the directory read-only so Save fails.
	os.Chmod(dir, 0o444)
	t.Cleanup(func() { os.Chmod(dir, 0o755) })

	_, err := s.Run(context.Background(), []int{22})
	if err == nil {
		t.Fatal("expected save error")
	}
}
