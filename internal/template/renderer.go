package template

import (
	"fmt"
	"regexp"
	"strings"
)

// placeholderRe matches {{ VAR_NAME }} style template expressions.
var placeholderRe = regexp.MustCompile(`\{\{\s*([A-Z_][A-Z0-9_]*)\s*\}\}`)

// RenderResult holds the output of a template render operation.
type RenderResult struct {
	Output   string
	Missing  []string
}

// Renderer expands template placeholders using a provided env map.
type Renderer struct {
	strict bool
}

// NewRenderer creates a Renderer. When strict is true, missing vars are errors.
func NewRenderer(strict bool) *Renderer {
	return &Renderer{strict: strict}
}

// Render replaces all {{ VAR }} placeholders in tmpl with values from env.
// Returns a RenderResult and an error if strict mode is enabled and vars are missing.
func (r *Renderer) Render(tmpl string, env map[string]string) (RenderResult, error) {
	var missing []string
	seen := map[string]bool{}

	output := placeholderRe.ReplaceAllStringFunc(tmpl, func(match string) string {
		inner := placeholderRe.FindStringSubmatch(match)
		if len(inner) < 2 {
			return match
		}
		key := inner[1]
		if val, ok := env[key]; ok {
			return val
		}
		if !seen[key] {
			missing = append(missing, key)
			seen[key] = true
		}
		return match
	})

	result := RenderResult{Output: output, Missing: missing}

	if r.strict && len(missing) > 0 {
		return result, fmt.Errorf("template: unresolved placeholders: %s", strings.Join(missing, ", "))
	}

	return result, nil
}
