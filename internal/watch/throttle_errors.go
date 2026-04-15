package watch

import "errors"

// Sentinel errors for Throttle construction.
var (
	// ErrInvalidWindow is returned when a non-positive window duration is provided.
	ErrInvalidWindow = errors.New("watch: throttle window must be greater than zero")

	// ErrInvalidMaxCount is returned when maxCount is less than one.
	ErrInvalidMaxCount = errors.New("watch: throttle maxCount must be at least 1")
)
