package omit

import (
	"fmt"
	"regexp"
)

// Options controls which keys are omitted from the source map.
type Options struct {
	// Keys is an explicit list of keys to remove.
	Keys []string
	// Pattern is a regex; any key matching it is removed.
	Pattern string
	// Prefix removes all keys that start with this string.
	Prefix string
}

// Omit returns a new map with the specified keys removed.
// The source map is never mutated.
func Omit(src map[string]string, opts Options) (map[string]string, error) {
	if src == nil {
		return map[string]string{}, nil
	}

	var re *regexp.Regexp
	if opts.Pattern != "" {
		var err error
		re, err = regexp.Compile(opts.Pattern)
		if err != nil {
			return nil, fmt.Errorf("omit: invalid pattern %q: %w", opts.Pattern, err)
		}
	}

	explicit := toSet(opts.Keys)

	out := make(map[string]string, len(src))
	for k, v := range src {
		if explicit[k] {
			continue
		}
		if re != nil && re.MatchString(k) {
			continue
		}
		if opts.Prefix != "" && len(k) >= len(opts.Prefix) && k[:len(opts.Prefix)] == opts.Prefix {
			continue
		}
		out[k] = v
	}
	return out, nil
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
