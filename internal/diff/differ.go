package diff

import (
	"fmt"
	"sort"
)

// ChangeType represents the kind of change between two contexts.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Change describes a single variable difference.
type Change struct {
	Key      string
	Type     ChangeType
	OldValue string
	NewValue string
}

// Result holds the full diff between two env maps.
type Result struct {
	Changes []Change
}

// HasChanges returns true if there are any non-unchanged entries.
func (r *Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Type != Unchanged {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary line.
func (r *Result) Summary() string {
	var added, removed, modified int
	for _, c := range r.Changes {
		switch c.Type {
		case Added:
			added++
		case Removed:
			removed++
		case Modified:
			modified++
		}
	}
	return fmt.Sprintf("+%d added, -%d removed, ~%d modified", added, removed, modified)
}

// Compare diffs two resolved env maps (from -> to).
func Compare(from, to map[string]string) *Result {
	keys := make(map[string]struct{})
	for k := range from {
		keys[k] = struct{}{}
	}
	for k := range to {
		keys[k] = struct{}{}
	}

	sorted := make([]string, 0, len(keys))
	for k := range keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	result := &Result{}
	for _, k := range sorted {
		oldVal, inFrom := from[k]
		newVal, inTo := to[k]

		switch {
		case inFrom && !inTo:
			result.Changes = append(result.Changes, Change{Key: k, Type: Removed, OldValue: oldVal})
		case !inFrom && inTo:
			result.Changes = append(result.Changes, Change{Key: k, Type: Added, NewValue: newVal})
		case oldVal != newVal:
			result.Changes = append(result.Changes, Change{Key: k, Type: Modified, OldValue: oldVal, NewValue: newVal})
		default:
			result.Changes = append(result.Changes, Change{Key: k, Type: Unchanged, OldValue: oldVal, NewValue: newVal})
		}
	}
	return result
}
