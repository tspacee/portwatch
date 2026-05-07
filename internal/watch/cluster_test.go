package watch

import (
	"testing"
)

func TestNewCluster_Empty(t *testing.T) {
	c := NewCluster()
	if len(c.Groups()) != 0 {
		t.Fatal("expected no groups on new cluster")
	}
}

func TestCluster_Add_Valid(t *testing.T) {
	c := NewCluster()
	if err := c.Add("web", 80); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !c.Contains("web", 80) {
		t.Fatal("expected port 80 in group 'web'")
	}
}

func TestCluster_Add_EmptyGroup(t *testing.T) {
	c := NewCluster()
	if err := c.Add("", 80); err == nil {
		t.Fatal("expected error for empty group name")
	}
}

func TestCluster_Add_InvalidPort(t *testing.T) {
	c := NewCluster()
	if err := c.Add("web", 0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := c.Add("web", 65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestCluster_Members_ReturnsAll(t *testing.T) {
	c := NewCluster()
	_ = c.Add("db", 5432)
	_ = c.Add("db", 5433)
	members := c.Members("db")
	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}
}

func TestCluster_Members_UnknownGroup(t *testing.T) {
	c := NewCluster()
	if m := c.Members("missing"); m != nil {
		t.Fatal("expected nil for unknown group")
	}
}

func TestCluster_Contains_False(t *testing.T) {
	c := NewCluster()
	if c.Contains("web", 443) {
		t.Fatal("expected false for unregistered port")
	}
}

func TestCluster_Remove_RemovesPort(t *testing.T) {
	c := NewCluster()
	_ = c.Add("web", 80)
	c.Remove("web", 80)
	if c.Contains("web", 80) {
		t.Fatal("expected port 80 to be removed")
	}
}

func TestCluster_Remove_EmptiesGroup(t *testing.T) {
	c := NewCluster()
	_ = c.Add("web", 80)
	c.Remove("web", 80)
	groups := c.Groups()
	for _, g := range groups {
		if g == "web" {
			t.Fatal("expected empty group to be removed from Groups()")
		}
	}
}

func TestCluster_Groups_ReturnsCopy(t *testing.T) {
	c := NewCluster()
	_ = c.Add("a", 1000)
	_ = c.Add("b", 2000)
	g1 := c.Groups()
	g1[0] = "mutated"
	for _, g := range c.Groups() {
		if g == "mutated" {
			t.Fatal("Groups() should return a copy")
		}
	}
}
