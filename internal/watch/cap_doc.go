// Package watch provides primitives for monitoring, controlling, and
// analyzing port activity within portwatch.
//
// Cap
//
// Cap enforces a hard upper bound on the number of concurrently tracked
// ports. It is useful when downstream consumers have finite capacity and
// must not be overwhelmed by a sudden burst of newly observed ports.
//
// Usage:
//
//	cap, err := watch.NewCap(100)
//	if err != nil { ... }
//
//	ok, err := cap.Acquire(port)
//	if ok {
//	    defer cap.Release(port)
//	    // process port
//	}
package watch
