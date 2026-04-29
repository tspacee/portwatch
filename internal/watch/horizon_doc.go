// Package watch provides primitives for monitoring, filtering, and
// analyzing port scan results over time.
//
// # Horizon
//
// Horizon tracks how long each port has been continuously open by
// recording the first time a port is observed. Once a port's age
// exceeds the configured cutoff duration, it is considered "beyond
// the horizon" and may warrant special attention or alerting.
//
// Typical usage:
//
//	h, _ := watch.NewHorizon(24 * time.Hour)
//	h.Observe(port)
//	if h.Beyond(port) {
//	    // port has been open for more than 24 hours
//	}
package watch
