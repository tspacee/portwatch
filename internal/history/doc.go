// Package history maintains a bounded, thread-safe ring buffer of port-scan
// events produced by the portwatch daemon.
//
// Each [Entry] captures the timestamp, the set of open ports observed, and any
// rule violations that were raised during that scan cycle.  Older entries are
// automatically evicted once the configured capacity is reached, keeping memory
// usage predictable regardless of how long the daemon runs.
//
// Typical usage:
//
//	h := history.New(200)          // keep last 200 scan results
//	h.Add(history.Entry{
//	    Timestamp:  time.Now(),
//	    OpenPorts:  ports,
//	    Violations: violationMessages,
//	})
//	fmt.Println(h.ViolationCount()) // total violations in window
package history
