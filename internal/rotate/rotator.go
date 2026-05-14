package rotate

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

// Entry represents a single rotation record for a key.
type Entry struct {
	Key       string    `json:"key"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	RotatedAt time.Time `json:"rotated_at"`
	Context   string    `json:"context"`
}

// Result summarises the outcome of a rotation run.
type Result struct {
	Rotated []Entry
	Skipped []string
	Errors  []error
}

func (r Result) Summary() string {
	return fmt.Sprintf("rotated=%d skipped=%d errors=%d", len(r.Rotated), len(r.Skipped), len(r.Errors))
}

// GeneratorFunc produces a new value for the given key.
type GeneratorFunc func(key string) (string, error)

// Options controls rotation behaviour.
type Options struct {
	Context   string
	Keys      []string // if empty, rotate all keys
	Overwrite bool
	Generator GeneratorFunc
}

// Rotate applies key rotation to the provided vars map according to opts.
// It returns a Result describing what happened.
func Rotate(vars map[string]string, opts Options) (Result, error) {
	if opts.Generator == nil {
		return Result{}, errors.New("rotate: generator func is required")
	}

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range vars {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}

	var res Result
	for _, k := range keys {
		oldVal, exists := vars[k]
		if !exists {
			res.Skipped = append(res.Skipped, k)
			continue
		}
		if !opts.Overwrite && oldVal != "" {
			res.Skipped = append(res.Skipped, k)
			continue
		}
		newVal, err := opts.Generator(k)
		if err != nil {
			res.Errors = append(res.Errors, fmt.Errorf("rotate: key %q: %w", k, err))
			continue
		}
		vars[k] = newVal
		res.Rotated = append(res.Rotated, Entry{
			Key:       k,
			OldValue:  oldVal,
			NewValue:  newVal,
			RotatedAt: time.Now().UTC(),
			Context:   opts.Context,
		})
	}
	return res, nil
}
