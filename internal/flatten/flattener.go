package flatten

import (
	"fmt"
	"sort"
	"strings"
)

// Options controls how flattening behaves.
type Options struct {
	// Separator is placed between nested key segments. Defaults to "_".
	Separator string
	// Prefix is prepended to every resulting key.
	Prefix string
	// UpperCase converts all keys to upper case.
	UpperCase bool
}

// Result holds the flattened key-value pairs and any warnings.
type Result struct {
	Vars     map[string]string
	Warnings []string
}

// Flatten collapses a nested map[string]any into a flat map[string]string.
// Nested maps are recursed; other non-string scalar values are converted via
// fmt.Sprintf. Slices and other complex types emit a warning and are skipped.
func Flatten(src map[string]any, opts Options) (*Result, error) {
	if src == nil {
		return &Result{Vars: map[string]string{}}, nil
	}
	if opts.Separator == "" {
		opts.Separator = "_"
	}
	res := &Result{Vars: make(map[string]string)}
	flatten(src, opts.Prefix, opts, res)
	return res, nil
}

func flatten(src map[string]any, prefix string, opts Options, res *Result) {
	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := src[k]
		full := join(prefix, k, opts.Separator)
		if opts.UpperCase {
			full = strings.ToUpper(full)
		}
		switch val := v.(type) {
		case map[string]any:
			flatten(val, full, opts, res)
		case string:
			res.Vars[full] = val
		case nil:
			res.Vars[full] = ""
		case bool, int, int64, float64:
			res.Vars[full] = fmt.Sprintf("%v", val)
		default:
			res.Warnings = append(res.Warnings, fmt.Sprintf("skipped %q: unsupported type %T", full, v))
		}
	}
}

func join(prefix, key, sep string) string {
	if prefix == "" {
		return key
	}
	return prefix + sep + key
}
