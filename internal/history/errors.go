package history

import "errors"

// ErrNilHistory is returned when a nil *History is passed where one is required.
var ErrNilHistory = errors.New("history: nil history provided")

// ErrInvalidLimit is returned when a non-positive limit is specified.
var ErrInvalidLimit = errors.New("history: limit must be greater than zero")
