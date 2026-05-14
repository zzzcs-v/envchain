package filter

import (
	"regexp"
	"strings"
)

// Options controls how filtering is applied.
type Options struct {
	// KeyPrefix filters vars whose key starts with the given prefix.
	KeyPrefix string
	// KeyPattern filters vars whose key matches the given regex.
	KeyPattern string
	// ExcludeKeys is a set of exact key names to exclude.
	ExcludeKeys []string
	// InvertMatch returns vars that do NOT match the criteria.
	InvertMatch bool
}

// Filter applies the given options to a map of env vars and returns a filtered copy.
func Filter(vars map[string]string, opts Options) (map[string]string, error) {
	var re *regexp.Regexp
	if opts.KeyPattern != "" {
		var err error
		re, err = regexp.Compile(opts.KeyPattern)
		if err != nil {
			return nil, err
		}
	}

	excluded := make(map[string]bool, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excluded[k] = true
	}

	result := make(map[string]string)
	for k, v := range vars {
		matched := matchesKey(k, opts.KeyPrefix, re) && !excluded[k]
		if opts.InvertMatch {
			matched = !matched
		}
		if matched {
			result[k] = v
		}
	}
	return result, nil
}

func matchesKey(key, prefix string, re *regexp.Regexp) bool {
	if prefix != "" && !strings.HasPrefix(key, prefix) {
		return false
	}
	if re != nil && !re.MatchString(key) {
		return false
	}
	return true
}
