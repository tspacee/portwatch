// Package watch provides primitives for monitoring, filtering, and
// controlling port scan behaviour in portwatch.
//
// Gate
//
// Gate is a simple allow-list controller that determines whether a port
// is permitted to proceed through the monitoring pipeline.
//
// When created with open=false (the default), only ports explicitly
// registered via Allow() will pass. When created with open=true, all
// ports pass unless removed via Deny().
//
// Example:
//
//	g := watch.NewGate(false)
//	_ = g.Allow(443)
//	_ = g.Allow(80)
//
//	if g.Pass(443) {
//		// port 443 is permitted
//	}
//
// Gate is safe for concurrent use.
package watch
