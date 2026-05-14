package compare

import (
	"fmt"
	"sort"
)

// Result holds the comparison outcome between two env maps.
type Result struct {
	OnlyInLeft  map[string]string
	OnlyInRight map[string]string
	Different   map[string][2]string // key -> [left, right]
	Shared      map[string]string
}

// Summary returns a human-readable summary of the diff.
func (r *Result) Summary() string {
	return fmt.Sprintf("+%d added, -%d removed, ~%d changed, %d shared",
		len(r.OnlyInRight), len(r.OnlyInLeft), len(r.Different), len(r.Shared))
}

// Keys returns all keys that differ between left and right, sorted.
func (r *Result) Keys() []string {
	seen := map[string]struct{}{}
	for k := range r.OnlyInLeft {
		seen[k] = struct{}{}
	}
	for k := range r.OnlyInRight {
		seen[k] = struct{}{}
	}
	for k := range r.Different {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Compare compares two environment variable maps and returns a Result.
func Compare(left, right map[string]string) *Result {
	res := &Result{
		OnlyInLeft:  make(map[string]string),
		OnlyInRight: make(map[string]string),
		Different:   make(map[string][2]string),
		Shared:      make(map[string]string),
	}

	for k, lv := range left {
		if rv, ok := right[k]; ok {
			if lv == rv {
				res.Shared[k] = lv
			} else {
				res.Different[k] = [2]string{lv, rv}
			}
		} else {
			res.OnlyInLeft[k] = lv
		}
	}

	for k, rv := range right {
		if _, ok := left[k]; !ok {
			res.OnlyInRight[k] = rv
		}
	}

	return res
}
