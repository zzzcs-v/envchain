package schema

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for an environment variable key.
type Rule struct {
	Pattern     string
	Required    bool
	Description string
}

// Schema holds a set of named rules.
type Schema struct {
	rules map[string]Rule
}

// Violation represents a single schema violation.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) String() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// NewSchema creates a Schema from a map of key -> Rule.
func NewSchema(rules map[string]Rule) (*Schema, error) {
	for key, rule := range rules {
		if rule.Pattern != "" {
			if _, err := regexp.Compile(rule.Pattern); err != nil {
				return nil, fmt.Errorf("invalid pattern for key %q: %w", key, err)
			}
		}
	}
	return &Schema{rules: rules}, nil
}

// Validate checks the given vars map against the schema rules.
// Returns a slice of violations (empty means valid).
func (s *Schema) Validate(vars map[string]string) []Violation {
	var violations []Violation

	for key, rule := range s.rules {
		val, exists := vars[key]
		if rule.Required && !exists {
			violations = append(violations, Violation{Key: key, Message: "required key is missing"})
			continue
		}
		if !exists {
			continue
		}
		if strings.TrimSpace(val) == "" {
			violations = append(violations, Violation{Key: key, Message: "value must not be empty"})
			continue
		}
		if rule.Pattern != "" {
			re := regexp.MustCompile(rule.Pattern)
			if !re.MatchString(val) {
				violations = append(violations, Violation{
					Key:     key,
					Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern),
				})
			}
		}
	}

	return violations
}
