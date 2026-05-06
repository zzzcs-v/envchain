package validate

import (
	"fmt"
	"regexp"
	"strings"
)

var envKeyRegex = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// Result holds validation findings for a set of env vars.
type Result struct {
	Warnings []string
	Errors   []string
}

// OK returns true when there are no errors.
func (r Result) OK() bool { return len(r.Errors) == 0 }

// Summary returns a human-readable summary string.
func (r Result) Summary() string {
	var b strings.Builder
	for _, e := range r.Errors {
		fmt.Fprintf(&b, "ERROR: %s\n", e)
	}
	for _, w := range r.Warnings {
		fmt.Fprintf(&b, "WARN:  %s\n", w)
	}
	return b.String()
}

// Vars validates a map of environment variable key/value pairs.
// It checks for naming conventions and flags empty values as warnings.
func Vars(vars map[string]string) Result {
	var res Result
	for k, v := range vars {
		if !envKeyRegex.MatchString(k) {
			res.Errors = append(res.Errors,
				fmt.Sprintf("key %q does not match expected pattern [A-Z_][A-Z0-9_]*", k))
		}
		if strings.TrimSpace(v) == "" {
			res.Warnings = append(res.Warnings,
				fmt.Sprintf("key %q has an empty value", k))
		}
	}
	return res
}
