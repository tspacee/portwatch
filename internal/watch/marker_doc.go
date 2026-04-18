// Package watch provides primitives for port monitoring, rate limiting,
// and scan lifecycle management.
//
// # Marker
//
// Marker allows operators to flag specific ports with a human-readable reason.
// Marked ports can be highlighted in status output or filtered by downstream
// pipeline stages. Marks are in-memory only and do not persist across restarts.
//
// Example:
//
//	m := watch.NewMarker()
//	_ = m.Mark(8080, "under investigation")
//	if m.IsMarked(8080) {
//		fmt.Println(m.Reason(8080))
//	}
package watch
