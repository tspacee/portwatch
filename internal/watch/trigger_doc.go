// Package watch provides primitives for port monitoring, rate limiting,
// and event control used by the portwatch daemon.
//
// Trigger
//
// Trigger fires a user-supplied callback once a configurable number of
// events (Record calls) accumulate within a sliding time window.
//
// After the callback fires the internal event list is cleared, so the
// trigger is immediately ready to arm again.
//
// Example:
//
//	tr, err := watch.NewTrigger(10*time.Second, 5, func() {
//		log.Println("threshold reached")
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	tr.Record() // call on each relevant event
package watch
