// Package freeze provides functionality to lock environment variable maps,
// preventing further modification by tracking a frozen snapshot.
package freeze

import (
	"errors"
	"fmt"
	"sort"
)

// ErrFrozen is returned when a mutation is attempted on a frozen map.
var ErrFrozen = errors.New("map is frozen and cannot be modified")

// Freezer holds an immutable snapshot of an env map.
type Freezer struct {
	original map[string]string
	frozen   map[string]string
}

// New creates a Freezer from the given source map.
// The source is deep-copied so external mutations don't affect the frozen state.
func New(source map[string]string) (*Freezer, error) {
	if source == nil {
		return nil, errors.New("source map must not be nil")
	}
	snap := make(map[string]string, len(source))
	for k, v := range source {
		snap[k] = v
	}
	return &Freezer{original: snap, frozen: snap}, nil
}

// Get returns the value for the given key from the frozen map.
func (f *Freezer) Get(key string) (string, bool) {
	v, ok := f.frozen[key]
	return v, ok
}

// Keys returns a sorted list of all keys in the frozen map.
func (f *Freezer) Keys() []string {
	keys := make([]string, 0, len(f.frozen))
	for k := range f.frozen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Snapshot returns a copy of the frozen map.
func (f *Freezer) Snapshot() map[string]string {
	copy := make(map[string]string, len(f.frozen))
	for k, v := range f.frozen {
		copy[k] = v
	}
	return copy
}

// DiffFrom compares the frozen map against a live map and returns changed keys.
func (f *Freezer) DiffFrom(live map[string]string) []string {
	var changed []string
	for k, frozenVal := range f.frozen {
		if liveVal, ok := live[k]; !ok || liveVal != frozenVal {
			changed = append(changed, k)
		}
	}
	for k := range live {
		if _, ok := f.frozen[k]; !ok {
			changed = append(changed, fmt.Sprintf("+%s", k))
		}
	}
	sort.Strings(changed)
	return changed
}
