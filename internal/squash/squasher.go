// Package squash merges multiple env maps into one, resolving conflicts
// by applying a configurable strategy (keep-first, keep-last, error).
package squash

import (
	"errors"
	"fmt"
	"sort"
)

// Strategy controls how key conflicts are resolved during squash.
type Strategy int

const (
	KeepFirst Strategy = iota // retain the value from the earliest source
	KeepLast                  // overwrite with the value from the latest source
	ErrorOnConflict           // return an error when a duplicate key is found
)

// Options configures the Squash operation.
type Options struct {
	Strategy Strategy
}

// Result holds the squashed map and any keys that were involved in conflicts.
type Result struct {
	Vars      map[string]string
	Conflicts []string
}

// Squash merges the provided sources left-to-right according to opts.
// Sources are applied in order; index 0 is considered the "first" source.
func Squash(sources []map[string]string, opts Options) (*Result, error) {
	out := make(map[string]string)
	conflictSet := make(map[string]struct{})

	for _, src := range sources {
		if src == nil {
			continue
		}
		for k, v := range src {
			if _, exists := out[k]; exists {
				switch opts.Strategy {
				case ErrorOnConflict:
					return nil, fmt.Errorf("squash: duplicate key %q", k)
				case KeepLast:
					out[k] = v
					conflictSet[k] = struct{}{}
				case KeepFirst:
					conflictSet[k] = struct{}{}
					// keep existing value
				}
			} else {
				out[k] = v
			}
		}
	}

	conflicts := make([]string, 0, len(conflictSet))
	for k := range conflictSet {
		conflicts = append(conflicts, k)
	}
	sort.Strings(conflicts)

	return &Result{Vars: out, Conflicts: conflicts}, nil
}

// ParseStrategy converts a string label to a Strategy constant.
func ParseStrategy(s string) (Strategy, error) {
	switch s {
	case "keep-first":
		return KeepFirst, nil
	case "keep-last":
		return KeepLast, nil
	case "error":
		return ErrorOnConflict, nil
	}
	return 0, errors.New("squash: unknown strategy " + s)
}
