// Package override applies key-level overrides to an env map from
// multiple sources such as flags, files, or inline pairs.
package override

import (
	"errors"
	"fmt"
	"strings"
)

// Entry represents a single key=value override pair.
type Entry struct {
	Key   string
	Value string
}

// Options controls how overrides are applied.
type Options struct {
	// AllowEmpty permits overriding a key with an empty string.
	AllowEmpty bool
	// Strict returns an error if an override key does not already exist in dst.
	Strict bool
}

// ParsePair parses a "KEY=VALUE" string into an Entry.
func ParsePair(s string) (Entry, error) {
	idx := strings.IndexByte(s, '=')
	if idx < 1 {
		return Entry{}, fmt.Errorf("override: invalid pair %q: must be KEY=VALUE", s)
	}
	return Entry{Key: s[:idx], Value: s[idx+1:]}, nil
}

// Apply merges entries into dst according to opts.
// dst must not be nil.
func Apply(dst map[string]string, entries []Entry, opts Options) error {
	if dst == nil {
		return errors.New("override: dst map must not be nil")
	}
	for _, e := range entries {
		if e.Key == "" {
			return errors.New("override: entry key must not be empty")
		}
		if !opts.AllowEmpty && e.Value == "" {
			return fmt.Errorf("override: empty value for key %q (use AllowEmpty to permit)", e.Key)
		}
		if opts.Strict {
			if _, ok := dst[e.Key]; !ok {
				return fmt.Errorf("override: key %q does not exist in destination (strict mode)", e.Key)
			}
		}
		dst[e.Key] = e.Value
	}
	return nil
}
