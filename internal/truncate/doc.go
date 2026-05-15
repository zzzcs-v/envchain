// Package truncate provides utilities for capping the length of environment
// variable values. It is useful when exporting configs to systems with strict
// value-length limits (e.g. certain CI providers or shell argument lists).
//
// Basic usage:
//
//	out, results, err := truncate.Truncate(src, truncate.Options{
//		MaxLen: 64,
//		Suffix: "...",
//	})
//
// When Keys is non-empty only those keys are considered for truncation;
// all other keys are passed through unchanged.
package truncate
