// Package watch provides utilities for monitoring port activity.
//
// # Registry
//
// Registry maintains a thread-safe map of observed ports and their
// metadata (protocol, observation count). It is intended to be used
// alongside the Watcher to accumulate state across scan cycles.
//
// Example usage:
//
//	reg := watch.NewRegistry()
//	_ = reg.Track(8080, "tcp")
//	meta, ok := reg.Get(8080)
//	if ok {
//		fmt.Println(meta.SeenCount) // 1
//	}
package watch
