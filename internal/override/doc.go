// Package override provides utilities for applying key-value overrides
// to environment variable maps.
//
// It supports:
//   - Parsing "KEY=VALUE" pair strings
//   - Applying multiple overrides to an existing map
//   - Strict mode, which rejects overrides for keys not already present
//   - AllowEmpty mode, which permits overriding a key with an empty string
//
// Example usage:
//
//	entries, _ := override.ParsePair("DEBUG=true")
//	override.Apply(envMap, []override.Entry{entries}, override.Options{Strict: true})
package override
