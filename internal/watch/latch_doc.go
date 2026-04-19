// Package watch provides primitives for monitoring and controlling port
// scanning behaviour in portwatch.
//
// # Latch
//
// A Latch records the first time a port event is observed and suppresses
// duplicate triggers until it is explicitly reset. This is useful when you
// want a one-shot alert — fire once when a port opens, then stay quiet until
// the condition is acknowledged and the latch is cleared.
//
// Basic usage:
//
//	latch := watch.NewLatch()
//
//	if armed, _ := latch.Arm(port); armed {
//		// first time we've seen this port — send alert
//	}
//
//	// later, after the operator acknowledges:
//	latch.Reset(port)
//
package watch
