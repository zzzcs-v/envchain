package trim

import (
	"fmt"
	"strings"
)

// Options controls trimming behaviour.
type Options struct {
	// Keys restricts trimming to these specific keys. If empty, all keys are trimmed.
	Keys []string
	// Prefix removes a leading prefix from values when present.
	Prefix string
	// Suffix removes a trailing suffix from values when present.
	Suffix string
	// Whitespace trims leading/trailing whitespace from every value.
	Whitespace bool
}

// Result holds the outcome of a trim operation.
type Result struct {
	Key      string
	Original string
	Trimmed  string
}

// Changed reports whether the value was actually modified.
func (r Result) Changed() bool { return r.Original != r.Trimmed }

// String returns a human-readable summary of the result.
func (r Result) String() string {
	if r.Changed() {
		return fmt.Sprintf("%s: %q -> %q", r.Key, r.Original, r.Trimmed)
	}
	return fmt.Sprintf("%s: unchanged", r.Key)
}

// Trim applies the given options to src and returns a new map plus per-key results.
// src is never mutated.
func Trim(src map[string]string, opts Options) (map[string]string, []Result, error) {
	if src == nil {
		return map[string]string{}, nil, nil
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	out := make(map[string]string, len(src))
	var results []Result

	for k, v := range src {
		trimmed := v

		if len(keySet) == 0 || func() bool { _, ok := keySet[k]; return ok }() {
			if opts.Whitespace {
				trimmed = strings.TrimSpace(trimmed)
			}
			if opts.Prefix != "" {
				trimmed = strings.TrimPrefix(trimmed, opts.Prefix)
			}
			if opts.Suffix != "" {
				trimmed = strings.TrimSuffix(trimmed, opts.Suffix)
			}
		}

		out[k] = trimmed
		results = append(results, Result{Key: k, Original: v, Trimmed: trimmed})
	}

	return out, results, nil
}
