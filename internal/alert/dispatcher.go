package alert

import "github.com/user/portwatch/internal/rules"

// Dispatcher receives rule violations and fans them out to registered Notifiers.
type Dispatcher struct {
	notifiers []Notifier
}

// NewDispatcher creates a Dispatcher with the given notifiers.
func NewDispatcher(notifiers ...Notifier) *Dispatcher {
	ns := make([]Notifier, len(notifiers))
	copy(ns, notifiers)
	return &Dispatcher{notifiers: ns}
}

// AddNotifier registers an additional Notifier.
func (d *Dispatcher) AddNotifier(n Notifier) {
	d.notifiers = append(d.notifiers, n)
}

// Dispatch converts each violation into an Alert and sends it to all notifiers.
// It returns the first error encountered, if any.
func (d *Dispatcher) Dispatch(violations []rules.Violation) error {
	for _, v := range violations {
		a := ViolationToAlert(v)
		for _, n := range d.notifiers {
			if err := n.Notify(a); err != nil {
				return err
			}
		}
	}
	return nil
}

// NotifierCount returns the number of registered notifiers.
func (d *Dispatcher) NotifierCount() int {
	return len(d.notifiers)
}
