package watch

import (
	"context"
	"errors"
	"testing"
)

func TestNewPipeline_NoStages(t *testing.T) {
	_, err := NewPipeline()
	if !errors.Is(err, ErrEmptyPipeline) {
		t.Fatalf("expected ErrEmptyPipeline, got %v", err)
	}
}

func TestNewPipeline_Valid(t *testing.T) {
	p, err := NewPipeline(func(ctx context.Context, ports []int) ([]int, error) {
		return ports, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Len() != 1 {
		t.Fatalf("expected 1 stage, got %d", p.Len())
	}
}

func TestPipeline_Run_PassesThroughStages(t *testing.T) {
	double := func(ctx context.Context, ports []int) ([]int, error) {
		out := make([]int, len(ports))
		for i, p := range ports {
			out[i] = p * 2
		}
		return out, nil
	}
	p, _ := NewPipeline(double, double)
	result, err := p.Run(context.Background(), []int{1, 2, 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0] != 4 || result[1] != 8 || result[2] != 12 {
		t.Fatalf("unexpected result: %v", result)
	}
}

func TestPipeline_Run_StopsOnError(t *testing.T) {
	sentinel := errors.New("stage error")
	called := false
	p, _ := NewPipeline(
		func(ctx context.Context, ports []int) ([]int, error) {
			return nil, sentinel
		},
		func(ctx context.Context, ports []int) ([]int, error) {
			called = true
			return ports, nil
		},
	)
	_, err := p.Run(context.Background(), []int{80})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if called {
		t.Fatal("second stage should not have been called")
	}
}

func TestPipeline_Run_StopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	p, _ := NewPipeline(
		func(ctx context.Context, ports []int) ([]int, error) {
			return ports, nil
		},
	)
	_, err := p.Run(ctx, []int{80})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
