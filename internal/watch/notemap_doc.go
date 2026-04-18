// Package watch provides port watching primitives for portwatch.
//
// NoteMap
//
// NoteMap allows operators to attach human-readable notes to specific port
// numbers. Notes are useful for documenting why a port is expected to be open
// or flagged, and can be surfaced in alerts and audit logs.
//
// Example:
//
//	nm := watch.NewNoteMap()
//	_ = nm.Set(8080, "internal proxy, approved 2024-01-15")
//	fmt.Println(nm.Get(8080))
package watch
