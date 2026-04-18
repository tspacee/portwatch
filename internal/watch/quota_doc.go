// Package watch provides primitives for controlling and observing port-scan
// cycles in portwatch.
//
// # Quota
//
// Quota limits the total number of scan cycles that may run within a rolling
// time window. It is complementary to RateLimiter (which enforces a minimum
// interval between scans) and Budget (which tracks CPU/time usage).
//
// Use Quota when you need a hard cap on how many scans can occur per hour or
// per day regardless of spacing:
//
//	q, err := watch.NewQuota(time.Hour, 60)
//	if err != nil { ... }
//	if q.Allow() {
//	    // run scan
//	}
//
// Remaining reports how many more scans are permitted before the oldest
// recorded timestamp leaves the window. Reset clears all state.
package watch
