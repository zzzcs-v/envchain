// Package defaults provides functionality for applying default values
// to environment variable maps when keys are missing or empty.
package defaults

import "fmt"

// Entry represents a single default value definition.
type Entry struct {
	Key      string
	Value    string
	Override bool // if true, replace even non-empty existing values
}

// Result holds the outcome of applying defaults.
type Result struct {
	Applied  []string
	Skipped  []string
	Overrode []string
}

// Summary returns a human-readable summary of the result.
func (r Result) Summary() string {
	return fmt.Sprintf("applied=%d skipped=%d overrode=%d",
		len(r.Applied), len(r.Skipped), len(r.Overrode))
}

// Apply merges default entries into dst, returning a Result.
// dst must not be nil. Entries with empty Key are skipped.
func Apply(dst map[string]string, entries []Entry) (Result, error) {
	if dst == nil {
		return Result{}, fmt.Errorf("defaults: dst map must not be nil")
	}

	var res Result
	for _, e := range entries {
		if e.Key == "" {
			continue
		}
		existing, exists := dst[e.Key]
		switch {
		case !exists || existing == "":
			dst[e.Key] = e.Value
			res.Applied = append(res.Applied, e.Key)
		case e.Override:
			dst[e.Key] = e.Value
			res.Overrode = append(res.Overrode, e.Key)
		default:
			res.Skipped = append(res.Skipped, e.Key)
		}
	}
	return res, nil
}
