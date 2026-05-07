// Package watch provides primitives for port monitoring, alerting, and
// behavioral analysis used by the portwatch daemon.
//
// # Intent
//
// Intent is a registry that lets operators record why a port is expected to be
// open. During evaluation, ports without a declared intent can be surfaced as
// candidates for investigation, even if no rule explicitly forbids them.
//
// Example:
//
//	intent := watch.NewIntent()
//	_ = intent.Declare(443, "public HTTPS endpoint")
//	_ = intent.Declare(22, "SSH management access")
//
//	if reason, ok := intent.Lookup(port); ok {
//		log.Printf("port %d is expected: %s", port, reason)
//	} else {
//		log.Printf("port %d has no declared intent", port)
//	}
package watch
