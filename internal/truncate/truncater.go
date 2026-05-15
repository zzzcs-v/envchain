package truncate

import (
	"errors"
	"strings"
)

// Options controls how truncation is applied.
type Options struct {
	// MaxLen is the maximum length of a value (in runes). Zero means no limit.
	MaxLen int
	// Suffix is appended when a value is truncated (e.g. "...").
	Suffix string
	// Keys restricts truncation to the given set of keys. Empty means all keys.
	Keys []string
}

// Result holds a single truncation event.
type Result struct {
	Key      string
	Original string
	Truncated string
}

// Truncate applies length truncation to values in src according to opts.
// It returns a new map and a slice of Result describing what was changed.
func Truncate(src map[string]string, opts Options) (map[string]string, []Result, error) {
	if src == nil {
		return map[string]string{}, nil, nil
	}
	if opts.MaxLen < 0 {
		return nil, nil, errors.New("truncate: MaxLen must be >= 0")
	}

	keySet := toSet(opts.Keys)
	out := make(map[string]string, len(src))
	var results []Result

	for k, v := range src {
		if opts.MaxLen == 0 || (!keySet[k] && len(keySet) > 0) {
			out[k] = v
			continue
		}
		runes := []rune(v)
		if len(runes) <= opts.MaxLen {
			out[k] = v
			continue
		}
		truncated := string(runes[:opts.MaxLen])
		if opts.Suffix != "" {
			truncated = strings.TrimRight(truncated, " ") + opts.Suffix
		}
		out[k] = truncated
		results = append(results, Result{Key: k, Original: v, Truncated: truncated})
	}
	return out, results, nil
}

func toSet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
