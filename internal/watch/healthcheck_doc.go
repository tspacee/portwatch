// Package watch provides the core port-watching loop and supporting primitives.
//
// # HealthChecker
//
// HealthChecker tracks the operational health of the watcher by recording
// scan successes and failures. After a configurable number of consecutive
// failures (maxFails), the checker transitions to an unhealthy state.
//
// A single RecordSuccess call resets the failure counter and restores healthy
// status, making it suitable for use with auto-recovery patterns.
//
// Usage:
//
//	hc, err := watch.NewHealthChecker(3)
//	if err != nil { ... }
//
//	// after each scan:
//	if scanErr != nil {
//	    hc.RecordFailure(scanErr)
//	} else {
//	    hc.RecordSuccess()
//	}
//
//	status := hc.Status()
//	if !status.Healthy {
//	    // trigger alert or circuit breaker
//	}
package watch
