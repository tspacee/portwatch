package alert

import "sync"

// Notifier is the interface implemented by all alert backends.
type Notifier interface {
	// Notify sends the given alert. Implementations should be safe to call
	// concurrently and must not modify the Alert value.
	Notify(a Alert) error
}

// NopNotifier is a no-op Notifier that silently discards every alert.
// Useful as a default when no notifier is configured.
type NopNotifier struct{}

// Notify implements Notifier. It always returns nil.
func (NopNotifier) Notify(_ Alert) error { return nil }

// MultiNotifier fans out each alert to a list of Notifiers in order.
// All notifiers are called even if one returns an error; the first
// non-nil error encountered is returned to the caller.
type MultiNotifier struct {
	mu        sync.RWMutex
	notifiers []Notifier
}

// NewMultiNotifier returns a MultiNotifier that dispatches to the given notifiers.
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

// Notify implements Notifier. It calls every contained Notifier and returns
// the first error encountered, continuing to notify remaining backends regardless.
func (m *MultiNotifier) Notify(a Alert) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var firstErr error
	for _, n := range m.notifiers {
		if err := n.Notify(a); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
