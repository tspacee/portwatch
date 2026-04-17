package watch

import (
	"testing"
	"time"
)

func TestNewDedupCache_InvalidTTL(t *testing.T) {
	_, err := NewDedupCache(0)
	if err == nil {
		t.Fatal("expected error for zero TTL")
	}
	if err != ErrInvalidTTL {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewDedupCache_Valid(t *testing.T) {
	c, err := NewDedupCache(time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
}

func TestDedupCache_Seen_FirstCall_ReturnsFalse(t *testing.T) {
	c, _ := NewDedupCache(time.Second)
	if c.Seen("port:8080") {
		t.Fatal("expected false on first call")
	}
}

func TestDedupCache_Seen_SecondCall_ReturnsTrue(t *testing.T) {
	c, _ := NewDedupCache(time.Second)
	c.Seen("port:8080")
	if !c.Seen("port:8080") {
		t.Fatal("expected true on second call within TTL")
	}
}

func TestDedupCache_Seen_AfterExpiry_ReturnsFalse(t *testing.T) {
	c, _ := NewDedupCache(20 * time.Millisecond)
	c.Seen("port:9090")
	time.Sleep(40 * time.Millisecond)
	if c.Seen("port:9090") {
		t.Fatal("expected false after TTL expiry")
	}
}

func TestDedupCache_Len_CountsActive(t *testing.T) {
	c, _ := NewDedupCache(time.Second)
	c.Seen("a")
	c.Seen("b")
	c.Seen("c")
	if c.Len() != 3 {
		t.Fatalf("expected 3, got %d", c.Len())
	}
}

func TestDedupCache_Evict_RemovesExpired(t *testing.T) {
	c, _ := NewDedupCache(20 * time.Millisecond)
	c.Seen("x")
	c.Seen("y")
	time.Sleep(40 * time.Millisecond)
	c.Evict()
	if c.Len() != 0 {
		t.Fatalf("expected 0 after eviction, got %d", c.Len())
	}
}

func TestDedupCache_Len_ExcludesExpired(t *testing.T) {
	c, _ := NewDedupCache(20 * time.Millisecond)
	c.Seen("exp")
	time.Sleep(40 * time.Millisecond)
	c.Seen("active")
	if c.Len() != 1 {
		t.Fatalf("expected 1 active entry, got %d", c.Len())
	}
}
