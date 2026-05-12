package clone

import (
	"errors"
	"fmt"
)

// Entry represents a single environment variable within a context.
type Entry struct {
	Key   string
	Value string
}

// Context is a named collection of environment entries.
type Context struct {
	Name    string
	Entries []Entry
}

// Cloner copies one named context into a new name, with optional key filtering.
type Cloner struct {
	source map[string][]Entry
}

// New creates a Cloner from a map of context name -> entries.
func New(source map[string][]Entry) (*Cloner, error) {
	if source == nil {
		return nil, errors.New("clone: source map must not be nil")
	}
	return &Cloner{source: source}, nil
}

// Clone copies the context identified by srcName into a new Context named dstName.
// If filterKeys is non-empty, only those keys are included in the clone.
// Returns an error if srcName does not exist or dstName is empty.
func (c *Cloner) Clone(srcName, dstName string, filterKeys []string) (*Context, error) {
	if dstName == "" {
		return nil, errors.New("clone: destination name must not be empty")
	}
	entries, ok := c.source[srcName]
	if !ok {
		return nil, fmt.Errorf("clone: source context %q not found", srcName)
	}

	filter := make(map[string]bool, len(filterKeys))
	for _, k := range filterKeys {
		filter[k] = true
	}

	var cloned []Entry
	for _, e := range entries {
		if len(filter) == 0 || filter[e.Key] {
			cloned = append(cloned, Entry{Key: e.Key, Value: e.Value})
		}
	}

	return &Context{
		Name:    dstName,
		Entries: cloned,
	}, nil
}
