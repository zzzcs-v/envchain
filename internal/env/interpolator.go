package env

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// placeholder pattern matches ${VAR_NAME} or ${VAR_NAME:-default}
var placeholderRe = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)(?::-(.*?))?\}`)

// Interpolator resolves variable references within env values.
type Interpolator struct {
	vars    map[string]string
	osLookup func(string) (string, bool)
}

// New creates an Interpolator seeded with the provided variable map.
func New(vars map[string]string) *Interpolator {
	return &Interpolator{
		vars:    vars,
		osLookup: os.LookupEnv,
	}
}

// Resolve expands all ${VAR} and ${VAR:-default} references in the given
// string. Resolution order: vars map → OS environment → inline default.
// Returns an error if a reference cannot be resolved and has no default.
func (i *Interpolator) Resolve(s string) (string, error) {
	var resolveErr error
	result := placeholderRe.ReplaceAllStringFunc(s, func(match string) string {
		if resolveErr != nil {
			return match
		}
		parts := placeholderRe.FindStringSubmatch(match)
		key := parts[1]
		defaultVal := parts[2]
		hasDefault := strings.Contains(match, ":-")

		if v, ok := i.vars[key]; ok {
			return v
		}
		if v, ok := i.osLookup(key); ok {
			return v
		}
		if hasDefault {
			return defaultVal
		}
		resolveErr = fmt.Errorf("unresolved variable: %s", key)
		return match
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}

// ResolveAll applies Resolve to every value in the map, returning a new map.
func (i *Interpolator) ResolveAll(vars map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(vars))
	for k, v := range vars {
		resolved, err := i.Resolve(v)
		if err != nil {
			return nil, fmt.Errorf("key %s: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}
