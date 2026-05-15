package normalize

import (
	"fmt"
	"regexp"
	"strings"
)

// Options controls normalization behavior.
type Options struct {
	// UpperCase converts all keys to uppercase.
	UpperCase bool
	// LowerCase converts all keys to lowercase (mutually exclusive with UpperCase).
	LowerCase bool
	// ReplacePattern is a regex matched against key characters to replace with ReplaceWith.
	ReplacePattern string
	// ReplaceWith is the replacement string for ReplacePattern matches.
	ReplaceWith string
	// StripNonAlnum removes any character that is not alphanumeric or underscore.
	StripNonAlnum bool
	// RestrictToKeys limits normalization to only these keys (nil = all keys).
	RestrictToKeys []string
}

// Normalize applies key normalization to src and returns a new map.
// Values are never modified. Keys are transformed according to opts.
// If two keys collide after normalization, the last one (in sorted iteration
// order) wins and an error is returned listing all collisions.
func Normalize(src map[string]string, opts Options) (map[string]string, error) {
	if src == nil {
		return map[string]string{}, nil
	}

	var re *regexp.Regexp
	if opts.ReplacePattern != "" {
		var err error
		re, err = regexp.Compile(opts.ReplacePattern)
		if err != nil {
			return nil, fmt.Errorf("normalize: invalid replace pattern: %w", err)
		}
	}

	restrict := toSet(opts.RestrictToKeys)

	result := make(map[string]string, len(src))
	collisions := map[string][]string{}

	for k, v := range src {
		newKey := k
		if len(restrict) == 0 || restrict[k] {
			newKey = transformKey(k, opts, re)
		}
		if existing, ok := result[newKey]; ok {
			collisions[newKey] = append(collisions[newKey], existing)
		}
		result[newKey] = v
	}

	if len(collisions) > 0 {
		keys := make([]string, 0, len(collisions))
		for k := range collisions {
			keys = append(keys, k)
		}
		return result, fmt.Errorf("normalize: key collisions after normalization: %s", strings.Join(keys, ", "))
	}

	return result, nil
}

func transformKey(k string, opts Options, re *regexp.Regexp) string {
	if opts.StripNonAlnum {
		k = regexp.MustCompile(`[^A-Za-z0-9_]`).ReplaceAllString(k, "")
	}
	if re != nil {
		k = re.ReplaceAllString(k, opts.ReplaceWith)
	}
	if opts.UpperCase {
		return strings.ToUpper(k)
	}
	if opts.LowerCase {
		return strings.ToLower(k)
	}
	return k
}

func toSet(keys []string) map[string]bool {
	if len(keys) == 0 {
		return nil
	}
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
