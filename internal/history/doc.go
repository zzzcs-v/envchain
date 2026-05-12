// Package history records and retrieves a log of past envchain context
// resolutions. Each time a context is exported or resolved, a history entry
// can be saved to a local directory for later inspection.
//
// Entries are stored as individual JSON files named by Unix nanosecond
// timestamp, making them easy to enumerate and sort without a database.
//
// Typical usage:
//
//	s, err := history.NewStore("~/.envchain/history")
//	if err != nil { ... }
//
//	// record after exporting
//	s.Record("dev", "dotenv", resolvedVars)
//
//	// list all past runs
//	entries, err := s.List()
package history
