package watch

import (
	"context"
	"sort"
)

// DeduplicateStage returns a Stage that removes duplicate ports.
func DeduplicateStage() Stage {
	return func(ctx context.Context, ports []int) ([]int, error) {
		seen := make(map[int]struct{}, len(ports))
		out := ports[:0:0]
		for _, p := range ports {
			if _, ok := seen[p]; !ok {
				seen[p] = struct{}{}
				out = append(out, p)
			}
		}
		return out, nil
	}
}

// SortStage returns a Stage that sorts ports in ascending order.
func SortStage() Stage {
	return func(ctx context.Context, ports []int) ([]int, error) {
		copy := append([]int(nil), ports...)
		sort.Ints(copy)
		return copy, nil
	}
}

// ExcludePortsStage returns a Stage that removes specified ports from the list.
func ExcludePortsStage(excluded []int) Stage {
	excludeSet := make(map[int]struct{}, len(excluded))
	for _, p := range excluded {
		excludeSet[p] = struct{}{}
	}
	return func(ctx context.Context, ports []int) ([]int, error) {
		out := ports[:0:0]
		for _, p := range ports {
			if _, skip := excludeSet[p]; !skip {
				out = append(out, p)
			}
		}
		return out, nil
	}
}
