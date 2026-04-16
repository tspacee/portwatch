// Package watch provides primitives for controlling scan execution.
//
// # Circuit Breaker
//
// CircuitBreaker protects the port scanner from repeated failures by
// temporarily halting scan attempts when consecutive errors exceed a
// configured threshold.
//
// States:
//
//	- Closed  — normal operation; calls are allowed through.
//	- Open    — breaker has tripped; calls return ErrCircuitOpen.
//	- HalfOpen — recovery window elapsed; one probe call is allowed.
//
// Usage:
//
//	cb, err := watch.NewCircuitBreaker(5, 30*time.Second)
//	if err != nil { /* handle */ }
//
//	if err := cb.Allow(); err != nil {
//	    // skip scan this cycle
//	}
//	// perform scan ...
//	cb.RecordSuccess() // or cb.RecordFailure()
package watch
