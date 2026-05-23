// Package watch provides building blocks for port observation and alerting.
//
// # Tripwire
//
// Tripwire monitors a fixed set of ports and fires a one-shot callback
// the moment any watched port disappears from an observed scan result.
//
// Once tripped the tripwire is silent until Reset is called, preventing
// alert storms when a port remains absent across multiple scan cycles.
//
// Typical usage:
//
//	tw, err := watch.NewTripwire([]int{22, 443}, func(p int) {
//		log.Printf("critical port %d vanished!", p)
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	// on each scan cycle:
//	tw.Observe(openPorts)
package watch
