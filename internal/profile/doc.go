// Package profile provides persistent named profiles for envchain.
//
// A Profile captures a context name and export format so that users
// can switch between common configurations (e.g. dev/dotenv, prod/export)
// without repeating flags on every invocation.
//
// Profiles are stored as JSON files under a configurable directory
// (default: .envchain/profiles). Each profile is saved as <name>.json.
//
// Example usage:
//
//	store, err := profile.NewStore(".envchain/profiles")
//	if err != nil { ... }
//
//	err = store.Save(profile.Profile{
//		Name:    "dev",
//		Context: "development",
//		Format:  "dotenv",
//	})
//
//	p, err := store.Load("dev")
package profile
