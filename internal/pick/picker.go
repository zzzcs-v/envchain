// Package pick provides functionality to select a subset of keys from an env map.
package pick

import (
	"fmt"
	"regexp"
)

// Options controls how keys are selected.
type Options struct {
	// Keys is an explicit list of keys to include.
	Keys []string
	// Pattern is a regex pattern; matching keys are included.
	Pattern string
}

// Pick returns a new map containing only the selected keys from src.
// If both Keys and Pattern are empty, an error is returned.
func Pick(src map[string]string, opts Options) (map[string]string, error) {
	if len(opts.Keys) == 0 && opts.Pattern == "" {
		return nil, fmt.Errorf("pick: at least one of Keys or Pattern must be specified")
	}

	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, fmt.Errorf("pick: invalid pattern %q: %w", opts.Pattern, err)
		}
	}

	explicit := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		explicit[k] = true
	}

	out := make(map[string]string)
	for k, v := range src {
		if explicit[k] {
			out[k] = v
			continue
		}
		if re != nil && re.MatchString(k) {
			out[k] = v
		}
	}
	return out, nil
}
