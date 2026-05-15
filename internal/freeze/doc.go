// Package freeze provides a Freezer type that captures an immutable snapshot
// of an environment variable map at a point in time.
//
// Use cases include:
//   - Detecting drift between a pinned config and a live environment
//   - Preventing accidental mutation of a baseline env state
//   - Auditing changes between two points in a pipeline
//
// Example:
//
//	f, err := freeze.New(envMap)
//	if err != nil {
//		log.Fatal(err)
//	}
//	// later...
//	changed := f.DiffFrom(currentEnvMap)
//	if len(changed) > 0 {
//		fmt.Println("env has drifted:", changed)
//	}
package freeze
