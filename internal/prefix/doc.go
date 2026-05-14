// Package prefix provides utilities for adding or removing a string prefix
// from environment variable keys in a map.
//
// It is commonly used to namespace variables before injection into a process
// (e.g. APP_) and to strip that namespace when reading them back.
//
// Example usage:
//
//	out, result, err := prefix.Apply(vars, prefix.Options{
//		Prefix: "APP_",
//	})
package prefix
