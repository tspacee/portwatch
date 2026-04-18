// Package watch provides the core port-watching primitives for portwatch.
//
// # Ticker
//
// Ticker is a jitter-aware interval timer used by the watcher loop to
// avoid thundering-herd effects when multiple portwatch instances run
// in the same environment.
//
// Basic usage:
//
//	t, err := watch.NewTicker(10*time.Second, 0.2)
//	if err != nil { ... }
//	defer t.Stop()
//	for range t.C {
//	    // perform scan
//	}
//
// A factor of 0 disables jitter and the ticker fires exactly at the
// base interval. A factor of 0.2 adds up to ±20 % random variance.
package watch
