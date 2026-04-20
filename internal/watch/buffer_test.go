package watch

import (
	"testing"
)

func TestNewBuffer_InvalidSize(t *testing.T) {
	_, err := NewBuffer(0)
	if err == nil {
		t.Fatal("expected error for size 0")
	}
}

func TestNewBuffer_Valid(t *testing.T) {
	b, err := NewBuffer(4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Cap() != 4 {
		t.Errorf("expected cap 4, got %d", b.Cap())
	}
	if b.Len() != 0 {
		t.Errorf("expected len 0, got %d", b.Len())
	}
}

func TestBuffer_Push_And_Len(t *testing.T) {
	b, _ := NewBuffer(3)
	b.Push(80)
	b.Push(443)
	if b.Len() != 2 {
		t.Errorf("expected len 2, got %d", b.Len())
	}
}

func TestBuffer_Push_EvictsOldest(t *testing.T) {
	b, _ := NewBuffer(2)
	b.Push(80)
	b.Push(443)
	b.Push(8080)
	ports := b.Peek()
	if len(ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(ports))
	}
	if ports[0] != 443 || ports[1] != 8080 {
		t.Errorf("unexpected ports after eviction: %v", ports)
	}
}

func TestBuffer_Peek_ReturnsCopy(t *testing.T) {
	b, _ := NewBuffer(3)
	b.Push(22)
	p1 := b.Peek()
	p1[0] = 9999
	p2 := b.Peek()
	if p2[0] == 9999 {
		t.Error("Peek should return an independent copy")
	}
}

func TestBuffer_Flush_ClearsBuffer(t *testing.T) {
	b, _ := NewBuffer(4)
	b.Push(80)
	b.Push(443)
	out := b.Flush()
	if len(out) != 2 {
		t.Fatalf("expected 2 flushed items, got %d", len(out))
	}
	if b.Len() != 0 {
		t.Errorf("expected buffer empty after flush, got %d", b.Len())
	}
}

func TestBuffer_Flush_ReturnsCopy(t *testing.T) {
	b, _ := NewBuffer(3)
	b.Push(22)
	out := b.Flush()
	out[0] = 9999
	// buffer is empty now, push again and verify independence
	b.Push(22)
	p := b.Peek()
	if p[0] == 9999 {
		t.Error("Flush result should be independent of internal state")
	}
}
