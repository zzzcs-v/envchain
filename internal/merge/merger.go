package merge

import (
	"fmt"
	"sort"
)

// MergeStrategy defines how conflicting keys are handled.
type MergeStrategy int

const (
	// StrategyOverwrite replaces existing values with incoming ones.
	StrategyOverwrite MergeStrategy = iota
	// StrategyKeepExisting preserves existing values on conflict.
	StrategyKeepExisting
	// StrategyError returns an error on any key conflict.
	StrategyError
)

// Merger combines multiple env var maps into one.
type Merger struct {
	strategy MergeStrategy
}

// NewMerger creates a Merger with the given strategy.
func NewMerger(strategy MergeStrategy) *Merger {
	return &Merger{strategy: strategy}
}

// Merge combines layers of env maps in order (later layers take precedence
// for StrategyOverwrite). Returns the merged map or an error.
func (m *Merger) Merge(layers ...map[string]string) (map[string]string, error) {
	result := make(map[string]string)

	for _, layer := range layers {
		// iterate keys in deterministic order
		keys := make([]string, 0, len(layer))
		for k := range layer {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := layer[k]
			if existing, exists := result[k]; exists {
				switch m.strategy {
				case StrategyOverwrite:
					result[k] = v
				case StrategyKeepExisting:
					_ = existing // keep what we have
				case StrategyError:
					return nil, fmt.Errorf("merge conflict: key %q already defined with value %q", k, existing)
				}
			} else {
				result[k] = v
			}
		}
	}

	return result, nil
}
