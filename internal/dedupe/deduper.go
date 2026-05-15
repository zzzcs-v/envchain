// Package dedupe provides utilities for removing duplicate environment variable
// entries across maps, preferring the first or last occurrence based on options.
package dedupe

import (
	"fmt"
	"sort"
)

// Strategy controls which value wins when a duplicate key is found.
type Strategy int

const (
	// KeepFirst retains the first occurrence of a duplicate key.
	KeepFirst Strategy = iota
	// KeepLast retains the last occurrence of a duplicate key.
	KeepLast
)

// Options configures deduplication behaviour.
type Options struct {
	Strategy Strategy
	// RestrictToKeys, if non-empty, only deduplicates the listed keys.
	RestrictToKeys []string
}

// Result holds the deduplicated map and a summary of removed entries.
type Result struct {
	Vars    map[string]string
	Removed []string
}

// Summary returns a human-readable description of the result.
func (r Result) Summary() string {
	if len(r.Removed) == 0 {
		return "no duplicates found"
	}
	return fmt.Sprintf("%d duplicate(s) removed: %v", len(r.Removed), r.Removed)
}

// Dedupe removes duplicate keys from the provided list of maps, merging them
// into a single map according to the given Options.
func Dedupe(sources []map[string]string, opts Options) (Result, error) {
	if len(sources) == 0 {
		return Result{Vars: map[string]string{}}, nil
	}

	restrict := toSet(opts.RestrictToKeys)
	seen := map[string]bool{}
	merged := map[string]string{}
	var removed []string

	process := func(src map[string]string) {
		keys := make([]string, 0, len(src))
		for k := range src {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := src[k]
			if len(restrict) > 0 && !restrict[k] {
				merged[k] = v
				continue
			}
			if seen[k] {
				if opts.Strategy == KeepLast {
					merged[k] = v
				}
				removed = append(removed, k)
				continue
			}
			seen[k] = true
			merged[k] = v
		}
	}

	for _, src := range sources {
		process(src)
	}

	sort.Strings(removed)
	return Result{Vars: merged, Removed: removed}, nil
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
