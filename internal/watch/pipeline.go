package watch

import (
	"context"
	"errors"
	"fmt"
)

// Stage is a function that processes a scan result and passes it along.
type Stage func(ctx context.Context, ports []int) ([]int, error)

// Pipeline chains multiple stages together, passing output of one into the next.
type Pipeline struct {
	stages []Stage
}

// ErrEmptyPipeline is returned when no stages are provided.
var ErrEmptyPipeline = errors.New("pipeline: at least one stage is required")

// NewPipeline constructs a Pipeline from the given stages.
func NewPipeline(stages ...Stage) (*Pipeline, error) {
	if len(stages) == 0 {
		return nil, ErrEmptyPipeline
	}
	return &Pipeline{stages: stages}, nil
}

// Run executes each stage in order, threading ports through the chain.
// If any stage returns an error, execution stops and the error is wrapped.
func (p *Pipeline) Run(ctx context.Context, ports []int) ([]int, error) {
	current := ports
	for i, stage := range p.stages {
		var err error
		current, err = stage(ctx, current)
		if err != nil {
			return nil, fmt.Errorf("pipeline stage %d: %w", i, err)
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}
	return current, nil
}

// Len returns the number of stages in the pipeline.
func (p *Pipeline) Len() int {
	return len(p.stages)
}
