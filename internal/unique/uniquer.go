package unique

import "strings"

// Options controls how uniqueness is determined.
type Options struct {
	// Keys restricts deduplication to specific keys.
	Keys []string
	// CaseSensitive controls whether values are compared case-sensitively.
	CaseSensitive bool
}

// Result holds the output of a Unique operation.
type Result struct {
	Kept    map[string]string
	Dropped map[string]string
}

// Unique removes entries whose values are duplicated across the map.
// When multiple keys share the same value, the first key (alphabetically) is kept.
func Unique(src map[string]string, opts Options) (Result, error) {
	if src == nil {
		return Result{
			Kept:    map[string]string{},
			Dropped: map[string]string{},
		}, nil
	}

	restrict := toSet(opts.Keys)

	// Collect keys in sorted order for determinism.
	keys := sortedKeys(src)

	seen := map[string]string{} // normalised value -> first key that claimed it
	kept := map[string]string{}
	dropped := map[string]string{}

	for _, k := range keys {
		if len(restrict) > 0 && !restrict[k] {
			kept[k] = src[k]
			continue
		}

		v := src[k]
		norm := v
		if !opts.CaseSensitive {
			norm = strings.ToLower(v)
		}

		if _, exists := seen[norm]; exists {
			dropped[k] = v
		} else {
			seen[norm] = k
			kept[k] = v
		}
	}

	return Result{Kept: kept, Dropped: dropped}, nil
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}

func sortedKeys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	// simple insertion sort — maps are small in practice
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}
