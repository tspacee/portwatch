// Package watch provides watcher utilities for portwatch.
//
// # Supervisor
//
// Supervisor wraps a WorkerFunc and automatically restarts it on failure
// using exponential backoff. It is useful for keeping long-running scan
// loops alive in the presence of transient errors.
//
// Basic usage:
//
//	backoff, _ := watch.NewBackoff(100*time.Millisecond, 10*time.Second, 2.0)
//	sup, _ := watch.NewSupervisor(myWorker, backoff, 5)
//	if err := sup.Run(ctx); err != nil {
//		log.Println("supervisor stopped:", err)
//	}
//
// A maxRestarts value of 0 or less means the supervisor will restart
// indefinitely until the context is cancelled.
package watch
