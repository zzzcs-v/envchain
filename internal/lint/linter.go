package lint

import (
	"fmt"
	"strings"
)

// Issue represents a single lint finding.
type Issue struct {
	Context string
	Key     string
	Message string
	Severity string // "error" | "warning"
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s.%s: %s", strings.ToUpper(i.Severity), i.Context, i.Key, i.Message)
}

// Result holds all issues found during a lint run.
type Result struct {
	Issues []Issue
}

func (r *Result) HasErrors() bool {
	for _, iss := range r.Issues {
		if iss.Severity == "error" {
			return true
		}
	}
	return false
}

func (r *Result) Summary() string {
	if len(r.Issues) == 0 {
		return "no issues found"
	}
	return fmt.Sprintf("%d issue(s) found", len(r.Issues))
}

// Run lints a map of context names to their key/value env vars.
// It checks for:
//   - keys that are not uppercase
//   - values that look like unresolved placeholders (${...})
//   - empty context names
func Run(contexts map[string]map[string]string) *Result {
	res := &Result{}

	for ctx, vars := range contexts {
		if strings.TrimSpace(ctx) == "" {
			res.Issues = append(res.Issues, Issue{
				Context:  "(unknown)",
				Key:      "-",
				Message:  "context name is empty or blank",
				Severity: "error",
			})
			continue
		}

		for k, v := range vars {
			if k != strings.ToUpper(k) {
				res.Issues = append(res.Issues, Issue{
					Context:  ctx,
					Key:      k,
					Message:  "key should be uppercase",
					Severity: "warning",
				})
			}

			if strings.Contains(v, "${") && strings.Contains(v, "}") {
				res.Issues = append(res.Issues, Issue{
					Context:  ctx,
					Key:      k,
					Message:  fmt.Sprintf("value appears to contain an unresolved placeholder: %q", v),
					Severity: "warning",
				})
			}
		}
	}

	return res
}
