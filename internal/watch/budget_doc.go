// Package watch provides Budget, a rolling-window scan-time governor.
//
// Budget tracks how much cumulative scan duration has been consumed within
// a configurable time window. When the total exceeds the configured maximum,
// Allow returns false, preventing new scans from starting until older entries
// expire out of the window.
//
// Typical usage:
//
//	b, err := watch.NewBudget(time.Minute, 10*time.Second)
//	if err != nil { ... }
//
//	if b.Allow() {
//		start := time.Now()
//		// ... perform scan ...
//		b.Record(time.Since(start))
//	}
package watch
