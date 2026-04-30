// Package watch provides primitives for monitoring, rate-limiting,
// and controlling port scan behaviour in portwatch.
//
// # Epoch
//
// Epoch tracks a monotonically incrementing scan generation counter.
// Each call to Advance increments the counter and records the wall-clock
// time of that generation. Downstream components can correlate port
// observations, diffs, and alerts back to a specific scan cycle by
// comparing epoch numbers.
//
// Typical usage:
//
//	epoch := watch.NewEpoch()
//
//	// at the start of each scan loop iteration:
//	gen := epoch.Advance()
//
//	// later, to measure how long ago that generation ran:
//	elapsed, err := epoch.Since(gen)
//
// Epoch is safe for concurrent use.
package watch
