// Package inject provides utilities for injecting environment variables
// into a target map from one or more named sources.
package inject

import (
	"errors"
	"fmt"
)

// Source represents a named set of key-value pairs.
type Source struct {
	Name string
	Vars map[string]string
}

// Options controls injection behaviour.
type Options struct {
	// Overwrite allows injected values to replace existing keys.
	Overwrite bool
	// Prefix is prepended to every injected key.
	Prefix string
}

// Result holds the outcome of a single injection.
type Result struct {
	Source   string
	Injected int
	Skipped  int
}

// Injector applies sources into a destination map.
type Injector struct {
	opts Options
}

// New returns an Injector configured with opts.
func New(opts Options) *Injector {
	return &Injector{opts: opts}
}

// Inject merges all sources into dst in order.
// Returns one Result per source and the first error encountered, if any.
func (inj *Injector) Inject(dst map[string]string, sources []Source) ([]Result, error) {
	if dst == nil {
		return nil, errors.New("inject: destination map must not be nil")
	}
	results := make([]Result, 0, len(sources))
	for _, src := range sources {
		if src.Name == "" {
			return results, errors.New("inject: source name must not be empty")
		}
		r := Result{Source: src.Name}
		for k, v := range src.Vars {
			key := inj.opts.Prefix + k
			if key == "" {
				return results, fmt.Errorf("inject: empty key in source %q", src.Name)
			}
			if _, exists := dst[key]; exists && !inj.opts.Overwrite {
				r.Skipped++
				continue
			}
			dst[key] = v
			r.Injected++
		}
		results = append(results, r)
	}
	return results, nil
}
