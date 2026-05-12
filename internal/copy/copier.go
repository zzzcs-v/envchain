package copy

import (
	"fmt"
	"sort"
)

// Entry represents a single environment variable in a context.
type Entry struct {
	Key   string
	Value string
}

// Result holds the outcome of a copy operation.
type Result struct {
	Copied    []string
	Skipped   []string
	Overwrite bool
}

// Summary returns a human-readable summary of the copy result.
func (r Result) Summary() string {
	return fmt.Sprintf("copied %d vars, skipped %d", len(r.Copied), len(r.Skipped))
}

// Copier copies environment variables from one context to another.
type Copier struct {
	overwrite bool
}

// New creates a new Copier. If overwrite is true, existing keys in dst are replaced.
func New(overwrite bool) *Copier {
	return &Copier{overwrite: overwrite}
}

// Copy merges vars from src into dst, returning a Result describing what happened.
func (c *Copier) Copy(src, dst map[string]string) (map[string]string, Result, error) {
	if src == nil {
		return nil, Result{}, fmt.Errorf("source context is nil")
	}
	if dst == nil {
		dst = make(map[string]string)
	}

	out := make(map[string]string, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	var result Result
	result.Overwrite = c.overwrite

	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if _, exists := out[k]; exists && !c.overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		out[k] = src[k]
		result.Copied = append(result.Copied, k)
	}

	return out, result, nil
}
