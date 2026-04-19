// Package watch provides primitives for monitoring, controlling, and
// analysing port scan activity.
//
// # Tally
//
// Tally maintains a cumulative hit-count per port number across multiple
// scan cycles. It is safe for concurrent use.
//
// Typical usage:
//
//	tal := watch.NewTally()
//	for _, port := range openPorts {
//		_ = tal.Inc(port)
//	}
//	topPort := tal.Top()
//
// Counts persist until explicitly cleared with Reset, making Tally
// suitable for long-running daemons that need to rank ports by
// observation frequency.
package watch
