package alert

import (
	"errors"
	"testing"
)

// fakeNotifier is a test double that records calls and optionally returns an error.
type fakeNotifier struct {
	called bool
	err    error
}

func (f *fakeNotifier) Notify(a Alert) error {
	f.called = true
	return f.err
}

func TestNewMultiNotifier_NoNotifiers(t *testing.T) {
	_, err := NewMultiNotifier()
	if err == nil {
		t.Fatal("expected error for empty notifier list, got nil")
	}
}

func TestNewMultiNotifier_StoresNotifiers(t *testing.T) {
	n1 := &fakeNotifier{}
	n2 := &fakeNotifier{}
	m, err := NewMultiNotifier(n1, n2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Len() != 2 {
		t.Errorf("expected 2 notifiers, got %d", m.Len())
	}
}

func TestMultiNotifier_Notify_AllCalled(t *testing.T) {
	n1 := &fakeNotifier{}
	n2 := &fakeNotifier{}
	m, _ := NewMultiNotifier(n1, n2)

	if err := m.Notify(Alert{Title: "test"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !n1.called || !n2.called {
		t.Error("expected both notifiers to be called")
	}
}

func TestMultiNotifier_Notify_PartialFailure(t *testing.T) {
	n1 := &fakeNotifier{err: errors.New("n1 failed")}
	n2 := &fakeNotifier{}
	m, _ := NewMultiNotifier(n1, n2)

	err := m.Notify(Alert{Title: "test"})
	if err == nil {
		t.Fatal("expected error when a notifier fails, got nil")
	}
	if !n2.called {
		t.Error("expected second notifier to be called despite first failure")
	}
}

func TestMultiNotifier_Notify_AllFail(t *testing.T) {
	n1 := &fakeNotifier{err: errors.New("n1 failed")}
	n2 := &fakeNotifier{err: errors.New("n2 failed")}
	m, _ := NewMultiNotifier(n1, n2)

	err := m.Notify(Alert{Title: "test"})
	if err == nil {
		t.Fatal("expected error when all notifiers fail, got nil")
	}
}
