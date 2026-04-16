package watch

import (
	"context"
	"testing"
)

func TestDeduplicateStage_RemovesDuplicates(t *testing.T) {
	stage := DeduplicateStage()
	result, err := stage(context.Background(), []int{80, 443, 80, 8080, 443})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 unique ports, got %d: %v", len(result), result)
	}
}

func TestDeduplicateStage_EmptyInput(t *testing.T) {
	stage := DeduplicateStage()
	result, err := stage(context.Background(), []int{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}
}

func TestSortStage_SortsPorts(t *testing.T) {
	stage := SortStage()
	result, err := stage(context.Background(), []int{8080, 443, 80})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0] != 80 || result[1] != 443 || result[2] != 8080 {
		t.Fatalf("unexpected order: %v", result)
	}
}

func TestSortStage_DoesNotMutateInput(t *testing.T) {
	input := []int{8080, 443, 80}
	stage := SortStage()
	_, _ = stage(context.Background(), input)
	if input[0] != 8080 {
		t.Fatal("original slice was mutated")
	}
}

func TestExcludePortsStage_RemovesExcluded(t *testing.T) {
	stage := ExcludePortsStage([]int{22, 80})
	result, err := stage(context.Background(), []int{22, 80, 443, 8080})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 ports, got %d: %v", len(result), result)
	}
	for _, p := range result {
		if p == 22 || p == 80 {
			t.Fatalf("excluded port %d still present", p)
		}
	}
}

func TestExcludePortsStage_EmptyExcludeList(t *testing.T) {
	stage := ExcludePortsStage([]int{})
	result, err := stage(context.Background(), []int{80, 443})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(result))
	}
}
