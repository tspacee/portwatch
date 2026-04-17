// Package watch provides utilities for monitoring port changes,
// including rate limiting, backoff, circuit breaking, and event deduplication.
//
// # DedupCache
//
// DedupCache suppresses duplicate events that occur within a configurable
// time-to-live (TTL) window. This is useful when the same port change is
// detected across multiple consecutive scan cycles before a rule action
// has been taken.
//
// Usage:
//
//	cache, err := watch.NewDedupCache(30 * time.Second)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	key := fmt.Sprintf("%s:%d", event.Protocol, event.Port)
//	if !cache.Seen(key) {
//		// handle event
//	}
//
// Call Evict() periodically to free memory from expired entries.
package watch
