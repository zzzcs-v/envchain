// Package defaults applies fallback values to environment variable maps.
//
// It is useful when loading configs that may have optional keys — callers
// define a list of Entry values specifying the key, default value, and
// whether existing non-empty values should be overridden.
//
// Example usage:
//
//	entries := []defaults.Entry{
//		{Key: "LOG_LEVEL", Value: "info"},
//		{Key: "PORT",      Value: "8080"},
//		{Key: "ENV",       Value: "production", Override: true},
//	}
//	res, err := defaults.Apply(envMap, entries)
//	fmt.Println(res.Summary())
package defaults
