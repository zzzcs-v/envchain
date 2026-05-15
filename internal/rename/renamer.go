package rename

import (
	"errors"
	"fmt"
	"regexp"
)

// Options configures rename behaviour.
type Options struct {
	// FromPattern is a regex applied to each key.
	FromPattern string
	// ToTemplate is the replacement string (supports $1, $2 capture groups).
	ToTemplate string
	// SkipConflicts silently drops a rename when the target key already exists.
	SkipConflicts bool
	// ErrorOnConflict returns an error when the target key already exists.
	ErrorOnConflict bool
}

// Result describes the outcome of a single key rename.
type Result struct {
	OldKey string
	NewKey string
	Skipped bool
}

// Rename applies regex-based key renaming to src and returns a new map plus
// a slice of Result describing every key that was considered.
func Rename(src map[string]string, opts Options) (map[string]string, []Result, error) {
	if src == nil {
		return map[string]string{}, nil, nil
	}
	if opts.FromPattern == "" {
		return nil, nil, errors.New("rename: FromPattern must not be empty")
	}
	re, err := regexp.Compile(opts.FromPattern)
	if err != nil {
		return nil, nil, fmt.Errorf("rename: invalid pattern: %w", err)
	}

	dst := make(map[string]string, len(src))
	var results []Result

	for k, v := range src {
		if !re.MatchString(k) {
			dst[k] = v
			continue
		}
		newKey := re.ReplaceAllString(k, opts.ToTemplate)
		if newKey == k {
			dst[k] = v
			continue
		}
		if _, exists := dst[newKey]; exists {
			if opts.ErrorOnConflict {
				return nil, nil, fmt.Errorf("rename: conflict: target key %q already exists", newKey)
			}
			if opts.SkipConflicts {
				results = append(results, Result{OldKey: k, NewKey: newKey, Skipped: true})
				dst[k] = v
				continue
			}
		}
		dst[newKey] = v
		results = append(results, Result{OldKey: k, NewKey: newKey})
	}
	return dst, results, nil
}
