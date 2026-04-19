// Package watch provides primitives for monitoring, filtering, and
// reacting to changes in open port sets.
//
// # Digest
//
// Digest computes a SHA-256 fingerprint over a sorted list of ports.
// It is useful for quickly determining whether a scan result has changed
// compared to the previous cycle without performing a full diff.
//
// Example usage:
//
//	d := watch.NewDigest()
//	if d.Changed(currentPorts) {
//		// port set has changed since last check
//	}
package watch
