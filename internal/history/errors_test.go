package history

import (
	"errors"
	"testing"
)

func TestErrNilHistory_IsDistinct(t *testing.T) {
	if errors.Is(ErrNilHistory, ErrInvalidLimit) {
		t.Fatal("ErrNilHistory should not match ErrInvalidLimit")
	}
}

func TestErrNilHistory_Message(t *testing.T) {
	if ErrNilHistory.Error() == "" {
		t.Fatal("expected non-empty error message for ErrNilHistory")
	}
}

func TestErrInvalidLimit_Message(t *testing.T) {
	if ErrInvalidLimit.Error() == "" {
		t.Fatal("expected non-empty error message for ErrInvalidLimit")
	}
}

func TestNewServer_NilReturnsErrNilHistory(t *testing.T) {
	_, err := NewServer(nil)
	if !errors.Is(err, ErrNilHistory) {
		t.Fatalf("expected ErrNilHistory, got %v", err)
	}
}

func TestNewExporter_NilReturnsErrNilHistory(t *testing.T) {
	_, err := NewExporter(nil)
	if !errors.Is(err, ErrNilHistory) {
		t.Fatalf("expected ErrNilHistory, got %v", err)
	}
}
