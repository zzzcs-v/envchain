// Package sort provides utilities for sorting environment variable maps
// by key, value, or custom ordering strategies.
package sort

import (
	"fmt"
	gosort "sort"
	"strings"
)

// Order defines the sort direction.
type Order string

const (
	Asc  Order = "asc"
	Desc Order = "desc"
)

// Options configures how sorting is applied.
type Options struct {
	By    string // "key" or "value"
	Order Order
	Keys  []string // if set, only sort these keys (others appended after)
}

// Result holds a sorted list of key-value pairs.
type Result struct {
	Pairs []Pair
}

// Pair represents a single key-value entry.
type Pair struct {
	Key   string
	Value string
}

// Sort takes a map and returns a sorted Result according to Options.
func Sort(src map[string]string, opts Options) (*Result, error) {
	if src == nil {
		return &Result{Pairs: []Pair{}}, nil
	}

	if opts.By == "" {
		opts.By = "key"
	}
	if opts.Order == "" {
		opts.Order = Asc
	}

	if opts.By != "key" && opts.By != "value" {
		return nil, fmt.Errorf("invalid sort field %q: must be \"key\" or \"value\"", opts.By)
	}

	pairs := make([]Pair, 0, len(src))
	for k, v := range src {
		pairs = append(pairs, Pair{Key: k, Value: v})
	}

	gosort.Slice(pairs, func(i, j int) bool {
		var a, b string
		if opts.By == "value" {
			a = strings.ToLower(pairs[i].Value)
			b = strings.ToLower(pairs[j].Value)
		} else {
			a = strings.ToLower(pairs[i].Key)
			b = strings.ToLower(pairs[j].Key)
		}
		if opts.Order == Desc {
			return a > b
		}
		return a < b
	})

	return &Result{Pairs: pairs}, nil
}

// ToMap converts a Result back to a plain map.
func (r *Result) ToMap() map[string]string {
	m := make(map[string]string, len(r.Pairs))
	for _, p := range r.Pairs {
		m[p.Key] = p.Value
	}
	return m
}
