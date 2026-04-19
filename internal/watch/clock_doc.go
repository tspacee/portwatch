// Package watch provides watcher primitives for portwatch.
//
// Clock
//
// Clock is a lightweight scan-timing tracker. It records when a scan
// begins and ends, computes elapsed duration, and maintains a count of
// completed scan cycles.
//
// Usage:
//
//	clock := watch.NewClock()
//	clock.Start()
//	// ... perform scan ...
//	clock.Stop()
//	fmt.Println(clock.Elapsed())  // duration of last scan
//	fmt.Println(clock.Count())    // total completed scans
package watch
