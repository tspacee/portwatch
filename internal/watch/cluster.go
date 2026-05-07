package watch

import (
	"errors"
	"sync"
)

// Cluster groups ports into named clusters and tracks membership.
// It is useful for correlating related ports (e.g., a service suite)
// and detecting when a cluster becomes partially or fully unavailable.
type Cluster struct {
	mu      sync.RWMutex
	groups  map[string]map[int]struct{}
}

// NewCluster returns an empty Cluster ready for use.
func NewCluster() *Cluster {
	return &Cluster{
		groups: make(map[string]map[int]struct{}),
	}
}

// Add registers a port under the named cluster group.
// Returns an error if the group name is empty or the port is out of range.
func (c *Cluster) Add(group string, port int) error {
	if group == "" {
		return errors.New("cluster: group name must not be empty")
	}
	if port < 1 || port > 65535 {
		return errors.New("cluster: port out of range")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.groups[group]; !ok {
		c.groups[group] = make(map[int]struct{})
	}
	c.groups[group][port] = struct{}{}
	return nil
}

// Members returns a copy of all ports registered under the named group.
// Returns nil if the group does not exist.
func (c *Cluster) Members(group string) []int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	set, ok := c.groups[group]
	if !ok {
		return nil
	}
	out := make([]int, 0, len(set))
	for p := range set {
		out = append(out, p)
	}
	return out
}

// Contains reports whether the given port is a member of the named group.
func (c *Cluster) Contains(group string, port int) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	set, ok := c.groups[group]
	if !ok {
		return false
	}
	_, found := set[port]
	return found
}

// Groups returns the names of all registered cluster groups.
func (c *Cluster) Groups() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	names := make([]string, 0, len(c.groups))
	for name := range c.groups {
		names = append(names, name)
	}
	return names
}

// Remove unregisters a port from the named cluster group.
func (c *Cluster) Remove(group string, port int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if set, ok := c.groups[group]; ok {
		delete(set, port)
		if len(set) == 0 {
			delete(c.groups, group)
		}
	}
}
