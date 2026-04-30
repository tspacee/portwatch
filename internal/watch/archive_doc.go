// Package watch provides primitives for monitoring and analyzing port activity.
//
// Archive
//
// Archive maintains a bounded, chronological log of port snapshots captured
// during successive scan cycles. Each entry pairs a timestamp with the list
// of open ports observed at that moment.
//
// Usage:
//
//	arch, err := watch.NewArchive(100)
//	if err != nil { ... }
//
//	// Record a snapshot after each scan.
//	arch.Store(openPorts)
//
//	// Retrieve all entries for trend analysis or export.
//	entries := arch.Entries()
//
The archive evicts the oldest entry when the capacity limit is reached,
ensuring memory usage remains bounded over long-running daemon sessions.
package watch
