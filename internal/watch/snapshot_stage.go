package watch

import (
	"context"
	"fmt"

	"github.com/user/portwatch/internal/snapshot"
)

// SnapshotStage compares a port list against a saved snapshot and returns only
// the ports that differ (added or removed). It updates the in-memory manager
// after each successful comparison.
type SnapshotStage struct {
	manager *snapshot.Manager
}

// NewSnapshotStage creates a SnapshotStage backed by the given manager.
func NewSnapshotStage(m *snapshot.Manager) (*SnapshotStage, error) {
	if m == nil {
		return nil, fmt.Errorf("snapshot manager must not be nil")
	}
	return &SnapshotStage{manager: m}, nil
}

// Run implements the pipeline Stage interface. It loads the previous snapshot,
// computes the diff, saves the new snapshot, and returns the added ports so
// downstream stages can act on changes only.
func (s *SnapshotStage) Run(ctx context.Context, ports []int) ([]int, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	diff := s.manager.Diff(ports)
	if err := s.manager.Save(ports); err != nil {
		return nil, fmt.Errorf("snapshot save: %w", err)
	}

	// Return only newly added ports for downstream alerting.
	return diff.Added, nil
}
