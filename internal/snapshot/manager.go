package snapshot

import "sync"

// Manager maintains the latest in-memory snapshot and persists it to disk.
type Manager struct {
	mu       sync.RWMutex
	filePath string
	latest   *Snapshot
}

// NewManager creates a Manager backed by the given file path.
// It attempts to load an existing snapshot from disk on creation.
func NewManager(filePath string) (*Manager, error) {
	s, err := Load(filePath)
	if err != nil {
		return nil, err
	}
	return &Manager{
		filePath: filePath,
		latest:   s,
	}, nil
}

// Update replaces the current snapshot with a new one built from ports,
// computes the diff against the previous snapshot, persists to disk,
// and returns the added and removed port lists.
func (m *Manager) Update(ports []int) (added, removed []int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	next := New(ports)
	added, removed = Diff(m.latest, next)

	if err = next.Save(m.filePath); err != nil {
		return nil, nil, err
	}
	m.latest = next
	return added, removed, nil
}

// Latest returns a copy of the most recent snapshot, or nil if none exists.
func (m *Manager) Latest() *Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.latest == nil {
		return nil
	}
	copy := *m.latest
	portsCopy := make([]int, len(m.latest.Ports))
	_ = append(portsCopy[:0], m.latest.Ports...)
	copy.Ports = portsCopy
	return &copy
}
