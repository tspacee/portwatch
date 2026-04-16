package watch

import "errors"

// ErrInvalidEventLogSize is returned when maxSize is less than 1.
var ErrInvalidEventLogSize = errors.New("eventlog: maxSize must be at least 1")
