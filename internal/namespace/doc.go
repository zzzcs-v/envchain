// Package namespace provides a store for managing named namespace entries.
//
// A namespace groups a set of environment contexts under a shared prefix,
// making it easy to scope and isolate variables across services or teams.
//
// Example usage:
//
//	store, err := namespace.NewStore("/home/user/.envchain/namespaces")
//	if err != nil { ... }
//
//	err = store.Save(namespace.Entry{
//		Name:     "payments",
//		Prefix:   "PAY_",
//		Contexts: []string{"dev", "staging", "prod"},
//	})
//
//	entries, err := store.List()
package namespace
