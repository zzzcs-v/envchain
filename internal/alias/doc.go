// Package alias provides a persistent store for short-name aliases that
// map user-defined labels to fully qualified context names used throughout
// envchain. Aliases are stored as individual JSON files in a configurable
// directory, making them easy to inspect and version-control if desired.
//
// Example usage:
//
//	store, _ := alias.NewStore("~/.envchain/aliases")
//	store.Set("prod", "production-us-east")
//	a, _ := store.Get("prod")
//	fmt.Println(a.Context) // production-us-east
package alias
