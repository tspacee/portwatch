package snapshot

import (
	"encoding/json"
	"os"
	"time"
)

// Snapshot holds a point-in-time record of open ports.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Ports     []int     `json:"ports"`
}

// New creates a new Snapshot from the given port list.
func New(ports []int) *Snapshot {
	copied := make([]int, len(ports))
	copy(copied, ports)
	return &Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     copied,
	}
}

// Save writes the snapshot as JSON to the given file path.
func (s *Snapshot) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads a snapshot from a JSON file at the given path.
// Returns nil and no error if the file does not exist.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Diff computes ports added and removed relative to a previous snapshot.
// If prev is nil, all current ports are treated as added.
func Diff(prev, curr *Snapshot) (added, removed []int) {
	if prev == nil {
		return curr.Ports, nil
	}
	prevSet := make(map[int]struct{}, len(prev.Ports))
	for _, p := range prev.Ports {
		prevSet[p] = struct{}{}
	}
	currSet := make(map[int]struct{}, len(curr.Ports))
	for _, p := range curr.Ports {
		currSet[p] = struct{}{}
	}
	for _, p := range curr.Ports {
		if _, ok := prevSet[p]; !ok {
			added = append(added, p)
		}
	}
	for _, p := range prev.Ports {
		if _, ok := currSet[p]; !ok {
			removed = append(removed, p)
		}
	}
	return added, removed
}
