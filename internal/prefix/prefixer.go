package prefix

import (
	"fmt"
	"strings"
)

// Options controls how keys are prefixed or stripped.
type Options struct {
	Prefix    string
	Strip     bool // if true, remove Prefix instead of adding it
	Overwrite bool // if true, re-prefix already-prefixed keys
}

// Result holds the outcome of a prefix operation.
type Result struct {
	Added   int
	Skipped int
	Stripped int
}

func (r Result) Summary() string {
	if r.Stripped > 0 {
		return fmt.Sprintf("stripped %d key(s), skipped %d", r.Stripped, r.Skipped)
	}
	return fmt.Sprintf("prefixed %d key(s), skipped %d", r.Added, r.Skipped)
}

// Apply adds or removes a prefix from every key in src and returns a new map.
func Apply(src map[string]string, opts Options) (map[string]string, Result, error) {
	if opts.Prefix == "" {
		return nil, Result{}, fmt.Errorf("prefix: prefix must not be empty")
	}
	if src == nil {
		return map[string]string{}, Result{}, nil
	}

	out := make(map[string]string, len(src))
	var res Result

	for k, v := range src {
		if opts.Strip {
			if strings.HasPrefix(k, opts.Prefix) {
				newKey := strings.TrimPrefix(k, opts.Prefix)
				out[newKey] = v
				res.Stripped++
			} else {
				out[k] = v
				res.Skipped++
			}
		} else {
			if strings.HasPrefix(k, opts.Prefix) && !opts.Overwrite {
				out[k] = v
				res.Skipped++
			} else {
				newKey := opts.Prefix + k
				if opts.Overwrite && strings.HasPrefix(k, opts.Prefix) {
					newKey = k
				}
				out[newKey] = v
				res.Added++
			}
		}
	}
	return out, res, nil
}
