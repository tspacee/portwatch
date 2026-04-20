package watch

import (
	"errors"
	"sync"
)

// Buffer holds a fixed-size sliding window of integer port values.
// It is safe for concurrent use and evicts the oldest entry when full.
type Buffer struct {
	mu   sync.Mutex
	data []int
	size int
}

// NewBuffer creates a Buffer with the given capacity.
// Returns an error if size is less than 1.
func NewBuffer(size int) (*Buffer, error) {
	if size < 1 {
		return nil, errors.New("buffer: size must be at least 1")
	}
	return &Buffer{
		data: make([]int, 0, size),
		size: size,
	}, nil
}

// Push appends a port value to the buffer.
// If the buffer is full the oldest value is evicted.
func (b *Buffer) Push(port int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.data) == b.size {
		b.data = b.data[1:]
	}
	b.data = append(b.data, port)
}

// Flush returns a copy of the current contents and resets the buffer.
func (b *Buffer) Flush() []int {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]int, len(b.data))
	copy(out, b.data)
	b.data = b.data[:0]
	return out
}

// Peek returns a copy of the current contents without clearing the buffer.
func (b *Buffer) Peek() []int {
	b.mu.Lock()
	defer b.mu.Unlock()
	out := make([]int, len(b.data))
	copy(out, b.data)
	return out
}

// Len returns the number of values currently stored.
func (b *Buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.data)
}

// Cap returns the maximum capacity of the buffer.
func (b *Buffer) Cap() int {
	return b.size
}
