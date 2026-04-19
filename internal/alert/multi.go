package alert

import (
	"errors"
	"fmt"
)

// MultiNotifier fans out alerts to multiple Notifier implementations.
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier creates a MultiNotifier from the provided notifiers.
// Returns an error if no notifiers are provided.
func NewMultiNotifier(notifiers ...Notifier) (*MultiNotifier, error) {
	if len(notifiers) == 0 {
		return nil, errors.New("multi notifier: at least one notifier is required")
	}
	return &MultiNotifier{notifiers: notifiers}, nil
}

// Notify sends the alert to all registered notifiers. It collects any errors
// and returns a combined error if one or more notifiers fail.
func (m *MultiNotifier) Notify(a Alert) error {
	var errs []error
	for _, n := range m.notifiers {
		if err := n.Notify(a); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("multi notifier: %d notifier(s) failed: %w", len(errs), errors.Join(errs...))
}

// Add appends one or more notifiers to the MultiNotifier.
func (m *MultiNotifier) Add(notifiers ...Notifier) {
	m.notifiers = append(m.notifiers, notifiers...)
}

// Len returns the number of registered notifiers.
func (m *MultiNotifier) Len() int {
	return len(m.notifiers)
}
