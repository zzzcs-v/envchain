// Package scope provides named groupings of environment variables that can be
// applied as boundaries within envchain configurations. A scope acts as a
// labelled container — for example "prod-secrets" or "shared-infra" — allowing
// operators to reason about which variables belong to which logical boundary
// independently of context inheritance.
//
// Scopes are stored as JSON files on disk and can be listed, loaded, saved, and
// deleted via the Store type.
package scope
