// Package watch provides primitives for monitoring, filtering, and controlling
// port scan behaviour in portwatch.
//
// # Suppress
//
// Suppress provides a time-bounded mute list for individual ports. When a port
// is muted it will not trigger alerts until the suppression window expires.
//
// Typical usage:
//
//	s, _ := watch.NewSuppress(5 * time.Minute)
//	s.Mute(8080)
//
//	if !s.IsMuted(8080) {
//		// send alert
//	}
//
// Suppressions are automatically cleaned up when IsMuted or Len is called,
// so no background goroutine is required.
package watch
