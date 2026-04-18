// Package watch provides primitives for controlling and monitoring
// the port-scanning loop.
//
// ScanLimiter combines two complementary rate-limiting strategies:
//
//  1. Minimum interval — enforces a hard floor between consecutive scans
//     so that a fast event loop cannot trigger back-to-back scans.
//
//  2. Rolling window — caps the total number of scans that may occur
//     within a sliding time window, preventing burst storms even when
//     individual scans are spaced apart.
//
// Both constraints must be satisfied for Allow to return true.
// Call Reset to clear accumulated state, e.g. after a reconfiguration.
package watch
