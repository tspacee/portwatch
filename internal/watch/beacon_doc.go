// Package watch provides utilities for monitoring, rate-limiting, and
// coordinating port scan activity within portwatch.
//
// # Beacon
//
// Beacon is a lightweight heartbeat emitter that fires a signal on a channel
// at a fixed interval. It is useful for triggering periodic actions such as
// flushing state, emitting health pings, or coordinating pipeline stages
// without coupling them to a specific ticker implementation.
//
// Example usage:
//
//	b, err := watch.NewBeacon(30 * time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
//	b.Start()
//	defer b.Stop()
//
//	for range b.C() {
//		// perform periodic work
//	}
package watch
