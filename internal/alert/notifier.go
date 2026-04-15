package alert

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
