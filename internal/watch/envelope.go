package watch

import (
	"errors"
	"time"
)

// Envelope wraps a scan result with metadata such as scan time,
// source identifier, and a sequence number for ordering.
type Envelope struct {
	Seq       uint64
	Source    string
	ScannedAt time.Time
	Ports     []int
}

// EnvelopeBuilder constructs Envelopes with an auto-incrementing sequence.
type EnvelopeBuilder struct {
	source string
	seq    uint64
}

// NewEnvelopeBuilder creates a new EnvelopeBuilder for the given source.
func NewEnvelopeBuilder(source string) (*EnvelopeBuilder, error) {
	if source == "" {
		return nil, errors.New("envelope: source must not be empty")
	}
	return &EnvelopeBuilder{source: source}, nil
}

// Wrap creates an Envelope for the given ports, incrementing the sequence.
func (b *EnvelopeBuilder) Wrap(ports []int) Envelope {
	b.seq++
	cp := make([]int, len(ports))
	copy(cp, ports)
	return Envelope{
		Seq:       b.seq,
		Source:    b.source,
		ScannedAt: time.Now(),
		Ports:     cp,
	}
}

// Seq returns the current sequence counter without incrementing.
func (b *EnvelopeBuilder) CurrentSeq() uint64 {
	return b.seq
}
