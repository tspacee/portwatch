// Package watch provides the core port-watching loop, status reporting,
// rate limiting, throttling, and retry backoff for portwatch.
//
// # Backoff
//
// Backoff implements exponential backoff to progressively delay retry
// attempts after scan errors or transient failures. The delay starts at
// a configurable base duration and doubles on each call to Next, up to
// a configured maximum.
//
// Example usage:
//
//	b, err := watch.NewBackoff(500*time.Millisecond, 30*time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for {
//		if err := doScan(); err != nil {
//			time.Sleep(b.Next())
//			continue
//		}
//		b.Reset()
//	}
package watch
